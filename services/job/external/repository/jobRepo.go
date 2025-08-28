package repository

import (
	"github.com/quiz_be/services/core/context"
	jobDomain "github.com/quiz_be/services/core/domain/job"
	"github.com/quiz_be/services/core/helper/sql/query"
	"gorm.io/datatypes"
	"time"
)

type jobGorm struct {
	ID            string `gorm:"primaryKey"`
	Name          string
	Data          datatypes.JSON
	Status        string `gorm:"default:PENDING;not null"`
	RetryCount    int    `gorm:"default:0;not null"`
	LastRunStatus string
	LastError     string
	LastRunAt     *time.Time
	TraceID       string
	Exchange      string
	RoutingKey    string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
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
		Status:        job.Status,
		RetryCount:    job.RetryCount,
		LastRunStatus: job.LastRunStatus,
		LastError:     job.LastError,
		LastRunAt:     job.LastRunAt,
		TraceID:       job.TraceID,
		Exchange:      job.Exchange,
		RoutingKey:    job.RoutingKey,
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
		Status:        jobG.Status,
		RetryCount:    jobG.RetryCount,
		LastRunStatus: jobG.LastRunStatus,
		LastError:     jobG.LastError,
		LastRunAt:     jobG.LastRunAt,
		TraceID:       jobG.TraceID,
		Exchange:      jobG.Exchange,
		RoutingKey:    jobG.RoutingKey,
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
