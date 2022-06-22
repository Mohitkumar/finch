package redis

import (
	"context"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence"
	"google.golang.org/protobuf/proto"
)

const WORKFLOW_KEY string = "WF"

var _ persistence.FlowDao = new(redisFlowDao)

type redisFlowDao struct {
	baseDao
}

func NewRedisFlowDao(conf Config) *redisFlowDao {
	return &redisFlowDao{
		baseDao: *newBaseDao(conf),
	}
}
func (rf *redisFlowDao) CreateAndSaveFlowContext(name string, action int, flow model.Flow) (*api.FlowContext, error) {
	key := rf.baseDao.getNamespaceKey(WORKFLOW_KEY, name)
	ctx := context.Background()
	flowCtx := &api.FlowContext{
		Id:                 flow.Id,
		WorkflowState:      api.FlowContext_RUNNING,
		CurrentActionState: api.FlowContext_A_RUNNING,
		CurrentAction:      int32(action),
	}
	data, err := proto.Marshal(flowCtx)
	if err != nil {
		return nil, err
	}
	if err := rf.baseDao.redisClient.HSet(ctx, key, []string{flow.Id, string(data)}).Err(); err != nil {
		return nil, err
	}
	return flowCtx, nil
}
