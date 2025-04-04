package db

import (
	"log"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)


func InitDB(dbpool *pgxpool.Pool) {
	_, err := dbpool.Exec(
			context.Background(), createQuestionsTable,
		)
	if err != nil {
		log.Fatalf("Error with accessing/creating table %s\n", err)
	}
	_, err = dbpool.Exec(
		context.Background(), createAnswersTable,
	)
	if err != nil {
		log.Fatalf("Error with accessing/creating table %s\n", err)
	}
	// Test data
	/*
	var question_id int
	err = dbpool.QueryRow(
		context.Background(), `INSERT INTO questions (title) VALUES ($1) RETURNING id`, "What animal is the fastest?",
	).Scan(&question_id)
	if err != nil {
		log.Fatalf("Error with adding question %s\n", err)
	}
	answers := []string{"Elephant", "Turtle", "Dog", "Gepard"}
	for i := 0; i < 4; i++ {
		_, err = dbpool.Exec(
			context.Background(),
			InsertAnswer,
			answers[i], question_id,
		)
		if err != nil {
			log.Fatalf("Error with adding answers %s\n", err)
		}
	}
	*/
}
