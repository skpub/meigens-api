package controller

import (
	"context"

	"github.com/uptrace/bun"

	"meigens-api/src/model"
)

func List(db *bun.DB) {
}

func Create(db *bun.DB) {
	_, err := db.NewCreateTable().Model((*model.Meigen)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
}