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
	"github.com/go-ceres/go-ceres/cache"
	"github.com/go-ceres/go-ceres/client/redis"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
)

// Config 配置信息
type Config struct {
	Prefix string `json:"prefix"`
	Type   string `json:"type"`
	*redis.Config
	logger logger.Interface
}

// DefaultConfig 默认配置文件
func DefaultConfig() *Config {
	return &Config{
		Prefix: "cerescache",
		Type:   "redis",
		Config: redis.DefaultConfig(),
		logger: logger.FrameLogger.With(logger.FieldMod(errors.ModCacheRedis)),
	}
}

// RawConfig ...
func RawConfig(key string) *Config {
	c := DefaultConfig()
	err := config.Get(key).Scan(c)
	if err != nil {
		c.logger.DPanicf("parse config", logger.FieldErr(err), logger.FieldAny("key", key), logger.FieldValue(c))
	}
	return c
}

// ScanConfig 扫描配置文件
func ScanConfig(name string) *Config {
	return RawConfig("ceres.cache." + name)
}

// Build 构建缓存组件
func (c *Config) Build() cache.Cache {
	return NewCacheRedis(c)
}
