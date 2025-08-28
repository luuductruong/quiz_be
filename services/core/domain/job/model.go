package job

import "time"

type Test struct {
	Name string
}

type Job struct {
	ID            string
	Name          string
	Data          []byte // sample
	LastRunStatus string
	LastRunAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
