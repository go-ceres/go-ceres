//   Copyright 2021 Go-Ceres
//   Author https://github.com/go-ceres/go-ceres
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package grpc

import (
	"fmt"
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
	"google.golang.org/grpc"
)

type Config struct {
	Debug               bool   // 是否开启调试
	Network             string // net.listen的network类型
	Host                string // 服务ip
	Port                int    // 服务端口
	PlainTextAddress    string // 注册中心显示的地址
	Version             string // 当前项目版本号
	Name                string // 服务名称
	TLS                 bool   // 是否使用tls连接
	CertFile            string // tls的cert文件路径
	KeyFile             string // tls的key文件路径
	ServerSlowThreshold int64  // 服务器素速度阈值
	serverOptions       []grpc.ServerOption
	streamInterceptors  []grpc.StreamServerInterceptor
	unaryInterceptors   []grpc.UnaryServerInterceptor
	logger              *logger.Logger // 日志组件
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Debug:               true,
		Network:             "tcp4",
		Host:                "0.0.0.0",
		Port:                5201,
		CertFile:            "",
		KeyFile:             "",
		ServerSlowThreshold: 500,
		logger:              logger.FrameLogger.With(logger.FieldMod("server.grpc")),
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
	return RawConfig("ceres.server." + name)
}

// WithServerOption 设置grpc服务参数
func (c *Config) WithServerOption(opts ...grpc.ServerOption) *Config {
	if c.serverOptions == nil {
		c.serverOptions = make([]grpc.ServerOption, 0)
	}
	c.serverOptions = append(c.serverOptions, opts...)
	return c
}

// WithStreamInterceptor 设置拦截 stream
func (c *Config) WithStreamInterceptor(opts ...grpc.StreamServerInterceptor) *Config {
	if c.streamInterceptors == nil {
		c.streamInterceptors = make([]grpc.StreamServerInterceptor, 0)
	}

	c.streamInterceptors = append(c.streamInterceptors, opts...)
	return c
}

// WithUnaryInterceptor 设置拦截 unary(一元) RPC
func (c *Config) WithUnaryInterceptor(opts ...grpc.UnaryServerInterceptor) *Config {
	if c.unaryInterceptors == nil {
		c.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0)
	}

	c.unaryInterceptors = append(c.unaryInterceptors, opts...)
	return c
}

// WithLogger 重新设置日志组件
func (c *Config) WithLogger(log *logger.Logger) *Config {
	c.logger = log
	return c
}

// Address 获取服务地址
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Build 构建grpc服务
func (c *Config) Build() *grpcServer {
	return newServer(c)
}
