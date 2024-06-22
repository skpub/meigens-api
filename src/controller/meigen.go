package controller

import (
	"context"
	"database/sql"

	"meigens-api/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddMeigen(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()
	tx, err := db_handle.BeginTx(ctx, nil)
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}
	queries := db.New(tx)

	user_id, _ := c.Get("user_id")
	poet := c.PostForm("poet")
	meigen := c.PostForm("meigen")

	// Get the Default group ID of the user.
	default_group_id, err := queries.GetDefaultGroupID(ctx, user_id.(string))
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	// Check if the meigen already exists.
	if count, err := queries.CheckMeigenExistsByMeigen(ctx, db.CheckMeigenExistsByMeigenParams{
		Meigen:  meigen,
		WhomID:  user_id.(string),
		GroupID: default_group_id,
		Name:    poet,
	}); err != nil || count > 0 {
		BadRequest(c, "Meigen already exists.")
		return
	}

	// Create poet if not exists.
	poet_id, err := queries.CheckPoetExists(ctx, db.CheckPoetExistsParams{
		Name: poet,
		GroupID: default_group_id,
	})
	if err != nil {
		poet_id, err = queries.CreatePoet(ctx, db.CreatePoetParams{
			Name: poet,
			GroupID: default_group_id,
		})
		if err != nil {
			InternalServerError(c, "DB error")
			tx.Rollback()
			return
		}
	}
	// Create meigen.
	meigen_id, err2 := queries.CreateMeigen(ctx, db.CreateMeigenParams{
		Meigen:  meigen,
		WhomID:  user_id.(string),
		GroupID: default_group_id,
		PoetID:  poet_id})

	if err2 != nil {
		InternalServerError(c, "DB error")
		tx.Rollback()
		return
	}
	tx.Commit()

	c.JSON(200, gin.H{
		"message":   "Successfully added the meigen.",
		"meigen_id": meigen_id,
	})
}

func AddMeigenToGroup(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	tx, err := db_handle.BeginTx(ctx, nil)
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	user_id, _ := c.Get("user_id")
	group_id, _ := uuid.Parse(c.PostForm("group_id"))
	poet := c.PostForm("poet")
	meigen := c.PostForm("meigen")

	queries := db.New(tx)

	// Check if the user is in the group.
	if count, err := queries.CheckUserExistsGroup(ctx, db.CheckUserExistsGroupParams{
		UserID: user_id.(string), GroupID: group_id}); err != nil {
		InternalServerError(c, "DB error")
	} else if count == 0 {
		BadRequest(c, "You are not in the group you specified.")
		return
	}
	// User has the permission to add the meigen to specified group.

	// Check if the meigen already exists.
	if count, err := queries.CheckMeigenExistsByMeigen(ctx, db.CheckMeigenExistsByMeigenParams{
		Meigen:  meigen,
		WhomID:  user_id.(string),
		GroupID: group_id,
		Name:    poet,
	}); err != nil || count > 0 {
		BadRequest(c, "Meigen already exists.")
		return
	}

	poet_id, err := queries.CheckPoetExists(ctx, db.CheckPoetExistsParams{
		Name: poet,
		GroupID: group_id,
	})
	if err != nil {
		poet_id, err = queries.CreatePoet(ctx, db.CreatePoetParams{
			Name: poet,
			GroupID: group_id,
		})
		if err != nil {
			tx.Rollback()
			InternalServerError(c, "DB error")
			return
		}
	}
	// Create Meigen.
	meigen_id, err2 := queries.CreateMeigen(ctx, db.CreateMeigenParams{
		Meigen:  meigen,
		WhomID:  user_id.(string),
		GroupID: group_id,
		PoetID:  poet_id,
	})
	if err2 != nil {
		InternalServerError(c, "DB error"+err2.Error())
		tx.Rollback()
		return
	}
	tx.Commit()

	c.JSON(200, gin.H{
		"message":   "Successfully added the meigen to the group.",
		"meigen_id": meigen_id,
	})
}
