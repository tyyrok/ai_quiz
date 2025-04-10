package main

import (
	"ai_quiz/db"
	"ai_quiz/routes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Println("Error loading .env file")
	}
	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	encodedPassword := os.Getenv("DB_PASSWORD")
	dbUrl := os.Getenv("DB_URL")
	connString := fmt.Sprintf(dbUrl, url.QueryEscape(encodedPassword))
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}
	defer dbpool.Close()

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error while initializing gocron %s", err)
	}
	defer func() { _ = s.Shutdown() }()
	_, err = s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(17, 50, 0),
				gocron.NewAtTime(9, 06, 0),
			),
		),
		gocron.NewTask(
			getDataFromAi,
			dbpool,
		),
	)	
	if err != nil {
		log.Fatalf("Error while gocron job %s", err)
	}
	s.Start()

	db.InitDB(dbpool)
	routes.Run(dbpool)
}
