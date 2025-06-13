package domain

import (
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/logger"
)

type QuizDomainParam struct {
	QuizRepo         quiz.QuizRepo
	QuestionRepo     quiz.QuestionRepo
	QuizQuestionRepo quiz.QuizQuestionRepo
}

type domain struct {
	logger           logger.Logger
	quizRepo         quiz.QuizRepo
	questionRepo     quiz.QuestionRepo
	quizQuestionRepo quiz.QuizQuestionRepo
}

func NewDomain(param *QuizDomainParam) quiz.Service {
	return &domain{
		logger:           logger.Default,
		quizRepo:         param.QuizRepo,
		questionRepo:     param.QuestionRepo,
		quizQuestionRepo: param.QuizQuestionRepo,
	}
}
