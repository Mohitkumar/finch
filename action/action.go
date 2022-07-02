package action

import (
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence"
)

type Action interface {
	GetId() int
	GetName() string
	GetType() string
	GetInputParams() map[string]any
	Execute(wfName string, flowContext *api.FlowContext) (*ActionResult, error)
}

var _ Action = new(baseAction)

type baseAction struct {
	id          int
	actType     string
	name        string
	inputParams map[string]any
	pFactory    persistence.PersistenceFactory
}

func NewBaseAction(id int, Type string, name string, inputParams map[string]any, pFactory persistence.PersistenceFactory) *baseAction {
	return &baseAction{
		id:          id,
		name:        name,
		inputParams: inputParams,
		actType:     Type,
		pFactory:    pFactory,
	}

}
func (ba *baseAction) GetId() int {
	return ba.id
}
func (ba *baseAction) GetName() string {
	return ba.name
}
func (ba *baseAction) GetType() string {
	return ba.actType
}
func (ba *baseAction) GetInputParams() map[string]any {
	return ba.inputParams
}

func (ba *baseAction) Execute(wfName string, flowContext *api.FlowContext) (*ActionResult, error) {
	return nil, fmt.Errorf("can not execute")
}

type ActionResult struct {
	NextAction int
	Data       map[string]any
}
