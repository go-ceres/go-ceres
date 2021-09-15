package ceres

import (
	"context"
	"github.com/go-ceres/go-ceres/cmd"
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/registry"
	"github.com/go-ceres/go-ceres/schedule"
	"github.com/go-ceres/go-ceres/server"
	"github.com/go-ceres/go-ceres/utils/cycle"
	"github.com/go-ceres/go-ceres/utils/signals"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
	"runtime"
	"strconv"
	"sync"
)

type Engine struct {
	isSetup      bool               // 是否设置
	rw           *sync.RWMutex      // 读写锁
	cycle        *cycle.Cycle       // 异步运行管理
	servers      []server.Server    // 服务
	schedule     *schedule.Schedule // 定时任务管理
	registry     registry.Registry  // 注册中心
	logger       *logger.Logger     // 日志框架
	beforeStarts []func() error     // 启动前回调
	beforeStops  []func() error     // 停止前回调
	afterStarts  []func() error     // 启动后回调
	afterStops   []func() error     // 停止后回调
	initOnce     sync.Once
	setupOnce    sync.Once
	stopOnce     sync.Once
	clear        func() // 程序结束后回调
}

// New 创建一个启动器
func New(fns ...func() (func(), error)) (*Engine, error) {
	eng := &Engine{}
	if err := eng.MustSetup(fns...); err != nil {
		return nil, err
	}
	return eng, nil
}

// initialize 初始化engine
func (eng *Engine) initialize() {
	eng.initOnce.Do(func() {
		eng.rw = &sync.RWMutex{}
		eng.cycle = cycle.NewCycle()
		eng.servers = make([]server.Server, 0)
		eng.beforeStarts = make([]func() error, 0)
		eng.beforeStops = make([]func() error, 0)
		eng.afterStarts = make([]func() error, 0)
		eng.afterStops = make([]func() error, 0)
		eng.logger = logger.FrameLogger.With(logger.FieldMod("app"))
	})
}

// setup 初始化组件
func (eng *Engine) setup() (err error) {
	eng.setupOnce.Do(func() {
		err = eng.serialUntilError(
			eng.initCmd,
			eng.printBanner,
			eng.initLogger,
			eng.initMaxProcs,
			eng.initCron,
		)
		eng.isSetup = true
	})
	return
}

// initCmd 初始化命令行
func (eng *Engine) initCmd() error {
	var opts []cmd.Option
	// 初始化插件命令行
	cmd.DefaultPluginManager.Range(func(n string, p cmd.Plugin) bool {
		// 获取该插件的命令行
		flags := p.Flags()
		opts = append(opts, cmd.WithFlags(flags))
		return true
	})
	return cmd.DefaultCmd.Init(opts...)
}

// MustSetup 必须初始化
func (eng *Engine) MustSetup(fns ...func() (func(), error)) error {
	eng.initialize()
	if err := eng.setup(); err != nil {
		return err
	}
	return eng.parallelUntilError(fns...)
}

// Server 设置服务
func (eng *Engine) Server(srv server.Server) error {
	if !eng.isSetup {
		return errors.New(errors.CodeAddServerErrorNoSetup, errors.MsgAddServerErrorNoSetup)
	}
	eng.rw.Lock()
	defer eng.rw.Unlock()
	eng.servers = append(eng.servers, srv)
	return nil
}

// Job 添加定时任务
func (eng *Engine) Job(job schedule.Job) (string, error) {
	err := eng.setup()
	if err != nil {
		return "", err
	}
	id, err := eng.schedule.AddJob(job)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(id)), nil
}

// SetRegistry 设置注册中心
func (eng *Engine) SetRegistry(registry registry.Registry) {
	eng.registry = registry
}

// Run 运行
func (eng *Engine) Run() error {
	// 启动前回调
	for _, fn := range eng.beforeStarts {
		if err := fn(); err != nil {
			return err
		}
	}
	// 等待退出信号
	eng.waitSignals()
	// 启动服务
	eng.cycle.Run(eng.startServer)
	// 启动定时任务
	eng.cycle.Run(eng.startCron)
	// 启动后回调
	for _, fn := range eng.afterStarts {
		if err := fn(); err != nil {
			return err
		}
	}
	// 阻止并等待退出
	if err := <-eng.cycle.Wait(); err != nil {
		eng.logger.Error("ceres shutdown with error", logger.FieldMod(errors.ModApp), logger.FieldErr(err))
		return err
	}
	eng.logger.Info("shutdown ceres, bye!", logger.FieldMod(errors.ModApp))
	return nil
}

// Stop 停止
func (eng *Engine) Stop() (err error) {
	eng.stopOnce.Do(func() {
		// 关闭前回调
		for _, fn := range eng.beforeStops {
			_ = fn()
		}
		if eng.registry != nil {
			err = eng.registry.Close()
			if err != nil {
				eng.logger.Error("stop registry close err", logger.FieldMod(errors.ModApp), logger.FieldErr(err))
			}
		}
		//stop servers
		eng.rw.RLock()
		for _, s := range eng.servers {
			func(s server.Server) {
				eng.cycle.Run(s.Stop)
			}(s)
		}
		eng.rw.RUnlock()

		<-eng.cycle.Done()
		// 关闭后回调
		for _, fn := range eng.afterStops {
			_ = fn()
		}
		eng.cycle.Close()
	})
	return
}

// parallelUntilError 并行方式运行，如果有错则返回错误
func (eng *Engine) parallelUntilError(fns ...func() (func(), error)) error {
	group := new(errgroup.Group)
	var clearArr []func()
	for _, fn := range fns {
		tfn := fn
		group.Go(func() error {
			clear, err := tfn()
			if err != nil {
				return err
			}
			clearArr = append(clearArr, clear)
			return nil
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

// printBanner 打印banner
func (eng *Engine) printBanner() error {
	// 如果自定义了banner
	if customBanner != "" {
		printBanner(customBanner)
		return nil
	}
	// 打印默认的banner
	printBanner(defaultBanner)
	return nil
}

// initLogger 初始化日志
func (eng *Engine) initLogger() error {
	// 框架日志
	if !config.Get("ceres.logger.frame").IsEmpty() {
		logger.FrameLogger = logger.ScanConfig("frame").Build()
	}
	// 项目日志
	if !config.Get("ceres.logger.default").IsEmpty() {
		logger.DefaultLogger = logger.ScanConfig("default").Build()
	}
	eng.logger = logger.FrameLogger.With(logger.FieldMod(errors.ModApp))
	return nil
}

// initMaxProcs 初始化MaxProcs
func (eng *Engine) initMaxProcs() error {
	if maxProcs := config.Get("ceres.application.maxProc").Int(0); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			eng.logger.Panicd("auto max procs", logger.FieldMod("engine"), logger.FieldErr(err))
		}
	}
	return nil
}

// initCron 初始化定时任务管理
func (eng *Engine) initCron() error {
	eng.schedule = schedule.ScanConfig("default").WithLogger(&schedule.Logger{Log: eng.logger.AddCallerSkip(1).With(logger.FieldMod("schedule"))}).Build()
	return nil
}

// waitSignals 等待退出信号
func (eng *Engine) waitSignals() {
	eng.logger.Infod("init listen signal", logger.FieldMod(errors.ModApp))
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		//todo: support timeout
		if grace {
			_ = eng.GracefulStop(context.TODO())
		} else {
			_ = eng.Stop()
		}
	})
}

// startServer 启动服务
func (eng *Engine) startServer() error {
	var eg errgroup.Group
	// start multi servers
	for _, s := range eng.servers {
		s := s
		eg.Go(func() (err error) {
			// 如果有注册中心,则注册服务
			if eng.registry != nil {
				err = eng.registerServer(s.Info())
				defer func(eng *Engine, info *server.ServiceInfo) {
					err = eng.unRegisterServer(info)
					if err != nil {
						eng.logger.Errord("unregister service", logger.FieldMod(errors.ModApp))
					}
				}(eng, s.Info())
			}
			eng.logger.Infod("start server", logger.FieldMod(errors.ModApp), logger.FieldValue(s.Info()))
			defer eng.logger.Infod("exit server", logger.FieldMod(errors.ModApp), logger.FieldValue(s.Info()))
			err = s.Start()
			return
		})
	}
	return eg.Wait()
}

// registerServer 注册服务
func (eng *Engine) registerServer(info *server.ServiceInfo) error {
	srv := &registry.Service{
		Name:    info.Name,
		Version: info.Version,
		Scheme:  info.Scheme,
		Nodes: []*registry.Node{
			{
				Id:       info.Id,
				Address:  info.Address,
				Weight:   info.Weight,
				Metadata: info.Metadata,
			},
		},
	}
	return eng.registry.Register(srv)
}

// unRegisterServer 取消注册
func (eng *Engine) unRegisterServer(info *server.ServiceInfo) error {
	srv := &registry.Service{
		Name:    info.Name,
		Version: info.Version,
		Scheme:  info.Scheme,
		Nodes: []*registry.Node{
			{
				Id:       info.Id,
				Address:  info.Address,
				Weight:   info.Weight,
				Metadata: info.Metadata,
			},
		},
	}
	return eng.registry.Deregister(srv)
}

// startCron 启动定时任务管理
func (eng *Engine) startCron() error {
	eng.schedule.Run()
	return nil
}

// GracefulStop 优雅地停止
func (eng *Engine) GracefulStop(todo context.Context) error {

	return nil
}
