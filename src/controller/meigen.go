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
	meigen_id, err2 := queries.CreateMeigen(ctx, db.CreateMeigenParams{
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
		"meigen_id": meigen_id,
	})
}

func AddMeigenToGroup(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id, _ 	:= c.Get("user_id")
	group_id, _ := uuid.Parse(c.PostForm("group_id"))
	poet 		:= c.PostForm("poet")
	meigen 		:= c.PostForm("meigen")

	queries := db.New(db_handle)

	// Check if the user is in the group.
	if count, err := queries.CheckUserExistsGroup(ctx, db.CheckUserExistsGroupParams{
		UserID: user_id.(string), GroupID: group_id}); err != nil {
		InternalServerError(c, "DB error")
	} else if count == 0 {
		BadRequest(c, "You are not in the group you specified.")
		return
	}
	// User has the permission to add the meigen to specified group.

	var poet_id uuid.UUID
	// Check if the poet already exists.
	if poet_ex, err := queries.CheckPoetExists(ctx, db.CheckPoetExistsParams{
		Name: poet, GroupID: group_id}); err != nil {
		InternalServerError(c, "DB error")
	} else if poet_ex == 0 {
		// Create Poet for specified group.
		poet_id, err = queries.CreatePoet(ctx, db.CreatePoetParams{Name: poet, GroupID: group_id})
		if err != nil {
			InternalServerError(c, "DB error")
			return
		}
	} else {
		// Get the Poet ID.
		poet_id, err = queries.GetPoetID(ctx, db.GetPoetIDParams{Name: poet, GroupID: group_id})
		if err != nil {
			InternalServerError(c, "DB error")
			return
		}
	}

	// Create Meigen.
	meigen_id, err2 := queries.CreateMeigen(ctx, db.CreateMeigenParams{
		Meigen: meigen,
		WhomID: user_id.(string),
		GroupID: group_id,
		PoetID: poet_id,
	})
	if err2 != nil {
		InternalServerError(c, "DB error" + err2.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully added the meigen to the group.",
		"meigen_id": meigen_id,
	})
}
