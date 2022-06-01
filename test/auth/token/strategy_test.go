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

package token

import (
	"encoding/json"
	"fmt"
	"github.com/go-ceres/go-ceres/auth/token"
	"github.com/go-ceres/go-ceres/cache/redis"
	"reflect"
	"testing"
)

type AAA struct {
	Name string
	Data string
}

//func registerType(elem interface{}) {
//	t := reflect.TypeOf(elem).Elem()
//	typeRegistry[t.Name()] = t
//}

func TestXuliehua(t *testing.T) {
	arr1 := string("111")
	arr2 := [6]string{}
	fmt.Printf("type of arr1 is %s, the kind is %s\n", reflect.TypeOf(arr1), reflect.TypeOf(arr1).Kind())
	fmt.Printf("type of arr2 is %s, the kind is %s\n", reflect.TypeOf(arr2), reflect.TypeOf(arr2).Kind())

}

func TestSetCreateToken(t *testing.T) {
	storage := redis.DefaultConfig().Build()
	tokenManager := token.ScanConfig("ceshi").Build("user").SetStorage(storage)
	login, err := tokenManager.Login("123456")
	if err != nil {
		return
	}
	fmt.Println(login)
	fmt.Println(tokenManager.IsLogin(login))
}

func TestCheckLogin(t *testing.T) {
	storage := redis.DefaultConfig().Build()
	tokenManager := token.ScanConfig("ceshi").Build("user").SetStorage(storage)
	marshal, err := json.Marshal(tokenManager.GetTokenInfo("a22a435b-0b1a-4321-8426-6ce0ea8ab40f"))
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}

func TestGetLoginId(t *testing.T) {
	storage := redis.DefaultConfig().Build()
	tokenManager := token.ScanConfig("ceshi").Build("user").SetStorage(storage)
	id := tokenManager.GetTokenInfo("a22a435b-0b1a-4321-8426-6ce0ea8ab40f")
	tokenManager.CheckLogin("a22a435b-0b1a-4321-8426-6ce0ea8ab40f")
	data, _ := json.Marshal(id)
	fmt.Println(string(data))
}
