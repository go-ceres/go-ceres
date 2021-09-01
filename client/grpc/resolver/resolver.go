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
	"sync"
)

const (
	EndpointSepChar int32 = ','
	subsetSize            = 32
)

type Resolver struct {
	stop chan struct{} // 停止监控通道
}

func (r Resolver) ResolveNow(_ resolver.ResolveNowOptions) {}

// Close 关闭
func (r Resolver) Close() {
	r.stop <- struct{}{}
}

// RegisterDiscover 注册带注册中心的
func RegisterDiscover(reg registry.Registry) {
	build := &DiscoverBuilder{
		registry: reg,
		rw:       &sync.RWMutex{},
	}
	// 注册builder
	resolver.Register(build)
}

// DeregisterDiscover 注销
func DeregisterDiscover(scheme string) {
	resolver.UnregisterForTesting(scheme)
}

// RegisterDirect 注册直连
func RegisterDirect() {
	if resolver.Get("direct") == nil {
		resolver.Register(
			&DirectBuilder{
				rw:    &sync.RWMutex{},
				cache: map[string]*registry.Service{},
			},
		)
	}
}
