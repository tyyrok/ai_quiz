package routes

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func setCookie(ctx *gin.Context, quizState *Cookie) {
	origin := os.Getenv("ORIGIN")
	jsonStr, err := json.Marshal(quizState)
	if err != nil {
		log.Printf("Error with marshalling cookie struct %s", err)
		return
	}
	ctx.SetCookie("quiz_state", string(jsonStr), 60*60*24, "/", origin, false, false)
}


func handleCookie(ctx *gin.Context, quizState *QuizState) {
	userCookie, err := ctx.Cookie("quiz_state")
	if err != nil {
		switch {
		case quizState.QuestionAnsweredId != 0:
			setCookie(ctx, &Cookie{Answered: []int{quizState.QuestionAnsweredId}})
		case quizState.QuestionLikedId != 0:
			setCookie(ctx, &Cookie{LikedQuestions: []int{quizState.QuestionLikedId}})
		case quizState.QuestionDislikedId != 0:
			setCookie(ctx, &Cookie{DislikedQuestions: []int{quizState.QuestionDislikedId}})
		case quizState.AnswerLikedId != 0:
			setCookie(ctx, &Cookie{LikedAnswers: []int{quizState.AnswerLikedId}})
		case quizState.AnswerDislikedId != 0:
			setCookie(ctx, &Cookie{DislikedAnswers: []int{quizState.AnswerDislikedId}})
		default:
			log.Printf("Error with setting cookie %s", err)
			return
		}
	}
	var cookie Cookie
	err = json.Unmarshal([]byte(userCookie), &cookie)
	if err != nil {
		log.Printf("Error with reading cookie %s", err)
		return
	}
	switch {
	case quizState.QuestionAnsweredId != 0:
		cookie.Answered = append(cookie.Answered, quizState.QuestionAnsweredId)
	case quizState.QuestionLikedId != 0:
		cookie.LikedQuestions = append(cookie.LikedQuestions, quizState.QuestionLikedId)
	case quizState.QuestionDislikedId != 0:
		cookie.DislikedQuestions = append(cookie.DislikedQuestions, quizState.QuestionDislikedId)
	case quizState.AnswerLikedId != 0:
		cookie.LikedAnswers = append(cookie.LikedAnswers, quizState.AnswerLikedId)
	case quizState.AnswerDislikedId != 0:
		cookie.DislikedAnswers = append(cookie.DislikedAnswers, quizState.AnswerDislikedId)
	default:
		log.Printf("Error with updating cookie %s", err)
		return
	}
	setCookie(ctx, &cookie)
}