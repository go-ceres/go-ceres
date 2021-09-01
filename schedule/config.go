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

package schedule

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/robfig/cron/v3"
)

type Config struct {
	Size int         `json:"size"` // 任务总数量
	log  cron.Logger // 日志组件
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Size: 100,
		log: &Logger{
			Log: logger.FrameLogger,
		},
	}
}

// RawConfig 完整key取配置
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	err := config.Get(key).Scan(conf)
	if err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 名称取配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.cron." + name)
}

// WithLogger 设置日志组件
func (c *Config) WithLogger(log cron.Logger) *Config {
	c.log = log
	return c
}

// Build 构建任务调度
func (c *Config) Build() *Schedule {
	s := newSchedule(c)
	return s
}
