/*
 * //  Copyright 2021 Go-Ceres
 * //  Author https://github.com/go-ceres/go-ceres
 * //
 * //  Licensed under the Apache License, Version 2.0 (the "License");
 * //  you may not use this file except in compliance with the License.
 * //  You may obtain a copy of the License at
 * //
 * //      http://www.apache.org/licenses/LICENSE-2.0
 * //
 * //  Unless required by applicable law or agreed to in writing, software
 * //  distributed under the License is distributed on an "AS IS" BASIS,
 * //  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * //  See the License for the specific language governing permissions and
 * //  limitations under the License.
 */

package server

import (
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/google/uuid"
)

// Server 服务接口
type Server interface {
	// Stop t停止
	Stop() error
	// Start 启动
	Start() error
	// Info 服务信息
	Info() *ServiceInfo
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Id       string            `json:"id"`       // 应用ID
	Name     string            `json:"name"`     // 应用名称
	Address  string            `json:"address"`  // 服务地址
	Version  string            `json:"version"`  // 当前版本
	Scheme   string            `json:"scheme"`   // 服务协议
	Weight   int               `json:"weight"`   // 权重
	Region   string            `json:"region"`   // 区域
	Zone     string            `json:"zone"`     // 地区
	Metadata map[string]string `json:"metadata"` // 元数据
}

// defaultInfo 默认服务信息
func defaultInfo() *ServiceInfo {
	info := &ServiceInfo{
		Id:       uuid.New().String(),
		Name:     cmd.GetAppName(),
		Scheme:   "",
		Weight:   100,
		Region:   cmd.GetRegion(),
		Zone:     cmd.GetZone(),
		Version:  cmd.GetAppVersion(),
		Metadata: make(map[string]string),
	}
	info.Metadata["hostname"] = cmd.GetHostname()
	info.Metadata["start_time"] = cmd.GetStartTime()
	info.Metadata["build_time"] = cmd.GetBuildTime()
	info.Metadata["build_user"] = cmd.GetBuildUser()
	info.Metadata["build_host"] = cmd.GetBuildHost()
	info.Metadata["ceres_version"] = cmd.FrameVersion()
	return info
}

// Option 服务信息参数
type Option func(info *ServiceInfo)

// ApplyOptions 添加服务信息
func ApplyOptions(opts ...Option) *ServiceInfo {
	info := defaultInfo()
	for _, opt := range opts {
		opt(info)
	}
	return info
}

// WithAddress 设置服务地址
func WithAddress(addr string) Option {
	return func(info *ServiceInfo) {
		info.Address = addr
	}
}

// WithScheme 设置服务协议
func WithScheme(scheme string) Option {
	return func(info *ServiceInfo) {
		info.Scheme = scheme
	}
}

// WithMetadata 添加元信息
func WithMetadata(key, value string) Option {
	return func(info *ServiceInfo) {
		info.Metadata[key] = value
	}
}
