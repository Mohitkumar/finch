package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	rd "github.com/go-redis/redis/v9"
	api_v1 "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence"
	"go.uber.org/zap"
)

type redisDelayQueue struct {
	baseDao
}

var _ persistence.DelayQueue = new(redisDelayQueue)

func NewRedisDelayQueue(config Config) *redisDelayQueue {
	return &redisDelayQueue{
		baseDao: *newBaseDao(config),
	}
}

func (rq *redisDelayQueue) Push(queueName string, message []byte) error {
	queueName = rq.getNamespaceKey(queueName)
	ctx := context.Background()
	currentTime := time.Now().UnixMilli()
	member := rd.Z{
		Score:  float64(currentTime),
		Member: message,
	}
	err := rq.redisClient.ZAdd(ctx, queueName, member).Err()
	if err != nil {
		logger.Error("error while push to redis list", zap.String("queue", queueName), zap.Error(err))
		return api_v1.StorageLayerError{}
	}
	return nil
}

func (rq *redisDelayQueue) PushWithDelay(queueName string, delay time.Duration, message []byte) error {
	queueName = rq.getNamespaceKey(queueName)
	ctx := context.Background()
	currentTime := time.Now().Add(delay).UnixMilli()
	member := rd.Z{
		Score:  float64(currentTime),
		Member: message,
	}
	err := rq.redisClient.ZAdd(ctx, queueName, member).Err()
	if err != nil {
		logger.Error("error while push to redis list", zap.String("queue", queueName), zap.Error(err))
		return api_v1.StorageLayerError{}
	}
	return nil
}

func (rq *redisDelayQueue) Pop(queueName string) ([]string, error) {
	queueName = rq.getNamespaceKey(queueName)
	ctx := context.Background()
	currentTime := time.Now().UnixMilli()
	pipe := rq.redisClient.Pipeline()

	zr := pipe.ZRange(ctx, queueName, 0, currentTime)
	pipe.ZRemRangeByScore(ctx, queueName, strconv.Itoa(0), strconv.FormatInt(currentTime, 10))

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error("error while pop from redis list", zap.String("queue", queueName), zap.Error(err))

		return nil, api_v1.StorageLayerError{}
	}

	res, err := zr.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, api_v1.PollError{QueueName: queueName}
		}
		logger.Error("error while pop from redis list", zap.String("queue", queueName), zap.Error(err))

		return nil, api_v1.StorageLayerError{}
	}
	return res, nil
}
