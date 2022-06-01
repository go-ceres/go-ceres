//    Copyright 2021. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package redis

import (
	"encoding/json"
	"github.com/go-ceres/go-ceres/client/redis"
	"time"
)

type CacheRedis struct {
	Config *Config
	client *redis.Redis
}

// NewCacheRedis 根据配置创建redis缓存
func NewCacheRedis(c *Config) *CacheRedis {
	redisClient := c.Config.Build()
	return &CacheRedis{
		Config: c,
		client: redisClient,
	}
}

// getSaveKey 获取实际存储key
func (c *CacheRedis) getSaveKey(key string) string {
	if len(c.Config.Prefix) > 0 {
		return c.Config.Prefix + ":" + key
	}
	return key
}

// Has 查询是否包含缓存
func (c *CacheRedis) Has(key string) bool {
	key = c.getSaveKey(key)
	for _, s := range c.client.Keys(key) {
		if key == s {
			return true
		}
	}
	return false
}

// Get 获取缓存，可以带默认值
func (c *CacheRedis) Get(key string, def ...string) string {
	key = c.getSaveKey(key)
	ret := ""
	if len(def) > 0 {
		ret = def[0]
	}
	if r := c.client.Get(key); r != "" {
		ret = r
	}
	return ret
}

// Set 设置缓存
func (c *CacheRedis) Set(key string, value string, timeout int64) bool {
	key = c.getSaveKey(key)
	return c.client.Set(key, value, time.Second*time.Duration(timeout))
}

// SetObject 设置对象
func (c *CacheRedis) SetObject(key string, value interface{}, timeout int64) bool {
	key = c.getSaveKey(key)
	marshal, err := json.Marshal(value)
	if err != nil {
		return false
	}
	return c.client.Set(key, marshal, time.Second*time.Duration(timeout))
}

// GetObject 获取obj
func (c *CacheRedis) GetObject(key string, obj interface{}) bool {
	key = c.getSaveKey(key)
	str := c.client.Get(key)
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return false
	}
	return true
}

// Update 修改数据,并且不修改过期时间
func (c *CacheRedis) Update(key string, value string) {
	expire, err := c.client.TTL(c.getSaveKey(key))
	if err != nil {
		return
	}
	c.Set(key, value, expire)
}

// UpdateObject 修改持久化数据
func (c *CacheRedis) UpdateObject(key string, value interface{}) bool {
	expire, err := c.client.TTL(c.getSaveKey(key))
	if err != nil {
		return false
	}
	return c.SetObject(key, value, expire)
}

// UpdateObjectTTl 修改持久化时间
func (c *CacheRedis) UpdateObjectTTl(key string, timeout int64) {
	key = c.getSaveKey(key)
	_, _ = c.client.Expire(key, time.Duration(timeout)*time.Second)
}

// Del 删除缓存
func (c *CacheRedis) Del(key string) bool {
	key = c.getSaveKey(key)
	return c.client.Del(key) == 0
}

// TTl 获取剩余过期时间
func (c *CacheRedis) TTl(key string) int64 {
	key = c.getSaveKey(key)
	ttl, _ := c.client.TTL(key)
	return ttl
}

// Clear 清除缓存
func (c *CacheRedis) Clear() bool {
	keys := c.client.Keys(c.Config.Prefix + "*")
	return c.client.Del(keys...) == 0
}
