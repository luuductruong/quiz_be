package grpchandler

import (
	"context"

	"github.com/quiz_be/services/core/application/quiz/dto"
	"github.com/quiz_be/services/core/application/quiz/service"
	appContext "github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/domain/quiz"
)

type handler struct {
	service.UnimplementedQuizServiceServer
	quizDomain quiz.Service
}

func NewHandler(quizDomain quiz.Service) service.QuizServiceServer {
	return &handler{
		quizDomain: quizDomain,
	}
}

func (h *handler) SubmitAnswer(ctx context.Context, req *dto.SubmitAnswerReq) (*dto.SubmitAnswerRes, error) {
	appCtx := appContext.FromContext(ctx)
	q, err := h.quizDomain.SubmitAnswer(appCtx, &quiz.SubmitAnswerReq{
		UserID:      req.UserId,
		QuizID:      req.QuizId,
		QuestionID:  req.QuestionId,
		AnswerTitle: req.AnswerTitle,
	})
	if err != nil {
		return nil, err
	}
	return &dto.SubmitAnswerRes{
		Quiz: dto.MapQuizFromDomain(q),
	}, nil
}

//func (h *handler) GetLeaderboard(ctx context.Context, req *dto.GetLeaderboardReq) (*dto.GetLeaderboardRes, error) {
//	appCtx := appContext.FromContext(ctx)
//	q, err := h.quizDomain.SubmitAnswer(appCtx, &quiz.SubmitAnswerReq{
//		UserID:      req.UserId,
//		QuizID:      req.QuizId,
//		QuestionID:  req.QuestionId,
//		AnswerTitle: req.AnswerTitle,
//	})
//	if err != nil {
//		return nil, err
//	}
//	return &dto.SubmitAnswerRes{
//		Quiz: dto.MapQuizFromDomain(q),
//	}, nil
//}
