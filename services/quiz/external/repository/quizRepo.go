package repository

import (
	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
)

type quizGorm struct {
	QuizID      string `gorm:"primary_key"`
	Content     string
	QuestionIDs []string
	Questions   []*questionGrom
}

func mapQuizToDomain(source *quizGorm) *quizDomain.Quiz {
	if source == nil {
		return nil
	}
	return &quizDomain.Quiz{
		QuizID:    source.QuizID,
		Content:   source.Content,
		Questions: nil,
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

func (q *quizQuery) Limit(limit int) quizDomain.QuizQuery {
	return query.Limit(q, limit)
}

func (q *quizQuery) Result() (*quizDomain.Quiz, error) {
	return query.Result(q, mapQuizToDomain)
}

func (q *quizQuery) ResultList() ([]*quizDomain.Quiz, error) {
	return query.ResultList(q, mapQuizToDomain)
}
