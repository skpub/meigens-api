package controller

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"meigens-api/db"

	"github.com/gin-gonic/gin"
)



func FetchTL(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)
	before := c.PostForm("before")
	var before_u int64
	if before == "" || before == "null" || before == "nil" {
		// before now+epsilon, that's mean all contents are expected to fetch.
		before_u = time.Now().Add(114514).Unix()
	} else {
		var err error
		before_u, err = strconv.ParseInt(before, 10, 64)
		if err != nil {
			BadRequest(c, "Invalid time format (before). Unixtime expected.")
			return
		}
	}
	before_t := time.Unix(before_u, 0)
	before_nulltime := sql.NullTime {
		Time: before_t,
		Valid: true,
	}

	queries := db.New(db_handle)

	contents, err := queries.FetchTL(ctx, db.FetchTLParams{
		FollowerID: user_id,
		Limit: 20,
		CreatedAt: before_nulltime,
	})
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}
	c.JSON(200, gin.H{
		"contents": contents,
	})
}
