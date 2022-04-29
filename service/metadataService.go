package service

import (
	"github.com/mohitkumar/finch/model"
)

type MetadataService interface {
	RegisterWorkflow(wf model.Workflow) error
	GetWorkflow(name string)
}
