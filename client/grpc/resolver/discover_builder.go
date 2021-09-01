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
	"encoding/json"
	"github.com/go-ceres/go-ceres/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"sync"
)

type DiscoverBuilder struct {
	registry registry.Registry // 注册中心
	rw       *sync.RWMutex     // 读写锁
}

// Build 构建方法
func (b DiscoverBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	authority := &Authority{}
	err := json.Unmarshal([]byte(target.Authority), authority)
	if err != nil {
		return nil, err
	}
	service, err := b.registry.WatchService(registry.WatchService(target.Endpoint), registry.WatchScheme("grpc"))
	if err != nil {
		return nil, err
	}
	var stop = make(chan struct{})
	go func() {
		for {
			select {
			case endpoint := <-service:
				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
				}
				for _, service := range endpoint.Services {
					// 如果存在版本号，并且版本号不为空的情况，直接跳过
					if authority.Version != "" && authority.Version != service.Version {
						continue
					}
					for _, node := range service.Nodes {
						// 如果地区不为空，并且当前节点不在所需的地区则跳过
						if authority.Zone != "" && node.Zone != authority.Zone {
							continue
						}
						// 如果区域不为空，并且当前节点不在所需的区域则跳过
						if authority.Region != "" && node.Region != authority.Region {
							continue
						}
						var address resolver.Address
						address.Addr = node.Address
						address.ServerName = target.Endpoint
						address.Attributes = attributes.New("__node_info_", node)
						state.Addresses = append(state.Addresses, address)
					}
				}
				cc.UpdateState(state)
			case <-stop:
				return
			}
		}
	}()
	return &Resolver{
		stop: stop,
	}, nil
}

// Scheme 获取协议
func (b DiscoverBuilder) Scheme() string {
	return b.registry.String()
}
