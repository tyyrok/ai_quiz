package routes

type Answer struct {
	Id int `json:"id"`
	Title string `json:"text"`
	Likes int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Users_answered int `json:"users_answered"`
}

type Question struct {
	Id int `json:"id"`
	Title string `json:"text"`
	Likes int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Answers []Answer `json:"answers"`
}

type Cookie struct {
	Answered []int `json:"answered"`
	LikedQuestions []int `json:"likedQuestions"`
	DislikedQuestions []int `json:"dislikedQuestions"`
	LikedAnswers []int `json:"likedAnswers"`
	DislikedAnswers []int `json:"dislikedAnswers"`
}

type QuizState struct {
	QuestionAnsweredId int
	QuestionLikedId int
	QuestionDislikedId int
	AnswerLikedId int
	AnswerDislikedId int
}