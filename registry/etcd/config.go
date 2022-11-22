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

package etcd

import (
	"github.com/go-ceres/go-ceres/client/etcd"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/registry"
	"time"
)

// Config 配置信息
type Config struct {
	*etcd.Config
	Prefix        string        `json:"prefix"`       // 前缀
	Namespace     string        `json:"namespace"`    // 服务空间
	ReadTimeout   time.Duration `json:"read_timeout"` // 请求超时时间
	ServiceTTL    time.Duration `json:"service_ttl"`  // 服务续约时间间隔
	EtcdConfigKey string        `json:"etcd_key"`     // etcd配置键
	etcdClient    *etcd.Client
	log           *logger.Logger
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Config:      etcd.DefaultConfig(),
		Prefix:      "/ceres/registry/",
		Namespace:   "go-ceres.com",
		ReadTimeout: time.Second * 3,
		ServiceTTL:  time.Second * 15,
		log:         logger.FrameLogger.With(logger.FieldMod(errors.ModRegistryEtcd)),
	}
}

// RawConfig 扫描配置
func RawConfig(key string) *Config {
	c := DefaultConfig()
	err := config.Get(key).Scan(c)
	if err != nil {
		c.log.Panicd("parse config", logger.FieldMod(errors.ModRegistryEtcd), logger.FieldErr(err), logger.FieldAny("key", key), logger.FieldValue(c))
	}
	return c
}

// ScanConfig 扫描配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.registry." + name)
}

// WithLogger 单独设置日志
func (c *Config) WithLogger(log *logger.Logger) *Config {
	c.log = log
	return c
}

// WithClient 单独设置etcd客户端
func (c *Config) WithClient(client *etcd.Client) *Config {
	c.etcdClient = client
	return c
}

// Build 构建etcd注册中心
func (c *Config) Build() registry.Registry {
	if len(c.EtcdConfigKey) > 0 {
		c.Config = etcd.RawConfig(c.EtcdConfigKey)
	}
	return newRegistry(c)
}
