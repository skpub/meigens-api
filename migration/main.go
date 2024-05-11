package main

import (
	"fmt"
	"os"
	"meigens-api/src/model"
	"meigens-api/src/db"
	"context"
	"github.com/uptrace/bun"
)

func main() {
	db, err := db.Conn()
	if err != nil {
		panic("failed to connect db.")
	}
	defer db.Close()

	models := []interface{}{
		&model.Users{},
		&model.Groups{},
		&model.Poets{},
		&model.Meigens{},
		&model.Relationships{},
		&model.Reactions{},
		&model.GroupPoets{},
		&model.UserGroupRels{},
		&model.UserPoets{},
	}

	switch arg := os.Args[1]; arg {
		case "migration":
			fmt.Println("execute migration.")
			destroy(db, models)
			initialization(db, models)

		case "destroy":
			fmt.Println("DESTROY")
			destroy(db, models)

		default:
			fmt.Println("specify the migration command.")
	}
}

func initialization(db *bun.DB, models []interface{}) {
	ctx := context.Background()
	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func destroy(db *bun.DB, models []interface{}) {
	ctx := context.Background()

	for _, model := range models {
		_, err := db.NewDropTable().
			Model(model).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {

			panic(err)
		}
	}
}
