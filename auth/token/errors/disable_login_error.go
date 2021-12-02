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

package errors

const (
	// =======账号封禁相关==========

	DisableLoginErrorValue   = "disable"
	DisableLoginErrorMessage = "此账号被封禁"
)

// DisableLoginError 账号封禁错误
type DisableLoginError struct {
	Message     string      // 提示语
	loginType   string      // 登录类型
	loginId     interface{} // 封禁用户id
	disableTime int64       // 封禁时间
}

// NewDisableLoginError 构建封禁错误
func NewDisableLoginError(loginType string, id interface{}, disableTime int64) *DisableLoginError {
	return &DisableLoginError{
		Message:     DisableLoginErrorMessage,
		loginType:   loginType,
		loginId:     id,
		disableTime: disableTime,
	}
}

// LoginType 获取登录类型
func (d *DisableLoginError) LoginType() string {
	return d.loginType
}

// LoginId 获取登录用户id
func (d *DisableLoginError) LoginId() interface{} {
	return d.loginId
}

// DisableTime 获取封禁时长
func (d *DisableLoginError) DisableTime() int64 {
	return d.disableTime
}

// Error 实现错误接口
func (d *DisableLoginError) Error() string {
	return d.Message
}
