package quiz

import (
	"github.com/quiz_be/services/core/context"
)

type QuizRepo interface {
	Query(ctx context.Context) QuizQuery
}

type QuizQuery interface {
	// query
	ByQuizID(quizID string) QuizQuery
	Limit(limit int) QuizQuery
	// result
	Result() (*Quiz, error)
	ResultList() ([]*Quiz, error)
}

type QuestionRepo interface {
	Query(ctx context.Context) QuestionQuery
}

type QuestionQuery interface {
	// query
	ByQuestionID(questionID string) QuestionQuery
	// result
	Result() (*Question, error)
}
