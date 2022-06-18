package redis

import (
	"fmt"

	rd "github.com/go-redis/redis/v9"
)

type baseDao struct {
	redisClient *rd.Client
}

func newBaseDao(conf Config) *baseDao {
	redisClient := rd.NewClient(&rd.Options{
		Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		DB:   0,
	})
	return &baseDao{
		redisClient: redisClient,
	}
}
