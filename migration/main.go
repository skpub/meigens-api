package main

import (
	"fmt"
	"os"
	"meigens-api/src/model"
	"meigens-api/src/db"
	"context"
)

func main() {
	switch arg := os.Args[1]; arg {
		case "migration":
			fmt.Println("execute migration.")

			db, err := db.Conn()
			if err != nil {
				panic("failed to connect db.")
			}
			defer db.Close()

			models := []interface{}{
				&model.Meigens{},
				&model.Users{},
				&model.Relationships{},
				&model.Reactions{},
			}

			ctx := context.Background()
			for _, model := range models {
				_, err = db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
				if err != nil {
					panic(err)
				}
			}

		default:
			fmt.Println("specify the migration command.")
	}
}
