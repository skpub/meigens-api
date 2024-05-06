package app

import (
	"context"
	"fmt"
	"meigens-api/src/model"
	"meigens-api/src/controller"
	"meigens-api/src/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	db_handle, _ := db.Conn()

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
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	secret := os.Getenv("SECRET")

	db_handle, _ := db.Conn()

	err := controller.Login(db_handle, username, password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	} else {
		var user model.Users
		err := db_handle.NewSelect().
			Model(&user).
			Where("name = ?", username).
			Scan(context.Background())
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
			"user_id": user.Id,
			"username": username,
			"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"error": "Failed to generate token.",
			})
			return
		}
		c.Header("Authorization", tokenString)
		c.JSON(200, gin.H{
			"message": "You got an access token.",
		})
	}
}

func AuthMiddleware (c *gin.Context) {
	secret := os.Getenv("SECRET")
	tokenString := c.GetHeader("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H {
			"error": "Unauthorized.",
		})
		c.Abort()
		return
	}

	c.Next()
}
