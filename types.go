package main

type PartInfo struct {
	Text string `json:"text"`
}

type ContentInfo struct {
	Parts []PartInfo `json:"parts"`
}

type GenerationConfigInfo struct {
	Temperature float32 `json:"temperature"`
	MaxOutputTokens int `json:"maxOutputTokens"`
	ResponseMimeType string `json:"responseMimeType"`
}

type PayloadHead struct {
	Content []ContentInfo `json:"contents"`
	GenerationConfig GenerationConfigInfo `json:"generationConfig"`
}

type AnswerResponse struct {
	AnswerText string `json:"answer_text"`
	IsCorrect bool `json:"is_correct"`
}

type QuestionResponse struct {
	QuestionText string `json:"question_text"`
	Answers []AnswerResponse `json:"answers"`
}

type QuizDataResponse struct {
	QuizData []QuestionResponse `json:"quiz_data"`
}

type PartResponse struct {
	Text string  `json:"text"`
}

type ContentResponseInfo struct {
	Parts []PartResponse `json:"parts"`
}

type Candidate struct {
	Content ContentResponseInfo `json:"content"`
}

type PartialResponse struct {
	Candidates []Candidate `json:"candidates"`
}