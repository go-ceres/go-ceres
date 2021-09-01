// Copyright 2020 Go-Ceres
// Author https://github.com/go-ceres/go-ceres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type Config struct {
	rw            *sync.RWMutex // 读写锁
	Debug         bool          `json:"debug"`       // 是否开启debug模式，默认false
	Stdout        bool          `json:"stdout"`      // 终端输出日志
	Level         string        `json:"level"`       // 日志等级
	Fields        []zap.Field   `json:"fields"`      // 初始化字段
	AddCaller     bool          `json:"add_caller"`  // 是否打印调用者信息，默认，true
	TimeFormat    string        `json:"time_format"` // 时间格式化
	CallerSkip    int           `json:"caller_skip"` // 表示输出当前栈帧，默认，1
	autoLevelKey  string        // 日志等级监听key
	Core          zapcore.Core
	EncoderConfig *zapcore.EncoderConfig   `json:"encoder_config"` // 日志编码设置
	writer        map[string]writer.Writer // 日志输出者
	Writers       map[string]interface{}   `json:"writers"` // 配置信息
}

// 获取一个默认的配置
func defaultConfig() *Config {
	return &Config{
		Level:         "info",
		Stdout:        true,
		AddCaller:     true,
		CallerSkip:    1,
		writer:        make(map[string]writer.Writer),
		Writers:       make(map[string]interface{}),
		EncoderConfig: defaultEncoderConfig(),
	}
}

// RawConfig 根据key构建配置
func RawConfig(key string) *Config {
	conf := defaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		panic(err)
	}
	return conf
}

// ScanConfig 根据name构建配置
func ScanConfig(name string) *Config {
	conf := RawConfig("ceres.logger." + name)
	if conf.TimeFormat != "" {
		conf.EncoderConfig.EncodeTime = timeEncoderStr(conf.TimeFormat)
	}
	return conf
}

// initialize 初始化
func (c *Config) initialize() {
	if c.rw == nil {
		c.rw = &sync.RWMutex{}
	}
	if c.EncoderConfig == nil {
		c.EncoderConfig = defaultEncoderConfig()
	}
	if c.Debug {
		c.EncoderConfig.EncodeLevel = debugEncodeLevel
	}
	if c.writer == nil {
		c.writer = make(map[string]writer.Writer)
	}
	// 如果有配置信息
	if len(c.Writers) > 0 {
		for name, conf := range c.Writers {
			// 如果有该配置文件的构造器，则初始化
			if build, ok := writer.Load(name); ok {
				c.AddWriter(build.Build(conf))
			}
		}
	}

	if c.rw == nil {
		c.rw = &sync.RWMutex{}
	}
}

// AddWriter 添加日志输出者
func (c *Config) AddWriter(writer writer.Writer) *Config {
	c.rw.Lock()
	c.writer[writer.Name()] = writer
	c.rw.Unlock()
	return c
}

// Build 创建logger
func (c Config) Build() *Logger {
	c.initialize()
	logger := newLogger(&c)
	if c.autoLevelKey != "" {
		logger.AutoLevel(c.autoLevelKey)
	}
	return logger
}
