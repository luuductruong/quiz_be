package domain

import (
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
)

type QuizDomainParam struct {
	QuizRepo         quiz.QuizRepo
	QuestionRepo     quiz.QuestionRepo
	QuizQuestionRepo quiz.QuizQuestionRepo
	UserRepo         quiz.UserRepo
	UserAnswerRepo   quiz.UserAnswerRepo
	ScoreRepo        quiz.ScoreRepo
	Publisher        pubsub.Publisher
}

type domain struct {
	logger           logger.Logger
	quizRepo         quiz.QuizRepo
	questionRepo     quiz.QuestionRepo
	quizQuestionRepo quiz.QuizQuestionRepo
	userRepo         quiz.UserRepo
	userAnswerRepo   quiz.UserAnswerRepo
	scoreRepo        quiz.ScoreRepo
	publisher        pubsub.Publisher
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
		publisher:        param.Publisher,
	}
}
