package routes

import (
	"context"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Answer struct {
	Id int `json:"id"`
	Title string `json:"text"`
	Likes int `json:"likes"`
	Users_answered int `json:"users_answered"`
}

type Question struct {
	Id int `json:"id"`
	Title string `json:"text"`
	Likes int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Answers []Answer `json:"answers"`
}


func Run(dbpool *pgxpool.Pool) {
	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", dbpool)
		ctx.Next()
	})
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func (ctx *gin.Context)  {
		dbctx, ok := ctx.Get("db")
		if !ok {
			log.Print("Failed to connect to db")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
			return
		}
		pool, ok := dbctx.(*pgxpool.Pool)
		if !ok {
			log.Print("Failed to connect to db")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
			return
		}
		questions, err := getTodayQuestions(pool)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
			return
		}
		questionsJSON, err := json.Marshal(questions)
		if err != nil {
			log.Println("Error marshalling questions:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{"Name": "Gin Framework", "questions": string(questionsJSON)})
	})

	router.Run(":8080")
}

func getTodayQuestions(pool *pgxpool.Pool) ([]Question, error) {
	rows, err := pool.Query(context.Background(), `SELECT questions.id, questions.title, questions.likes, questions.dislikes, answers.id, answers.title, answers.likes, answers.users_answered FROM questions INNER JOIN answers ON answers.question_id = questions.id WHERE questions.date = $1;`, "2025-04-04")
	if err != nil {
		log.Printf("Failed to retrieve data %s", err)
		return nil, err
	}
	defer rows.Close()
	var questions []Question
	for rows.Next() {
		var q Question
		var a Answer
		is_new := true
		err := rows.Scan(
			&q.Id, &q.Title, &q.Likes, &q.Dislikes, &a.Id, &a.Title, &a.Likes, &a.Users_answered)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		for i, elem := range questions {
			if elem.Id == q.Id {
				questions[i].Answers = append(questions[i].Answers, a)
				is_new = false
				break
			}
		}
		if is_new {
			q.Answers = append(q.Answers, a)
			questions = append(questions, q)
		}
	}
	return questions, nil
}