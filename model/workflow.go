package model

type Flow struct {
	Id         string
	RootAction int
	Actions    map[int]Action
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

func (wf *Workflow) Convert(id string) Flow {
	actionMap := make(map[int]Action)
	for _, actionDef := range wf.Actions {
		if actionDef.Type == string(WF_TYPE_SYSTEM) {

		} else {
			flAct := NewUserAction(actionDef.Id, actionDef.Type, actionDef.Name, actionDef.InputParams)
			actionMap[actionDef.Id] = flAct
		}
	}
	flow := Flow{
		Id:         id,
		RootAction: wf.RootAction,
		Actions:    actionMap,
	}
	return flow
}

type ActionDef struct {
	Id          int            `json:"id"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	InputParams map[string]any `json:"inputParams"`
	Next        int            `json:"next"`
	Expression  string         `json:"expression"`
	Cases       map[string]int `json:"cases"`
	Forks       []int          `json:"forks"`
	Join        int            `json:"join"`
}
