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
	"encoding/json"
	"github.com/go-ceres/go-ceres/client/etcdv3"
	"github.com/go-ceres/go-ceres/errors"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/registry"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"path"
	"strings"
	"sync"
)

type etcdRegistry struct {
	client   *etcdv3.Client // 客户端
	rw       *sync.RWMutex  // 读写锁
	cache    sync.Map       // 注册缓存信息
	sessions map[string]*concurrency.Session
	*Config
}

// registerStore 注册服务缓存信息
type registerStore struct {
	service *registry.Service
	opts    *registry.RegisterOptions
}

// newRegistry 新建etcd注册中心
func newRegistry(c *Config) *etcdRegistry {
	var client *etcdv3.Client
	if c.etcdClient != nil {
		client = c.etcdClient
	} else {
		if c.Config != nil {
			client = c.Config.Build()
		}
	}
	return &etcdRegistry{
		client:   client,
		rw:       &sync.RWMutex{},
		Config:   c,
		sessions: make(map[string]*concurrency.Session),
		cache:    sync.Map{},
	}
}

// String 获取注册中心实例名称
func (e *etcdRegistry) String() string {
	return "etcd"
}

// Register 注册服务
func (e *etcdRegistry) Register(srv *registry.Service, opts ...registry.RegisterOption) error {
	if len(srv.Nodes) == 0 {
		return errors.New(errors.CodeRegisterServerErrorNoNode, errors.MsgRegisterServerErrorNoNode)
	}
	var gerr error
	// register each node individually
	for _, node := range srv.Nodes {
		err := e.registerNode(srv, node, opts...)
		if err != nil {
			gerr = err
		}
	}
	return gerr
}

// RegisterNode 注册节点
func (e *etcdRegistry) registerNode(srv *registry.Service, node *registry.Node, opts ...registry.RegisterOption) error {
	// 服务信息
	service := &registry.Service{
		Name:     srv.Name,
		Version:  srv.Version,
		Metadata: srv.Metadata,
		Scheme:   srv.Scheme,
		Nodes:    []*registry.Node{node},
	}
	// 操作参数
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	// 存储信息
	store := registerStore{
		service: service,
		opts:    &options,
	}
	// 组装请求参数
	opOptions := make([]clientv3.OpOption, 0)
	// 如果参数包含ttl
	if ttl := options.TTL.Seconds(); ttl > 0 {
		sess, err := e.getSession(nodePath(e.Prefix, srv.Name, srv.Scheme, node.Id), concurrency.WithTTL(int(ttl)))
		if err != nil {
			return err
		}
		opOptions = append(opOptions, clientv3.WithLease(sess.Lease()))
	} else if ttl := e.Config.ServiceTTL.Seconds(); ttl > 0 {
		sess, err := e.getSession(nodePath(e.Prefix, srv.Name, srv.Scheme, node.Id), concurrency.WithTTL(int(ttl)))
		if err != nil {
			return err
		}
		opOptions = append(opOptions, clientv3.WithLease(sess.Lease()))
	}
	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), e.Config.ReadTimeout)
	defer cancel()
	// etcd存储键
	key := nodePath(e.Prefix, srv.Name, srv.Scheme, node.Id)
	val := encode(service)
	// 发起请求
	_, err := e.client.Put(ctx, key, val, opOptions...)
	if err != nil {
		e.log.Errord("register service", logger.FieldMod(errors.ModRegistryEtcd), logger.FieldErr(err), logger.FieldAny("etcd.key", key), logger.FieldValue(service))
		return err
	}
	e.log.Infod("register service node", logger.FieldAny("etcd.key", key), logger.FieldValue(service))
	e.cache.Store(key, store)
	return nil
}

// Deregister 注销服务
func (e *etcdRegistry) Deregister(srv *registry.Service, opts ...registry.DeRegisterOption) error {
	if len(srv.Nodes) == 0 {
		return errors.New(errors.CodeRegisterServerErrorNoNode, errors.MsgRegisterServerErrorNoNode)
	}
	// 操作参数
	options := registry.DeRegisterOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	if _, ok := options.Context.Deadline(); !ok {
		var cancel context.CancelFunc
		options.Context, cancel = context.WithTimeout(options.Context, e.Config.ReadTimeout)
		defer cancel()
	}
	// 循环注销节点
	for _, node := range srv.Nodes {
		// etcd存储键
		key := nodePath(e.Prefix, srv.Scheme, srv.Name, node.Id)
		// 删除session
		if err := e.delSession(key); err != nil {
			return err
		}
		_, err := e.client.Delete(options.Context, key)
		if err != nil {
			e.log.Errord("deregister service node", logger.FieldMod(errors.ModRegistryEtcd), logger.FieldErr(err), logger.FieldAny("etcd.key", key), logger.FieldValue(srv))
			return err
		}
		e.cache.Delete(key)
		e.log.Infod("deregister service", logger.FieldAny("etcd.key", key), logger.FieldValue(srv))
	}
	return nil
}

// GetService 根据名称获取服务列表
func (e *etcdRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	opt := &registry.GetOptions{}
	for _, option := range opts {
		option(opt)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.ReadTimeout)
	defer cancel()
	rsp, err := e.client.Get(ctx, servicePath(e.Prefix, name, opt.Scheme)+"/", clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	if len(rsp.Kvs) == 0 {
		return nil, registry.ErrNotFound
	}
	serviceMap := map[string]*registry.Service{}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			s, ok := serviceMap[sn.Version]
			if !ok {
				s = &registry.Service{
					Name:     sn.Name,
					Version:  sn.Version,
					Metadata: sn.Metadata,
					Scheme:   sn.Scheme,
				}
				serviceMap[s.Version] = s
			}

			s.Nodes = append(s.Nodes, sn.Nodes...)
		}
	}

	services := make([]*registry.Service, 0, len(serviceMap))
	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services, nil
}

// WatchService 监控服务
func (e *etcdRegistry) WatchService(opts ...registry.WatchOption) (chan registry.WatchServiceResult, error) {
	// 获取当前存在的列表
	opt := &registry.WatchOptions{}
	for _, o := range opts {
		o(opt)
	}
	if len(opt.Service) == 0 {
		return nil, registry.ErrNoServiceName
	}
	// 获取列表
	var getOptions []registry.GetOption
	if len(opt.Scheme) > 0 {
		getOptions = append(getOptions, registry.GetScheme(opt.Scheme))
	}
	services, err := e.GetService(opt.Service, getOptions...)
	if err != nil {
		return nil, err
	}
	var addresses = make(chan registry.WatchServiceResult, 10)
	var result = &registry.WatchServiceResult{
		Services: []registry.Service{},
	}
	for _, service := range services {
		e.update(result, &registry.Result{Service: service, Action: "create"})
	}

	addresses <- *result.DeepCopy()

	watcher, err := newEtcdWatcher(e, opts...)
	if err != nil {
		return nil, err
	}
	ch := make(chan bool)

	go func() {
		select {
		case <-ch:
			watcher.Stop()
		}
	}()

	go func() {
		for {
			res, err := watcher.Next()
			if err != nil {
				close(ch)
			}
			e.update(result, res)
			out := result.DeepCopy()
			select {
			// case addresses <- snapshot:
			case addresses <- *out:
			default:
				e.log.Warnf("invalid")
			}
		}
	}()
	// 获取现有的服务
	return addresses, nil
}

// update 根据操作修改监听结果集
func (e *etcdRegistry) update(result *registry.WatchServiceResult, res *registry.Result) {
	if res == nil || res.Service == nil || result == nil {
		return
	}
	e.rw.Lock()
	e.rw.Unlock()

	var service *registry.Service
	var index int
	for i, s := range result.Services {
		if s.Version == res.Service.Version {
			service = &s
			index = i
		}
	}

	switch res.Action {
	case "create", "update":
		// 如果没有找到
		if service == nil {
			result.Services = append(result.Services, *res.Service)
			return
		}
		for _, cur := range service.Nodes {
			var seen bool
			for _, node := range res.Service.Nodes {
				if cur.Id == node.Id {
					seen = true
					break
				}
			}
			if !seen {
				res.Service.Nodes = append(res.Service.Nodes, cur)
			}
		}
		result.Services[index] = *res.Service
	case "delete":
		if service == nil {
			return
		}

		var nodes []*registry.Node

		// filter cur nodes to remove the dead one
		for _, cur := range service.Nodes {
			var seen bool
			for _, del := range res.Service.Nodes {
				if del.Id == cur.Id {
					seen = true
					break
				}
			}
			if !seen {
				nodes = append(nodes, cur)
			}
		}

		// still got nodes, save and return
		if len(nodes) > 0 {
			service.Nodes = nodes
			result.Services[index] = *service
			return
		}
	}
}

// ListService 获取服务列表
func (e *etcdRegistry) ListService() ([]*registry.Service, error) {
	panic("implement me")
}

// Close 关闭连接
func (e *etcdRegistry) Close() error {
	panic("implement me")
}

// getSession 获取session
func (e *etcdRegistry) getSession(key string, opts ...concurrency.SessionOption) (*concurrency.Session, error) {
	e.rw.RLock()
	sess, ok := e.sessions[key]
	e.rw.RUnlock()
	if ok {
		return sess, nil
	}
	sess, err := concurrency.NewSession(e.client.Client, opts...)
	if err != nil {
		return sess, err
	}
	e.rw.Lock()
	e.sessions[key] = sess
	e.rw.Unlock()
	return sess, nil
}

// delSession 删除session
func (e *etcdRegistry) delSession(key string) error {
	e.rw.RLock()
	sess, ok := e.sessions[key]
	e.rw.RUnlock()
	if ok {
		e.rw.Lock()
		delete(e.sessions, key)
		e.rw.Unlock()
		if err := sess.Close(); err != nil {
			return err
		}
	}
	return nil
}

// nodePath 组装节点注册路径
func nodePath(prefix, name, scheme, id string) string {
	service := strings.Replace(name, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(prefix, service, scheme, node)
}

// servicePath 组装服务注册路径
func servicePath(prefix, name, scheme string) string {
	pathStr := path.Join(prefix, strings.Replace(name, "/", "-", -1))
	if len(scheme) > 0 {
		pathStr = path.Join(pathStr, scheme)
	}
	return pathStr
}

// decode 解码
func decode(ds []byte) *registry.Service {
	var s *registry.Service
	_ = json.Unmarshal(ds, &s)
	return s
}

// encode 编码
func encode(s *registry.Service) string {
	b, _ := json.Marshal(s)
	return string(b)
}
