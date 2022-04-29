package persistence

import "github.com/mohitkumar/finch/model"

type Persistence interface {
	Save(workflow model.Workflow) error
}
