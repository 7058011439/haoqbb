package DataBase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/go-redis/redis"
	"reflect"
)

type RedisDB struct {
	*redis.Client
}

func NewRedisDB(ip string, port int, passWord string, dbIndex int) *RedisDB {
	ret := &RedisDB{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", ip, port),
			Password: passWord,
			DB:       dbIndex,
		}),
	}
	if ret.Client == nil {
		Log.ErrorLog("Failed to NewRedisDB, client is nil")
		return nil
	}
	if _, err := ret.Client.Ping().Result(); err != nil {
		Log.ErrorLog("Failed to redis ping, err = %v", err)
		return nil
	}
	return ret
}

func (r *RedisDB) HGetString(key string, field string) string {
	if ret, err := r.Client.HGet(key, field).Result(); err == nil {
		return ret
	}
	return ""
}

func (r *RedisDB) HGetValue(key string, field string, data interface{}) error {
	if d := r.HGetString(key, field); d != "" {
		return json.Unmarshal([]byte(d), data)
	} else {
		return errors.New("未能获取到数据")
	}
}

func (r *RedisDB) HGetAll(key string) map[string]string {
	if ret, err := r.Client.HGetAll(key).Result(); err == nil {
		return ret
	}
	return nil
}

func (r *RedisDB) HSetValue(key string, field string, value interface{}) error {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Ptr:
		value, _ = json.Marshal(value)
	}
	return r.Client.HSet(key, field, value).Err()
}

func (r *RedisDB) HIncBy(key string, field string, value int64) int64 {
	if ret, err := r.Client.HIncrBy(key, field, value).Result(); err == nil {
		return ret
	}
	return 0
}

func (r *RedisDB) IncBy(key string, value int64) int64 {
	if ret, err := r.Client.IncrBy(key, value).Result(); err == nil {
		return ret
	}
	return 0
}

func (r *RedisDB) IsKeyExist(key string) bool {
	if ret, err := r.Client.Exists(key).Result(); err != nil {
		return false
	} else {
		return ret > 0
	}
}

func (r *RedisDB) IsFieldExist(key string, field string) bool {
	return r.Client.HExists(key, field).Val()
}
