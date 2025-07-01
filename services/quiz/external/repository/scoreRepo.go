package repository

import (
	"time"

	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
)

type scoreGorm struct {
	ID        string `gorm:"primary_key"`
	UserID    string
	QuizID    string
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *userGorm `gorm:"foreignKey:UserID;references:UserID"`
	Quiz      *quizGorm `gorm:"foreignKey:QuizID;references:QuizID"`
}

func mapScoreToDomain(source *scoreGorm) *quizDomain.Score {
	if source == nil {
		return nil
	}
	return &quizDomain.Score{
		ID:        source.ID,
		UserID:    source.UserID,
		QuizID:    source.QuizID,
		Score:     source.Score,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		User:      mapUserToDomain(source.User),
		Quiz:      mapQuizToDomain(source.Quiz),
	}
}

func mapScoreFromDomain(source *quizDomain.Score) *scoreGorm {
	if source == nil {
		return nil
	}
	return &scoreGorm{
		ID:        source.ID,
		UserID:    source.UserID,
		QuizID:    source.QuizID,
		Score:     source.Score,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}

func (s scoreGorm) TableName() string {
	return "scores"
}

func NewScoreRepo() quizDomain.ScoreRepo {
	return &scoreRepo{}
}

type scoreRepo struct {
}

type scoreQuery struct {
	query.BaseQuery
}

func (s *scoreRepo) Query(ctx context.Context) quizDomain.ScoreQuery {
	return &scoreQuery{query.NewBQ(ctx.GetDbTx().Model(&scoreGorm{}))}
}

func (s *scoreRepo) Upsert(ctx context.Context, score *quizDomain.Score) error {
	return query.Upsert(ctx.GetDbTx(), score, mapScoreFromDomain)
}

func (s *scoreQuery) ByID(id string) quizDomain.ScoreQuery {
	return query.Where(s, "id = ?", id)
}

func (s *scoreQuery) ByUserID(userID string) quizDomain.ScoreQuery {
	return query.Where(s, "user_id = ?", userID)
}

func (s *scoreQuery) ByQuizID(quizID string) quizDomain.ScoreQuery {
	return query.Where(s, "quiz_id = ?", quizID)
}
func (s *scoreQuery) WithUser() quizDomain.ScoreQuery {
	s.SetDB(s.GetDB().Preload("User"))
	return s
}

// paging
func (s *scoreQuery) Limit(limit int) quizDomain.ScoreQuery {
	return query.Limit(s, limit)
}

func (s *scoreQuery) Offset(offset int) quizDomain.ScoreQuery {
	return query.Offset(s, offset)
}

// ordering
func (s *scoreQuery) OrderByScore(desc bool) quizDomain.ScoreQuery {
	ord := ""
	if desc {
		ord = " desc"
	}
	s.SetDB(s.GetDB().Order("score" + ord))
	return s
}

func (s *scoreQuery) Result() (*quizDomain.Score, error) {
	return query.Result(s, mapScoreToDomain)
}

func (s *scoreQuery) ResultList() ([]*quizDomain.Score, error) {
	return query.ResultList(s, mapScoreToDomain)
}

func (s *scoreQuery) Count() (int, error) {
	return query.Count(s)
}

func (s *scoreQuery) Delete() error {
	return query.Delete(s)
}
