package persistence

import "github.com/mohitkumar/finch/model"

const WF_PREFIX string = "WF_"
const METADATA_CF string = "METADATA_"

type WorkflowDao interface {
	Save(wf model.Workflow) (bool, error)

	Delete(name string) (bool, error)

	Get(name string) (*model.Workflow, error)
}
