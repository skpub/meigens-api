package dbconn

import (
	"database/sql"
	"fmt"
	"os"
	"log"

	// "github.com/uptrace/bun"
	// "github.com/uptrace/bun/dialect/pgdialect"
	// "github.com/uptrace/bun/driver/pgdriver"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Conn() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	dbname := os.Getenv("PG_DBNAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	// pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	// db := bun.NewDB(pgdb, pgdialect.New())

	return db, nil
}