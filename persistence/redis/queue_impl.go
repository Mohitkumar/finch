package redis

import (
	"context"

	rd "github.com/go-redis/redis/v9"
	"github.com/mohitkumar/finch/persistence"
)

type redisQueue struct {
	redisClient *rd.Client
}

var _ persistence.Queue = new(redisQueue)

func (rq *redisQueue) Push(queueName string, mesage []byte) error {
	ctx := context.Background()

	return rq.redisClient.LPush(ctx, queueName, mesage).Err()
}

func (rq *redisQueue) Pop(queuName string) ([]byte, error) {
	ctx := context.Background()
	res, err := rq.redisClient.LPop(ctx, queuName).Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}
