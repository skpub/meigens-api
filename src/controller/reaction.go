package controller

import (
	"context"
	"database/sql"
	"meigens-api/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Reaction(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)

	meigen_id := c.PostForm("meigen_id")
	meigen_id_uuid, err := uuid.Parse(meigen_id)
	if err != nil {
		BadRequest(c, "Invalid meigen_id.")
		return
	}

	reaction := c.PostForm("reaction")
	reaction_num_, err2 := strconv.Atoi(reaction)
	reaction_num := int32(reaction_num_)
	if err2 != nil {
		BadRequest(c, "Invalid reaction.")
		return
	}

	queries := db.New(db_handle)
	// Check if the reaction already exists.
	if count, err := queries.CheckReactionExists(ctx, db.CheckReactionExistsParams{
		MeigenID: meigen_id_uuid,
		UserID: user_id,
		Reaction: reaction_num}); count > 0 || err != nil {
		BadRequest(c, "Reaction already exists.")
		return
	}
	// Check if the meigen exists.
	if count, err := queries.CheckMeigenExists(ctx, meigen_id_uuid); count < 1 || err != nil {
		BadRequest(c, "Meigen does not exist.")
		return
	}	


	reaction_id, err3 := queries.CreateReaction(ctx, db.CreateReactionParams{
		MeigenID: meigen_id_uuid,
		UserID: user_id,
		Reaction: reaction_num,
	})
	if err3 != nil {
		InternalServerError(c, "DB error")
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully added the reaction.",
		"reaction_id": reaction_id,
	})
}
