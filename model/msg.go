package model

type FlowContextMessage struct {
	WorkflowName string `json:"wfName"`
	FlowId       string `json:"flowId"`
	ActionId     int    `json:"actionId"`
}
