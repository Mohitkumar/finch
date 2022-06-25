package model

type Workflow struct {
	Name       string         `json:"name"`
	Data       map[string]any `json:"data"`
	RootAction int            `json:"rootAction"`
	Actions    []ActionDef    `json:"actions"`
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
