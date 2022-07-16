package action

import (
	"encoding/json"
	"time"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence/factory"
)

var _ Action = new(delayAction)

type delayAction struct {
	baseAction
	delay      time.Duration
	nextAction int
}

func NewDelayAction(id int, Type ActionType, name string, delaySeconds int, nextAction int, pFactory *factory.PersistenceFactory) *delayAction {
	inputParams := map[string]any{}
	return &delayAction{
		baseAction: *NewBaseAction(id, Type, name, inputParams, pFactory),
		delay:      time.Duration(delaySeconds) * time.Second,
		nextAction: nextAction,
	}
}
func (d *delayAction) Execute(wfName string, flowContext *api.FlowContext) error {
	msg := &model.FlowContextMessage{
		WorkflowName: wfName,
		FlowId:       flowContext.Id,
		ActionId:     d.nextAction,
	}
	data, _ := json.Marshal(msg)
	err := d.pFactory.GetDelayQueue().PushWithDelay("delay_action", d.delay, data)
	if err != nil {
		return err
	}
	return d.pFactory.GetFlowDao().UpdateFlowStatus(wfName, flowContext.Id, flowContext, api.FlowContext_DELAY_WATING)
}
