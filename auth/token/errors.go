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

// TokenErrorType 定义token错误类型
type TokenErrorType string

const (
	NEVER_EXPIRE     int64 = -1 //常量，表示一个key永不过期 (在一个key被标注为永远不过期时返回此值)
	NOT_VALUE_EXPIRE int64 = -2 //常量，表示系统中不存在这个缓存 (在对不存在的key获取剩余存活时间时返回此值)
)

// 定义所有的错误类型
const (
	// DisableLoginErrorType 封禁账号错误
	DisableLoginErrorType TokenErrorType = "disable-login-error"
	NotLoginErrorType     TokenErrorType = "not-login-error"
)

type tokenError struct {
	errorType TokenErrorType
	message   string
}

// Type 错误类型
func (t *tokenError) Type() TokenErrorType {
	return t.errorType
}

// Error 错误信息
func (t *tokenError) Error() string {
	return t.message
}

// NewTokenError 创建一个token
func NewTokenError(message string) TokenError {
	return &tokenError{
		errorType: "default",
		message:   message,
	}
}

type TokenError interface {
	// Type 错误类型
	Type() TokenErrorType
	// Error 错误信息
	Error() string
}
