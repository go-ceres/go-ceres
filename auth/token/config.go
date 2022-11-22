//    Copyright 2022. Go-Ceres
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

package token

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
)

// Config 配置信息
type Config struct {
	TokenName       string     `json:"token_name" yaml:"TokenName" toml:"token_name"`                   // token名称，提交的时候按照此名称获取token数据
	Timeout         int64      `json:"timeout" yaml:"Timeout" toml:"timeout"`                           // userToken 过期时间
	ActivityTimeout int64      `json:"activity_timeout" yaml:"ActivityTimeout" toml:"activity_timeout"` // 临时过期时间，用于（超过多少时间不能再操作）
	IsConcurrent    bool       `json:"is_concurrent" yaml:"IsConcurrent" toml:"is_concurrent"`          // 是否支持多账号登录
	IsShare         bool       `json:"is_share" yaml:"IsShare" toml:"is_share"`                         // 是否共享token
	TokenStyle      TokenStyle `json:"token_style" yaml:"TokenStyle" toml:"token_style"`                // token的样式
	AutoRenew       bool       `json:"auto_renew" yaml:"AutoRenew"  toml:"auto_renew"`                  // 自动续签
	TokenPrefix     string     `json:"token_prefix" yaml:"TokenPrefix" toml:"token_prefix"`             // token前缀
	IsLog           bool       `json:"is_log" yaml:"IsLog" toml:"is_log"`                               // 是否打印日志
	CheckLogin      bool       `json:"check_login" yaml:"CheckLogin" toml:"check_login"`                // 检查是否登录
	logger          Logger     // 日志组件
}

func DefaultConfig() *Config {
	return &Config{
		TokenName:       "ceres-token",
		Timeout:         2592000,
		ActivityTimeout: -1,
		IsConcurrent:    true,
		IsShare:         true,
		TokenStyle:      TOKEN_STYLE_UUID,
		AutoRenew:       true,
		TokenPrefix:     "Bearer",
		IsLog:           true,
		CheckLogin:      true,
		logger:          logger.FrameLogger.With(logger.FieldMod(errors.ModAuthToken)),
	}
}

// RawConfig 根据配置键扫描配置
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	err := config.Get(key).Scan(conf)
	if err != nil {
		conf.logger.Panicd("parse config", logger.FieldErr(err), logger.FieldAny("key", key), logger.FieldValue(conf))
	}
	return conf
}

// ScanConfig 根据配置名扫描配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.auth.token." + name)
}

// WithLogger 设置日志组件
func (c *Config) WithLogger(log *logger.Logger) *Config {
	c.logger = log
	return c
}

// Build 构建token
func (c *Config) Build(loginType string) *Logic {
	return NewLogic(loginType, c)
}
