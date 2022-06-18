package redis

import (
	"fmt"
	"strings"

	rd "github.com/go-redis/redis/v9"
)

type baseDao struct {
	redisClient *rd.Client
	namespace   string
}

func newBaseDao(conf Config) *baseDao {
	redisClient := rd.NewClient(&rd.Options{
		Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		DB:   0,
	})
	return &baseDao{
		redisClient: redisClient,
		namespace:   conf.Namespace,
	}
}

func (bs *baseDao) getNamespaceKey(args ...string) string {
	return fmt.Sprintf("%s:%s", bs.namespace, strings.Join(args, ":"))
}
