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

package token

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
)

const (
	TokenStyleUuid        = "uuid"        // uuid样式
	TokenStyleSimpleUuid  = "simple-uuid" // uuid不带下划线
	TOKEN_STYLE_RANDOM_32 = "random-32"   // 随机32字符串
)

type Config struct {
	TokenName       string         `json:"token_name" yaml:"token_name" toml:"token_name"`                   // token的名称，同时提交的token名称也是如此
	Timeout         int64          `json:"timeout" yaml:"timeout" toml:"timeout"`                            // usertoken的过期时间
	ActivityTimeout int64          `json:"activity_timeout" yaml:"activity_timeout" toml:"activity_timeout"` // 临时过期时间（适用场景为，多少时间后就不允许操作）
	IsConcurrent    bool           `json:"is_concurrent" yaml:"is_concurrent" toml:"is_concurrent"`          // 是否支持同时在线（为false时会挤掉旧登录）
	IsShare         bool           `json:"is_share" yaml:"is_share" toml:"is_share"`                         // 是否共享token，
	TokenStyle      string         `json:"token_style" yaml:"token_style" toml:"token_style"`                // token样式
	AutoRenew       bool           `json:"auto_renew" yaml:"auto_renew" toml:"auto_renew"`                   // 自动续签
	TokenPrefix     string         `json:"token_prefix" yaml:"token_prefix" toml:"token_prefix"`             // token前缀
	IsLog           bool           `json:"is_log" yaml:"is_log" toml:"is_log"`                               // 是否打印日志
	CheckLogin      bool           `json:"check_login" yaml:"check_login" toml:"check_login"`                // 获取tokensession时是否检查登录，
	logger          *logger.Logger // 日志组件
}

func DefaultConfig() *Config {
	return &Config{
		TokenName:       "cerestoken",
		Timeout:         2592000,
		ActivityTimeout: -1,
		IsConcurrent:    true,
		IsShare:         true,
		TokenStyle:      TokenStyleUuid,
		AutoRenew:       true,
		TokenPrefix:     "",
		IsLog:           true,
		CheckLogin:      true,
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
func (c *Config) Build(loginType string) *Token {
	return NewToken(c, loginType)
}
