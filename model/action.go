package model

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type Action interface {
	GetId() int
	Execute(flowContext *api.FlowContext) error
}

var _ Action = new(UserAction)

func NewUserAction(id int, Type string, name string, data map[string]any) *UserAction {
	act := &UserAction{
		Id:   id,
		Name: name,
		Data: data,
		Type: Type,
	}
	return act
}

type UserAction struct {
	Id   int
	Type string
	Name string
	Data map[string]any
}

func (ba *UserAction) GetId() int {
	return ba.Id
}

func (ba *UserAction) Execute(flowContext *api.FlowContext) error {

}
