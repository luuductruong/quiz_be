package repository

import (
	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
	"time"
)

type userAnswerGorm struct {
	ID             string `gorm:"primary_key"`
	UserID         string
	QuizID         string
	QuestionID     string
	SelectedAnswer string
	IsCorrect      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	User           *userGorm     `gorm:"foreignKey:UserID;references:UserID"`
	Quiz           *quizGorm     `gorm:"foreignKey:QuizID;references:QuizID"`
	Question       *questionGrom `gorm:"foreignKey:QuestionID;references:QuestionID"`
}

func mapUserAnswerToDomain(source *userAnswerGorm) *quizDomain.UserAnswer {
	if source == nil {
		return nil
	}
	return &quizDomain.UserAnswer{
		ID:             source.ID,
		UserID:         source.UserID,
		QuizID:         source.QuizID,
		QuestionID:     source.QuestionID,
		SelectedAnswer: source.SelectedAnswer,
		IsCorrect:      source.IsCorrect,
		CreatedAt:      source.CreatedAt,
		UpdatedAt:      source.UpdatedAt,
		User:           mapUserToDomain(source.User),
		Quiz:           mapQuizToDomain(source.Quiz),
		Question:       mapQuestionToDomain(source.Question),
	}
}

func mapUserAnswerFromDomain(source *quizDomain.UserAnswer) *userAnswerGorm {
	if source == nil {
		return nil
	}
	return &userAnswerGorm{
		ID:             source.ID,
		UserID:         source.UserID,
		QuizID:         source.QuizID,
		QuestionID:     source.QuestionID,
		SelectedAnswer: source.SelectedAnswer,
		IsCorrect:      source.IsCorrect,
		CreatedAt:      source.CreatedAt,
		UpdatedAt:      source.UpdatedAt,
	}
}

func (ua userAnswerGorm) TableName() string {
	return "user_answer"
}

func NewUserAnswerRepo() quizDomain.UserAnswerRepo {
	return &userAnswerRepo{}
}

type userAnswerRepo struct {
}

type userAnswerQuery struct {
	query.BaseQuery
}

func (ua *userAnswerRepo) Query(ctx context.Context) quizDomain.UserAnswerQuery {
	return &userAnswerQuery{query.NewBQ(ctx.GetDbTx().Model(&userAnswerGorm{}))}
}

func (ua *userAnswerRepo) Upsert(ctx context.Context, usAnswer *quizDomain.UserAnswer) error {
	return query.Upsert(ctx.GetDbTx(), usAnswer, mapUserAnswerFromDomain)
}

func (ua *userAnswerQuery) ByID(id string) quizDomain.UserAnswerQuery {
	return query.Where(ua, "id = ?", id)
}

func (ua *userAnswerQuery) ByUserID(userID string) quizDomain.UserAnswerQuery {
	return query.Where(ua, "user_id = ?", userID)
}

func (ua *userAnswerQuery) ByQuizID(quizID string) quizDomain.UserAnswerQuery {
	return query.Where(ua, "quiz_id = ?", quizID)
}

func (ua *userAnswerQuery) ByQuestionID(questionID string) quizDomain.UserAnswerQuery {
	return query.Where(ua, "question_id = ?", questionID)
}

func (ua *userAnswerQuery) ByListQuestionID(questionIDs ...string) quizDomain.UserAnswerQuery {
	return query.Where(ua, "question_id in ?", questionIDs)
}

func (ua *userAnswerQuery) Result() (*quizDomain.UserAnswer, error) {
	return query.Result(ua, mapUserAnswerToDomain)
}

func (ua *userAnswerQuery) ResultList() ([]*quizDomain.UserAnswer, error) {
	return query.ResultList(ua, mapUserAnswerToDomain)
}

func (ua *userAnswerQuery) Delete() error {
	return query.Delete(ua)
}
