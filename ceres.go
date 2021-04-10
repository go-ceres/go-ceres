package ceres

import (
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/go-ceres/go-ceres/registry"
	"github.com/go-ceres/go-ceres/server"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Engine struct {
	server			server.Server		// 服务
	registry		registry.Registry	// 注册中心
	cmd 			*cmd.Command		// 命令管理
	beforeStarts 	[]func() error		// 启动前回调
	beforeStops  	[]func() error		// 停止前回调
	afterStarts  	[]func() error		// 启动后回调
	afterStops   	[]func() error		// 停止后回调
	initOnce 		sync.Once
	setupOnce		sync.Once
	stopOnce		sync.Once
}

// New 创建一个启动器
func New(fns ...func() error) (*Engine,error) {
	eng := &Engine{}
	if err := eng.Setup(fns...); err != nil {
		return nil, err
	}
	return eng, nil
}

// initialize 初始化engine
func (eng *Engine) initialize()  {
	eng.initOnce.Do(func() {
		eng.beforeStarts	=	make([]func() error,0)
		eng.beforeStops		=	make([]func() error,0)
		eng.afterStarts		=	make([]func() error,0)
		eng.afterStops		=	make([]func() error,0)
	})
}

// setup 初始化组件
func (eng *Engine) setup() (err error) {
	eng.setupOnce.Do(func() {
		err = eng.serialUntilError(


		)
	})
	return
}

// Setup 额外注册设置
func (eng *Engine) Setup(fns ...func() error) error {
	eng.initialize()
	if err := eng.setup(); err != nil {
		return err
	}
	return eng.parallelUntilError(fns...)
}

// Run 运行
func (eng *Engine) Run() error {

}

// Stop 停止
func (eng *Engine) Stop() (err error) {

}

// parallelUntilError 并行方式运行，如果有错则返回错误
func (eng *Engine) parallelUntilError(fns ...func() error) error {
	group := new(errgroup.Group)
	for _, fn := range fns {
		tfn := fn
		group.Go(func() error {
			return tfn()
		})
	}
	return group.Wait()
}

// serialUntilError 串行方式运行方法
func (eng *Engine) serialUntilError(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}