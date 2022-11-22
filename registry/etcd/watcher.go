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

package etcd

import (
	"context"
	"github.com/go-ceres/go-ceres/client/etcd"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/registry"
	"go.etcd.io/etcd/client/v3"
	"path"
	"strings"
)

// etcdWatcher 服务观察者
type etcdWatcher struct {
	stop      chan bool          // 停止监控通道
	watchChan clientv3.WatchChan // etcd的监控通道
	client    *etcd.Client       // etcd客户端
}

// NewWatch 创建服务观察者
func NewWatch(r *etcdRegistry, opts ...registry.WatchOption) (registry.Watcher, error) {
	return newEtcdWatcher(r, opts...)
}

// newEtcdWatcher 创建服务观察者
func newEtcdWatcher(r *etcdRegistry, opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}
	// etcd请求上下文
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan bool, 1)
	// 开启线程监控停止监控通道
	go func() {
		<-stop
		cancel()
	}()

	watchPath := r.Config.Prefix
	if len(wo.Service) > 0 {
		watchPath = path.Join(watchPath, strings.Replace(wo.Service, "/", "-", -1))
	}

	if len(wo.Scheme) > 0 {
		watchPath = path.Join(watchPath, wo.Scheme)
	}

	watchPath = watchPath + "/"

	return &etcdWatcher{
		stop:      stop,
		watchChan: r.client.Watch(ctx, watchPath, clientv3.WithPrefix(), clientv3.WithPrevKV()),
		client:    r.client,
	}, nil

}

// Next 获取下一个结果
func (ew *etcdWatcher) Next() (*registry.Result, error) {
	for resp := range ew.watchChan {
		if resp.Err() != nil {
			return nil, resp.Err()
		}
		if resp.Canceled {
			return nil, errors.New(errors.CodeWatcherServiceErrorCanceled, errors.MsgWatcherServiceErrorCanceled)
		}
		for _, ev := range resp.Events {
			service := decode(ev.Kv.Value)
			var action string
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					action = "create"
				} else if ev.IsModify() {
					action = "update"
				}
			case clientv3.EventTypeDelete:
				action = "delete"
				// get service from prevKv
				service = decode(ev.PrevKv.Value)
			}
			if service == nil {
				continue
			}
			return &registry.Result{
				Action:  action,
				Service: service,
			}, nil
		}
	}
	return nil, errors.New(errors.CodeWatcherServiceErrorPassFor, errors.MsgWatcherServiceErrorPassFor)
}

// Stop 停止监控
func (ew *etcdWatcher) Stop() {
	select {
	case <-ew.stop:
		return
	default:
		close(ew.stop)
	}
}
