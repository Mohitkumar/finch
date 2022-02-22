package coordinator

import (
	"strings"

	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/protobuf/proto"
)

func getFlowKeyPrefix(name string) string {
	return strings.Join([]string{"flow", name}, "_")
}

func (c *Coordinator) CreateFlow(flow *api.Flow) (*api.FlowCreateResponse, error) {
	key := getFlowKeyPrefix(flow.Name)
	flowBytes, err := proto.Marshal(flow)
	if err != nil {
		return &api.FlowCreateResponse{Status: api.FlowCreateResponse_FAILED}, err
	}
	kvItem := &api.KVItem{
		Key:   []byte(key),
		Value: flowBytes,
	}
	if _, err := c.apply(DBPutRequestType, kvItem); err != nil {
		return &api.FlowCreateResponse{Status: api.FlowCreateResponse_FAILED}, err
	}
	return &api.FlowCreateResponse{Status: api.FlowCreateResponse_SUCCESS}, nil
}

func (c *Coordinator) GetServers() ([]*api.Server, error) {
	future := c.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}
	var servers []*api.Server
	for _, server := range future.Configuration().Servers {
		servers = append(servers, &api.Server{
			Id:       string(server.ID),
			RpcAddr:  string(server.Address),
			IsLeader: c.raft.Leader() == server.Address,
		})
	}
	return servers, nil
}
