package cache

import(
	"time"
	"gopkg.in/redis.v3"
)

var(
	redisCli *redis.Client
)

func RedisConfig(address, password string, db int64) error{
	redisCli = redis.NewClient(&redis.Options{
		Addr: address,
		Password: password,
		DB: db,
	})

	return redisCli.Ping().Err()
}

func RedisExists(key string)(bool, error){
	return redisCli.Exists(key).Result()
}


func RedisExpire(key string, exp time.Duration) error{
	if exp > 0{
		return redisCli.Expire(key, exp).Err()
	}
	return nil
}

func RedisDel(keys ...string) error {
	return redisCli.Del(keys...).Err()
}

func RedisSet(key, value string) error{
	return RedisSetWithExp(key, 0, value)
}

func RedisSetWithExp(key string, exp time.Duration, value string) error{
	return redisCli.Set(key, value, exp).Err()
}

func RedisGet(key string)(string, error){
	value, err := redisCli.Get(key).Result()
	if err == redis.Nil {
		//key does not exist
		return "", nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

////////////////
// Redis List //
////////////////

func RedisLPush(key string, values ...string) error{
	return RedisLPushWithExp(key, 0, values...)
}

func RedisLPushWithExp(key string, exp time.Duration, values ...string) error{
	err := redisCli.LPush(key, values...).Err()
	if err != nil {
		return err
	}
	return RedisExpire(key, exp)
}

func RedisRPush(key string, values ...string) error{
	return RedisRPushWithExp(key, 0, values...)
}

func RedisRPushWithExp(key string, exp time.Duration, values ...string) error{
	err := redisCli.RPush(key, values...).Err()
	if err != nil{
		return err
	}
	return RedisExpire(key, exp)
}

func RedisLRangeAll(key string)([]string, error){
	length, err := redisCli.LLen(key).Result()
	if err != nil {
		return nil, err
	}

	return redisCli.LRange(key, 0, length).Result()
}

func RedisLRem(key string, count int64, value string) error {
	return redisCli.LRem(key, count, value).Err()
}

//////////////////
// Redis Hashes //
//////////////////

func RedisHSet(key, field, value string) error{
	return RedisHSetWithExp(key, 0, field, value)
}

func RedisHSetWithExp(key string, exp time.Duration, field, value string) error{
	err := redisCli.HSet(key, field, value).Err()
	if err != nil{
		return err
	}
	return RedisExpire(key, exp)
}

func RedisHGet(key, field string) (string, error){
	value, err := redisCli.HGet(key, field).Result()
	if err == redis.Nil {
		//key does not exist
		return "", nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

func RedisHGetAll(key string)([]string, error){
	values, err := redisCli.HGetAll(key).Result()
	if err == redis.Nil {
		//key does not exist
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return values, nil
}

func RedisHExists(key, field string) (bool, error) {
	return redisCli.HExists(key, field).Result()
}


func RedisHSetMap(key string, values map[string]string) error {
	return RedisHSetMapWithExp(key, 0, values)
}

func RedisHSetMapWithExp(key string, exp time.Duration, values map[string]string) error {
	for field, value := range values {
		err := redisCli.HSet(key, field, value).Err()
		if err != nil {
			return err
		}
	}
	return RedisExpire(key, exp)
}

func RedisHGetMap(key string) (map[string]string, error) {
	return redisCli.HGetAllMap(key).Result()
}

func RedisHDel(key, field string) error {
	return redisCli.HDel(key, field).Err()
}

func RedisHLen(key string) (int64, error) {
	return redisCli.HLen(key).Result()
}

////////////////
// Redis Sets //
////////////////

func RedisSAdd(key string, values ...string) error{
	return RedisSAddWithExp(key, 0, values...)
}

func RedisSAddWithExp(key string, exp time.Duration, values ...string) error{
	err := redisCli.SAdd(key, values...).Err()
	if err != nil{
		return err
	}
	return RedisExpire(key, exp)
}

func RedisSMembers(key string) ([]string, error){
	return redisCli.SMembers(key).Result()
}

///////////////////////
// Redis Sorted Sets //
///////////////////////

func RedisZAdd(key, value string, score float64)error{
	return RedisZAddWithExp(key, 0, value, score)
}

func RedisZAddWithExp(key string, exp time.Duration, value string, score float64) error{
	z := redis.Z{
		Score: score,
		Member: value,
	}
	err := redisCli.ZAdd(key, z).Err()
	if err != nil{
		return err
	}
	return RedisExpire(key, exp)
}

func RedisZAddList(key string, values ...redis.Z) error{
	return RedisZAddListWithExp(key, 0, values...)
}

func RedisZAddListWithExp(key string, exp time.Duration, values ...redis.Z)error{
	if err := redisCli.ZAdd(key, values...).Err(); err != nil{
		return err
	}
	return RedisExpire(key, exp)
}

func RedisZIncr1(key, value string) error{
	return redisCli.ZIncrBy(key, 1.0, value).Err()
}

func RedisZCount(key, min, max string) (int64, error){
	return redisCli.ZCount(key, min, max).Result()
}

func RedisZCountAll(key string)(int64, error){
	return RedisZCount(key, "-inf", "+inf")
}

func RedisZDecsLimit(key string, offset, count int64)([]string, error){
	return RedisZRevRangeByScore(key, "-inf", "+inf", offset, count)
}

func RedisZRevRangeByScore(key, min, max string, offset, count int64)([]string, error){
	opt := redis.ZRangeByScore{
		Min: min,
		Max: max,
	}

	if offset > 0{
		opt.Offset = offset
	}

	if count > 0{
		opt.Count = count
	}
	return redisCli.ZRevRangeByScore(key, opt).Result()
}
