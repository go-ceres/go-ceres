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

package manager

import (
	"github.com/go-ceres/go-ceres/auth/token/stp"
	"github.com/go-ceres/go-ceres/cache"
)

var defaultManager Manager

type Manager struct {
	Store cache.Cache   // 持久化接口
	Stp   stp.Interface // 权限接口
}

func SetStorage(store cache.Cache) {
	defaultManager.Store = store
}

func SetStp(stp stp.Interface) {
	defaultManager.Stp = stp
}

func Storage() cache.Cache {
	return defaultManager.Store
}

func Stp() stp.Interface {
	return defaultManager.Stp
}
