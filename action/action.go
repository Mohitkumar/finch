package action

import (
	"fmt"
	"strings"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence/factory"
	"github.com/mohitkumar/finch/util"
	"github.com/oliveagle/jsonpath"
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
	Execute(wfName string, flowContext *api.FlowContext) error
}

var _ Action = new(baseAction)

type baseAction struct {
	id          int
	actType     ActionType
	name        string
	inputParams map[string]any
	pFactory    *factory.PersistenceFactory
}

func NewBaseAction(id int, Type ActionType, name string, inputParams map[string]any, pFactory *factory.PersistenceFactory) *baseAction {
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

func (ba *baseAction) Execute(wfName string, flowContext *api.FlowContext) error {
	return fmt.Errorf("can not execute")
}

func (ba *baseAction) ResolveInputParams(flowContext *api.FlowContext) map[string]any {
	flowData := util.ConvertFromProto(flowContext.Data)
	data := make(map[string]any)
	ba.resolveParams(flowData, ba.inputParams, data)
	return data
}

func (ba *baseAction) resolveParams(flowData map[string]any, params map[string]any, output map[string]any) {
	for k, v := range params {
		switch v.(type) {
		case map[string]any:
			out := make(map[string]any)
			output[k] = out
			ba.resolveParams(flowData, v.(map[string]any), out)
		case string:
			if strings.HasPrefix(v.(string), "$") {
				value, _ := jsonpath.JsonPathLookup(flowData, v.(string))
				output[k] = value
			} else {
				output[k] = v
			}
		default:
			output[k] = v
		}
	}
}
