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

package writer

import (
	"io"
	"sync"
)

var DefaultBuilders = newBuilders()

// Manager 构造器管理
type Manager struct {
	builders sync.Map
}

// Builder 构造器接口
type Builder interface {
	Build(c interface{}) Writer
}

// Writer 日志写入者
type Writer interface {
	io.WriteCloser
	Name() string
	Rotate() error
}

//
func newBuilders() *Manager {
	return &Manager{
		builders: sync.Map{},
	}
}

// Register 注册一个writer的构造器
func (m *Manager) Register(name string, build Builder) {
	m.builders.Store(name, build)
}

// UnRegister 删除一个已经存在的构造器
func (m *Manager) UnRegister(name string) {
	m.builders.Delete(name)
}

// Range 循环
func (m *Manager) Range(fn func(key string, build Builder) bool) {
	m.builders.Range(func(key, value interface{}) bool {
		strKey := key.(string)
		build := value.(Builder)
		return fn(strKey, build)
	})
}

// Load 读取
func (m *Manager) Load(key string) (Builder, bool) {
	val, ok := m.builders.Load(key)
	if ok {
		value, ok := val.(Builder)
		return value, ok
	}
	return nil, false
}

// Register 注册构造器
func Register(name string, builder Builder) {
	DefaultBuilders.Register(name, builder)
}

// UnRegister 删除一个已经存在的构造器
func UnRegister(name string) {
	DefaultBuilders.UnRegister(name)
}

// Load 获取构造器
func Load(name string) (Builder, bool) {
	return DefaultBuilders.Load(name)
}

// Range 循环构造器
func Range(fn func(key string, build Builder) bool) {
	DefaultBuilders.Range(fn)
}
