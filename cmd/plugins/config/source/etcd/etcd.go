//    Copyright 2022. Go-Ceres
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

package etcd

import (
	"github.com/go-ceres/cli/v2"
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/source/etcd"
)

type etcdPlugin struct {
}

// Name 插件名称
func (f *etcdPlugin) Name() string {
	return "config.source.etcd"
}

// Flags 需要注册的Flag命令
func (f *etcdPlugin) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "endpoints",
			Aliases: []string{"e"},
			Usage:   "configuration file path",
			EnvVars: []string{"CERES_CONFIG_ENDPOINTS"},
		}, &cli.StringFlag{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "profile decoder",
			EnvVars: []string{"CERES_CONFIG_DECODE"},
		}, &cli.BoolFlag{
			Name:    "prefix",
			Aliases: []string{"p"},
			Usage:   "etcd path prefix",
			EnvVars: []string{"CERES_CONFIG_PREFIX"},
		},
	}
}

// Init 初始化方法
func (f *etcdPlugin) Init(ctx *cli.Context) error {
	conf := etcd.DefaultConfig()
	conf.Endpoints = ctx.StringSlice("endpoints")
	conf.Encoding = ctx.String("decode")
	err := config.Load(conf.Build())
	if err != nil {
		return err
	}
	return nil
}

// Config 当配置组件初始化完成后
func (f *etcdPlugin) Config() error {
	panic("implement me")
}

// Destroy 当服务销毁时调用
func (f *etcdPlugin) Destroy() {
	panic("implement me")
}

func init() {
	p := &etcdPlugin{}
	err := cmd.RegisterPlugin(p)
	if err != nil {
		logger.FrameLogger.Panicd("register plugin", logger.FieldErr(err))
	}
}
