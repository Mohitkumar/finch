package persistance

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type (
	QueueDao interface {
		Enqueue(*api.Function) error
		Dequeue() (*api.Function, error)
	}
)
