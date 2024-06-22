package socket

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"meigens-api/db"
	"strings"
	"sync"

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

func TLSocket(c *gin.Context) {
	user_id := c.MustGet("user_id").(string)
	conn, err := socketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	mtx.Lock()
	connections, not_first := clients.Get(user_id)
	// connections: nil or hashset of conn
	if !not_first {
		var connections_set = hashset.New()
		connections_set.Add(conn)
		clients.Put(user_id, connections_set)
	} else {
		connections.(*hashset.Set).Add(conn)
		// conections is a pointer to hashset.
		// so this operation is reflected to the clients.
	}
	mtx.Unlock()
	db_handle := c.MustGet("db").(*sql.DB)
	queries := db.New(db_handle)
	ctx := context.Background()
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %+v", err)
			removeClientFromClients(user_id, conn)
			break
		}

		// parse
		inst_json := strings.SplitN(string(msg), ",", 2)
		json_str := inst_json[1]
		switch instruction := inst_json[0]; instruction {
			case "0":
				// receive client TL state.
				var jsonData MsgTLState
				err := json.Unmarshal([]byte(inst_json[1]), &jsonData)
				if err != nil {
					conn.WriteMessage(t, []byte("Invalid JSON format."))
					continue
				}
				//
				// TODO: implement.
				//
			case "1":
				// create meigen.
				var jsonData MsgMeigen
				err := json.Unmarshal([]byte(inst_json[1]), &jsonData)
				if err != nil {
					conn.WriteMessage(t, []byte("Invalid JSON format."))
					continue
				}
				followers, err := queries.GetFollowers(ctx, user_id)
				if err != nil {
					log.Printf("Failed to get followers: %+v", err)
					continue
				}
				//
				SendMessage(followers, []byte(json_str))
				//
			case "2":
				// create meigen to group.
				var jsonData MsgMeigenGroup
				err := json.Unmarshal([]byte(inst_json[1]), &jsonData)
				if err != nil {
					conn.WriteMessage(t, []byte("Invalid JSON format."))
					continue
				}
				//
				// TODO: implement.
				//
		}
		// end parse

		mtx.Lock()
		it := clients.Iterator()
		// it.Key(): user_id, it.Value(): hashset of conn
		for it.Next() {
			for _, conn := range it.Value().(*hashset.Set).Values() {
				conn := conn.(*websocket.Conn)
				conn.WriteMessage(t, msg)
			}
		}
		mtx.Unlock()
	}
	// clean
	removeClientFromClients(user_id, conn)
}

func SendMessage(recipients_candidate []string, msg []byte) {
	// Resipients = INTERSECTION of (recipients_candidate, LOGGED_IN_USER)
	// Both are sorted.
	// So this algorithm can be used. O(max(len(A), len(B)))
	/*
	    A, B: sorted.
		====
		for i = 0 to len(A)
			for j = i to len(B)
			if A[i] == B[j] then C.append(v)
			else if A[i] < B[j] then break
		====
		C: result.
	*/
	resipients := hashset.New()
	it := clients.Iterator()
	for it.Next() {
		for candidate := range recipients_candidate {
			if it.Key() == candidate {
				// it.Values are the (user's) hashset of connections.
				// user may have multiple connections.
				for _, conn := range it.Value().(*hashset.Set).Values() {
					// conn: *websocket.Conn
					resipients.Add(conn)
				}
			}
		}
	}
	for _, conn := range resipients.Values() {
		go func(conn *websocket.Conn) {
			conn.WriteMessage(websocket.TextMessage, msg)
		} (conn.(*websocket.Conn))
	}
}
