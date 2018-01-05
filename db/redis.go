package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisHandler struct {
	Pool *redis.Pool
}

func (h *RedisHandler) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := h.Pool.Get()
	defer conn.Close()

	data, err := conn.Do(cmd, args...)
	if err != nil {
		return nil, fmt.Errorf("%s args %v: %v", cmd, args, err)
	}
	return data, nil
}

func (h *RedisHandler) Init(addr string) {
	h.Close()
	h.Pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 120 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr,
				redis.DialConnectTimeout(time.Duration(100)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(100)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(100)*time.Millisecond))
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (h *RedisHandler) Close() {
	if h.Pool != nil {
		h.Pool.Close()
		h.Pool = nil
	}
}

func (h *RedisHandler) Ping() error {

	conn := h.Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

func (h *RedisHandler) Get(key string) (string, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, nil
}

func (h *RedisHandler) Set(key string, value string) error {

	conn := h.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := value
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func (h *RedisHandler) Exists(key string) (bool, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func (h *RedisHandler) Delete(key string) error {

	conn := h.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (h *RedisHandler) GetKeys(pattern string) ([]string, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (h *RedisHandler) Incr(key string) (int, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

func (h *RedisHandler) Llen(key string) (int64, error) {
	data, err := h.do("LLEN", key)
	if err != nil {
		return int64(0), err
	}
	return data.(int64), nil
}

func (h *RedisHandler) Lpop(key string) (string, error) {
	data, err := redis.String(h.do("LPOP", key))
	if err != nil {
		return "", err
	}
	return data, nil
}

func (h *RedisHandler) Rpush(key string, value string) error {
	_, err := h.do("RPUSH", key, value)
	return err
}

func (h *RedisHandler) Zadd(key string, value string, score int64) error {
	conn := h.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error zadd key %s to %s: err  %v", key, v, conn.Err())
	}
	return err
}

func (h *RedisHandler) Zcard(key string) (int, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZCARD", key)
	if err != nil {
		return -1, fmt.Errorf("error zcard key %s %v", key, err)
	}
	return int(res.(int64)), nil
}

func (h *RedisHandler) mapZrangeResults(ress []interface{}) ([]map[string]int64, error) {
	results := []map[string]int64{}
	for idx, _ := range ress {
		if idx&0x1 == 0 {
			score, err := strconv.ParseInt(string(ress[idx+1].([]byte)), 10, 64)
			if err != nil {
				return nil, err
			}
			val := map[string]int64{
				string(ress[idx].([]byte)): score,
			}
			results = append(results, val)
		}
	}
	return results, nil
}

func (h *RedisHandler) Zrange(key string, offset, limit int) ([]map[string]int64, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZRANGE", key, offset, offset*limit+limit, "WITHSCORES")
	if err != nil {
		return nil, fmt.Errorf("error zrange key %s %v", key, err)
	}
	return h.mapZrangeResults(res.([]interface{}))
}

func (h *RedisHandler) Zrangebyscore(key string, min, max int64, offset, limit int) ([]map[string]int64, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZRANGEBYSCORE", key, "("+strconv.FormatInt(min, 10), strconv.FormatInt(max, 10), "WITHSCORES", "LIMIT", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error zrangebyscore key %s %v", key, err)
	}
	return h.mapZrangeResults(res.([]interface{}))
}

func (h *RedisHandler) ZrangebyscoreInf(key string, offset, limit int) ([]map[string]int64, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZRANGEBYSCORE", key, "-inf", "+inf", "WITHSCORES", "LIMIT", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error zrangebyscore key %s %v", key, err)
	}
	return h.mapZrangeResults(res.([]interface{}))
}

func (h *RedisHandler) Zrevrangebyscore(key string, min, max int64, offset, limit int) ([]map[string]int64, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZREVRANGEBYSCORE", key, strconv.FormatInt(max, 10), "("+strconv.FormatInt(min, 10), "WITHSCORES", "LIMIT", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error zrangebyscore key %s %v", key, err)
	}
	return h.mapZrangeResults(res.([]interface{}))
}

func (h *RedisHandler) ZrevrangebyscoreInf(key string, offset, limit int) ([]map[string]int64, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES", "LIMIT", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error zrangebyscore key %s %v", key, err)
	}
	return h.mapZrangeResults(res.([]interface{}))
}

func (h *RedisHandler) Zcount(key string, min, max int64) (int, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZCOUNT", key, "("+strconv.FormatInt(min, 10), strconv.FormatInt(max, 10))
	if err != nil {
		return -1, fmt.Errorf("error zcount key %s %v", key, err)
	}
	return int(res.(int64)), nil
}

func (h *RedisHandler) ZcountInf(key string) (int, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("ZCOUNT", key, "-inf", "+inf")
	if err != nil {
		return -1, fmt.Errorf("error zcount key %s %v", key, err)
	}
	return int(res.(int64)), nil
}

func (h *RedisHandler) Zremrangebyrank(key string, start, stop int64) error {

	conn := h.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZREMRANGEBYRANK", key, start, stop)
	if err != nil {
		return fmt.Errorf("error zremrangebyrank key %s :%v", key, err)
	}
	return err
}

func (h *RedisHandler) Zincrby(key string, value string, score int64) error {
	conn := h.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", key, score, value)
	if err != nil {
		v := value
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error zadd key %s to %s: err  %v", key, v, conn.Err())
	}
	return err
}

func (h *RedisHandler) Expire(key string, seconds int64) error {

	conn := h.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return fmt.Errorf("error expire key %s to %d: %v", key, seconds, err)
	}
	return err
}

func (h *RedisHandler) Ttl(key string) (int, error) {

	conn := h.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("TTL", key)
	if err != nil {
		return -1, fmt.Errorf("error ttl key %s %v", key, err)
	}
	return int(res.(int64)), nil
}
