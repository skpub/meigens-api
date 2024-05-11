package controller

import (
	"github.com/gin-gonic/gin"
)

func InternalServerError(c *gin.Context, er_str string) {
	c.JSON(500, gin.H{
		"message": "Internal Server Error: " + er_str,
	})
	c.Abort()
}

func BadRequest(c *gin.Context, er_str string) {
	c.JSON(400, gin.H{
		"message": "Bad Request: " + er_str,
	})
	c.Abort()
}
