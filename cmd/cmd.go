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
	"os"
)

var (
	DefaultCmd   = newCmd()
	DefaultFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "show version",
		},
		&cli.StringFlag{
			Name:    "region",
			Usage:   "service region",
			EnvVars: []string{"CERES_REGION"},
		},
		&cli.StringFlag{
			Name:    "zone",
			Usage:   "service zone",
			EnvVars: []string{"CERES_ZONE"},
		},
	}
)

type Command struct {
	ctx  *cli.Context
	app  *cli.App
	opts Options
}

// newCmd 创建命令交互
func newCmd(opts ...Option) *Command {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	// 应用名称
	if len(options.Name) == 0 {
		options.Name = appName
	}

	// 名称后的简介
	if len(options.Usage) == 0 {
		options.Usage = "a service for go-ceres"
		if len(usage) > 0 {
			options.Usage = usage
		}
	}

	// 应用描述
	if len(options.Description) == 0 {
		options.Description = "a service for go-ceres"
		if len(usage) > 0 {
			options.Description = description
		}
	}

	cmd := new(Command)
	cmd.opts = options
	cmd.app = cli.NewApp()

	cmd.app.Name = appName
	if len(options.Name) > 0 {
		cmd.app.Name = options.Name
	}
	cmd.app.HideVersion = true
	cmd.app.Description = options.Description
	cmd.app.Usage = options.Usage
	cmd.app.Before = cmd.Before
	cmd.app.Flags = append(cmd.app.Flags, DefaultFlags...)
	cmd.app.Action = func(context *cli.Context) error {
		return nil
	}
	cmd.app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "show version",
			Action: func(context *cli.Context) error {
				ShowVersion()
				os.Exit(0)
				return nil
			},
		},
	}
	return cmd
}

// Init 初始化命令行
func (c *Command) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	c.app.Flags = append(c.app.Flags, c.opts.Flags...)
	// 运行
	c.app.RunAndExitOnError()
	return nil
}

// Before 应用运行前
func (c Command) Before(ctx *cli.Context) (err error) {
	showVersion := ctx.Bool("version")
	if showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// 获取数据中心
	appRegion = ctx.String("region")
	// 获取区域
	appZone = ctx.String("zone")
	// 设置context
	c.ctx = ctx
	// 初始化插件
	DefaultPluginManager.Range(func(n string, p Plugin) bool {
		err = p.Init(ctx)
		if err != nil {
			return false
		}
		return true
	})
	return
}

// App 获取应用信息
func (c *Command) App() *cli.App {
	return c.app
}

// Context 获取应用
func (c *Command) Context() *cli.Context {
	return c.ctx
}
