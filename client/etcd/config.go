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
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"google.golang.org/grpc"
	"time"
)

// Config 配置信息
type Config struct {
	Endpoints            []string      `json:"endpoints"`
	AutoSyncInterval     time.Duration `json:"auto_sync_interval"`
	DialTimeout          time.Duration `json:"dial_timeout"`
	DialKeepAliveTime    time.Duration `json:"dial_keep_alive_time"`
	DialKeepAliveTimeout time.Duration `json:"dial_keep_alive_timeout"`
	CertFile             string        `json:"cert_file"`
	KeyFile              string        `json:"key_file"`
	CaCert               string        `json:"ca_cert"`
	Username             string        `json:"username"`
	Password             string        `json:"password"`
	RejectOldCluster     bool          `json:"reject_old_cluster"`
	PermitWithoutStream  bool          `json:"permit_without_stream"`
	Secure               bool          `json:"secure"`
	logger               *logger.Logger
	DialOptions          []grpc.DialOption
}

// DefaultConfig 默认的配置
func DefaultConfig() *Config {
	return &Config{
		Endpoints:            []string{"127.0.0.1:2379"},
		DialTimeout:          time.Second * 5,
		DialKeepAliveTime:    time.Second * 10,
		DialKeepAliveTimeout: time.Second * 3,
		logger:               logger.FrameLogger.With(logger.FieldMod(errors.ModClientEtcd)),
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
		},
	}
}

// RawConfig 扫描配置
func RawConfig(key string) *Config {
	conf := DefaultConfig()
	err := config.Get(key).Scan(conf)
	if err != nil {
		conf.logger.Panicd("etcd client parse config panic", logger.FieldErr(err), logger.FieldValue(conf))
	}
	return conf
}

// ScanConfig 扫描配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.etcd." + name)
}

// WithLogger 单独设置日志组件
func (c *Config) WithLogger(log *logger.Logger) *Config {
	c.logger = log
	return c
}

// WithDialOption 配置额外的连接信息
func (c *Config) WithDialOption(opts ...grpc.DialOption) *Config {
	c.DialOptions = append(c.DialOptions, opts...)
	return c
}

// Build 构建etcd客户端
func (c *Config) Build() *Client {
	return newClient(c)
}
