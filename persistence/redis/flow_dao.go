package redis

import (
	"context"
	"fmt"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence"
	"github.com/mohitkumar/finch/util"
	"go.uber.org/zap"
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
func (rf *redisFlowDao) CreateAndSaveFlowContext(wFname string, flowId string, action int, input map[string]any) (*api.FlowContext, error) {
	dataMap := make(map[string]any)
	dataMap["input"] = input
	flowCtx := &api.FlowContext{
		Id:            flowId,
		WorkflowState: api.FlowContext_RUNNING,
		CurrentAction: int32(action),
		Data:          util.ConvertToProto(dataMap),
	}
	if err := rf.saveFlowContext(wFname, flowId, flowCtx); err != nil {
		return nil, err
	}

	return flowCtx, nil
}

func (rf *redisFlowDao) AddActionOutputToFlowContext(wFname string, flowId string, action int, dataMap map[string]any) (*api.FlowContext, error) {
	flowCtx, err := rf.GetFlowContext(wFname, flowId)
	if err != nil {
		return nil, err
	}
	data := flowCtx.GetData()
	output := make(map[string]any)
	output["output"] = dataMap
	data[fmt.Sprintf("%d", action)] = util.ConvertMapToStructPb(output)
	if err := rf.saveFlowContext(wFname, flowId, flowCtx); err != nil {
		return nil, err
	}
	return flowCtx, nil
}

func (rf *redisFlowDao) saveFlowContext(wfName string, flowId string, flowCtx *api.FlowContext) error {
	key := rf.baseDao.getNamespaceKey(WORKFLOW_KEY, wfName)
	ctx := context.Background()
	data, err := proto.Marshal(flowCtx)
	if err != nil {
		return err
	}
	if err := rf.baseDao.redisClient.HSet(ctx, key, []string{flowId, string(data)}).Err(); err != nil {
		logger.Error("error in saving flow context", zap.String("flowName", wfName), zap.String("flowId", flowId), zap.Error(err))
		return api.StorageLayerError{}
	}
	return nil
}

func (rf *redisFlowDao) UpdateFlowContextNextAction(wfName string, flowId string, flowCtx *api.FlowContext, nextAction int) error {
	flowCtx.NextAction = int32(nextAction)
	return rf.saveFlowContext(wfName, flowId, flowCtx)
}

func (rf *redisFlowDao) UpdateFlowStatus(wfName string, flowId string, flowCtx *api.FlowContext, flowState api.FlowContext_WorkflowState) error {
	flowCtx.WorkflowState = flowState
	return rf.saveFlowContext(wfName, flowId, flowCtx)
}
func (rf *redisFlowDao) GetFlowContext(wfName string, flowId string) (*api.FlowContext, error) {
	key := rf.baseDao.getNamespaceKey(WORKFLOW_KEY, wfName)
	ctx := context.Background()
	flowCtxStr, err := rf.baseDao.redisClient.HGet(ctx, key, flowId).Result()
	if err != nil {
		logger.Error("error in getting flow context", zap.String("flowName", wfName), zap.String("flowId", flowId), zap.Error(err))
		return nil, api.StorageLayerError{}
	}
	flowCtx := &api.FlowContext{}
	if err := proto.Unmarshal([]byte(flowCtxStr), flowCtx); err != nil {
		return nil, err
	}
	return flowCtx, nil
}
