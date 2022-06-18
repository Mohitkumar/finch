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
	Actions    []ActionDef            `json:"actions"`
}

type ActionDef struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Data       map[string]any `json:"data"`
	Next       int            `json:"next"`
	Expression string         `json:"expression"`
	Cases      map[string]int `json:"cases"`
	Forks      []int          `json:"forks"`
	Join       int            `json:"join"`
}
