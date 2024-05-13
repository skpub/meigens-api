package controller

import (
	"database/sql"
	"context"
	"meigens-api/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchGroups(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	user_id_uuid, _ := uuid.Parse(user_id.(string))

	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	queries := db.New(db_handle)
	if group_ids, err := queries.GetGroupsParticipated(ctx, user_id_uuid); err != nil {
		InternalServerError(c, "failed to fetch groups.")
		return
	} else {
		c.JSON(200, gin.H{
			"groups": group_ids,
		})
	}
}

func AddGroup(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	user_id_uuid, _ := uuid.Parse(user_id.(string))
	group_name := c.PostForm("group_name")

	queries := db.New(db_handle)

	group_ex := db.GroupEXParams {
		UserID: user_id_uuid,
		Name: group_name,
	}
	if count, err := queries.GroupEX(ctx, group_ex); err != nil {
		InternalServerError(c, "failed to check the group.")
	} else if count > 0 {
		BadRequest(c, "The group you specified already exists.")
	} else {
		// Not exist: then add the group user specified.
		new_group_params := db.CreateGroupParams {
			Name: group_name,
			OwnerID: user_id_uuid,
		}
		if new_group_id, err := queries.CreateGroup(ctx, new_group_params); err != nil {
			InternalServerError(c, "failed to add the group.")
			return
		} else {
			if err := queries.AddUserToGroup(ctx, db.AddUserToGroupParams {
				UserID: user_id_uuid,
				GroupID: new_group_id}); err != nil {
				// Strange error !!!!
				InternalServerError(c, "failed to add the user to the group.")
				// TODO: delete the group.
				return 
			} else {
			// Successfully added.
				c.JSON(200, gin.H{
					"message": "Successfully added the group. and you are the owner of the group.",
				})
			}
		}
	}
}