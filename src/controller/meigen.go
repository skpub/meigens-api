package controller

import (
	"github.com/uptrace/bun"
	"github.com/gin-gonic/gin"
	"context"
	"meigens-api/src/model"
	"github.com/google/uuid"
)

func AddGroup(c *gin.Context) {
	db := c.MustGet("db").(*bun.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	group_name := c.PostForm("group_name")

	new_group := &model.Groups {
		Name: group_name,
	}

	db.NewInsert().Model(new_group).Exec(ctx)

	// Add the user to the group.
	user_id_uuid, _ := uuid.Parse(user_id.(string))
	new_user_group_rel := &model.UserGroupRels {
		GroupID: new_group.Id,
		UserID: user_id_uuid,
	}
	db.NewInsert().
		Model(new_user_group_rel).
		Exec(ctx)
	
	defer db.Close()
	c.JSON(200, gin.H{
		"message": "Successfully added the group, and you took part in it.",
	})
}

func AddMeigenToGroup(c *gin.Context) {
	db := c.MustGet("db").(*bun.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	group_id := c.PostForm("group_id")
	poet := c.PostForm("poet")
	meigen := c.PostForm("meigen")

	// Check if the user is in the group.
	ugRel := new(model.UserGroupRels)
	err := db.NewSelect().
		Model(&ugRel).
		Where("user_id = ?", user_id).
		Where("group_id = ?", group_id).
		Scan(ctx)
	if err != nil {
		InternalServerError(c, "DB error")
	}
	if ugRel == nil {
		BadRequest(c, "You are not in the group you specified.")
	}

	group_id_uuid, _ := uuid.Parse(group_id)

	// Check if the poet exists.
	g_poet := new(model.GroupPoets)

	if err := db.NewSelect().
		Model(&g_poet).
		Join("poets ON poets.id = group_poets.poet_id").
		Scan(ctx); err != nil {
		// Not exist: then add the poet user specified.
		new_poet := model.Poets {
			Name: poet,
		}
		db.NewInsert().
			Model(new_poet).
			Exec(ctx)

		new_poet_g := model.GroupPoets {
			GroupID: group_id_uuid,
			PoetID: new_poet.Id,
		}
		db.NewInsert().
			Model(new_poet_g).
			Exec(ctx)
	} else {

	}

	// Exist: then get the poet_id.
		poet_id := g_poet.PoetID


	new_column := model.Meigens {
		Meigen: meigen,
		WhomID: user_id.(uuid.UUID),
		GroupID: group_id_uuid,
		PoetID: poet_id,
	}
	db.NewAddColumn().
		Model(new_column).
		Exec(ctx)
	
	defer db.Close()

	c.JSON(200, gin.H{
		"message": "Successfully added the meigen to the group.",
	})
}