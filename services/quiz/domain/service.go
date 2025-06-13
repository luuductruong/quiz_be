package domain

import (
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/logger"
)

type QuizDomainParam struct {
	QuizRepo     quiz.QuizRepo
	QuestionRepo quiz.QuestionRepo
}

type domain struct {
	logger       logger.Logger
	quizRepo     quiz.QuizRepo
	questionRepo quiz.QuestionRepo
}

func NewDomain(param *QuizDomainParam) quiz.Service {
	return &domain{
		logger:       logger.Default,
		quizRepo:     param.QuizRepo,
		questionRepo: param.QuestionRepo,
	}
}
