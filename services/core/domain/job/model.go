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
	Status        string
	RetryCount    int
	LastError     string
	TraceID       string
	Exchange      string
	RoutingKey    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
