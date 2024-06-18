package app

import (
	"database/sql"
	"log"
	"meigens-api/src/controller"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = treemap.NewWithStringComparator()

func SetRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/signup", Signup)
	r.POST("/login", Login)

	authGroup := r.Group("/auth")
	authGroup.Use(AuthMiddleware)
	authGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	socketUpgrader := websocket.Upgrader {
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	authGroup.GET("/socket", func(ctx *gin.Context) {
		user_id := ctx.MustGet("user_id").(string)
		conn, err := socketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %+v", err)
			return
		}
		clients.Put(user_id, conn)
		// 以下の処理についていずれはMutexを考慮して書く必要があるよ。
		for {
			t, msg, err := conn.ReadMessage()	
			if err != nil {
				log.Printf("Failed to read message: %+v", err)
				break
			}
			it := clients.Iterator()
			for it.Next() {
				it.Value().(*websocket.Conn).WriteMessage(t, msg)
			}
		}
	})

	authGroup.GET("/fetch_group_ids", controller.FetchGroups)
	authGroup.GET("/fetch_tl", controller.FetchTL)
	authGroup.POST("/search_users", controller.SearchUsers)
	authGroup.POST("/add_group", controller.AddGroup)
	authGroup.POST("/add_meigen_to_group", controller.AddMeigenToGroup)
	authGroup.POST("/add_meigen", controller.AddMeigen)
	authGroup.POST("/follow", controller.Follow)
	authGroup.POST("/reaction", controller.Reaction)
	authGroup.PATCH("/patch_user_image", controller.PatchUserImage)
	authGroup.PATCH("/patch_group_image", controller.PatchGroupImage)

	// defer db.Close()

	return r
}
