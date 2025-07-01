package repository

import (
	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
)

type quizQuestionGorm struct {
	QuizID     string        `gorm:"primaryKey;column:quiz_id"`
	QuestionID string        `gorm:"primaryKey;column:question_id"`
	Question   *questionGrom `gorm:"foreignKey:QuestionID;references:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func mapQuizQuestionToDomain(source *quizQuestionGorm) *quizDomain.QuizQuestion {
	if source == nil {
		return nil
	}
	return &quizDomain.QuizQuestion{
		QuizID:     source.QuizID,
		QuestionID: source.QuestionID,
		Question:   mapQuestionToDomain(source.Question),
	}
}

func mapQuizQuestionFromDomain(source *quizDomain.QuizQuestion) *quizQuestionGorm {
	if source == nil {
		return nil
	}
	return &quizQuestionGorm{
		QuizID:     source.QuizID,
		QuestionID: source.QuestionID,
		//Question:   mapQuestionToDomain(source.Question),
	}
}

func (q quizQuestionGorm) TableName() string {
	return "quiz_question"
}

func NewQuizQuestionRepo() quizDomain.QuizQuestionRepo {
	return &quizQuestionRepo{}
}

type quizQuestionRepo struct {
}

type quizQuestionQuery struct {
	query.BaseQuery
}

func (q *quizQuestionRepo) BulkUpsert(ctx context.Context, quizQuestion []*quizDomain.QuizQuestion) error {
	if len(quizQuestion) == 0 {
		return nil
	}
	return query.BatchUpsert(ctx.GetDbTx(), quizQuestion, mapQuizQuestionFromDomain)
}

func (q *quizQuestionRepo) Query(ctx context.Context) quizDomain.QuizQuestionQuery {
	return &quizQuestionQuery{query.NewBQ(ctx.GetDbTx().Model(&quizQuestionGorm{}))}
}

func (q *quizQuestionQuery) ByQuizID(quizID string) quizDomain.QuizQuestionQuery {
	return query.Where(q, "quiz_id = ?", quizID)
}

func (q *quizQuestionQuery) ByQuestionID(questionID string) quizDomain.QuizQuestionQuery {
	return query.Where(q, "question_id = ?", questionID)
}

func (q *quizQuestionQuery) Limit(limit int) quizDomain.QuizQuestionQuery {
	return query.Limit(q, limit)
}

func (q *quizQuestionQuery) Result() (*quizDomain.QuizQuestion, error) {
	return query.Result(q, mapQuizQuestionToDomain)
}

func (q *quizQuestionQuery) ResultList() ([]*quizDomain.QuizQuestion, error) {
	return query.ResultList(q, mapQuizQuestionToDomain)
}
