package coordinator

import (
	"strings"

	api "github.com/mohitkumar/finch/api/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func getFlowKeyPrefix(name string) string {
	return strings.Join([]string{"flow", name}, "_")
}

func (c *Coordinator) CreateFlow(flow *api.Flow) (*api.FlowCreateResponse, error) {
	c.logger.Debug("create flow request recevied")
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

func (c *Coordinator) GetFlow(req *api.FlowGetRequest) (*api.FlowGetResponse, error) {
	c.logger.Debug("get flow request recevied")
	key := []byte(getFlowKeyPrefix(req.Name))
	res, err := c.get(key)
	if err != nil {
		return nil, err
	}
	flow := &api.Flow{}
	if err := proto.Unmarshal(res, flow); err != nil {
		return nil, err
	}
	return &api.FlowGetResponse{Flow: flow}, nil
}

func (c *Coordinator) GetServers() ([]*api.Server, error) {
	future := c.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return nil, err
	}
	var servers []*api.Server
	var srvs []string
	for _, server := range future.Configuration().Servers {
		servers = append(servers, &api.Server{
			Id:       string(server.ID),
			RpcAddr:  string(server.Address),
			IsLeader: c.raft.Leader() == server.Address,
		})
		srvs = append(srvs, string(server.ID))
	}
	c.logger.Info("get server result", zap.String("servers", strings.Join(srvs, ",")))
	return servers, nil
}
