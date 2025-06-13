package repository

import (
	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
)

type questionGrom struct {
	QuestionID    string
	Content       string
	CorrectAnswer string
	Answers       []*answerGorm `json:"answers"`
}

type answerGorm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func mapQuestionToDomain(source *questionGrom) *quizDomain.Question {
	if source == nil {
		return nil
	}
	return &quizDomain.Question{
		QuestionID:    source.QuestionID,
		Content:       source.Content,
		CorrectAnswer: source.CorrectAnswer,
		Answers:       nil,
	}
}

func (q questionGrom) TableName() string {
	return "question"
}

func NewQuestionRepo() quizDomain.QuestionRepo {
	return &questionRepo{}
}

type questionRepo struct {
}

type questionQuery struct {
	query.BaseQuery
}

func (q *questionRepo) Query(ctx context.Context) quizDomain.QuestionQuery {
	return &questionQuery{query.NewBQ(ctx.GetDbTx().Model(&questionGrom{}))}
}

func (q *questionQuery) ByQuestionID(questionID string) quizDomain.QuestionQuery {
	return query.Where(q, "question_id = ?", questionID)
}

func (q *questionQuery) Result() (*quizDomain.Question, error) {
	return query.Result(q, mapQuestionToDomain)
}
