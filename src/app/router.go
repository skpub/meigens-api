package app

import (
	"database/sql"
	"meigens-api/src/controller"
	"meigens-api/src/controller/socket"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	front := os.Getenv("FRONT_ORIGIN")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{front},
		AllowHeaders: []string{"Authorization"},
	}))
	r.POST("/signup", Signup)
	r.POST("/login", Login)

	authGroup := r.Group("/auth")
	authGroup.Use(AuthMiddleware)
	authGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/socket", socket.TLSocket)
	authGroup.GET("/fetch_group_ids", controller.FetchGroups)
	authGroup.GET("/fetch_tl", controller.FetchTL)
	authGroup.GET("/fetch_user_imgs", controller.FetchUserImgs)
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
