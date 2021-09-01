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

package grpc

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/registry"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Config struct {
	Debug         bool                        // 是否开启调试模式
	Block         bool                        // 等待连接启动再返回
	ReadTimeout   time.Duration               // 调用超时时间
	DialTimeout   time.Duration               // 调用超时时间
	SlowThreshold time.Duration               // 超时阈值
	Balancer      string                      // 负载均衡策略
	Secure        bool                        // 安全链接
	KeepAlive     *keepalive.ClientParameters // 存活策略
	registry      registry.Registry           // 注册中心
	logger        *logger.Logger              // 日志组件
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Block:         true,
		ReadTimeout:   time.Second,
		DialTimeout:   time.Second * 3,
		SlowThreshold: time.Duration(0),
		Balancer:      roundrobin.Name,
		Debug:         true,
		Secure:        false,
		logger:        logger.FrameLogger.With(logger.FieldMod(errors.ModClientGrpc)),
	}
}

// RawConfig 读取配置
func RawConfig(key string) *Config {
	c := DefaultConfig()
	err := config.Get(key).Scan(c)
	if err != nil {
		c.logger.Panicd("parse config", logger.FieldErr(err), logger.FieldAny("key", key), logger.FieldValue(c))
	}
	return c
}

// ScanConfig 根据name扫描配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.client.grpc." + name)
}

// WithLogger 设置日志组件
func (c *Config) WithLogger(log *logger.Logger) *Config {
	c.logger = log
	return c
}

// WithRegistry 设置注册中心
func (c *Config) WithRegistry(registry registry.Registry) *Config {
	c.registry = registry
	return c
}

// Build 构建grpc客户端
func (c *Config) Build() *Client {
	return newClient(c)
}
