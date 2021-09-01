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
	"fmt"
	"github.com/go-ceres/go-ceres/registry"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestLeijia(t *testing.T) {
	fmt.Println("111111111")
	fmt.Println("111111111")
	fmt.Println("111111111")
}

// 测试注册服务
func TestEtcdRegistry_Register(t *testing.T) {
	conf := DefaultConfig()
	conf.Endpoints = []string{"127.0.0.1:2379"}
	etcdReg := conf.Build()

	srv := &registry.Service{
		Name:    "category",
		Version: "1.0.0",
		Scheme:  "grpc",
		Nodes: []*registry.Node{
			{
				Id:      uuid.NewString(),
				Address: "127.0.0.1:5098",
				Metadata: map[string]string{
					"aaa": "bbb",
				},
			},
		},
	}
	err := etcdReg.Register(srv, registry.RegisterTTl(time.Second*30))
	if err != nil {
		t.Error(err)
	}

	timer := time.NewTimer(time.Second * 20)
	select {
	case l := <-timer.C:
		fmt.Print("当前的时间为：", l)
		err = etcdReg.Deregister(srv)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestEtcdRegistry_Watch(t *testing.T) {

}
