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

type UserRepo interface {
	Query(ctx context.Context) UserQuery
}

type UserQuery interface {
	// query
	ByUserID(userID string) UserQuery
	ByPhoneNumber(phoneNumber string) UserQuery
	// result
	Result() (*User, error)
	ResultList() ([]*User, error)
}

type UserAnswerRepo interface {
	Upsert(ctx context.Context, usAnswer *UserAnswer) error
	Query(ctx context.Context) UserAnswerQuery
}

type UserAnswerQuery interface {
	// query
	ByID(id string) UserAnswerQuery
	ByUserID(userID string) UserAnswerQuery
	ByQuizID(quizID string) UserAnswerQuery
	ByQuestionID(questionID string) UserAnswerQuery
	// result
	Result() (*UserAnswer, error)
	ResultList() ([]*UserAnswer, error)
}

type ScoreRepo interface {
	Upsert(ctx context.Context, score *Score) error
	Query(ctx context.Context) ScoreQuery
}

type ScoreQuery interface {
	// query
	ByID(id string) ScoreQuery
	ByUserID(userID string) ScoreQuery
	ByQuizID(quizID string) ScoreQuery
	WithUser() ScoreQuery
	// paging
	Limit(limit int) ScoreQuery
	Offset(offset int) ScoreQuery
	// ordering
	OrderByScore(desc bool) ScoreQuery
	// result
	Result() (*Score, error)
	ResultList() ([]*Score, error)
	Count() (int, error)
}
