package job

import "github.com/quiz_be/services/core/context"

type JobRepo interface {
	Upsert(ctx context.Context, job *Job) error
	Query(ctx context.Context) JobQuery
}

type JobQuery interface {
	// query
	ByName(name string) JobQuery
	ByTopic(topic string) JobQuery
	Limit(limit int) JobQuery
	Offset(offset int) JobQuery
	// result
	Result() (*Job, error)
	ResultList() ([]*Job, error)
	Count() (int, error)
}
