package domain

import (
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/logger"
)

type QuizDomainParam struct {
	QuizRepo         quiz.QuizRepo
	QuestionRepo     quiz.QuestionRepo
	QuizQuestionRepo quiz.QuizQuestionRepo
	UserRepo         quiz.UserRepo
	UserAnswerRepo   quiz.UserAnswerRepo
	ScoreRepo        quiz.ScoreRepo
}

type domain struct {
	logger           logger.Logger
	quizRepo         quiz.QuizRepo
	questionRepo     quiz.QuestionRepo
	quizQuestionRepo quiz.QuizQuestionRepo
	userRepo         quiz.UserRepo
	userAnswerRepo   quiz.UserAnswerRepo
	scoreRepo        quiz.ScoreRepo
}

func NewDomain(param *QuizDomainParam) quiz.Service {
	return &domain{
		logger:           logger.Default,
		quizRepo:         param.QuizRepo,
		questionRepo:     param.QuestionRepo,
		quizQuestionRepo: param.QuizQuestionRepo,
		userRepo:         param.UserRepo,
		userAnswerRepo:   param.UserAnswerRepo,
		scoreRepo:        param.ScoreRepo,
	}
}
