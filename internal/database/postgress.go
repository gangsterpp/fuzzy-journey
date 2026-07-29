package database

import (
	"context"
	"log"
	"time"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPgPool(ctx context.Context, url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close()

	// Verify the connection is alive
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}
	// Maximum number of open connections to the database.
	db.SetMaxOpenConns(25)

	// Maximum number of idle connections retained in the pool.
	db.SetMaxIdleConns(10)

	// Maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(5 * time.Minute)

	// Maximum amount of time a connection may sit idle before being closed.
	db.SetConnMaxIdleTime(1 * time.Minute)

	return db

}
