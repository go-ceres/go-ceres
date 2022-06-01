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

const (
	NOT_TOKEN             = "-1"
	NOT_TOKEN_MESSAGE     = "未提供Token"
	INVALID_TOKEN         = "-2"
	INVALID_TOKEN_MESSAGE = "Token无效"
	TOKEN_TIMEOUT         = "-3"
	TOKEN_TIMEOUT_MESSAGE = "Token已过期"
	BE_REPLACED           = "-4"
	BE_REPLACED_MESSAGE   = "Token已被顶下线"
	KICK_OUT              = "-5"
	KICK_OUT_MESSAGE      = "Token已被踢下线"
	DEFAULT_MESSAGE       = "当前会话未登录"
)

// ABNORMAL_LIST 异常map
var ABNORMAL_LIST = map[string]bool{
	NOT_TOKEN:     true,
	INVALID_TOKEN: true,
	TOKEN_TIMEOUT: true,
	BE_REPLACED:   true,
	KICK_OUT:      true,
}

type NotLoginError struct {
	LogicType string
	Code      string
	Message   string
}

// NewNotLoginError 创建一个未登录错误
func NewNotLoginError(logicType string, code string, token string) *NotLoginError {
	var msg = ""
	switch code {
	case NOT_TOKEN:
		msg = NOT_TOKEN_MESSAGE
	case INVALID_TOKEN:
		msg = INVALID_TOKEN_MESSAGE
	case TOKEN_TIMEOUT:
		msg = TOKEN_TIMEOUT_MESSAGE
	case BE_REPLACED:
		msg = BE_REPLACED_MESSAGE
	case KICK_OUT:
		msg = KICK_OUT_MESSAGE
	default:
		msg = DEFAULT_MESSAGE
	}
	if len(token) == 0 {
		msg = msg + ":" + token
	}
	return &NotLoginError{
		LogicType: logicType,
		Code:      code,
		Message:   msg,
	}
}

// Error 获取错误信息
func (n *NotLoginError) Error() string {
	return n.Message
}

// Type 获取错误类型
func (n *NotLoginError) Type() TokenErrorType {
	return NotLoginErrorType
}
