package db

import(
	"fmt"
	"time"
	"math"

	"github.com/seefan/gossdb"
)

type SSDBHandler struct{
	conn *gossdb.Connectors
}

func SSDBNewHandlerDefault(host string, port int) (*SSDBHandler, error){
	return SSDBNewHandler(&gossdb.Config{
		Host:             host,
		Port:             port,
	})
}

func SSDBNewHandler(conf *gossdb.Config) (*SSDBHandler, error){
	conn, err := gossdb.NewPool(conf)
	if err != nil{
		return nil, err
	}

	//test connection
	handler := &SSDBHandler{
		conn: conn,
	}
	cli, err := handler.newClient()
	if err != nil{
		return nil, err
	}
	defer cli.Close()
	if !cli.Ping() {
		return nil, fmt.Errorf("Failed to ping ssdb host: %s:%d", conf.Host, conf.Port)
	}
	return handler, nil
}

func (h *SSDBHandler)newClient()(*gossdb.Client, error){
	if h.conn == nil{
		return nil, fmt.Errorf("Failed to initialize ssdb connection pool")
	}
	c, err := h.conn.NewClient()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ssdbFormatSec(dur time.Duration) int64 {
	return int64(dur / time.Second)
}

func (h *SSDBHandler)Exists(key string)(bool, error){
	cli, err := h.newClient()
	if err != nil{
		return false, err
	}
	defer cli.Close()

	return cli.Exists(key)
}

func (h *SSDBHandler)Expire(key string, exp time.Duration) error{
	if exp > 0{
		cli, err := h.newClient()
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

func (h *SSDBHandler)Del(key string) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Del(key)
}

func (h *SSDBHandler)Set(key string, value string) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Set(key, value)
}

func (h *SSDBHandler)SetWithExp(key string, value string, exp time.Duration) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Set(key, value, ssdbFormatSec(exp))
}

func (h *SSDBHandler)Get(key string) (string, error){
	cli, err := h.newClient()
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
func (h *SSDBHandler)QPushBack(key string, values ...string) (int64, error){
	return h.QPushBackWithExpire(key, 0, values...)
}

func (h *SSDBHandler)QPushBackWithExpire(key string, exp time.Duration, values ...string) (int64, error){
	cli, err := h.newClient()
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
	return size, h.Expire(key, exp)
}

func (h *SSDBHandler)QRangeAll(key string)([]string, error){
	cli, err := h.newClient()
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

func (h *SSDBHandler)QClear(key string) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Qclear(key)
}

//////////////////
// SSDB Hashmap //
//////////////////
func (h *SSDBHandler)HSet(key, field string, value string) error{
	return h.HSetWithExp(key, 0, field, value)
}

func (h *SSDBHandler)HSetWithExp(key string, exp time.Duration, field string, value string) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()
	err = cli.Hset(key, field, value)
	if err != nil{
		return err
	}
	return h.Expire(key, exp)
}

func (h *SSDBHandler)HGet(key, field string)(string, error){
	cli, err := h.newClient()
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

func (h *SSDBHandler)HSetMap(key string, values map[string]string) error{
	return h.HSetMapWithExp(key, 0, values)
}

func (h *SSDBHandler)HSetMapWithExp(key string, exp time.Duration, values map[string]string)(error){
	cli, err := h.newClient()
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

	return h.Expire(key, exp)
}

func (h *SSDBHandler)HGetMap(key string) (map[string]string, error){
	cli, err := h.newClient()
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

func (h *SSDBHandler)HClear(key string) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Hclear(key)
}

func (h *SSDBHandler)HDel(key, field string) error{
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Hdel(key, field)
}

func (h *SSDBHandler)HSize(key string)(int64, error){
	cli, err := h.newClient()
	if err != nil{
		return 0, err
	}
	defer cli.Close()
	return cli.Hsize(key)
}

func (h *SSDBHandler)HAllFields(key string)([]string, error){
	return h.HFields(key, "", "", math.MaxInt64)
}

func (h *SSDBHandler)HFields(key, fieldStart, fieldEnd string, limit int64)([]string, error){
	cli, err := h.newClient()
	if err != nil{
		return nil, err
	}
	defer cli.Close()
	return cli.Hkeys(key, fieldStart, fieldEnd, limit)
}

//////////////////////
// SSDB Sorted Sets //
//////////////////////
func (h *SSDBHandler)ZSet(key string, field string, score int64)(error){
	return h.ZSetWithExp(key, 0, field, score)
}

func (h *SSDBHandler)ZSetWithExp(key string, exp time.Duration, field string, score int64)(error){
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	err = cli.Zset(key, field, score)
	if err != nil{
		return err
	}
	return h.Expire(key, exp)
}

func (h *SSDBHandler)ZGet(key, field string)(int64, error){
	cli, err := h.newClient()
	if err != nil{
		return 0, err
	}
	defer cli.Close()

	return cli.Zget(key, field)
}

func (h *SSDBHandler)ZTopX(key, keyStart string, start, limit int64)([]string,  []int64, error){
	return h.ZRScan(key, keyStart, start, "", limit)
}

func (h *SSDBHandler)ZRScan(key, fieldStart string, start, end interface{}, limit int64) (keys []string, scores []int64, err error){
	cli, err := h.newClient()
	if err != nil{
		return nil, nil, err
	}
	defer cli.Close()
	return cli.Zrscan(key, fieldStart, start, end, limit)
}

func (h *SSDBHandler)ZExists(key, field string) (bool, error){
	cli, err := h.newClient()
	if err != nil{
		return false, err
	}
	defer cli.Close()

	return cli.Zexists(key, field)
}

func (h *SSDBHandler)ZDel(key, field string)(error){
	cli, err := h.newClient()
	if err != nil{
		return err
	}
	defer cli.Close()

	return cli.Zdel(key, field)
}
