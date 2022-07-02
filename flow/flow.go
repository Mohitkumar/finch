package flow

import (
	"strings"

	"github.com/mohitkumar/finch/action"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence"
)

type FlowType string

const (
	FLOW_TYPE_SYSTEM FlowType = "SYSTEM"
	FLOW_TYPE_USER   FlowType = "USER"
)

type Flow struct {
	Id         string
	RootAction int
	Actions    map[int]action.Action
}

func Convert(wf *model.Workflow, id string, pFactory persistence.PersistenceFactory) Flow {
	actionMap := make(map[int]action.Action)
	for _, actionDef := range wf.Actions {
		var flAct action.Action = action.NewBaseAction(actionDef.Id, actionDef.Type,
			actionDef.Name, actionDef.InputParams, pFactory)
		if actionDef.Type == string(FLOW_TYPE_SYSTEM) {
			if strings.EqualFold(actionDef.Name, "decision") {
				flAct = action.NewDecisionAction(actionDef.Id, actionDef.Type,
					actionDef.Name, actionDef.InputParams, actionDef.Expression, pFactory)
			}
		} else {
			flAct = action.NewUserAction(actionDef.Id, actionDef.Type,
				actionDef.Name, actionDef.InputParams, actionDef.Next, pFactory)
		}
		actionMap[actionDef.Id] = flAct
	}
	flow := Flow{
		Id:         id,
		RootAction: wf.RootAction,
		Actions:    actionMap,
	}
	return flow
}
