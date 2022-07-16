package flow

import (
	"strings"

	"github.com/mohitkumar/finch/action"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence/factory"
)

type Flow struct {
	Id         string
	RootAction int
	Actions    map[int]action.Action
}

func Convert(wf *model.Workflow, id string, pFactory *factory.PersistenceFactory) Flow {
	actionMap := make(map[int]action.Action)
	for _, actionDef := range wf.Actions {
		actionType := action.ToActionType(actionDef.Type)
		var flAct action.Action = action.NewBaseAction(actionDef.Id, actionType,
			actionDef.Name, actionDef.InputParams, pFactory)
		if actionType == action.ACTION_TYPE_SYSTEM {
			if strings.EqualFold(actionDef.Name, "switch") {
				flAct = action.NewSwitchAction(actionDef.Id, actionType,
					actionDef.Name, actionDef.Expression, actionDef.Cases, pFactory)
			} else if strings.EqualFold(actionDef.Name, "delay") {
				flAct = action.NewDelayAction(actionDef.Id, actionType,
					actionDef.Name, actionDef.DelaySeconds, actionDef.Next, pFactory)
			}
		} else {
			flAct = action.NewUserAction(actionDef.Id, actionType,
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
