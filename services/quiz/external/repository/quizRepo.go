package repository

import (
	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper"
	"github.com/quiz_be/services/core/helper/sql/query"
	"gorm.io/gorm"
	"time"
)

type quizGorm struct {
	QuizID        string `gorm:"primary_key"`
	Title         string
	QuizQuestions []*quizQuestionGorm `gorm:"foreignKey:QuizID;references:QuizID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func mapQuizToDomain(source *quizGorm) *quizDomain.Quiz {
	if source == nil {
		return nil
	}
	return &quizDomain.Quiz{
		QuizID:        source.QuizID,
		Title:         source.Title,
		QuizQuestions: helper.MapList(source.QuizQuestions, mapQuizQuestionToDomain),
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}
}

func (q quizGorm) TableName() string {
	return "quiz"
}

func NewQuizRepo() quizDomain.QuizRepo {
	return &quizRepo{}
}

type quizRepo struct {
}

type quizQuery struct {
	query.BaseQuery
}

func (q *quizRepo) Query(ctx context.Context) quizDomain.QuizQuery {
	return &quizQuery{query.NewBQ(ctx.GetDbTx().Model(&quizGorm{}))}
}

func (q *quizQuery) ByQuizID(quizID string) quizDomain.QuizQuery {
	return query.Where(q, "quiz_id = ?", quizID)
}

func (q *quizQuery) WithQuizQuestion(questionID string) quizDomain.QuizQuery {
	q.SetDB(q.GetDB().Preload("QuizQuestions", func(db *gorm.DB) *gorm.DB {
		if len(questionID) > 0 {
			return db.Where("question_id = ?", questionID).Preload("Question")
		}
		return db.Preload("Question")
	}))
	return q
}

func (q *quizQuery) Limit(limit int) quizDomain.QuizQuery {
	return query.Limit(q, limit)
}

func (q *quizQuery) Offset(offset int) quizDomain.QuizQuery {
	return query.Offset(q, offset)
}

func (q *quizQuery) Result() (*quizDomain.Quiz, error) {
	return query.Result(q, mapQuizToDomain)
}

func (q *quizQuery) ResultList() ([]*quizDomain.Quiz, error) {
	return query.ResultList(q, mapQuizToDomain)
}

func (q *quizQuery) Count() (int, error) {
	return query.Count(q)
}
