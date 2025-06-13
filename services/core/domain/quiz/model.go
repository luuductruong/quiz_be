package quiz

import "time"

type Quiz struct {
	QuizID        string `json:"quiz_id"`
	Title         string `json:"title"`
	QuizQuestions []*QuizQuestion
}

type Question struct {
	QuestionID    string    `json:"question_id"`
	Content       string    `json:"content"`
	Score         uint      `json:"score"`
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

type User struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
}

type Score struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	QuizID    string    `json:"quiz_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User
	Quiz      *Quiz
}

type UserAnswer struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	QuizID         string    `json:"quiz_id"`
	QuestionID     string    `json:"question_id"`
	SelectedAnswer string    `json:"selected_answer"`
	IsCorrect      bool      `json:"is_correct"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           *User
	Quiz           *Quiz
	Question       *Question
}
