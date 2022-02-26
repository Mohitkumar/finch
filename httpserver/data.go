package httpserver

type Flow struct {
	Name          string                 `json:"name"`
	Data          map[string]interface{} `json:"data"`
	StartActionId uint64                 `json:"startAction"`
	Actions       []Action               `json:"actions"`
}

type Action struct {
	Id              uint64                 `json:"id"`
	Name            string                 `json:"name"`
	InputParameters map[string]interface{} `json:"inputParameters"`
	Next            []ActionNode           `json:"next"`
}

type ActionNode struct {
	Id    uint64 `json:"id"`
	Event string `json:"event"`
}
