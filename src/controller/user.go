package controller

import (
	"bytes"
	"context"
	"database/sql"
	"image"
	"image/png"
	_ "image/jpeg"

	"github.com/gin-gonic/gin"

	"meigens-api/db"
)

func PatchUserImage(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)

	// encode image file to be required format.
	// 256x256 png, small size.
	img_obj, err := c.FormFile("image")
	if err != nil {
		BadRequest(c, "No image file.")
		return
	}
	img, err2 := img_obj.Open()

	if err2 != nil {
		BadRequest(c, "Failed to open image file.")
		return
	}

	img_bin, _, err3 := image.Decode(img)
	if err3 != nil {
		BadRequest(c, "Invalid image file.")
		return
	}
	if img_bin.Bounds().Dx() != 256 || img_bin.Bounds().Dy() != 256 {
		BadRequest(c, "Invalid image file. 256x256 size is expected.")
		return
	}

	defer img.Close()

	var img_png bytes.Buffer
	err4 := png.Encode(&img_png, img_bin)
	if err4 != nil {
		InternalServerError(c, "Failed to encode image.")
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
