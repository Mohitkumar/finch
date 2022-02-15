package shard

import (
	"strings"

	"github.com/golang/protobuf/proto"
	api "github.com/mohitkumar/finch/api/v1"
)

type FlowCreateStatus int16

const (
	FlowCreateStatusFailed  FlowCreateStatus = -1
	FlowCreateStatusSuccess FlowCreateStatus = 1
)

func getFlowKeyPrefix(name string) string {
	return strings.Join([]string{"flow", name}, "_")
}

func getFlowShardPrefix(name string) string {
	return strings.Join([]string{"flow_shard", name}, "_")
}

func (shard *Shard) CreateFlow(flow *api.Flow) (FlowCreateStatus, error) {
	key := getFlowKeyPrefix(flow.Name)
	flowBytes, err := proto.Marshal(flow)
	if err != nil {
		return FlowCreateStatusFailed, err
	}
	if err := shard.kvStore.Put([]byte(key), flowBytes); err != nil {
		return FlowCreateStatusFailed, err
	}
	if err := shard.kvStore.Put([]byte(getFlowShardPrefix(flow.Name)), []byte(shard.ID)); err != nil {
		return FlowCreateStatusFailed, err
	}

	return FlowCreateStatusSuccess, nil
}
