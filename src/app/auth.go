package app

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"meigens-api/db"
	"meigens-api/src/auth"
	"meigens-api/src/controller"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Signup(c *gin.Context) {
	user_id := c.PostForm("user_id")
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	ctx := context.Background()

	db_handle := c.MustGet("db").(*sql.DB)
	tx, err := db_handle.BeginTx(ctx, nil)
	if err != nil {
		controller.InternalServerError(c, "DB error")
		return
	}
	queries := db.New(tx)

	// Check if user already exists.
	if count_users, err := queries.CheckUserExists(ctx, user_id); err != nil {
		controller.InternalServerError(c, "DB error")
	} else if count_users > 0 {
		c.JSON(400, gin.H{
			"message": "User already exists.",
		})
		c.Abort()
		return
	}

	// Create Default group for the user.
	if _, err := queries.CheckGroupExists(ctx, db.CheckGroupExistsParams{
		UserID: user_id,
		Name:   user_id + "_DEFAULT",
	}); err != nil {
		controller.InternalServerError(c, "DB error")
	}

	group_id, err := queries.CreateGroup(ctx, user_id+"_DEFAULT")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to create default group.",
		})
		c.Abort()
		tx.Rollback()
		return
	}

	// Create new user.
	password_hash := sha256.Sum256([]byte(password))
	new_user_params := db.CreateUserParams{
		ID:             user_id,
		Name:           username,
		Email:          email,
		Password:       hex.EncodeToString(password_hash[:]),
		DefaultGroupID: group_id,
	}

	_, err2 := queries.CreateUser(context.Background(), new_user_params)
	if err2 != nil {
		c.JSON(500, gin.H{
			"message": "failed to create user.",
		})
		c.Abort()
		tx.Rollback()
		return
	}

	queries.InitDefaultUG(ctx, db.InitDefaultUGParams{UserID: user_id, GroupID: group_id})
	tx.Commit()

	// Successfully added.
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("added %s", username),
	})
}

func Login(c *gin.Context) {
	user_id := c.PostForm("user_id")
	password := c.PostForm("password")

	secret := os.Getenv("SECRET")

	db_handle, _ := c.MustGet("db").(*sql.DB)

	queries := db.New(db_handle)

	password_hash := sha256.Sum256([]byte(password))
	user_params := db.LoginParams{
		ID:       user_id,
		Password: hex.EncodeToString(password_hash[:]),
	}

	if _, err := queries.Login(
		context.Background(), user_params); err != nil {
		// invalid username or password
		c.JSON(400, gin.H{
			"message": "invalid username or password.",
		})
		c.Abort()
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(time.Hour * 24 * 3).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		c.Abort()
		return
	} else {
		c.JSON(200, gin.H{
			"message": "You got an access token.",
			"token":   tokenString,
		})
	}
}

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user_id, err := auth.Auth(tokenString)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.Set("user_id", user_id)
}
