package app_config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Connection_string = ""
	DB_url            = ""
	API_port          = ""
	SECRET_KEY        []byte
)

func SetupEnvironments() {

	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	API_port = os.Getenv("API_PORT")
	if err != nil {
		API_port = "3000"
	}

	Connection_string =
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL"),
		)
	DB_url =
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL"),
		)

	SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
}
