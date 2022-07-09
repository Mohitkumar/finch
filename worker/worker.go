package worker

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type Worker interface {
	Execute(*api.Task) (*api.TaskResult, error)
	GetName() string
	GetPollInterval() int
}
