package controller

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"

	"meigens-api/db"
)

func SearchUsers(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	query := c.PostForm("query")

	query0 := "%" + query + "%"

	queries := db.New(db_handle)

	users, err := queries.SearchUsers(ctx, query0)
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	c.JSON(200, gin.H{
		"found_users": users,
	})

}

func Follow(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	target_id := c.PostForm("target_id")
	user_id, _ := c.Get("user_id")

	queries := db.New(db_handle)

	if _, err := queries.CheckUserExists(ctx, target_id); err != nil {
		BadRequest(c, "The target user does not exist.")
	}

	if err := queries.Follow(ctx, db.FollowParams{
		FollowerID: user_id.(string),
		FolloweeID: target_id,
	}); err != nil {
		BadRequest(c, "Already followed.")
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully followed \"" + target_id + "\".",
	})
}
