package domain

import (
	"context"
	appCtx "github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/errors"
	repo "github.com/quiz_be/services/quiz/external/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mockQuizRepo         = repo.NewQuizRepo()
	mockQuestionRepo     = repo.NewQuestionRepo()
	mockQuizQuestionRepo = repo.NewQuizQuestionRepo()
	mockUserRepo         = repo.NewUserRepo()
	mockUserAnswerRepo   = repo.NewUserAnswerRepo()
	mockScoreRepo        = repo.NewScoreRepo()
	mockQuizDomain       = &QuizDomainParam{
		QuizRepo:         mockQuizRepo,
		QuestionRepo:     mockQuestionRepo,
		QuizQuestionRepo: mockQuizQuestionRepo,
		UserRepo:         mockUserRepo,
		UserAnswerRepo:   mockUserAnswerRepo,
		ScoreRepo:        mockScoreRepo,
	}
)

func TestSubmitAnswer_InvalidInput(t *testing.T) {
	d := NewDomain(mockQuizDomain) // domain với repo mock
	ctx := appCtx.FromContext(context.Background())

	tests := []struct {
		name string
		inp  *model.SubmitAnswerReq
	}{
		{"empty userID", &model.SubmitAnswerReq{UserID: "", QuizID: "q1", SelectAnswers: []*model.SelectedQuestionAnswer{{}}}},
		{"empty quizID", &model.SubmitAnswerReq{UserID: "u1", QuizID: "", SelectAnswers: []*model.SelectedQuestionAnswer{{}}}},
		{"empty answers", &model.SubmitAnswerReq{UserID: "u1", QuizID: "q1", SelectAnswers: []*model.SelectedQuestionAnswer{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quiz, err := d.SubmitAnswer(ctx, tt.inp)
			assert.Nil(t, quiz)
			assert.Error(t, err)
			assert.Equal(t, errors.Code(err), errors.InvalidArgumentCode())
		})
	}
}
