package action

import (
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence"
)

type ActionType string

const ACTION_TYPE_SYSTEM ActionType = "SYSTEM"
const ACTION_TYPE_USER ActionType = "USER"

func ToActionType(at string) ActionType {
	if at == "SYSTEM" {
		return ACTION_TYPE_SYSTEM
	}
	return ACTION_TYPE_USER
}

type Action interface {
	GetId() int
	GetName() string
	GetType() ActionType
	GetInputParams() map[string]any
	Execute(wfName string, flowContext *api.FlowContext) (*ActionResult, error)
}

var _ Action = new(baseAction)

type baseAction struct {
	id          int
	actType     ActionType
	name        string
	inputParams map[string]any
	pFactory    persistence.PersistenceFactory
}

func NewBaseAction(id int, Type ActionType, name string, inputParams map[string]any, pFactory persistence.PersistenceFactory) *baseAction {
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
func (ba *baseAction) GetType() ActionType {
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
	ActionType ActionType
}
