package cache

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	defaultMaxIdle     = 5
	defaultMaxActive   = 20
	defaultIdleTimeout = time.Minute * 5
)

// Rediscache - объект для работы с redis
type Rediscache struct {
	pool *redis.Pool
}

// NewRediscache - конструктор Rediscache
func NewRediscache(addr string, maxIdle, maxActive int, idleTimeout time.Duration) *Rediscache {
	if maxIdle == 0 {
		maxIdle = defaultMaxIdle
	}

	if maxActive == 0 {
		maxActive = defaultMaxActive
	}

	if idleTimeout == 0 {
		idleTimeout = defaultIdleTimeout
	}

	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}

	return &Rediscache{
		pool: pool,
	}
}

func (r *Rediscache) GetConnection() redis.Conn {
	return r.pool.Get()
}

func (r *Rediscache) SetWithTTL(key, val string, ttl int) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val, "PX", ttl)
	return err
}

func (r *Rediscache) Get(key string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (r *Rediscache) Publish(channel, val string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PUBLISH", channel, val)
	return err

}

func (r *Rediscache) Add(key, val string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	return err

}

func (r *Rediscache) Read(key string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("SPOP", key))
}

func (r *Rediscache) Close() {
	r.pool.Close()
}
