package redis

import (
	"context"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/persistence"
	"github.com/mohitkumar/finch/util"
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
func (rf *redisFlowDao) CreateAndSaveFlowContext(wFname string, flowId string, action int, dataMap map[string]any) (*api.FlowContext, error) {
	key := rf.baseDao.getNamespaceKey(WORKFLOW_KEY, wFname)
	ctx := context.Background()
	flowCtx := &api.FlowContext{
		Id:                 flowId,
		WorkflowState:      api.FlowContext_RUNNING,
		CurrentActionState: api.FlowContext_A_RUNNING,
		CurrentAction:      int32(action),
		Data:               util.ConvertToProto(dataMap),
	}
	data, err := proto.Marshal(flowCtx)
	if err != nil {
		return nil, err
	}
	if err := rf.baseDao.redisClient.HSet(ctx, key, []string{flowId, string(data)}).Err(); err != nil {
		return nil, err
	}
	return flowCtx, nil
}
