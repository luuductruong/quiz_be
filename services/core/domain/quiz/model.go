package quiz

type Quiz struct {
	QuizID    string      `json:"quiz_id"`
	Content   string      `json:"content"`
	Questions []*Question `json:"questions"`
}

type Question struct {
	QuestionID    string    `json:"question_id"`
	Content       string    `json:"content"`
	Answers       []*Answer `json:"answers"`
	CorrectAnswer string
}

type Answer struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
