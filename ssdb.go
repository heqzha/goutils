package utils

import(
	"fmt"
	"time"

	"github.com/seefan/gossdb"
)

var(
	ssdbConn *gossdb.Connectors
)

func ssdbNewClient() (*gossdb.Client, error){
	if ssdbConn == nil{
		return nil, fmt.Errorf("Failed to initialize ssdb connection pool")
	}
	c, err := ssdbConn.NewClient()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func SSDBConfig(host string, port int, minPoolSize int, maxPoolSize, acquireIncr int) error{
	var err error
	ssdbConn, err = gossdb.NewPool(&gossdb.Config{
		Host:             host,
		Port:             port,
		MinPoolSize:      minPoolSize,
		MaxPoolSize:      maxPoolSize,
		AcquireIncrement: acquireIncr,
	})
	if err != nil{
		return err
	}

	//test connection
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()
	if !cli.Ping() {
		return fmt.Errorf("Failed to ping ssdb host: %s:%d", host, port)
	}
	return nil
}

func ssdbFormatSec(dur time.Duration) int64 {
	return int64(dur / time.Second)
}

func SSDBExists(key string)(bool, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return false, err
	}
	defer cli.Close()

	return cli.Exists(key)
}

func SSDBExpire(key string, exp time.Duration) error{
	if exp > 0{
		cli, err := ssdbNewClient()
		if err != nil{
			return err
		}
		defer cli.Close()

		success, err := cli.Expire(key, ssdbFormatSec(exp))
		if err != nil{
			return err
		}else if !success{
			return fmt.Errorf("Set expire to key %s failed", key)
		}
	}
	return nil
}

func SSDBDel(key string) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Del(key)
}

func SSDBSet(key string, value string) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Set(key, value)
}

func SSDBSetWithExp(key string, value string, exp time.Duration) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Set(key, value, ssdbFormatSec(exp))
}

func SSDBGet(key string) (string, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return "", err
	}
	defer cli.Close()

	value, err := cli.Get(key)
	if err != nil{
		return "", err
	}
	return value.String(), nil
}

////////////////
// SSDB Queue //
////////////////
func SSDBQPushBack(key string, values ...string) (int64, error){
	return SSDBQPushBackWithExpire(key, 0, values...)
}

func SSDBQPushBackWithExpire(key string, exp time.Duration, values ...string) (int64, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return 0, err
	}
	defer cli.Close()

	interfs := make([]interface{}, len(values))
	for i, v := range values {
		interfs[i] = v
	}
	size, err := cli.Qpush_back(key, interfs)
	if err != nil{
		return 0, err
	}
	return size, SSDBExpire(key, exp)
}

func SSDBQRangeAll(key string)([]string, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return nil, err
	}
	defer cli.Close()

	results, err := cli.Qslice(key, 0, -1)
	if err != nil{
		return nil, err
	}
	values := make([]string, len(results))
	for i, r := range results{
		values[i] = r.String()
	}
	return values, nil
}

func SSDBQClear(key string) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Qclear(key)
}

//////////////////
// SSDB Hashmap //
//////////////////
func SSDBHSet(key, field string, value string) error{
	return SSDBHSetWithExp(key, 0, field, value)
}

func SSDBHSetWithExp(key string, exp time.Duration, field string, value string) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()
	err = cli.Hset(key, field, value)
	if err != nil{
		return err
	}
	return SSDBExpire(key, exp)
}

func SSDBHGet(key, field string)(string, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return "", err
	}
	defer cli.Close()

	value, err := cli.Hget(key, field)
	if err != nil{
		return "", err
	}
	return value.String(), nil
}

func SSDBHSetMap(key string, values map[string]string) error{
	return SSDBHSetMapWithExp(key, 0, values)
}

func SSDBHSetMapWithExp(key string, exp time.Duration, values map[string]string)(error){
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	for k, v := range values{
		err := cli.Hset(key, k, v)
		if err != nil{
			return err
		}
	}

	return SSDBExpire(key, exp)
}

func SSDBHGetMap(key string) (map[string]string, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return nil, err
	}
	defer cli.Close()

	m, err := cli.HgetAll(key)
	if err != nil{
		return nil, err
	}
	values := map[string]string{}
	for k, v := range m{
		values[k] = v.String()
	}
	return values, nil
}

func SSDBHClear(key string) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Hclear(key)
}

func SSDBHDel(key, field string) error{
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Hdel(key, field)
}

func SSDBHSize(key string)(int64, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return 0, err
	}
	defer cli.Close()
	return cli.Hsize(key)
}

//////////////////////
// SSDB Sorted Sets //
//////////////////////
func SSDBZSet(key string, field string, score int64)(error){
	return SSDBZSetWithExp(key, 0, field, score)
}

func SSDBZSetWithExp(key string, exp time.Duration, field string, score int64)(error){
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	err = cli.Zset(key, field, score)
	if err != nil{
		return err
	}
	return SSDBExpire(key, exp)
}

func DBZGet(key, field string)(int64, error){
	cli, err := ssdbNewClient()
	if err != nil{
		return 0, err
	}
	defer cli.Close()

	return cli.Zget(key, field)
}

func SSDBZRScan(name, keyStart string, scoreStart, scoreEnd interface{}, limit int64) (keys []string, scores []int64, err error){
	cli, err := ssdbNewClient()
	if err != nil{
		return nil, nil, err
	}
	defer cli.Close()

	return cli.Zrscan(name, keyStart, scoreStart, scoreEnd, limit)
}

func SSDBZDel(name, key string)(error){
	cli, err := ssdbNewClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Zdel(name, key)
}
