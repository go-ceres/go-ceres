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
	"context"
	"encoding/json"
	"github.com/go-ceres/go-ceres/client/grpc/resolver"
	"google.golang.org/grpc"
	"strings"
	"time"
)

// Client 客户端
type Client struct {
	dialOptions []grpc.DialOption
	config      *Config
}

// newClient 创建一个grpc链接
func newClient(c *Config) *Client {
	var dialOptions []grpc.DialOption

	// 是否使用安全连接
	if !c.Secure {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}
	// 存活策略
	if c.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*c.KeepAlive))
	}
	// 负载均衡策略
	dialOptions = append(dialOptions, grpc.WithBalancerName(c.Balancer))
	// 注册非注册中心的调用
	resolver.RegisterDirect()
	// 注册有注册中心的
	if c.registry != nil {
		resolver.RegisterDiscover(c.registry)
	}
	// 设置连接
	client := &Client{
		config:      c,
		dialOptions: dialOptions,
	}
	return client
}

// Close 关闭客户端
func (c *Client) Close() {
	if c.config.registry != nil {
		// 注销
		resolver.DeregisterDiscover(c.config.registry.String())
	}
}

// Conn 获取连接
func (c *Client) Conn(service string, opts ...Option) (*grpc.ClientConn, error) {
	// 解析服务
	target := parseTarget(service)
	if target.Scheme == "" {
		if c.config.registry != nil {
			target.Scheme = c.config.registry.String()
		} else if len(strings.Split(target.Endpoint, ",")) > 1 {
			target.Scheme = "direct"
		}
	}
	// 解析附加参数
	var ctx = context.Background()
	var dialOptions = c.dialOptions
	// 调用的额外参数
	var options = &Options{
		DialOptions: []grpc.DialOption{},
		Authority:   &resolver.Authority{},
	}
	var bytes []byte
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(options)
		}
		var err error
		bytes, err = json.Marshal(options.Authority)
		if err != nil {
			return nil, err
		}
	}
	target.Authority = string(bytes)
	// 设置调用target
	service = stringTarget(target)
	// 默认配置使用block
	if c.config.Block {
		if c.config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	// 如果有附加调用参数
	if len(options.DialOptions) > 0 {
		dialOptions = append(dialOptions, options.DialOptions...)
	}
	// 如果配置了debug
	if c.config.Debug {
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(c.config.logger, service)))
	}
	// 如果配置了超时阈值
	if c.config.SlowThreshold > time.Duration(0) {
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(c.config.logger, c.config.ReadTimeout, c.config.SlowThreshold)))
	}
	// 获取ClientConn
	conn, err := grpc.DialContext(ctx, service, dialOptions...)

	if err != nil {
		return nil, err
	}
	return conn, nil
}
