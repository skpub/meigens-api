package app

import (
	"database/sql"
	"meigens-api/src/controller"

	"github.com/gin-gonic/gin"
)

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

	authGroup.GET("/fetch_group_ids", controller.FetchGroups)
	authGroup.POST("/search_users", controller.SearchUsers)
	authGroup.POST("/add_group", controller.AddGroup)
	authGroup.POST("/add_meigen_to_group", controller.AddMeigenToGroup)
	authGroup.POST("/add_meigen", controller.AddMeigen)
	authGroup.POST("/follow", controller.Follow)
	authGroup.PATCH("/patch_user_image", controller.PatchUserImage)
	authGroup.PATCH("/patch_group_image", controller.PatchGroupImage)

	// defer db.Close()

	return r
}
