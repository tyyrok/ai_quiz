package routes

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func mainPageHandler(ctx *gin.Context) {
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
	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{"Name": "Gin Framework", "questions": template.JS(questionsJSON)})
}

func answerHandler(ctx *gin.Context) {
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
	question_id := ctx.Params.ByName("question_id")
	answer_id := ctx.Params.ByName("answer_id")
	if question_id == " " || answer_id == " " {
		log.Print("Not found question_id/answer_id")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	correctAnswer, err := answerQuestion(pool, question_id, answer_id)
	if err != nil {
		log.Printf("Error while processing %s\n", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	answerJSON, err := json.Marshal(correctAnswer)
	if err != nil {
		log.Println("Error marshalling questions:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	q_id, _ := strconv.Atoi(question_id)
	handleCookie(ctx, &QuizState{QuestionAnsweredId: q_id})
	ctx.JSON(http.StatusOK, gin.H{"correctAnswer": string(answerJSON)})
}

func questionLikeHandler(ctx *gin.Context) {
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
	question_id := ctx.Params.ByName("question_id")
	if question_id == " " {
		log.Print("Not found question_id")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	q_id, err := strconv.Atoi(question_id)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	_, is_dislike := ctx.GetQuery("is_dislike")
	tx, err := pool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
		return
	}
	defer tx.Rollback(ctx)

	if is_dislike {
		row := tx.QueryRow(
			ctx,
			`SELECT dislikes FROM questions WHERE id=$1 FOR UPDATE;`,
			q_id,
		)
		var dislikes int
		err = row.Scan(&dislikes)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		_, err = tx.Exec(
			ctx,
			`UPDATE questions SET dislikes=$1 WHERE id=$2;`,
			dislikes+1, q_id,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		tx.Commit(ctx)
		handleCookie(ctx, &QuizState{QuestionDislikedId: q_id})
		ctx.JSON(http.StatusOK, gin.H{"dislikes": dislikes+1})
	} else {
		row := tx.QueryRow(
			ctx,
			`SELECT likes FROM questions WHERE id=$1 FOR UPDATE;`,
			q_id,
		)
		var likes int
		err = row.Scan(&likes)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		_, err = tx.Exec(
			ctx,
			`UPDATE questions SET likes=$1 WHERE id=$2;`,
			likes+1, q_id,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		tx.Commit(ctx)
		handleCookie(ctx, &QuizState{QuestionLikedId: q_id})
		ctx.JSON(http.StatusOK, gin.H{"likes": likes+1})
	}
}

func answerLikeHandler(ctx *gin.Context) {
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
	question_id := ctx.Params.ByName("question_id")
	answer_id := ctx.Params.ByName("answer_id")
	if question_id == " " || answer_id == " "{
		log.Print("Not found question_id or answer_id")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	_, err := strconv.Atoi(question_id)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	a_id, err := strconv.Atoi(answer_id)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong request"})
		return
	}
	_, is_dislike := ctx.GetQuery("is_dislike")
	tx, err := pool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
		return
	}
	defer tx.Rollback(ctx)

	if is_dislike {
		row := tx.QueryRow(
			ctx,
			`SELECT dislikes FROM answers WHERE id=$1 FOR UPDATE;`,
			a_id,
		)
		var dislikes int
		err = row.Scan(&dislikes)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		_, err = tx.Exec(
			ctx,
			`UPDATE answers SET dislikes=$1 WHERE id=$2;`,
			dislikes+1, a_id,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		tx.Commit(ctx)
		handleCookie(ctx, &QuizState{AnswerDislikedId: a_id})
		ctx.JSON(http.StatusOK, gin.H{"dislikes": dislikes+1})
	} else {
		row := tx.QueryRow(
			ctx,
			`SELECT likes FROM answers WHERE id=$1 FOR UPDATE;`,
			a_id,
		)
		var likes int
		err = row.Scan(&likes)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		_, err = tx.Exec(
			ctx,
			`UPDATE answers SET likes=$1 WHERE id=$2;`,
			likes+1, a_id,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Smt wrong"})
			return
		}
		tx.Commit(ctx)
		handleCookie(ctx, &QuizState{AnswerLikedId: a_id})
		ctx.JSON(http.StatusOK, gin.H{"likes": likes+1})
	}
}

func getTodayQuestions(pool *pgxpool.Pool) ([]Question, error) {
	rows, err := pool.Query(
		context.Background(),
		`SELECT questions.id, questions.title, questions.likes, questions.dislikes, answers.id, answers.title, answers.likes, answers.dislikes, answers.users_answered FROM questions INNER JOIN answers ON answers.question_id = questions.id WHERE questions.date = $1;`,
		time.Now().Format("2006-01-02"))
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
			&q.Id, &q.Title, &q.Likes, &q.Dislikes, &a.Id, &a.Title, &a.Likes, &a.Dislikes, &a.Users_answered)
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

func answerQuestion(pool *pgxpool.Pool, question_id, answer_id string) (Answer, error) {
	q_id, err := strconv.Atoi(question_id)
	if err != nil {
		return Answer{}, err
	}
	a_id, err := strconv.Atoi(answer_id)
	if err != nil {
		return Answer{}, err
	}
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Answer{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(
		ctx,
		`SELECT id, users_answered from answers WHERE id=$1 AND question_id=$2 FOR UPDATE;`,
		a_id, q_id,
	)
	var id, users_answered int
	err = row.Scan(&id, &users_answered)
	if err != nil {
		return Answer{}, err
	}
	_, err = tx.Exec(
		ctx,
		`UPDATE answers SET users_answered = $1 WHERE id = $2;`,
		users_answered+1, id,
	)
	if err != nil {
		return Answer{}, err
	}
	tx.Commit(ctx)
	row = pool.QueryRow(
		ctx,
		`SELECT id, title, likes, users_answered FROM answers WHERE question_id = $1 AND is_correct = True`,
		q_id,
	)
	var answer Answer
	err = row.Scan(&answer.Id, &answer.Title, &answer.Likes, &answer.Users_answered)
	if err != nil {
		log.Printf("FUCK %s", err)
		return Answer{}, err
	}
	return answer, nil
}