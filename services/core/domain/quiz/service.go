package quiz

import (
	"github.com/quiz_be/services/core/context"
)

const (
	QuizTopic            = "quiz"
	PushJobGetQuizDetail = "getQuizDetail"
)

type Service interface {
	// admin
	ManageQuiz(ctx context.Context, inp *ManageQuizReq) (*Quiz, error)
	ManageQuestion(ctx context.Context, inp *ManageQuestionReq) (*Question, error)
	GetListQuestion(ctx context.Context, page int, limit int) ([]*Question, int32, error)
	// user
	SubmitAnswer(ctx context.Context, inp *SubmitAnswerReq) (*Quiz, error)
	GetLeaderboard(ctx context.Context, inp *GetLeaderboardReq) ([]*Score, int32, error)
	GetListQuiz(ctx context.Context, page int, limit int) ([]*Quiz, int32, error)
	GetQuizDetail(ctx context.Context, quizID string) (*Quiz, error)
	GetUser(ctx context.Context) ([]*User, error)
}
