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

package cmd

import (
	"fmt"
	"github.com/go-ceres/cli/v2"
	"github.com/go-ceres/go-ceres/errors"
	"sync"
)

var (
	DefaultPluginManager = &PluginManager{
		smu:     sync.RWMutex{},
		plugins: make(map[string]Plugin),
	}
)

// PluginManager 插件管理器定义
type PluginManager struct {
	smu     sync.RWMutex
	plugins map[string]Plugin
}

// Plugin 插件接口定义
type Plugin interface {
	// Name 获取插件名称
	Name() string
	// Flags 返回该插件需要注入的flags
	Flags() []cli.Flag
	// Init 初始化插件时会调用该方法
	Init(ctx *cli.Context) error
	// Config 配置组件加载完毕后调用该方法
	Config() error
	// Destroy 应用退出时调用，用于释放资源
	Destroy()
}

// Register 注册插件
func (m *PluginManager) Register(p Plugin) error {
	m.smu.Lock()
	defer m.smu.Unlock()
	name := p.Name()
	// 不允许注册插件名为空的插件
	if name == "" {
		return errors.New(500, "plugin name is empty").WithMod("cmd")
	}
	// 判断该插件是否已经注册
	if _, ok := m.plugins[name]; ok {
		return errors.New(500, fmt.Sprintf("Plugin with name %s already registered", name)).WithMod("cmd")
	}
	// 注册插件
	m.plugins[name] = p
	return nil
}

// IsRegistered 判断插件是否已经注册
func (m *PluginManager) IsRegistered(p Plugin) bool {
	m.smu.Lock()
	defer m.smu.Unlock()
	name := p.Name()
	if _, ok := m.plugins[name]; !ok {
		return false
	}
	return true
}

// Range 循环插件
func (m *PluginManager) Range(fn func(n string, p Plugin) bool) {
	for s, plugin := range m.plugins {
		m.smu.Lock()
		b := fn(s, plugin)
		m.smu.Unlock()
		if !b {
			break
		}
	}
}

// RegisterPlugin 注册插件
func RegisterPlugin(p Plugin) error {
	return DefaultPluginManager.Register(p)
}

// IsRegisteredPlugin 判断是否已经注册插件
func IsRegisteredPlugin(p Plugin) bool {
	return DefaultPluginManager.IsRegistered(p)
}

// RangePlugins 循环插件
func RangePlugins(fn func(n string, p Plugin) bool) {
	DefaultPluginManager.Range(fn)
}
