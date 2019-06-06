package cache

import (
	"github.com/stretchr/testify/require"
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

var rc *Rediscache

func init() {
	rc = &Rediscache{
		pool: &redis.Pool{
			MaxIdle:     10,
			MaxActive:   50,
			IdleTimeout: 200,
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:6379") },
		},
	}
}

func TestCache(t *testing.T) {

	t.Run("check set and get value with ttl", func(t *testing.T) {
		err := rc.SetWithTTL("test", "value1", 500)
		require.NoError(t, err)

		val, err := rc.Get("test")
		require.NoError(t, err)
		require.Equal(t, "value1", val)
		time.Sleep(time.Second)

		val, err = rc.Get("testval")

		require.Error(t, err)
		require.Empty(t, val)

	})

	t.Run("add, read value and delete", func(t *testing.T) {
		err := rc.Add("test", "value2")
		require.NoError(t, err)

		val, err := rc.Read("test")
		require.NoError(t, err)
		require.Equal(t, "value2", val)

		val, err = rc.Read("test")
		require.Error(t, err)
		require.Empty(t, val)

	})

}

