package redis

import (
	"context"
	"strconv"

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
	flowCtx := &api.FlowContext{
		Id:                 flowId,
		WorkflowState:      api.FlowContext_RUNNING,
		CurrentActionState: api.FlowContext_A_RUNNING,
		CurrentAction:      int32(action),
		Data:               util.ConvertToProto(dataMap),
	}
	if err := rf.SaveFlowContext(wFname, flowId, flowCtx); err != nil {
		return nil, err
	}

	return flowCtx, nil
}

func (rf *redisFlowDao) UpdateFlowContextData(wFname string, flowId string, action int, dataMap map[string]any) (*api.FlowContext, error) {
	flowCtx, err := rf.GetFlowContext(wFname, flowId)
	if err != nil {
		return nil, err
	}
	data := flowCtx.GetData()
	data[strconv.Itoa(action)] = util.ConvertMapToStructPb(dataMap)
	if err := rf.SaveFlowContext(wFname, flowId, flowCtx); err != nil {
		return nil, err
	}
	return flowCtx, nil
}

func (rf *redisFlowDao) SaveFlowContext(wfName string, flowId string, flowCtx *api.FlowContext) error {
	key := rf.baseDao.getNamespaceKey(WORKFLOW_KEY, wfName)
	ctx := context.Background()
	data, err := proto.Marshal(flowCtx)
	if err != nil {
		return err
	}
	if err := rf.baseDao.redisClient.HSet(ctx, key, []string{flowId, string(data)}).Err(); err != nil {
		return err
	}
	return nil
}

func (rf *redisFlowDao) GetFlowContext(wfName string, flowId string) (*api.FlowContext, error) {
	key := rf.baseDao.getNamespaceKey(WORKFLOW_KEY, wfName)
	ctx := context.Background()
	flowCtxStr, err := rf.baseDao.redisClient.HGet(ctx, key, flowId).Result()
	if err != nil {
		return nil, err
	}
	var flowCtx *api.FlowContext
	if err := proto.Unmarshal([]byte(flowCtxStr), flowCtx); err != nil {
		return nil, err
	}
	return flowCtx, nil
}
