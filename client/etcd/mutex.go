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
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

// Mutex 互斥锁
type Mutex struct {
	s *concurrency.Session
	m *concurrency.Mutex
}

// NewMutex 创建一个互斥锁
func (c *Client) NewMutex(key string, opts ...concurrency.SessionOption) (*Mutex, error) {
	mutex := &Mutex{}
	// 默认session ttl = 60s
	session, err := concurrency.NewSession(c.Client, opts...)
	if err != nil {
		return nil, err
	}
	mutex.s = session
	mutex.m = concurrency.NewMutex(mutex.s, key)
	return mutex, nil
}

// Lock 获取锁
func (mutex *Mutex) Lock(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.Lock(ctx)
}

// TryLock 尝试获取锁
func (mutex *Mutex) TryLock(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.TryLock(ctx)
}

// Unlock 解除锁
func (mutex *Mutex) Unlock() error {
	return mutex.m.Unlock(context.TODO())
}

// Close 关闭会话
func (mutex *Mutex) Close() error {
	return mutex.s.Close()
}
