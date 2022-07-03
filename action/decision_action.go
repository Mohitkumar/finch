package action

import (
	"strconv"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence/factory"
	"github.com/mohitkumar/finch/util"
	"github.com/yalp/jsonpath"
)

var _ Action = new(decisionAction)

type decisionAction struct {
	baseAction
	expression string
	cases      map[string]int
}

func NewDecisionAction(id int, Type ActionType, name string, inputParams map[string]any, expression string, pFactory *factory.PersistenceFactory) *decisionAction {
	return &decisionAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		expression: expression,
	}
}

func (d *decisionAction) Execute(wfName string, flowContext *api.FlowContext) error {
	dataMap := util.ConvertFromProto(flowContext.Data)
	expressionValue, err := jsonpath.Read(dataMap, d.expression)
	if err != nil {
		return err
	}
	var nextAction int
	switch expValue := expressionValue.(type) {
	case int:
		nextAction = d.cases[strconv.Itoa(nextAction)]
	case string:
		nextAction = d.cases[expValue]
	}
	flowContext.NextAction = int32(nextAction)
	d.pFactory.GetFlowDao().SaveFlowContext(wfName, flowContext.Id, flowContext)
	return nil
}
