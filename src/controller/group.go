package controller

import (
	"context"
	"fmt"
	"meigens-api/src/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func FetchGroups(c *gin.Context) {
	db := c.MustGet("db").(*bun.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	// user_id_uuid, _ := uuid.Parse(user_id.(string))

	var g_r_rels []model.UserGroupRels
	err := db.NewSelect().
		Model(&g_r_rels).
		Where("user_id = ?", user_id).
		Scan(ctx)
	if err != nil {
		fmt.Println(err)
		InternalServerError(c, "DB error")
	}

	group_ids := make([]uuid.UUID, len(g_r_rels))
	for i := 0; i < len(g_r_rels); i++ {
		group_ids[i] = g_r_rels[i].GroupID
	}

	c.JSON(200, gin.H{
		"groups": group_ids,
	})
}

func AddGroup(c *gin.Context) {
	db := c.MustGet("db").(*bun.DB)
	ctx := context.Background()

	user_id, _ := c.Get("user_id")
	group_name := c.PostForm("group_name")

	group_stil_ex := []model.Groups {}
	if count, err := db.NewSelect().
		Model(&group_stil_ex).
		Where("name = ?", group_name).
		Count(ctx); err != nil || count == 0 {
		// Not exist: then add the group user specified.
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
	
		c.JSON(200, gin.H{
			"message": "Successfully added the group, and you took part in it.",
		})
		return
	} else {
		BadRequest(c, "The group you specified already exists.")
	}
}