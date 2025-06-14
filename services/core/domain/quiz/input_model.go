package quiz

type SubmitAnswerReq struct {
	UserID      string `json:"user_id"`
	QuizID      string `json:"quiz_id"`
	QuestionID  string `json:"question_id"`
	AnswerTitle string `json:"answer_title"`
}

type GetLeaderboardReq struct {
	QuizID string `json:"quiz_id"`
}
