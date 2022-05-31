package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RCache struct {
	cache *redis.Client
}

// InitRedis 获取redis对象
func InitRedis(addr, password string, dbNum int) (*RCache, error) {
	var rcache RCache
	rcache.cache = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		PoolSize: 200,
		DB:       dbNum,
	})
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := rcache.cache.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &rcache, nil
}

func (rcache *RCache) IsExist(key string) bool {
	_, err := rcache.cache.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return true
}

func (rcache *RCache) Delete(key string) error {
	_, err := rcache.cache.Del(ctx, key).Result()
	return err
}

// GetSeqence 自增序列
func (rcache *RCache) GetSeqence(key string) (int64, error) {
	seq, err := rcache.cache.Incr(ctx, key).Result()
	return seq, err
}

// GetString
func (rcache *RCache) GetString(key string) (string, error) {
	return rcache.cache.Get(ctx, key).Result()
}

// GetInt
func (rcache *RCache) GetInt(key string) (int, error) {
	return rcache.cache.Get(ctx, key).Int()
}

// GetInt64
func (rcache *RCache) GetInt64(key string) (int64, error) {
	return rcache.cache.Get(ctx, key).Int64()
}

// GetFloat64
func (rcache *RCache) GetFloat64(key string) (float64, error) {
	return rcache.cache.Get(ctx, key).Float64()
}

// SetString 设置KEY值 永不过期
func (rcache *RCache) SetString(key string, val string) error {
	err := rcache.cache.Set(ctx, key, val, 0).Err()
	return err
}

// SetStringExpire 设置KEY值
func (rcache *RCache) SetStringExpire(key string, val string, expiration time.Duration) error {
	err := rcache.cache.Set(ctx, key, val, expiration).Err()
	return err
}

// PushList 存入队列(左进)
func (rcache *RCache) PushList(key string, val string) error {
	err := rcache.cache.LPush(ctx, key, val).Err()
	return err
}

// PopList 移除并获取列表最后一个元素(右出)
func (rcache *RCache) PopList(key string) string {
	return rcache.cache.RPop(ctx, key).Val()
}

// BRPopList 阻塞式的移除并获取列表最后一个元素(右出)
func (rcache *RCache) BRPopList(key string) string {
	rec := rcache.cache.BRPop(ctx, time.Duration(0), key).Val()
	if len(rec) == 2 {
		return rec[1]
	}
	return ""
}

// GetIndexList 通过索引获取列表中的元素
func (rcache *RCache) GetIndexList(key string, dx int64) string {
	return rcache.cache.LIndex(ctx, key, dx).Val()
}

// Hset
func (rcache *RCache) Hset(key string, fields string, val string) error {
	return rcache.cache.HSet(ctx, key, fields, val).Err()
}

// HGet
func (rcache *RCache) HGet(key string, fields string) string {
	val := rcache.cache.HGet(ctx, key, fields).Val()
	return val
}

// HGetAll
func (rcache *RCache) HGetAll(key string) map[string]string {
	val := rcache.cache.HGetAll(ctx, key).Val()
	return val
}

// HDel
func (rcache *RCache) HDel(key string, fields string) error {
	return rcache.cache.HDel(ctx, key, fields).Err()
}

// BRPopLPush 队列转移
func (rcache *RCache) BRPopLPush(source, destination string) string {
	return rcache.cache.BRPopLPush(ctx, source, destination, 0).Val()
}

// RPopLPush 队列转移
func (rcache *RCache) RPopLPush(source, destination string) string {
	return rcache.cache.RPopLPush(ctx, source, destination).Val()
}
