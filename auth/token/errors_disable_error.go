//    Copyright 2022. Go-Ceres
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

import "fmt"

type DisableLoginError struct {
	LoginId     string `json:"login_id"`
	LogicType   string `json:"logic_type"`
	DisableTime int64  `json:"disable_time"`
}

func (d *DisableLoginError) Type() TokenErrorType {
	return DisableLoginErrorType
}

// NewDisableLoginError 创建封禁错误
func NewDisableLoginError(LoginId string, LogicType string, DisableTime int64) TokenError {
	return &DisableLoginError{
		LoginId:     LoginId,
		LogicType:   LogicType,
		DisableTime: DisableTime,
	}
}

// Error 详细的错误信息
func (d *DisableLoginError) Error() string {
	return fmt.Sprintf("loginId:%s账号已经被封禁，", d.LoginId)
}

// Message 错误提示信息
func (d *DisableLoginError) Message() string {
	return "该账号被封禁"
}
