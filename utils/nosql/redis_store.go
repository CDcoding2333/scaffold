package nosql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const redisMaxIdleConn = 20
const redisMaxActive = 20

type RedisStore struct {
	pool        *redis.Pool
	host        string
	port        int
	db          int
	withTimeout bool
}

func newRedisStore(host, password string, port int, db int, withTimeout bool) (*RedisStore, error) {
	f := func() (redis.Conn, error) {
		var c redis.Conn
		var err error
		if withTimeout {
			c, err = redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second*20, time.Second*10, time.Second*10)
		} else {
			c, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		}
		if err != nil {
			return nil, err
		}
		if password != "" {
			if _, err := c.Do("AUTH", password); err != nil {
				return nil, err
			}
		}

		if _, err := c.Do("SELECT", db); err != nil {
			return nil, err
		}

		return c, err
	}
	pool := redis.NewPool(f, redisMaxIdleConn)
	pool.MaxActive = redisMaxActive
	pool.Wait = true

	store := &RedisStore{pool: pool, host: host, port: port, db: db, withTimeout: withTimeout}
	return store, nil
}

func NewRedisStore(host, password string, port int, db int) (*RedisStore, error) {
	return newRedisStore(host, password, port, db, true)
}

func NewRedisStoreWithoutTimeout(host, password string, port int, db int) (*RedisStore, error) {
	return newRedisStore(host, password, port, db, false)
}

func (r *RedisStore) SetMaxIdle(maxIdle int) {
	r.pool.MaxIdle = maxIdle
}

func (r *RedisStore) SetMaxActive(maxActive int) {
	r.pool.MaxActive = maxActive
}

func (r *RedisStore) Get(key string) (interface{}, error) {
	return r.do("GET", key)
}

func (r *RedisStore) Set(key string, value interface{}) error {
	_, err := r.do("SET", key, value)
	return err
}

func (r *RedisStore) SetNx(key string, value interface{}) (int64, error) {
	return r.Int64(r.do("SETNX", key, value))
}

func (r *RedisStore) SetEx(key string, value interface{}, timeout int64) error {
	_, err := r.do("SETEX", key, timeout, value)
	return err
}

// nil表示成功，ErrNil表示数据库内已经存在这个key，其他表示数据库发生错误
func (r *RedisStore) SetNxEx(key string, value interface{}, timeout int64) error {
	_, err := r.String(r.do("SET", key, value, "NX", "EX", timeout))
	return err
}

func (r *RedisStore) ExpireAt(key string, timestamp int64) (int64, error) {
	return r.Int64(r.do("EXPIREAT", key, timestamp))
}

func (r *RedisStore) Del(keys ...string) (int64, error) {
	ks := make([]interface{}, len(keys))
	for i, key := range keys {
		ks[i] = key
	}
	return r.Int64(r.do("DEL", ks...))
}

func (r *RedisStore) Incr(key string) (int64, error) {
	return r.Int64(r.do("INCR", key))
}

func (r *RedisStore) IncrBy(key string, delta int64) (int64, error) {
	return r.Int64(r.do("INCRBY", key, delta))
}

func (r *RedisStore) Expire(key string, duration int64) (int64, error) {
	return r.Int64(r.do("EXPIRE", key, duration))
}

func (r *RedisStore) HGet(key string, field string) (interface{}, error) {
	return r.do("HGET", key, field)
}

func (r *RedisStore) HLen(key string) (int64, error) {
	return r.Int64(r.do("HLEN", key))
}

func (r *RedisStore) HSet(key string, field string, val interface{}) error {
	_, err := r.do("HSET", key, field, val)
	return err
}

func (r *RedisStore) HDel(key string, fields ...string) (int64, error) {
	ks := make([]interface{}, len(fields)+1)
	ks[0] = key
	for i, key := range fields {
		ks[i+1] = key
	}
	return r.Int64(r.do("HDEL", ks...))
}

func (r *RedisStore) HClear(key string) error {
	_, err := r.do("DEL", key)
	return err
}

func (r *RedisStore) HMGet(key string, fields ...string) (interface{}, error) {
	if len(fields) == 0 {
		return nil, ErrNil
	}
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	return r.do("HMGET", args...)

}

func (r *RedisStore) HMSet(key string, kvs ...interface{}) error {
	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		return ErrWrongArgsNum
	}
	args := make([]interface{}, len(kvs)+1)
	args[0] = key
	for i := 0; i < len(kvs); i += 2 {
		if _, ok := kvs[i].(string); !ok {
			return errors.New("field must be string")
		}
		args[i+1] = kvs[i]
		args[i+2] = kvs[i+1]
	}
	_, err := r.do("HMSET", args...)
	return err
}

func (r *RedisStore) HExpire(key string, duration int64) error {
	_, err := r.do("EXPIRE", key, duration)
	return err
}

func (r *RedisStore) HGetAll(key string) (map[string]interface{}, error) {
	//TODO
	vals, err := r.Values(r.do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	num := len(vals) / 2
	result := make(map[string]interface{}, num)
	for i := 0; i < num; i++ {
		key, _ := r.String(vals[2*i], nil)
		result[key] = vals[2*i+1]
	}
	return result, nil
}

func (r *RedisStore) HIncrBy(key, field string, delta int64) (int64, error) {
	return r.Int64(r.do("HINCRBY", key, field, delta))
}

func (r *RedisStore) ZAdd(key string, kvs ...interface{}) (int64, error) {
	if len(kvs) == 0 {
		return 0, nil
	}
	if len(kvs)%2 != 0 {
		return 0, errors.New("args num error")
	}
	args := make([]interface{}, len(kvs)+1)
	args[0] = key
	for i := 0; i < len(kvs); i += 2 {
		args[i+1] = kvs[i]
		args[i+2] = kvs[i+1]
	}
	return r.Int64(r.do("ZAdd", args...))
}

func (r *RedisStore) ZRem(key string, members ...string) (int64, error) {
	args := make([]interface{}, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	return r.Int64(r.do("ZREM", args...))
}

func (r *RedisStore) ZClear(key string) error {
	_, err := r.do("DEL", key)
	return err
}

func (r *RedisStore) ZExpire(key string, duration int64) error {
	_, err := r.do("EXPIRE", key, duration)
	return err
}

func (r *RedisStore) ZRange(key string, min, max int64, withScores bool) (interface{}, error) {
	if withScores {
		return r.do("ZRANGE", key, min, max, "WITHSCORES")
	} else {
		return r.do("ZRANGE", key, min, max)
	}
}

func (r *RedisStore) ZRangeByScoreWithScore(key string, min, max int64) (map[string]int64, error) {
	vals, err := r.Values(r.do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	n := len(vals) / 2
	result := make(map[string]int64, n)
	for i := 0; i < n; i++ {
		key, _ := r.String(vals[2*i], nil)
		score, _ := r.String(vals[2*i+1], nil)
		v, _ := strconv.ParseFloat(score, 64)
		result[key] = int64(v)
	}
	return result, nil
}

func (r *RedisStore) LRange(key string, start, stop int64) (interface{}, error) {
	return r.do("LRANGE", key, start, stop)
}

func (r *RedisStore) LRem(key string, count int64, value interface{}) (int64, error) {
	return r.Int64(r.do("LREM", key, count, value))
}

func (r *RedisStore) LLen(key string) (int64, error) {
	return r.Int64(r.do("LLEN", key))
}

func (r *RedisStore) TTl(key string) (int64, error) {
	return r.Int64(r.do("TTL", key))
}

func (r *RedisStore) LPop(key string) (interface{}, error) {
	return r.do("LPOP", key)
}

func (r *RedisStore) RPop(key string) (interface{}, error) {
	return r.do("RPOP", key)
}

func (r *RedisStore) BLPop(key string, timeout int) (interface{}, error) {
	vals, err := r.Values(r.do("BLPOP", key, timeout))
	if err != nil {
		return nil, err
	}
	return vals[1], nil
}

func (r *RedisStore) BRPop(key string, timeout int) (interface{}, error) {
	vals, err := r.Values(r.do("BRPOP", key, timeout))
	if err != nil {
		return nil, err
	}
	return vals[1], nil
}

func (r *RedisStore) LPush(key string, value ...interface{}) error {
	args := make([]interface{}, len(value)+1)
	args[0] = key
	for i, v := range value {
		args[i+1] = v
	}
	_, err := r.do("LPUSH", args...)
	return err
}

func (r *RedisStore) RPush(key string, value ...interface{}) error {
	args := make([]interface{}, len(value)+1)
	args[0] = key
	for i, v := range value {
		args[i+1] = v
	}
	_, err := r.do("RPUSH", args...)
	return err
}

func (r *RedisStore) MGet(key ...interface{}) ([]interface{}, error) {
	values, err := r.do("MGET", key...)
	if err != nil {
		return nil, err
	}

	return values.([]interface{}), nil
}

func (r *RedisStore) SAdd(key string, members ...interface{}) (int64, error) {
	args := make([]interface{}, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	return r.Int64(r.do("SADD", args...))
}

func (r *RedisStore) SCard(key string) (int64, error) {
	return r.Int64(r.do("SCARD", key))
}

func (r *RedisStore) SIsMember(key string, member interface{}) (bool, error) {
	return r.Bool(r.do("SISMEMBER", key, member))
}

func (r *RedisStore) SMembers(key string) (interface{}, error) {
	values, err := r.do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (r *RedisStore) SRem(key string, members ...interface{}) (int64, error) {
	args := make([]interface{}, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[1+i] = m
	}
	return r.Int64(r.do("SREM", args...))
}

func (r *RedisStore) Subscribe(ctx context.Context, onMessage func(channel string, data []byte) error, channels ...interface{}) error {
	const healthCheckPeriod = time.Second * 10

	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", r.host, r.port),
		// Read timeout on server should be greater than ping period.
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

loop:
	for err == nil {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}
	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	psc.Unsubscribe()

	// Wait for goroutine to complete.
	return <-done
}

func (r *RedisStore) Publish(key string, members ...interface{}) (int64, error) {
	args := make([]interface{}, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	return r.Int64(r.do("PUBLISH", args...))
}

func (r *RedisStore) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	res, err := conn.Do(cmd, args...)
	if err == redis.ErrNil {
		return nil, ErrNil
	}
	return res, err
}

func (r *RedisStore) GetPool() *redis.Pool {
	return r.pool
}
