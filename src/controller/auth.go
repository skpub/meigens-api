package controller

import (
	"context"
	"database/sql"
	"fmt"
	"encoding/hex"

	"github.com/uptrace/bun"

	"crypto/sha256"
	"meigens-api/src/model"
)

func Login(db *bun.DB, username string, password string) error {
	pwhash := sha256.Sum256([]byte(password))

	var user model.Users
	err := db.NewSelect().
		Model(&user).
		Where("name = ?", username).
		Scan(context.Background())
	if err != nil {
		return err
	} else {
		if user.Password == hex.EncodeToString(pwhash[:]) {
			return nil
		} else {
			return fmt.Errorf("password is incorrect.")
		}
	}
}

func CreateUser(db *bun.DB, username string, password string, email string) error {
	user := new(model.Users)
	exists, err := db.NewSelect().
		Model(user).
		Where("name = ?", username).Exists(context.Background())
	if err != nil && err != sql.ErrNoRows {
		return err
	} else
	if exists {
		return fmt.Errorf("user already exists.")
	} else {
		pw_hash := sha256.Sum256([]byte(password))

		user := model.Users {
			Name: username,
			Bio: "",
			Email: email,
			Password: hex.EncodeToString(pw_hash[:]),
		}
		_, err := db.NewInsert().Model(&user).Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
}