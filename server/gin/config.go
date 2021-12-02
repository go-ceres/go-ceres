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

package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
)

type Config struct {
	Host                string // 服务ip
	Port                int    // 服务端口
	Mode                string // 运行模式
	PlainTextAddress    string // 注册中心显示地址
	Version             string // 当前项目版本号
	Name                string // 服务名称
	ServerSlowThreshold int64  // 服务器超时阈值
	logger              *logger.Logger
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:                "0.0.0.0",
		Port:                5202,
		ServerSlowThreshold: 500,
		Mode:                gin.ReleaseMode,
		logger:              logger.FrameLogger.With(logger.FieldMod("server.gin")),
		Name:                cmd.DefaultCmd.App().Name,
		Version:             cmd.DefaultCmd.App().Version,
	}
}

// RawConfig 读取配置
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	if err := config.Get(key).Scan(conf); err != nil {
		conf.logger.Panicd(
			"grpc server parse config panic",
			logger.FieldErr(err),
			logger.FieldValue(conf),
		)
	}
	return conf
}

// ScanConfig 从config组件读取配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.server.gin." + name)
}

// WithLogger 重新设置日志
func (c *Config) WithLogger(log *logger.Logger) *Config {
	c.logger = log
	return c
}

// WithHost 设置主机名
func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

// WithPort 设置端口
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// Address 获取服务地址
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Build 构建gin的服务
func (c *Config) Build() *Server {
	// 新建服务
	server := newGinServer(c)
	// 日志中间件
	server.Use(loggerMiddleware(c.logger, c.ServerSlowThreshold))

	return server
}
