package app

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"meigens-api/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	db_handle := c.MustGet("db").(*sql.DB)

	queries := db.New(db_handle)

	// Check if user already exists.
	if name, _ := queries.GetUserByName(context.Background(), username); name == username {
		c.JSON(400, gin.H{
			"message": "User already exists.",
		})
		c.Abort()
		return
	}

	// Create new user.
	new_user_params := db.CreateUserParams {
		Name: username,
		Email: email,
		Password: password,
	}

	err2 := queries.CreateUser(context.Background(), new_user_params)
	if err2 != nil {
		c.JSON(500, gin.H{
			"message": "failed to create user.",
		})
		c.Abort()
		return
	}
	
	// Successfully added.
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("added %s", username),
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	secret := os.Getenv("SECRET")

	db_handle, _ := c.MustGet("db").(*sql.DB)

	queries := db.New(db_handle)

	password_hash := sha256.Sum256([]byte(password))
	user_params := db.LoginParams {
		Name: username,
		Password: hex.EncodeToString(password_hash[:]),
	}

	if user, err := queries.Login(
		context.Background(), user_params);
		err != nil {
		// invalid username or password
		c.JSON(400, gin.H{
			"message": "invalid username or password.",
		})
		c.Abort()
		return
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
			"user_id": user.ID,
			"username": user.Name,
			"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"error": "Failed to generate token.",
			})
			c.Abort()
			return
		} else {
			c.JSON(200, gin.H{
				"message": "You got an access token.",
				"token": tokenString,
			})
		}
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

			db_handle := c.MustGet("db").(*sql.DB)
			queries := db.New(db_handle)
			if name, err := queries.GetUsernameByID(context.Background(), user_id);
				err != nil || name != claims["username"] {
				// invalid username
				c.JSON(http.StatusUnauthorized, gin.H {
					"error": "Unauthorized. (token is valid but user not found)",
				})
				c.Abort()
				return
			} else {
				// Authorized
				c.Set("user_id", claims["user_id"].(string))
				c.Set("username", claims["username"].(string))
				c.Next()
			}
		}
	}
}
