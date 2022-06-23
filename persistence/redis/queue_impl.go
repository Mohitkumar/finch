package redis

import (
	"context"

	"github.com/mohitkumar/finch/persistence"
)

type redisQueue struct {
	baseDao
}

var _ persistence.Queue = new(redisQueue)

func NewRedisQueue(config Config) *redisQueue {
	return &redisQueue{
		baseDao: *newBaseDao(config),
	}
}
func (rq *redisQueue) Push(queueName string, mesage []byte) error {
	queueName = rq.getNamespaceKey(queueName)
	ctx := context.Background()

	return rq.redisClient.LPush(ctx, queueName, mesage).Err()
}

func (rq *redisQueue) Pop(queueName string) ([]byte, error) {
	queueName = rq.getNamespaceKey(queueName)
	ctx := context.Background()
	res, err := rq.redisClient.LPop(ctx, queueName).Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}
