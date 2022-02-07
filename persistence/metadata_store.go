package persistence

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type (
	MetadataManager interface {
		SaveStateMachine(m *api.StateMachine) error
		GetStateMachine(name string) (*api.StateMachine, error)
	}
)
