package main

import (
	"ai_quiz/db"
	"ai_quiz/routes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	encodedPassword := os.Getenv("DB_PASSWORD")
	dbUrl := os.Getenv("DB_URL")
	
	connString := fmt.Sprintf(dbUrl, url.QueryEscape(encodedPassword))
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}
	defer dbpool.Close()

	db.InitDB(dbpool)
	routes.Run(dbpool)
}
