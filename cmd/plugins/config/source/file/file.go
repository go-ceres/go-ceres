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

package file

import (
	"github.com/go-ceres/cli/v2"
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/source/file"
)

// filePlugin
type filePlugin struct {
	source config.Source
}

// Name 插件名称
func (f *filePlugin) Name() string {
	return "config.source.file"
}

// Flags 需要注册的Flag命令
func (f *filePlugin) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "configuration file path",
			EnvVars: []string{"CERES_CONFIG_FILE"},
		}, &cli.StringFlag{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "profile decoder",
			EnvVars: []string{"CERES_CONFIG_DECODE"},
		}, &cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Whether to monitor configuration changes",
			EnvVars: []string{"CERES_CONFIG_WATCH"},
		},
	}
}

// Init 初始化方法
func (f *filePlugin) Init(ctx *cli.Context) error {
	path := ctx.String("file")
	if path == "" {
		path = "./config/config.toml"
	}
	var opts []file.Option
	decode := ctx.String("decode")
	if decode != "" {
		opts = append(opts, file.Unmarshal(decode))
	}
	f.source = file.NewSource(path, opts...)
	err := config.Load(f.source)
	if err != nil {
		return err
	}
	return nil
}

// Config 当配置组件初始化完成后
func (f *filePlugin) Config() error {
	panic("implement me")
}

// Destroy 当服务销毁时调用
func (f *filePlugin) Destroy() {
	panic("implement me")
}

func init() {
	p := &filePlugin{}
	err := cmd.RegisterPlugin(p)
	if err != nil {
		logger.FrameLogger.Panicd("register plugin", logger.FieldErr(err))
	}
}
