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
	"github.com/go-ceres/go-ceres/client/grpc/resolver"
	"google.golang.org/grpc"
)

type Options struct {
	// 额外的调用参数
	DialOptions []grpc.DialOption
	// 附加信息
	Authority *resolver.Authority
}

type Option func(opt *Options)

// Version 调用服务的版本
func Version(v string) Option {
	return func(opt *Options) {
		opt.Authority.Version = v
	}
}

// Region 指定地域
func Region(region string) Option {
	return func(opt *Options) {
		opt.Authority.Region = region
	}
}

// Zone 指定地区
func Zone(zone string) Option {
	return func(opt *Options) {
		opt.Authority.Zone = zone
	}
}

// DialOption 设置调用参数
func DialOption(dialOption ...grpc.DialOption) Option {
	return func(opt *Options) {
		opt.DialOptions = append(opt.DialOptions, dialOption...)
	}
}
