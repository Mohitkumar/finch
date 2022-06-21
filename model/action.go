package model

type Action interface {
	GetId() int
	Execute() error
}

var _ Action = new(BaseAction)

func NewBaseAction(id int, Type string, name string, data map[string]any) *BaseAction {
	act := &BaseAction{
		Id:   id,
		Name: name,
		Data: data,
		Type: Type,
	}
	return act
}

type BaseAction struct {
	Id   int
	Type string
	Name string
	Data map[string]any
}

func (ba *BaseAction) GetId() int {
	return ba.Id
}

func (ba *BaseAction) Execute() error {

}
