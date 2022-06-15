package model

type Flow struct {
	Id         uint64
	RootAction Action
	Actions    map[int]Action
	Data       map[string]interface{}
}

type Workflow struct {
	Name       string                 `json:"name"`
	Data       map[string]interface{} `json:"data"`
	RootAction int                    `json:"rootAction"`
	Actions    []Action               `json:"actions"`
}
