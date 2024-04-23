package redis

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	redis        *redis.Client
	closeChannel chan struct{}
	subscription map[string]*SubscriptionInfo
}

var (
	instance *Cache
	once     sync.Once
)

// New represents a new redis client
func New(addr string, port int, password string) *Cache {
	if instance != nil {
		return instance
	}
	options, err := redis.ParseURL(fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:%d", addr, port),
			Password: password,
		}
	}

	if options.TLSConfig != nil {
		options.TLSConfig.InsecureSkipVerify = true
	}

	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	if err != nil {
		_ = fmt.Errorf("redis connection error: %s", err)
	}
	once.Do(func() {
		instance = &Cache{
			redis:        client,
			closeChannel: make(chan struct{}),
			subscription: make(map[string]*SubscriptionInfo),
		}
	})
	return instance
}

// Get retrieves value at key from cache
func (c *Cache) Get(key string, value interface{}) (err error) {
	var data []byte
	if err := c.redis.Get(key).Scan(&data); err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// Has checks if key is available in cache
func (c *Cache) Has(key string) (ok bool) {
	item, err := c.redis.Keys(key).Result()
	if err != nil {
		return false
	}
	if len(item) != 0 {
		return true
	}
	return false
}

// Set stores a key with a given life time. 0 for permanent
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) (err error) {
	raw, _ := json.Marshal(value)
	_, err = c.redis.Set(key, raw, ttl).Result()
	if err != nil {
		_ = fmt.Errorf("write error: %s", err)
	}
	return err
}

// Del removes a value from redis
func (c *Cache) Del(key string) (err error) {
	_, err = c.redis.Del(key).Result()
	return err
}

// Keys list all available cache keys
func (c *Cache) Keys(pattern string) (available []string) {
	result, err := c.redis.Keys(pattern).Result()
	if err != nil {
		panic(err)
	}
	return result

}

// Type returns the type of the cache
func (c *Cache) Type() string {
	return "redis"
}

// Clear removes all keys and closes the client
func (c *Cache) Clear() {
	defer func() {
		_ = recover()
	}()
	c.redis.FlushAll()
	close(c.closeChannel)
}
func (c *Cache) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.redis.SetNX(key, value, expiration)
}
