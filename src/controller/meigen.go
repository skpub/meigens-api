package controller

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"meigens-api/db"
)

func AddMeigen(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	poet := c.PostForm("poet")
	meigen := c.PostForm("meigen")

	queries := db.New(db_handle)

	// Get the Default group ID of the user.
	default_group_id, err := queries.GetDefaultGroupID(ctx, user_id.(string))
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	// Create poet if not exists.
	poet_column_id, err1 := queries.CreatePoet(
		ctx, db.CreatePoetParams{Name: poet, GroupID: default_group_id})
	
	if err1 != nil {
		InternalServerError(c, "DB error")
		return
	}

	// Create meigen.
	_, err2 := queries.CreateMeigen(ctx, db.CreateMeigenParams{
		Meigen: meigen,
		WhomID: user_id.(string),
		GroupID: default_group_id,
		PoetID: poet_column_id})
	
	if err2 != nil {
		InternalServerError(c, "DB error")
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully added the meigen.",
	})
}

// TODO: FIX
func AddMeigenToGroup(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	user_id_uuid, _ := uuid.Parse(user_id.(string))
	group_id := c.PostForm("group_id")
	group_id_uuid, _ := uuid.Parse(group_id)
	poet := c.PostForm("poet")
	meigen := c.PostForm("meigen")

	queries := db.New(db_handle)

	// Check if the user is in the group.
	gu := db.UserEXGroupParams {
		UserID: user_id_uuid,
		GroupID: group_id_uuid,
	}
	if count, err := queries.UserEXGroup(ctx, gu); err != nil {
		InternalServerError(c, "DB error")
		return
	} else if count == 0 {
		BadRequest(c, "You are not in the group you specified.")
		return
	} else {
		// User has the permission to add the meigen to specified group.

		// Check if the poet exists in the group.
		poet_ex_param := db.PoetExGroupParams {
			Name: poet,
			GroupID: group_id_uuid,
		}
		if count, err := queries.PoetExGroup(ctx, poet_ex_param); err != nil {
			InternalServerError(c, "DB error")
		} else if count == 0 {
			if poet_column_id, err := queries.CreatePoet(ctx, poet); err != nil {
				InternalServerError(c, "DB error")
				return
			} else {
				poet_group_rel := db.CreatePoetGroupRelParams {
					PoetID: poet_column_id,
					GroupID: group_id_uuid,
				}
				if err := queries.CreatePoetGroupRel(ctx, poet_group_rel); err != nil {
					InternalServerError(c, "DB error")
					return
				} else {
					// Successfully added the poet to the group.
					// Then INSERT the meigen.
					meigen_params := db.CreateMeigenParams {
						Meigen: meigen,
						WhomID: user_id_uuid,
						GroupID: uuid.NullUUID{ UUID: group_id_uuid, Valid: true },
						PoetID: poet_column_id,
					}
					if err := queries.CreateMeigen(ctx, meigen_params); err != nil {
						InternalServerError(c, "DB error")
						return
					}
				}
			}
		}

	}
	c.JSON(200, gin.H{
		"message": "Successfully added the meigen to the group.",
	})
}
