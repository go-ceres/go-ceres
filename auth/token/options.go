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

import "github.com/go-ceres/go-ceres/auth/token/entity"

type loginOptions struct {
	device  string // 设备标识
	timeout int64  // 当前此次登录的过期时间
}

func (l *loginOptions) Device() string {
	return l.device
}

func (l *loginOptions) SetDevice(device string) {
	l.device = device
}

func (l *loginOptions) Timeout() int64 {
	return l.timeout
}

func (l *loginOptions) SetTimeout(timeout int64) {
	l.timeout = timeout
}

type LoginOption func(o *loginOptions)

// DefaultOption 默认的登录额外参数
func DefaultOption(timeout int64) *loginOptions {
	opts := &loginOptions{
		timeout: timeout,
		device:  entity.DefaultLoginDevice,
	}
	return opts
}

// Device 客户端
func Device(device string) LoginOption {
	return func(o *loginOptions) {
		o.device = device
	}
}

// Timeout 当前此次登录的过期时间
func Timeout(timeout int64) LoginOption {
	return func(o *loginOptions) {
		o.timeout = timeout
	}
}
