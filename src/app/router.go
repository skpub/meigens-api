package app

import "github.com/gin-gonic/gin"

func SetRouter() *gin.Engine {
	r:= gin.Default()
	r.GET("test", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "test OK.",
		})
	})
	return r
}