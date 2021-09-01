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

package registry

import "github.com/go-ceres/go-ceres/errors"

var (
	ErrNotFound      = errors.New(errors.CodeGetServiceErrorNotFound, errors.MsgGetServiceErrorNotFound)
	ErrNoServiceName = errors.New(errors.CodeWatchServiceErrorNoServiceName, errors.MsgWatchServiceErrorNoServiceName)
)

// Registry 注册中心接口
type Registry interface {
	// Register 注册服务
	Register(srv *Service, opts ...RegisterOption) error
	// Deregister 注销服务
	Deregister(srv *Service, opts ...DeRegisterOption) error
	// GetService 获取服务
	GetService(name string, opts ...GetOption) ([]*Service, error)
	// WatchService 监听服务
	WatchService(opts ...WatchOption) (chan WatchServiceResult, error)
	// ListService 服务列表
	ListService() ([]*Service, error)
	// String 获取注册中心类型
	String() string
	// Close 关闭注册中心
	Close() error
}

// WatchServiceResult 监听服务返回
type WatchServiceResult struct {
	Services []Service // 服务列表
}

func (w *WatchServiceResult) DeepCopy() *WatchServiceResult {
	if w == nil {
		return nil
	}

	res := &WatchServiceResult{Services: []Service{}}
	res.Services = w.Services
	return res
}

// Node 节点信息
type Node struct {
	Id       string            `json:"id"`       // 节点id
	Address  string            `json:"address"`  // 节点地址
	Weight   int               `json:"weight"`   // 节点权重
	Region   string            `json:"region"`   // 区域
	Zone     string            `json:"zone"`     // 地区
	Metadata map[string]string `json:"metadata"` // 节点元信息
}

// Service 服务
type Service struct {
	Name     string            `json:"name"`     // 服务名称
	Version  string            `json:"version"`  // 版本
	Scheme   string            `json:"scheme"`   // 服务协议
	Nodes    []*Node           `json:"nodes"`    // 所有节点信息
	Metadata map[string]string `json:"metadata"` // 元信息
}
