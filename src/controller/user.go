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

	queries := db.New(db_handle)
	if permission, err := queries.CheckUserExistsGroup(ctx, db.CheckUserExistsGroupParams{
		UserID:  user_id,
		GroupID: group_id,
	}); err != nil {
		InternalServerError(c, "DB error. Can't check if the user is in the group.")
		return
	} else if permission&1 == 0 { // check if the user has permission to patch the group image(WRITE).
		BadRequest(c, "You don't have the permission to patch the group image.")
		return
	}
	var img_png bytes.Buffer
	if err2 := imgEncode(c, &img_png); err2 != nil {
		return
	}
	// Finaly, we got the image binary 'img_png' to be stored in the database.
	if err3 := queries.PatchGroupImage(ctx, db.PatchGroupImageParams{
		ID:  group_id,
		Img: img_png.Bytes(),
	}); err3 != nil {
		InternalServerError(c, "Failed to store image. :"+err3.Error())
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
		InternalServerError(c, "Failed to store image. :"+err5.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully patched the user image.",
	})
}

func PatchUserName(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)

	name := c.PostForm("name")

	queries := db.New(db_handle)
	err := queries.PatchUserName(ctx, db.PatchUserNameParams{
		ID:   user_id,
		Name: name,
	})
	if err != nil {
		InternalServerError(c, "Failed to patch name.")
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully patched the name.",
		"name":    name,
	})
}

func PatchUserBio(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	user_id := c.MustGet("user_id").(string)

	bio := c.PostForm("bio")

	queries := db.New(db_handle)
	bio_nullstr := sql.NullString{String: bio, Valid: true}
	err := queries.PatchUserBio(ctx, db.PatchUserBioParams{
		ID:  user_id,
		Bio: bio_nullstr,
	})

	if err != nil {
		InternalServerError(c, "Failed to patch name.")
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully patched the name.",
		"bio":     bio,
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

type UserProfile struct {
	Username    string `json:"name"`
	Bio         string `json:"bio"`
	IsFollowing bool   `json:"is_following"`
}

func FetchUserProfile(c *gin.Context) {
	db_handle := c.MustGet("db").(*sql.DB)
	ctx := context.Background()

	me, _ := c.Get("user_id")
	user_id := c.Query("user_id")
	queries := db.New(db_handle)

	user, err := queries.GetUserProfile(ctx, user_id)
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	is_following, err := queries.CheckFollowing(ctx, db.CheckFollowingParams{
		FollowerID: user_id,
		FolloweeID: me.(string),
	})
	if err != nil {
		InternalServerError(c, "DB error")
		return
	}

	var bio string
	if user.Bio.Valid {
		bio = user.Bio.String
	} else {
		bio = ""
	}

	user_profile := UserProfile{
		Username:    user.Name,
		Bio:         bio,
		IsFollowing: is_following > 0,
	}

	c.JSON(200, gin.H{
		"contents": user_profile,
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
