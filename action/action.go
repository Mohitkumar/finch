package action

import (
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
)

type Action interface {
	GetId() int
	GetName() string
	GetType() string
	GetInputParams() map[string]any
	GetExpression() string
	Execute(wfName string, flowContext *api.FlowContext) error
}

var _ Action = new(baseAction)

type baseAction struct {
	id          int
	actType     string
	name        string
	inputParams map[string]any
}

func NewBaseAction(id int, Type string, name string, inputParams map[string]any) *baseAction {
	return &baseAction{
		id:          id,
		name:        name,
		inputParams: inputParams,
		actType:     Type,
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

func (ba *baseAction) GetExpression() string {
	return "nil"
}

func (ba *baseAction) Execute(wfName string, flowContext *api.FlowContext) error {
	return fmt.Errorf("can not execute")
}
