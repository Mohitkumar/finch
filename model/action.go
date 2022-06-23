package model

type Action interface {
	GetId() int
	GetName() string
	GetType() string
	GetInputParams() map[string]any
}

var _ Action = new(UserAction)

func NewUserAction(id int, Type string, name string, inputParams map[string]any) *UserAction {
	act := &UserAction{
		Id:          id,
		Name:        name,
		InputParams: inputParams,
		Type:        Type,
	}
	return act
}

type UserAction struct {
	Id          int
	Type        string
	Name        string
	InputParams map[string]any
}

func (ba *UserAction) GetId() int {
	return ba.Id
}
func (ba *UserAction) GetName() string {
	return ba.Name
}
func (ba *UserAction) GetType() string {
	return ba.Type
}
func (ba *UserAction) GetInputParams() map[string]any {
	return ba.InputParams
}
