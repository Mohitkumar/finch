package action

import (
	"strconv"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence/factory"
	"github.com/mohitkumar/finch/util"
	"github.com/oliveagle/jsonpath"
)

var _ Action = new(switchAction)

type switchAction struct {
	baseAction
	expression string
	cases      map[string]int
}

func NewSwitchAction(id int, Type ActionType, name string, inputParams map[string]any, expression string, pFactory *factory.PersistenceFactory) *decisionAction {
	return &switchAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		expression: expression,
	}
}

func (d *switchAction) Execute(wfName string, flowContext *api.FlowContext) error {
	dataMap := util.ConvertFromProto(flowContext.Data)
	expressionValue, err := jsonpath.JsonPathLookup(dataMap, d.expression)
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
	err = d.pFactory.GetFlowDao().UpdateFlowContextNextAction(wfName, flowContext.Id, flowContext, nextAction)
	if err != nil {
		return err
	}
	return nil
}
