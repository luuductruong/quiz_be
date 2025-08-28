package repository

import (
	"github.com/quiz_be/services/core/context"
	jobDomain "github.com/quiz_be/services/core/domain/job"
	"github.com/quiz_be/services/core/helper/sql/query"
	"time"
)

type jobGorm struct {
	ID            string `gorm:"primary_key"`
	Name          string
	Data          []byte
	LastRunStatus string
	LastRunAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (j jobGorm) TableName() string {
	return "job"
}

func mapJobFromDomain(job *jobDomain.Job) *jobGorm {
	if job == nil {
		return nil
	}
	return &jobGorm{
		ID:            job.ID,
		Name:          job.Name,
		Data:          job.Data,
		LastRunStatus: job.LastRunStatus,
		LastRunAt:     job.LastRunAt,
		CreatedAt:     job.CreatedAt,
		UpdatedAt:     job.UpdatedAt,
	}
}

func mapJobToDomain(jobG *jobGorm) *jobDomain.Job {
	if jobG == nil {
		return nil
	}
	return &jobDomain.Job{
		ID:            jobG.ID,
		Name:          jobG.Name,
		Data:          jobG.Data,
		LastRunStatus: jobG.LastRunStatus,
		LastRunAt:     jobG.LastRunAt,
		CreatedAt:     jobG.CreatedAt,
		UpdatedAt:     jobG.UpdatedAt,
	}
}

func NewJobRepo() jobDomain.JobRepo {
	return &jobRepo{}
}

type jobRepo struct {
}

func (j *jobRepo) Query(ctx context.Context) jobDomain.JobQuery {
	return &jobQuery{query.NewBQ(ctx.GetDbTx().Model(&jobGorm{}))}
}

func (q *jobRepo) Upsert(ctx context.Context, jobD *jobDomain.Job) error {
	return query.Upsert(ctx.GetDbTx(), jobD, mapJobFromDomain)
}

type jobQuery struct {
	query.BaseQuery
}

func (j *jobQuery) ByName(name string) jobDomain.JobQuery {
	return query.Where(j, "name = ?", name)
}

func (j *jobQuery) ByTopic(topic string) jobDomain.JobQuery {
	return query.Where(j, "topic = ?", topic)
}

func (j *jobQuery) Limit(limit int) jobDomain.JobQuery {
	return query.Limit(j, limit)
}

func (j *jobQuery) Offset(offset int) jobDomain.JobQuery {
	return query.Offset(j, offset)
}

func (j *jobQuery) Result() (*jobDomain.Job, error) {
	return query.Result(j, mapJobToDomain)
}

func (j *jobQuery) ResultList() ([]*jobDomain.Job, error) {
	return query.ResultList(j, mapJobToDomain)

}

func (j *jobQuery) Count() (int, error) {
	return query.Count(j)
}
