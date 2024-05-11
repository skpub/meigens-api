package app

import (
	"context"
	"encoding/json"
	"fmt"
	"meigens-api/src/controller"
	"meigens-api/src/db"
	"meigens-api/src/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
		// c.Header("Authorization", tokenString)
		c.JSON(200, gin.H{
			"message": "You got an access token.",
			"token": tokenString,
		})
	}
}

func AuthMiddleware (c *gin.Context) {
	secret := os.Getenv("SECRET")
	tokenString := c.GetHeader("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithJSONNumber())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H {
			"error": "Unauthorized. (invalid token)",
		})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, _ := claims["exp"].(json.Number).Int64()
		if exp < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H {
				"error": "Unauthorized. (your token is expired)",
			})
			c.Abort()
			return
		} else {
			user_id, _ := uuid.Parse(claims["user_id"].(string))
			user := model.Users {
				Id: user_id,
				Name: claims["username"].(string),
			}
			db := c.MustGet("db").(*bun.DB)
			if err := db.NewSelect().Model(&user).Scan(context.Background()); err != nil && user.Id == uuid.Nil {
				c.JSON(http.StatusUnauthorized, gin.H {
					"error": "Unauthorized. (token is valid but user not found)",
				})
				c.Abort()
				return
			}
			c.Set("user_id", claims["user_id"].(string))
			c.Set("username", claims["username"].(string))
			c.Next()
		}
	}
}
