package quiz

import "github.com/quiz_be/services/core/context"

type Service interface {
	SubmitAnswer(ctx context.Context, inp *SubmitAnswerReq) (*Quiz, error)
}
