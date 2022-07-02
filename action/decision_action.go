package action

import (
	"strconv"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence"
	"github.com/mohitkumar/finch/util"
	"github.com/yalp/jsonpath"
)

var _ Action = new(decisionAction)

type decisionAction struct {
	baseAction
	expression string
	cases      map[string]int
}

func NewDecisionAction(id int, Type ActionType, name string, inputParams map[string]any, expression string, pFactory persistence.PersistenceFactory) *decisionAction {
	return &decisionAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		expression: expression,
	}
}

func (d *decisionAction) Execute(wfName string, flowContext *api.FlowContext) (*ActionResult, error) {
	dataMap := util.ConvertFromProto(flowContext.Data)
	expressionValue, err := jsonpath.Read(dataMap, d.expression)
	if err != nil {
		return nil, err
	}
	var nextAction int
	switch expValue := expressionValue.(type) {
	case int:
		nextAction = d.cases[strconv.Itoa(nextAction)]
	case string:
		nextAction = d.cases[expValue]
	}
	result := &ActionResult{
		NextAction: nextAction,
		ActionType: ACTION_TYPE_SYSTEM,
	}
	return result, nil
}
