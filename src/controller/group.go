package controller

import (
	"database/sql"
	"context"
	"meigens-api/db"

	"github.com/gin-gonic/gin"
)

func FetchGroups(c *gin.Context) {
	user_id, _ := c.Get("user_id")

	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	queries := db.New(db_handle)
	if group_ids, err := queries.GetGroupsParticipated(ctx, user_id.(string)); err != nil {
		InternalServerError(c, "failed to fetch groups.")
		return
	} else {
		c.JSON(200, gin.H{
			"group_ids": group_ids,
		})
	}
}

func AddGroup(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	group_name := c.PostForm("group_name")

	queries := db.New(db_handle)

	// Check if the group already exists.
	if count_groups, err := queries.CheckGroupExists(ctx, db.CheckGroupExistsParams{
		UserID: user_id.(string),
		Name: group_name,
	}); err != nil {
		InternalServerError(c, "failed to check if the group exists.")
		return
	} else if count_groups > 0 {
		BadRequest(c, "Group already exists.")
		return
	}

	new_group_id, err := queries.CreateGroup(ctx, group_name)
	if err != nil {
		InternalServerError(c, "failed to add the group.")
		return
	}
	master_permission := 0xffff;
	err2 := queries.AddUserToGroup(ctx, db.AddUserToGroupParams {
		UserID: user_id.(string),
		GroupID: new_group_id,
		Permission: int16(master_permission),
	})
	if err2 != nil {
		InternalServerError(c, "failed to add the user to the group.")
		queries.DeleteGroup(ctx, new_group_id)
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully added the group.",
		"group_id": new_group_id,
	})
}
