package socket

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"meigens-api/db"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"meigens-api/src/auth"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var mtx sync.Mutex
var clients = treemap.NewWithStringComparator()

// user_id: hashset of conn
// user couuld have multiple connections so we need to store them in a slice.

var socketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		front_origin := os.Getenv("FRONT_ORIGIN")
		origin := r.Header.Get("Origin")
		return origin == front_origin
	},
}

func removeClientFromClients(user_id string, conn *websocket.Conn) {
	mtx.Lock()
	connections, not_found := clients.Get(user_id)
	if !not_found {
		connections.(*hashset.Set).Remove(conn)
	}
	if connections.(*hashset.Set).Size() == 0 {
		clients.Remove(user_id)
	}
	mtx.Unlock()
}

type connection_state struct {
	conn  *websocket.Conn
	state uint8
}

func TLSocket(c *gin.Context) {
	conn, err := socketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	mtx.Lock()
	user_id := c.Query("user_id")
	state := c.Query("state")
	connections, not_first := clients.Get(user_id)
	var state_num int64
	state_num, err = strconv.ParseInt(state, 10, 8)
	if err != nil {
		SendErrorMessage(conn, "Invalid state.")
		return
	}
	if state_num < 0 || state_num > 1 {
		SendErrorMessage(conn, "Invalid state.")
		return
	}
	// connections: nil or hashset of conn
	if !not_first {
		var connections_set = hashset.New()
		connections_set.Add(connection_state{conn, uint8(state_num)})
		clients.Put(user_id, connections_set)
	} else {
		connections.(*hashset.Set).Add(connection_state{conn, uint8(state_num)})
		// conections is a pointer to hashset.
		// so this operation is reflected to the clients.
	}
	mtx.Unlock()
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			removeClientFromClients(user_id, conn)
			break
		}

		// parse
		inst_json := strings.SplitN(string(msg), ",", 3)
		if len(inst_json) != 3 {
			SendErrorMessage(conn, "Invalid instruction format.")
			continue
		}
		tokenString := inst_json[1]
		json_str := inst_json[2]

		// Token validation
		var jsonData MsgToken
		err = json.Unmarshal([]byte(json_str), &jsonData)
		if err != nil {
			SendErrorMessage(conn, "Invalid JSON format. (Unauthorized)")
			continue
		}
		_, err = auth.Auth(tokenString)
		if err != nil {
			SendErrorMessage(conn, "Invalid Token. (Unauthorized)")
			// Unauthorized
			continue
		}

		switch instruction := inst_json[0]; instruction {

		case "0":
			// receive client TL state.
			var jsonData MsgTLState
			err := json.Unmarshal([]byte(json_str), &jsonData)
			if err != nil {
				SendErrorMessage(conn, "Invalid JSON format.")
				continue
			}
			state := jsonData.State
			connections, not_first := clients.Get(user_id)
			if !not_first {
				SendErrorMessage(conn, "Invalid user ID.")
				continue
			}
			set := connections.(*hashset.Set)
			set.Add(connection_state{conn, state})
		case "1":
			// create meigen.
			var jsonData MsgMeigen
			err := json.Unmarshal([]byte(json_str), &jsonData)
			if err != nil {
				SendErrorMessage(conn, "Invalid JSON format.")
				continue
			}
			tx, err := db_handle.BeginTx(ctx, nil)
			if err != nil {
				log.Printf("Failed to create tx: %+v", err)
				continue
			}
			queries := db.New(tx)
			followers, err := queries.GetFollowers(ctx, user_id)
			if err != nil {
				log.Printf("Failed to get followers: %+v", err)
				continue
			}
			// Create Meigen.
			// 1. Get default_group_id
			def_grp_id, err := queries.GetDefaultGroupID(ctx, user_id)
			if err != nil {
				log.Printf("Failed to get default group ID: %+v", err)
				continue
			}
			// 2. Create Poet if not exists.
			poet_id, err := queries.CheckPoetExists(ctx, db.CheckPoetExistsParams{
				Name:    jsonData.Poet,
				GroupID: def_grp_id,
			})
			if err != nil {
				poet_id, err = queries.CreatePoet(ctx, db.CreatePoetParams{
					Name:    jsonData.Poet,
					GroupID: def_grp_id,
				})
				if err != nil {
					log.Printf("Failed to create poet: %+v", err)
					continue
				}
			}
			if err != nil {
				log.Printf("Failed to get poet ID: %+v", err)
				tx.Rollback()
				continue
			}
			// 4. Create Meigee.
			meigen_id, err := queries.CreateMeigen(ctx, db.CreateMeigenParams{
				Meigen:  jsonData.Meigen,
				WhomID:  user_id,
				GroupID: def_grp_id,
				PoetID:  poet_id,
			})
			if err != nil {
				log.Printf("Failed to create meigen: %+v", err)
				tx.Rollback()
				continue
			}

			record, _ := queries.GetMeigenContent(ctx, meigen_id)
			meigen, _ := json.Marshal(record)
			tx.Commit()
			//
			SendMeigenToFollowers(followers, []byte("1"+","+string(meigen)), user_id)
			//
		case "2":
			// create meigen to group.
			var jsonData MsgMeigenGroup
			err := json.Unmarshal([]byte(json_str), &jsonData)
			if err != nil {
				SendErrorMessage(conn, "Invalid JSON format.")
				continue
			}
			//
			// TODO: implement.
			//
		}
	}
	// clean
	removeClientFromClients(user_id, conn)
}

func SendErrorMessage(conn *websocket.Conn, msg string) {
	conn.WriteMessage(websocket.TextMessage, []byte("0,"+msg))
}

func SendMeigenToFollowers(recipients_candidate_ []string, msg []byte, user_id string) {
	// Resipients = INTERSECTION of (recipients_candidate, LOGGED_IN_USER)
	// Both are sorted.
	// So this algorithm can be used. O(max(len(A), len(B)))
	/*
		A, B: sorted.
		====
		for i = 0 to len(A)
			for j = i to len(B)
			if A[i] == B[j] then C.append(v)
			else if A[i] < B[j] then continue
		====
		C: result.
	*/
	recipients_candidate := append(recipients_candidate_, user_id)
	resipients := hashset.New()
	it := clients.Iterator()
	for it.Next() {
		for _, pair := range it.Value().(*hashset.Set).Values() {
			if pair.(connection_state).state == 1 {
				resipients.Add(pair.(connection_state).conn)
			}
		}
		for _, candidate := range recipients_candidate {
			if it.Key() == candidate {
				// it.Values are the (user's) hashset of connections.
				// user may have multiple connections.
				for _, pair := range it.Value().(*hashset.Set).Values() {
					// conn: *websocket.Conn
					conn := pair.(connection_state).conn
					resipients.Add(conn)
				}
			} else if it.Key().(string) > candidate {
				continue
			}
		}
	}
	for _, conn := range resipients.Values() {
		go func(conn *websocket.Conn) {
			conn.WriteMessage(websocket.TextMessage, msg)
		}(conn.(*websocket.Conn))
	}
}
