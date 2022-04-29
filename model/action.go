package model

import (
	api "github.com/mohitkumar/finch/api/v1"
)

type Action interface {
	GetId() uint64
	Init(id uint64, data string)
	GetNext() map[string]Action
	Process(context api.FlowContext) error
	Unprocess(context api.FlowContext) error
}
