package coordinator

import (
	"strings"

	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/protobuf/proto"
)

type FlowCreateStatus int16

const (
	FlowCreateStatusFailed  FlowCreateStatus = -1
	FlowCreateStatusSuccess FlowCreateStatus = 1
)

func getFlowKeyPrefix(name string) string {
	return strings.Join([]string{"flow", name}, "_")
}

func (c *Coordinator) CreateFlow(flow *api.Flow) (FlowCreateStatus, error) {
	key := getFlowKeyPrefix(flow.Name)
	flowBytes, err := proto.Marshal(flow)
	if err != nil {
		return FlowCreateStatusFailed, err
	}
	kvItem := &api.KVItem{
		Key:   []byte(key),
		Value: flowBytes,
	}
	if _, err := c.apply(DBPutRequestType, kvItem); err != nil {
		return FlowCreateStatusFailed, err
	}
	return FlowCreateStatusSuccess, nil
}
