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
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-redis/redis"
	"time"
)

type Mode string

const (
	ClusterMode Mode = "cluster"
	SimpleMode  Mode = "simple"
)

type Config struct {
	Addrs        []string       `json:"addrs"`          // 连接地址
	Mode         Mode           `json:"mode"`           // 模式（cluster,simple）
	Password     string         `json:"password"`       // 密码
	DB           int            `json:"db"`             // DB，默认为0, 一般应用不推荐使用DB分片
	PoolSize     int            `json:"pool_size"`      // 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	MaxRetries   int            `json:"maxRetries"`     //网络相关的错误最大重试次数 默认5次
	MinIdleConns int            `json:"min_idle_conns"` // 最小空闲连接数,默认100
	DialTimeout  time.Duration  `json:"dial_timeout"`   // 连接超时
	ReadTimeout  time.Duration  `json:"read_timeout"`   //读取超时 默认3s
	WriteTimeout time.Duration  `json:"write_timeout"`  // 写入超时 默认3s
	IdleTimeout  time.Duration  `json:"idle_timeout"`   // 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	Debug        bool           `json:"debug"`          // 是否开启debug模式
	ReadOnly     bool           `json:"readOnly"`       // 集群模式中在从属节点上启用读模式
	logger       *logger.Logger // 日志组件
}

func DefaultConfig() *Config {
	return &Config{
		Addrs:        []string{"127.0.0.1:6379"},
		DB:           0,
		PoolSize:     10,
		MaxRetries:   5,
		MinIdleConns: 100,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 60,
		Debug:        false,
		ReadOnly:     false,
		logger:       logger.FrameLogger.With(logger.FieldMod(errors.ModClientRedis)),
	}
}

// RawConfig 根据key扫描配置
func RawConfig(key string) *Config {
	c := DefaultConfig()
	err := config.Get(key).Scan(c)
	if err != nil {
		c.logger.Panicd("parse config", logger.FieldErr(err), logger.FieldAny("key", key), logger.FieldValue(c))
	}
	return c
}

// ScanConfig 根据名称扫描配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.client.redis." + name)
}

// buildClusterClient 构建redis集群客户端
func (c *Config) buildClusterClient() *redis.ClusterClient {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.Addrs,
		MaxRedirects: c.MaxRetries,
		ReadOnly:     c.ReadOnly,
		Password:     c.Password,
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		IdleTimeout:  c.IdleTimeout,
	})
	if err := clusterClient.Ping().Err(); err != nil {
		c.logger.Panic("start cluster redis", logger.FieldErr(err))
	}
	return clusterClient
}

// buildSimpleClient 构建单节点redis客户端
func (c *Config) buildSimpleClient() *redis.Client {
	simpleClient := redis.NewClient(&redis.Options{
		Addr:         c.Addrs[0],
		Password:     c.Password,
		DB:           c.DB,
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		IdleTimeout:  c.IdleTimeout,
	})
	if err := simpleClient.Ping().Err(); err != nil {
		c.logger.Panic("dial redis fail", logger.FieldErr(err), logger.FieldAny("config", c))
	}
	return simpleClient
}

// Build 构建redis
func (c *Config) Build() *Redis {
	// 如果没有配置模式
	if len(c.Mode) == 0 {
		c.Mode = SimpleMode
		if len(c.Addrs) > 1 {
			c.Mode = ClusterMode
		}
	}
	var client redis.UniversalClient
	switch c.Mode {
	case ClusterMode:
		if len(c.Addrs) == 1 {
			c.logger.Warn("redis config has only 1 address but with cluster mode")
		}
		client = c.buildClusterClient()
	case SimpleMode:
		if len(c.Addrs) > 1 {
			c.logger.Warn("redis config has more than 1 address but with simple mode")
		}
		client = c.buildSimpleClient()
	default:
		c.logger.Panic("redis mode must be one of (simple, cluster)")
	}
	return &Redis{
		Config: c,
		client: client,
	}
}
