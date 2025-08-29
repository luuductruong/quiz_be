package grpchandler

import (
	"context"
	"github.com/quiz_be/services/core/application/quiz/model"
	"github.com/quiz_be/services/core/helper"

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
func (h *handler) ManageQuiz(ctx context.Context, req *dto.ManageQuizReq) (*model.Quiz, error) {
	appCtx := appContext.FromContext(ctx)
	quiz, err := h.quizDomain.ManageQuiz(appCtx, &quiz.ManageQuizReq{
		QuizID:      req.QuizId,
		Title:       req.Title,
		QuestionIDs: req.QuestionIds,
	})
	if err != nil {
		return nil, err
	}
	return dto.MapQuizFromDomain(quiz), nil
}

func (h *handler) ManageQuestion(ctx context.Context, req *dto.ManageQuestionReq) (*model.Question, error) {
	appCtx := appContext.FromContext(ctx)
	question, err := h.quizDomain.ManageQuestion(appCtx, &quiz.ManageQuestionReq{
		QuestionID:    req.QuestionId,
		Content:       req.Content,
		Score:         uint(req.Score),
		CorrectAnswer: req.CorrectAnswer,
		Answers:       helper.MapList(req.Answers, dto.MapAnswerToDomain),
	})
	if err != nil {
		return nil, err
	}
	return dto.AdminMapQuestionFromDomain(question), nil
}

func (h *handler) SubmitAnswer(ctx context.Context, req *dto.SubmitAnswerReq) (*dto.SubmitAnswerRes, error) {
	appCtx := appContext.FromContext(ctx)
	q, err := h.quizDomain.SubmitAnswer(appCtx, &quiz.SubmitAnswerReq{
		UserID: req.UserId,
		QuizID: req.QuizId,
		SelectAnswers: helper.MapList(req.SelectAnswers, func(s *dto.SelectedQuestionAnswer) *quiz.SelectedQuestionAnswer {
			return &quiz.SelectedQuestionAnswer{
				QuestionID:  s.QuestionId,
				AnswerTitle: s.AnswerTitle,
			}
		}),
	})
	if err != nil {
		return nil, err
	}
	return &dto.SubmitAnswerRes{
		Quiz: dto.MapQuizFromDomain(q),
	}, nil
}

func (h *handler) GetLeaderboard(ctx context.Context, req *dto.GetLeaderboardReq) (*dto.GetLeaderboardRes, error) {
	appCtx := appContext.FromContext(ctx)
	page, limit := appCtx.GetPageAndLimit()
	listScore, total, err := h.quizDomain.GetLeaderboard(appCtx, &quiz.GetLeaderboardReq{
		QuizID: req.QuizId,
		Limit:  limit,
		Page:   page,
	})
	if err != nil {
		return nil, err
	}
	return &dto.GetLeaderboardRes{
		Leaderboard: helper.MapList(listScore, dto.MapScoreFromDomain),
		Page:        int32(page),
		Total:       total,
	}, nil
}

func (h *handler) GetListQuiz(ctx context.Context, req *dto.GetListQuizReq) (*dto.GetListQuizRes, error) {
	appCtx := appContext.FromContext(ctx)
	page, limit := appCtx.GetPageAndLimit()
	listQuiz, total, err := h.quizDomain.GetListQuiz(appCtx, page, limit)
	if err != nil {
		return nil, err
	}
	return &dto.GetListQuizRes{
		ListQuiz: helper.MapList(listQuiz, dto.MapQuizFromDomain),
		Page:     int32(page),
		Total:    total,
	}, nil
}

func (h *handler) GetListQuestionByAdmin(ctx context.Context, req *dto.GetListQuestionReq) (*dto.GetListQuestionRes, error) {
	appCtx := appContext.FromContext(ctx)
	page, limit := appCtx.GetPageAndLimit()
	//TODO validate admin here
	listQuestion, total, err := h.quizDomain.GetListQuestion(appCtx, page, limit)
	if err != nil {
		return nil, err
	}
	return &dto.GetListQuestionRes{
		QuestionList: helper.MapList(listQuestion, dto.AdminMapQuestionFromDomain),
		Page:         int32(page),
		Total:        total,
	}, nil
}

func (h *handler) GetQuizDetail(ctx context.Context, req *dto.GetQuizDetailReq) (*model.Quiz, error) {
	appCtx := appContext.FromContext(ctx)
	q, err := h.quizDomain.GetQuizDetail(appCtx, req.QuizId)
	if err != nil {
		return nil, err
	}
	return dto.MapQuizFromDomain(q), nil
}

func (h *handler) GetUser(ctx context.Context, req *dto.GetUserReq) (*dto.GetUserRes, error) {
	appCtx := appContext.FromContext(ctx)
	listUser, err := h.quizDomain.GetUser(appCtx)
	if err != nil {
		return nil, err
	}
	return &dto.GetUserRes{
		ListUser: helper.MapList(listUser, dto.MapUserFromDomain),
	}, nil
}
