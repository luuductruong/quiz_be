package quiz

import (
	"github.com/quiz_be/services/core/context"
)

type QuizRepo interface {
	Upsert(ctx context.Context, quiz *Quiz) error
	Query(ctx context.Context) QuizQuery
}

type QuizQuery interface {
	// query
	ByQuizID(quizID string) QuizQuery
	Limit(limit int) QuizQuery
	Offset(offset int) QuizQuery
	// load
	WithQuizQuestion(questionID string) QuizQuery
	// result
	Result() (*Quiz, error)
	ResultList() ([]*Quiz, error)
	Count() (int, error)
}

type QuestionRepo interface {
	Upsert(ctx context.Context, question *Question) error
	Query(ctx context.Context) QuestionQuery
}

type QuestionQuery interface {
	// query
	ByQuestionID(questionID string) QuestionQuery
	ByQuestionIDs(questionIDs []string) QuestionQuery
	Limit(limit int) QuestionQuery
	Offset(offset int) QuestionQuery
	// result
	OrderByUpdatedAt(desc bool) QuestionQuery
	Result() (*Question, error)
	ResultList() ([]*Question, error)
	Count() (int, error)
}

type QuizQuestionRepo interface {
	BulkUpsert(ctx context.Context, quizQuestion []*QuestionQuiz) error
	Query(ctx context.Context) QuizQuestionQuery
}

type QuizQuestionQuery interface {
	// query
	ByQuizID(quizID string) QuizQuestionQuery
	ByQuestionID(questionID string) QuizQuestionQuery
	Limit(limit int) QuizQuestionQuery
	// result
	Result() (*QuestionQuiz, error)
	ResultList() ([]*QuestionQuiz, error)
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
	ByListQuestionID(questionIDs ...string) UserAnswerQuery
	// result
	Result() (*UserAnswer, error)
	ResultList() ([]*UserAnswer, error)
	Delete() error
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
	Delete() error
}
