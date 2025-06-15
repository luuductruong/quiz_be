package quiz

import "github.com/quiz_be/services/core/context"

type Service interface {
	SubmitAnswer(ctx context.Context, inp *SubmitAnswerReq) (*Quiz, error)
	GetLeaderboard(ctx context.Context, inp *GetLeaderboardReq) ([]*Score, int32, error)
	GetListQuiz(ctx context.Context, page int, limit int) ([]*Quiz, int32, error)
	GetQuizDetail(ctx context.Context, quizID string) (*Quiz, error)
	GetUser(ctx context.Context) ([]*User, error)
}
