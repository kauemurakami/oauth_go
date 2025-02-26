package db

import (
	"context"
	"fmt"
	"log"
	app_config "oauth2/core/config"
	"os"

	"github.com/jackc/pgx/v5"
)

func SetupDB() *pgx.Conn {
	connStr := app_config.Connection_string
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Database connection established")
	return db
}
