package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun/extra/bundebug"

	"meigens-api/src/controller"
	"meigens-api/src/db"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	db_handle, _ := db.Conn()

	db_handle.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	r.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")

		err := controller.CreateUser(db_handle, username, password, email)
		if err != nil {
			// c.Error(err)
			// User already exists.
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else {
			// Successfully added.
			c.JSON(200, gin.H{
				"message": fmt.Sprintf("added %s", username),
			})
		}

	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		err := controller.Login(db_handle, username, password)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "this is token",
				// TODO implement to generate token.
			})
		}
	})

	return r
}