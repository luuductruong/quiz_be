package repository

import (
	"encoding/json"
	"gorm.io/datatypes"
	"time"

	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
)

type questionGrom struct {
	QuestionID    string `gorm:"primary_key"`
	Content       string
	CorrectAnswer string
	Score         uint
	Answers       datatypes.JSON `json:"answers"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type answerGorm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func mapQuestionToDomain(source *questionGrom) *quizDomain.Question {
	if source == nil {
		return nil
	}
	var answers []*quizDomain.Answer
	_ = json.Unmarshal(source.Answers, &answers)
	return &quizDomain.Question{
		QuestionID:    source.QuestionID,
		Content:       source.Content,
		CorrectAnswer: source.CorrectAnswer,
		Score:         source.Score,
		Answers:       answers,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}
}

func mapQuestionFromDomain(source *quizDomain.Question) *questionGrom {
	if source == nil {
		return nil
	}
	answers, _ := json.Marshal(source.Answers)
	return &questionGrom{
		QuestionID:    source.QuestionID,
		Content:       source.Content,
		CorrectAnswer: source.CorrectAnswer,
		Score:         source.Score,
		Answers:       answers,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
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

func (q *questionRepo) Upsert(ctx context.Context, question *quizDomain.Question) error {
	return query.Upsert(ctx.GetDbTx(), question, mapQuestionFromDomain)
}

func (q *questionRepo) Query(ctx context.Context) quizDomain.QuestionQuery {
	return &questionQuery{query.NewBQ(ctx.GetDbTx().Model(&questionGrom{}))}
}

func (q *questionQuery) ByQuestionID(questionID string) quizDomain.QuestionQuery {
	return query.Where(q, "question_id = ?", questionID)
}

func (q *questionQuery) ByQuestionIDs(questionIDs []string) quizDomain.QuestionQuery {
	return query.Where(q, "question_id in ?", questionIDs)
}

func (q *questionQuery) Limit(limit int) quizDomain.QuestionQuery {
	return query.Limit(q, limit)
}

func (q *questionQuery) Offset(offset int) quizDomain.QuestionQuery {
	return query.Offset(q, offset)
}

func (q *questionQuery) OrderByUpdatedAt(desc bool) quizDomain.QuestionQuery {
	ord := ""
	if desc {
		ord = " desc"
	}
	q.SetDB(q.GetDB().Order("updated_at" + ord))
	return q
}

func (q *questionQuery) Result() (*quizDomain.Question, error) {
	return query.Result(q, mapQuestionToDomain)
}

func (q *questionQuery) ResultList() ([]*quizDomain.Question, error) {
	return query.ResultList(q, mapQuestionToDomain)
}

func (q *questionQuery) Count() (int, error) {
	return query.Count(q)
}
