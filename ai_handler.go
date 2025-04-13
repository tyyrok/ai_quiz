package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


func getDataFromAi(dbpool *pgxpool.Pool) {
	ok, _ := checkTodayQuestionInDB(dbpool)
	if ok {
		log.Println("Found questions for today in db")
		return
	}
	api_key := os.Getenv("G_API_KEY")
	raw_url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="
	url := fmt.Sprintf("%s%s", raw_url, api_key)
	cl := &http.Client{
		Timeout: 120 * time.Second,
	}
	payload := PayloadHead{
		Content: []ContentInfo{
			{
				Parts: []PartInfo{
					{
						Text: "Generate exactly 4 distinct IT development multiple-choice questions suitable for typical junior and middle level developers. Try to be creative and try to bind current datetime to your questions. Structure the entire output as a single JSON object with a top-level key named `quiz_data`. The value of `quiz_data` must be an array. Each element in this array represents one question and must be an object containing two keys: `question_text` (string containing the question text) and `answers` (an array). The `answers` array must contain exactly 4 answer objects. Each answer object must have two keys: `answer_text` (string containing the answer choice) and `is_correct` (a boolean value: `true` or `false`). For each question, exactly one answer object in its `answers` array must have `is_correct` set to `true`, while the other two must have `is_correct` set to `false`. Ensure the output is only the requested JSON object.",
					},
				},
			},
		},
		GenerationConfig: GenerationConfigInfo{
			Temperature: 0.7,
			MaxOutputTokens: 1500,
			ResponseMimeType: "application/json",
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := cl.Post(
		url, "application/json", bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Printf("Problem with request, status code: %d", resp.StatusCode)
		return
	}
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var partialRes PartialResponse
	if err = json.Unmarshal(body, &partialRes); err != nil {
		log.Println(err)
		return
	}
	var quizData QuizDataResponse
	if err = json.Unmarshal([]byte(partialRes.Candidates[0].Content.Parts[0].Text), &quizData); err != nil {
		log.Println(err)
		return
	}
	log.Println(quizData)
	saveQuestionsToDB(dbpool, &quizData)
}

func checkTodayQuestionInDB(dbpool *pgxpool.Pool) (bool, error) {
	var id int
	err := dbpool.QueryRow(
		context.Background(),
		`SELECT questions.id FROM questions WHERE questions.date = $1;`,
		time.Now().Format("2006-01-02")).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		log.Printf("Failed to retrieve data %s", err)
		return false, err
	} else {
		return true, nil
	}
}

func saveQuestionsToDB(dbpool *pgxpool.Pool, quizData *QuizDataResponse) {
	ctx := context.Background()
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		log.Printf("Problem with pool connection %s", err)
		return
	}
	defer tx.Rollback(ctx)
	for _, q := range quizData.QuizData {
		var q_id int
		err := tx.QueryRow(
			ctx,
			`INSERT INTO questions (title) VALUES ($1) RETURNING id`,
			q.QuestionText,
		).Scan(&q_id)
		if err != nil {
			log.Printf("Error with adding question %s\n", err)
			continue
		}
		for _, a := range q.Answers {
			_, err := tx.Exec(
				ctx,
				`INSERT INTO answers (title, question_id, is_correct) VALUES ($1, $2, $3) RETURNING id`,
				a.AnswerText, q_id, a.IsCorrect,
			)
			if err != nil {
				log.Printf("Error with adding answer %s\n", err)
				continue
			}
		}
	}
	tx.Commit(ctx)
}