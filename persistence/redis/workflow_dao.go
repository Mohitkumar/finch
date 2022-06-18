package redis

import (
	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence"
)

var _ persistence.WorkflowDao = new(redisWorkflowDao)

type redisWorkflowDao struct {
	baseDao
}

func NewRedisWorkflowDao(conf Config) *redisWorkflowDao {
	return &redisWorkflowDao{
		baseDao: *newBaseDao(conf),
	}
}

func (rfd *redisWorkflowDao) Save(wf model.Workflow) (bool, error) {

}

func (rfd *redisWorkflowDao) Delete(name string) (bool, error) {

}

func (rfd *redisWorkflowDao) Get(name string) (*model.Workflow, error) {

}
