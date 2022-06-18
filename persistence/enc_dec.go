package persistence

import (
	"encoding/json"

	"github.com/mohitkumar/finch/model"
)

type serializable interface {
	model.Workflow | model.Flow
}

type JsonEncDec[T serializable] struct{}

func (encdec *JsonEncDec[T]) Encode(value T) ([]byte, error) {
	res, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (encdec *JsonEncDec[T]) Decode(data []byte) (*T, error) {
	var res T
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
