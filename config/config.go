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

package config

var DefaultConfig = NewConfig()

// Get 获取配置信息
func Get(path string) Value {
	return DefaultConfig.Get(path)
}

func Root() Values {
	return DefaultConfig.Root()
}

// Set 设置配置信息
func Set(path string, data interface{}) error {
	return DefaultConfig.Set(path, data)
}

// OnChange 添加一个配置文件更改监听
func OnChange(fn ChangeFunc) {
	DefaultConfig.OnChange(fn)
}

// Load 从数据源获取配置信息
func Load(sources Source) error {
	return DefaultConfig.LoadSource(sources)
}

// Watch 开启配置文件监听
func Watch() {
	DefaultConfig.Watch()
}

// UnWatch 取消配置文件监听
func UnWatch() {
	DefaultConfig.UnWatch()
}
