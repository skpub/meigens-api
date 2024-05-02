package app

import (
	"github.com/gin-gonic/gin"

	"meigens-api/src/controller"
	"meigens-api/src/db"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	dbb, _ := db.Conn()


	r.GET("/list", func(c *gin.Context){
		controller.Create(dbb)
		c.JSON(200, gin.H{
			"message": "test OK.",
		})
	})
	return r
}