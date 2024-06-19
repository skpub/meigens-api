package controller

import (
	"context"
	"database/sql"
	"log"
	"sync"

	// "math"
	"meigens-api/db"
	"strconv"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var mtx sync.Mutex
var clients = treemap.NewWithStringComparator()
// user_id: hashset of conn
// user couuld have multiple connections so we need to store them in a slice.

var socketUpgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
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

func TLSocket(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	conn, err := socketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
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
	for {
		t, msg, err := conn.ReadMessage()	
		if err != nil {
			log.Printf("Failed to read message: %+v", err)
			removeClientFromClients(user_id, conn)
			break
		}
		mtx.Lock()
		it := clients.Iterator()
		// it.Key(): user_id, it.Value(): hashset of conn
		mtx.Unlock()
		for it.Next() {
			for _, conn := range it.Value().(*hashset.Set).Values() {
				conn := conn.(*websocket.Conn)
				conn.WriteMessage(t, msg)
			}
		}
	}
	// clean
	removeClientFromClients(user_id, conn)
}

func FetchTL(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)
	before := c.PostForm("before")
	var before_u int64
	if before == "" || before == "null" || before == "nil" {
		// before now+epsilon, that's mean all contents are expected to fetch.
		before_u = time.Now().Add(114514).Unix()
	} else {
		var err error
		before_u, err = strconv.ParseInt(before, 10, 64)
		if err != nil {
			BadRequest(c, "Invalid time format (before). Unixtime expected.")
			return
		}
	}
	before_t := time.Unix(before_u, 0)
	before_nulltime := sql.NullTime {
		Time: before_t,
		Valid: true,
	}

	queries := db.New(db_handle)

	contents, err := queries.FetchTL(ctx, db.FetchTLParams{
		FollowerID: user_id,
		Limit: 20,
		CreatedAt: before_nulltime,
	})
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}
	c.JSON(200, gin.H{
		"contents": contents,
	})
}
