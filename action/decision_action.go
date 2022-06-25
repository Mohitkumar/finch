package action

import (
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
)

var _ Action = new(decisionAction)

type decisionAction struct {
	baseAction
	expression string
}

func NewDecisionAction(id int, Type string, name string, inputParams map[string]any, expression string) *decisionAction {
	return &decisionAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams),
		expression: expression,
	}
}

func (d *decisionAction) GetExpression() string {
	return d.expression
}

func (d *decisionAction) Execute(wfName string, flowContext *api.FlowContext) error {
	return fmt.Errorf("can not execute")
}
