package quiz

type SubmitAnswerReq struct {
	UserID        string                    `json:"user_id"`
	QuizID        string                    `json:"quiz_id"`
	SelectAnswers []*SelectedQuestionAnswer `json:"select_answers"`
}

type SelectedQuestionAnswer struct {
	QuestionID  string `json:"question_id"`
	AnswerTitle string `json:"answer_title"`
}

type GetLeaderboardReq struct {
	QuizID string `json:"quiz_id"`
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
}
