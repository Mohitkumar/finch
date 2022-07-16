package redis

import (
	"testing"
	"time"

	api_v1 "github.com/mohitkumar/finch/api/v1"
	"github.com/stretchr/testify/require"
)

func TestDelayQueue(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, queue *redisDelayQueue,
	){
		"test simple push":     testPushPop,
		"test push with delay": testPushPopDelay,
	} {
		t.Run(scenario, func(t *testing.T) {
			conf := &Config{
				Host:      "localhost",
				Port:      6379,
				Namespace: "test",
			}
			queue := NewRedisDelayQueue(*conf)

			fn(t, queue)
		})
	}
}

func testPushPop(t *testing.T, queue *redisDelayQueue) {
	err := queue.Push("test-delay", []byte("test_msg1"))
	require.NoError(t, err)
	time.Sleep(1 * time.Second)

	res, err := queue.Pop("test-delay")
	require.NoError(t, err)

	require.Equal(t, "test_msg1", res[0])

	_, err = queue.Pop("test-delay")
	_, ok := err.(api_v1.PollError)
	require.True(t, ok)
}

func testPushPopDelay(t *testing.T, queue *redisDelayQueue) {
	err := queue.PushWithDelay("test-delay", 5*time.Second, []byte("test_msg2"))
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	res, err := queue.Pop("test-delay")
	require.Error(t, err)
	_, ok := err.(api_v1.PollError)
	require.True(t, ok)

	time.Sleep(1 * time.Second)
	res, err = queue.Pop("test-delay")
	require.Error(t, err)
	_, ok = err.(api_v1.PollError)
	require.True(t, ok)

	time.Sleep(4 * time.Second)
	res, err = queue.Pop("test-delay")
	require.NoError(t, err)
	require.Equal(t, "test_msg2", res[0])
}
