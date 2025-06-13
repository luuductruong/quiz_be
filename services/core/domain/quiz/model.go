package quiz

type Quiz struct {
	QuizID        string `json:"quiz_id"`
	Content       string `json:"content"`
	QuizQuestions []*QuizQuestion
}

type Question struct {
	QuestionID    string    `json:"question_id"`
	Content       string    `json:"content"`
	Answers       []*Answer `json:"answers"`
	CorrectAnswer string
}

type QuizQuestion struct {
	QuizID     string
	QuestionID string
	Question   *Question
}

type Answer struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
