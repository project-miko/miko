package models

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

// completed commands
// string set / get
// hash hset / hget / hmset / hmget / hgetall
// list lpush / rpop / brpop

type RedisClient struct {
	client *redis.Pool
}

func (rc *RedisClient) Get() redis.Conn {
	return rc.client.Get()
}

// block pop an element from the right side of the queue (supports monitoring multiple queues)
// note: the timeout setting needs to reference the redis server keep alive setting, if it exceeds this value, the timeout error will be triggered due to the server actively disconnecting
// self-judge redis.ErrNil redis empty value error
func (rc *RedisClient) Brpop(keys []string, blockFor int64) (k, result string, e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	r, e := redis.Strings(client.Do("BRPOP", redis.Args{}.AddFlat(keys).Add(blockFor)...))
	if e != nil {
		return "", "", e
	}

	return r[0], r[1], nil
}

// pop an element from the right side of the queue
func (rc *RedisClient) Rpop(key string) (result string, e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	result, e = redis.String(client.Do("RPOP", key))
	return result, e
}

// push a batch of elements to the left side of the queue
func (rc *RedisClient) LBatchPush(key string, value []string) (e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	_, e = client.Do("LPUSH", redis.Args{}.Add(key).AddFlat(value)...)
	return e
}

// push an element to the left side of the queue
func (rc *RedisClient) Lpush(key, value string) (e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	_, e = client.Do("LPUSH", key, value)
	return e
}

// get all values of a hash structure
func (rc *RedisClient) Hgetall(key string) (result map[string]string, e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	result = make(map[string]string)

	r, e := redis.Strings(client.Do("HGETALL", key))

	for i, c := 0, len(r); i < c; i += 2 {
		result[r[i]] = r[i+1]
	}

	return result, e
}

// get multiple fields of a hash structure
func (rc *RedisClient) Hmget(key string, fields []string) (result map[string]string, e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	result = make(map[string]string)

	r, e := redis.Strings(client.Do("HMGET", redis.Args{}.Add(key).AddFlat(fields)...))

	for i, v := range r {
		result[fields[i]] = v
	}

	return result, e
}

// set multiple key-value pairs for a hash structure
func (rc *RedisClient) Hmset(key string, kvs map[string]string) (e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	_, e = client.Do("HMSET", redis.Args{}.Add(key).AddFlat(kvs)...)
	return e
}

// get a single field value of a hash structure
func (rc *RedisClient) Hget(key, field string) (result string, e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	result, e = redis.String(client.Do("HGET", key, field))
	return result, e
}

// set a single key-value pair for a hash structure
func (rc *RedisClient) Hset(key, field, value string) (e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	_, e = client.Do("HSET", key, field, value)
	return e
}

// get command, get a string from redis
func (rc *RedisClient) GetString(key string) (str string, e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	_d, e := redis.String(client.Do("GET", key))
	if e != nil {
		return "", e
	}

	return _d, nil
}

// mget command, get multiple strings from redis, will filter out non-existent keys
func (rc *RedisClient) MGet(keys []string) (map[string]string, error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	result := make(map[string]string)

	r, e := redis.Strings(client.Do("MGET", redis.Args{}.AddFlat(keys)...))
	if e != nil {
		return nil, e
	}

	for i, v := range r {
		if v != "" {
			result[keys[i]] = v
		}
	}

	return result, nil
}

// same as set command, set a string to redis
func (rc *RedisClient) SetString(key, value string, expire int64) (e error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	if expire > 0 {
		_, e = client.Do("SETEX", key, expire, value)
	} else {
		_, e = client.Do("SET", key, value)
	}

	return e
}

// get a struct from redis, and parse it to the corresponding struct pointer
// data is a pointer to a new struct, the data will be unmarshalled into this memory
func (rc *RedisClient) GetStruct(key string, data interface{}) error {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	_d, e := redis.Bytes(client.Do("GET", key))
	if e != nil {
		return e
	}

	e = json.Unmarshal(_d, data)
	if e != nil {
		return e
	}

	return nil
}

// save a struct to redis
// data : struct pointer
// automatically marshal the struct and save it to redis
func (rc *RedisClient) SetStruct(key string, data interface{}, expire int64) error {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	d, e := json.Marshal(data)
	if e != nil {
		return e
	}

	if expire > 0 {
		_, e = client.Do("SETEX", key, expire, d)
	} else {
		_, e = client.Do("SET", key, d)
	}

	if e != nil {
		return e
	}

	return nil
}

func (rc *RedisClient) IncrWithExpire(key string, expire int64) (int, error) {
	client := rc.Get()
	defer func() {
		_ = client.Close()
	}()

	exists, err := redis.Bool(client.Do("EXISTS", key))
	if err != nil {
		return 0, err
	}

	if !exists { // if the key does not exist, set the initial value to 1 and set the expiration time
		err = rc.SetString(key, "1", expire)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

	// if the key exists, only execute the INCR operation
	return redis.Int(client.Do("INCR", key))
}
