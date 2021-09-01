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

package resolver

import (
	"github.com/go-ceres/go-ceres/registry"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
)

type DirectBuilder struct {
	rw    *sync.RWMutex                // 读写锁
	cache map[string]*registry.Service // 本地缓存服务
}

// Build 构建方法
func (b DirectBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var addrs []resolver.Address
	hosts := strings.FieldsFunc(target.Endpoint, func(r rune) bool {
		return r == EndpointSepChar
	})
	for _, val := range subset(hosts, subsetSize) {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	cc.UpdateState(resolver.State{
		Addresses: addrs,
	})

	return &Resolver{
		stop: nil,
	}, nil
}

// Scheme 获取协议
func (b DirectBuilder) Scheme() string {
	return "direct"
}
