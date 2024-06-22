package controller

import (
	"bytes"
	"context"
	"database/sql"
	"image"
	_ "image/jpeg"
	"image/png"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"meigens-api/db"
)

func imgEncode(c *gin.Context, img_png *bytes.Buffer) error {
	// encode image file to be required format.
	// 256x256 png, small size.
	img_obj, err := c.FormFile("image")
	if err != nil {
		BadRequest(c, "No image file.")
		return err
	}
	img, err1 := img_obj.Open()
	if err1 != nil {
		BadRequest(c, "No image file.")
		return err1
	}
	img_bin, _, err2 := image.Decode(img)
	if err2 != nil {
		BadRequest(c, "Invalid image file.")
		return err
	}
	if img_bin.Bounds().Dx() != 256 || img_bin.Bounds().Dy() != 256 {
		BadRequest(c, "Invalid image file. 256x256 size is expected.")
		return err
	}
	defer img.Close()
	err3 := png.Encode(img_png, img_bin)
	if err3 != nil {
		InternalServerError(c, "Failed to encode image.")
		return err
	}
	return nil
}

func PatchGroupImage(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)
	group_id := c.PostForm("group_id")

	group_id_uuid, err := uuid.Parse(group_id)
	if err != nil {
		BadRequest(c, "Invalid group_id.")
		return
	}

	queries := db.New(db_handle)
	if permission, err := queries.CheckUserExistsGroup(ctx, db.CheckUserExistsGroupParams{
		UserID: user_id,
		GroupID: group_id_uuid,
	}); err != nil {
		InternalServerError(c, "DB error. Can't check if the user is in the group.")
		return
	} else if permission & 1 == 0 { // check if the user has permission to patch the group image(WRITE).
		BadRequest(c, "You don't have the permission to patch the group image.")
		return
	}
	var img_png bytes.Buffer
	if err2 := imgEncode(c, &img_png); err2 != nil {
		return
	}
	// Finaly, we got the image binary 'img_png' to be stored in the database.
	if err3 := queries.PatchGroupImage(ctx, db.PatchGroupImageParams{
		ID: group_id_uuid,
		Img: img_png.Bytes(),
	}); err3 != nil {
		InternalServerError(c, "Failed to store image. :" + err3.Error())
	}

	c.JSON(200, gin.H{
		"message": "Successfully patched the group image.",
	})
}

func PatchUserImage(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)

	var img_png bytes.Buffer
	if err2 := imgEncode(c, &img_png); err2 != nil {
		return
	}
	// Finally, we got the image binary 'img_png' to be stored in the database.
	queies := db.New(db_handle)
	_, err5 := queies.PatchUserImage(ctx, db.PatchUserImageParams{
		ID:  user_id,
		Img: img_png.Bytes(),
	})
	if err5 != nil {
		InternalServerError(c, "Failed to store image. :" + err5.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully patched the user image.",
	})
}

func SearchUsers(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	query := c.PostForm("query") // search query
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

func FetchUserImgs(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()
	queries := db.New(db_handle)

	user_ids_str := c.Query("user_ids")
	user_ids := strings.Split(user_ids_str, ",")
	user_imgs := make(map[string][]byte)
	for _, user_id := range user_ids {
		img, err := queries.GetUserImg(ctx, user_id)
		if err == nil {
			user_imgs[user_id] = img
		}
	}
	c.JSON(200, gin.H{
		"contents": user_imgs,
	})
}
