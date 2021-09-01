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

package registry

import (
	"context"
	"time"
)

// DeRegisterOptions 注销服务参数
type DeRegisterOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// DeRegisterOption 注册服务参数
type DeRegisterOption func(opt *DeRegisterOptions)

// DeregisterContext 注销服务的上下文
func DeregisterContext(ctx context.Context) DeRegisterOption {
	return func(opt *DeRegisterOptions) {
		opt.Context = ctx
	}
}

// RegisterOptions 注册服务时的额外参数
type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// RegisterOption 注册服务参数
type RegisterOption func(opt *RegisterOptions)

// RegisterTTl 设置租约时间
func RegisterTTl(t time.Duration) RegisterOption {
	return func(opt *RegisterOptions) {
		opt.TTL = t
	}
}

// RegisterContext 设置注册服务上下文
func RegisterContext(ctx context.Context) RegisterOption {
	return func(opt *RegisterOptions) {
		opt.Context = ctx
	}
}

// WatchOptions 监听服务参数
type WatchOptions struct {
	Service string          // 监听的服务名
	Scheme  string          // 服务协议
	Context context.Context // 上下文
}

// WatchOption 参数
type WatchOption func(opt *WatchOptions)

// WatchService 设置监听的服务
func WatchService(name string) WatchOption {
	return func(opt *WatchOptions) {
		opt.Service = name
	}
}

// WatchScheme 设置监听的服务
func WatchScheme(scheme string) WatchOption {
	return func(opt *WatchOptions) {
		opt.Scheme = scheme
	}
}

// WatchContext 设置监听服务的上下文
func WatchContext(ctx context.Context) WatchOption {
	return func(opt *WatchOptions) {
		opt.Context = ctx
	}
}

// GetOptions 获取服务参数
type GetOptions struct {
	Context context.Context // 上下文
	Scheme  string          // 服务协议
}

type GetOption func(opt *GetOptions)

func GetContext(ctx context.Context) GetOption {
	return func(opt *GetOptions) {
		opt.Context = ctx
	}
}

func GetScheme(scheme string) GetOption {
	return func(opt *GetOptions) {
		opt.Scheme = scheme
	}
}
