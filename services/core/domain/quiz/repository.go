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
	// load
	WithQuizQuestion(questionID string) QuizQuery
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

type QuizQuestionRepo interface {
	Query(ctx context.Context) QuizQuestionQuery
}

type QuizQuestionQuery interface {
	// query
	ByQuizID(quizID string) QuizQuestionQuery
	ByQuestionID(questionID string) QuizQuestionQuery
	Limit(limit int) QuizQuestionQuery
	// result
	Result() (*QuizQuestion, error)
	ResultList() ([]*QuizQuestion, error)
}
