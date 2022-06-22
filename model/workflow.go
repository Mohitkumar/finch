package model

type Flow struct {
	Id         string
	RootAction int
	Actions    map[int]Action
	Data       map[string]any
}

type WorkflowType string

const (
	WF_TYPE_SYSTEM WorkflowType = "SYSTEM"
	WF_TYPE_USER   WorkflowType = "USER"
)

type Workflow struct {
	Name       string         `json:"name"`
	Data       map[string]any `json:"data"`
	RootAction int            `json:"rootAction"`
	Actions    []ActionDef    `json:"actions"`
}

func (wf *Workflow) Convert(id string, data map[string]any) Flow {
	actionMap := make(map[int]Action)
	for _, action := range wf.Actions {
		if action.Type == string(WF_TYPE_SYSTEM) {

		} else {
			flAct := NewUserAction(action.Id, action.Type, action.Name, action.Data)
			actionMap[action.Id] = flAct
		}
	}
	flow := Flow{
		Id:         id,
		Data:       data,
		RootAction: wf.RootAction,
		Actions:    actionMap,
	}
	return flow
}

type ActionDef struct {
	Id         int            `json:"id"`
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Data       map[string]any `json:"data"`
	Next       int            `json:"next"`
	Expression string         `json:"expression"`
	Cases      map[string]int `json:"cases"`
	Forks      []int          `json:"forks"`
	Join       int            `json:"join"`
}
