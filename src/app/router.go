package app

import (
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", Signup)
	r.POST("/login", Login)

	authGroup := r.Group("/auth")
	authGroup.Use(AuthMiddleware)
	authGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return r
}
