package app

import (
	"github.com/gin-gonic/gin"
	"meigens-api/src/db"
	"meigens-api/src/controller"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	db, err := db.Conn()
	if err != nil {
		panic("failed to connect db.")
	}

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

	authGroup.POST("/add_group", controller.AddGroup)
	authGroup.POST("/add_meigen", controller.AddMeigenToGroup)

	return r
}
