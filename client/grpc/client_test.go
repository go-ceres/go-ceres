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

package grpc

import (
	"fmt"
	"testing"
)

func TestGetConn(t *testing.T) {
	conf := DefaultConfig()
	client := conf.Build()
	conn, err := client.Conn("etcd_server,aaa,ddd", Version("1.0.0"))
	if err != nil {
		return
	}
	fmt.Println(conn)

}

func TestParseTarget(t *testing.T) {
	service := "main"
	data := parseTarget(service)
	fmt.Println(data.Scheme)
	fmt.Println(data.Endpoint)
	fmt.Println(data.Authority)
}
