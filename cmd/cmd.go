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

import (
	"github.com/go-ceres/cli/v2"
	"sync"
)

var (
	DefaultCmd = newCmd()
	DefaultPluginManager = &PluginManager{
		smu: sync.RWMutex{},
		plugins: make(map[string]Plugin),
	}
)

type Command struct {
	ctx *cli.Context
	app *cli.App
	opts Options
}

// newCmd 创建命令交互
func newCmd(opts ...Option) *Command {
	options := Options{
	}
	for _, o := range opts {
		o(&options)
	}

	if len(options.Description) == 0 {
		options.Description = "a service for go-ceres"
	}

	cmd := new(Command)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = cmd.opts.Name
	cmd.app.Version = cmd.opts.Version
	cmd.app.Description = cmd.opts.Description
	cmd.app.Before = cmd.Before
	cmd.app.Action = cmd.Action

	return cmd
}

// Init 初始化命令行
func (c *Command) Init(opts ...Option) error {

}

// Before 应用运行前
func (c Command) Before(ctx *cli.Context) error {
	
}

// Action 在启动时进入
func (c Command) Action(ctx *cli.Context) error {
	
}


func (c *Command) App() *cli.App {
	return c.app
}

func (c *Command) Context() *cli.Context {
	return c.ctx
}