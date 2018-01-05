package test

import "testing"
import "github.com/heqzha/goutils/db"
import "strconv"

import "time"

func TestRedisHandlerInit(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")
	t.Log(handler.Pool)

}

func TestRedisHandlerPing(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")

	if err := handler.Ping(); err != nil {
		t.Error(err)
		return
	}
	t.Log("Ping")
}

func TestRedisHandlerSetGet(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")
	key := "test_set"
	exist, err := handler.Exists(key)
	if err != nil {
		t.Error(err)
		return
	} else if exist {
		if err := handler.Delete(key); err != nil {
			t.Error(err)
			return
		}
		t.Logf("exist %s, deleted!", key)
	}

	if err := handler.Set(key, "123"); err != nil {
		t.Error(err)
		return
	}
	val, err := handler.Get(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(val))
	handler.Delete(key)
}

func TestRedisHandlerGetKeys(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")
	key := "test_set"

	if err := handler.Set(key, "123"); err != nil {
		t.Error(err)
		return
	}

	keys, err := handler.GetKeys("test_*")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(keys)

	handler.Delete(key)
}

func TestRedisHandlerIncr(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")
	key := "test_incr"
	cnter, err := handler.Incr(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cnter)
}

func TestRedisHandlerZ(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")
	key := "test_z"

	for index := 0; index < 3; index++ {
		if err := handler.Zadd(key, strconv.Itoa(index), int64(index)); err != nil {
			t.Error(err)
			return
		}
	}
	if err := handler.Zincrby(key, "test", 3); err != nil {
		t.Error(err)
		return
	}

	card, err := handler.Zcard(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(card)

	r, err := handler.Zrange(key, 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r)

	r, err = handler.Zrangebyscore(key, -1, 1, 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r)

	r, err = handler.ZrangebyscoreInf(key, 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r)

	r, err = handler.Zrevrangebyscore(key, -1, 1, 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r)

	r, err = handler.ZrevrangebyscoreInf(key, 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r)

	cnt, err := handler.Zcount(key, -1, 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cnt)

	cnt, err = handler.ZcountInf(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cnt)

	if err := handler.Zremrangebyrank(key, 0, 1); err != nil {
		t.Error(err)
		return
	}
	r, err = handler.Zrange(key, 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r)
}

func TestRedisHandlerExpire(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9600")
	key := "test_set"

	if err := handler.Set(key, "123"); err != nil {
		t.Error(err)
		return
	}
	if err := handler.Expire(key, 3); err != nil {
		t.Error(err)
		return
	}
	ttl, err := handler.Ttl(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ttl)
	val, err := handler.Get(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(val))

	// for exist, _ := handler.Exists(key); exist; {
	// 	fmt.Println("exist")
	// 	time.Sleep(1 * time.Second)
	// }
	time.Sleep(5 * time.Second)
	exist, err := handler.Exists(key)
	if err != nil {
		t.Error(err)
		return
	} else if !exist {
		t.Log("not exist")
		return
	}
	val, err = handler.Get(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(val))
}

func TestRedisHandlerL(t *testing.T) {
	handler := &db.RedisHandler{}
	defer handler.Close()
	handler.Init("127.0.0.1:9601")
	key := "test_set"

	for _, v := range []string{"1", "2", "safsdf"} {
		if err := handler.Rpush(key, v); err != nil {
			t.Error(err)
			return
		}
	}

	l, err := handler.Llen(key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("length of list: %d", l)
	for ; err == nil && l > 0; l, err = handler.Llen(key) {
		data, err := handler.Lpop(key)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(data)
	}
}
