package persistence

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type (
	ExecutionDao interface {
		Save(sm *api.StateMachineContext) error
	}
)
