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

package cmd

import "github.com/go-ceres/cli/v2"

type Options struct {
	Name string
	Description string
	Version     string
	Flags		[]cli.Flag
}

type Option func(o *Options)

// Name 设置应用名称
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Name 设置应用介绍
func Description(desc string) Option {
	return func(o *Options) {
		o.Description = desc
	}
}

// Name 设置应用介绍
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Name 添加命令参数
func Flags(flags []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flags...)
	}
}
