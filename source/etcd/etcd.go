package etcd

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/errors"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var DefaultPrefix = "/ceres/config/"

type etcdSource struct {
	client  *clientv3.Client
	config  *Config
	changed chan struct{}
	err     error
}

// Read 读取配置
func (e *etcdSource) Read() (*config.DataSet, error) {
	if e.err != nil {
		return nil, e.err
	}
	kv := clientv3.NewKV(e.client)
	res, err := kv.Get(e.config.Ctx, e.config.Prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	data := makeMapData(res.Kvs, e.config.TrimPrefix)
	b, err := Marshals[e.config.Encoding](data)
	if err != nil {
		return nil, errors.New(500, "error reading source: "+err.Error())
	}

	cs := &config.DataSet{
		Format:    e.getUnmarshal(),
		Source:    e.String(),
		Timestamp: time.Now(),
		Data:      b,
	}
	return cs, nil
}

// Write 写入配置
func (e *etcdSource) Write(set *config.DataSet) error {
	panic("implement me")
}

// IsChanged 配置变化通道
func (e *etcdSource) IsChanged() <-chan struct{} {
	return e.changed
}

// Watch 开启监听
func (e *etcdSource) Watch() {
	if e.changed == nil {
		e.changed = make(chan struct{})
	}
	go func() {
		ch := e.client.Watch(e.config.Ctx, e.config.Prefix, clientv3.WithPrefix())
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
				select {
				case e.changed <- struct{}{}:
				default:
				}
			}
		}
	}()
}

// String 获取监控名称
func (e *etcdSource) String() string {
	return "etcd"
}

// UnWatch 取消监控
func (e *etcdSource) UnWatch() {
	close(e.changed)
	e.changed = nil
	e.client.Watcher.Close()
}

// getUnmarshal 获取解码
func (e *etcdSource) getUnmarshal() string {
	if e.config.Encoding != "" {
		return e.config.Encoding
	}
	return "json"
}

// NewSource 创建一个资源
func NewSource(c *Config) config.Source {
	// 处理前缀
	if c.TrimPrefix == "" {
		c.TrimPrefix = c.Prefix
	}
	cli, err := clientv3.New(*c.Config)
	return &etcdSource{
		client: cli,
		config: c,
		err:    err,
	}
}
