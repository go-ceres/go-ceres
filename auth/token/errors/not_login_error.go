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

import "strings"

const (
	// 没有提供token

	NotToken        string = "-1"
	NotTokenMessage string = "未提供Token"

	// 提供token无效

	InvalidToken        string = "-2"
	InvalidTokenMessage string = "Token无效"

	// 表示token已过期

	TokenTimeout        string = "-3"
	TokenTimeoutMessage string = "Token已过期"

	BeReplaced        string = "-4"
	BeReplacedMessage string = "Token已被顶下线"

	KickOut        string = "-5"
	KickOutMessage string = "Token已被踢下线"

	DefaultMessage string = "当前会话未登录"
)

var AbnormalList = []string{NotToken, InvalidToken, TokenTimeout, BeReplaced, KickOut}

type NotLoginError struct {
	message   string
	Type      string
	LoginType string
}

func (n *NotLoginError) Error() string {
	return n.message
}

type NotLoginErrorOptions struct {
	LoginType string
	ErrType   string
	Token     string
}

type NotLoginErrorOption func(o *NotLoginErrorOptions)

func NotLoginErrorLoginType(loginType string) NotLoginErrorOption {
	return func(o *NotLoginErrorOptions) {
		o.LoginType = loginType
	}
}

func NotLoginErrorErrType(errType string) NotLoginErrorOption {
	return func(o *NotLoginErrorOptions) {
		o.ErrType = errType
	}
}

func NotLoginErrorToken(token string) NotLoginErrorOption {
	return func(o *NotLoginErrorOptions) {
		o.Token = token
	}
}

// NewNotLoginError 建立错误
func NewNotLoginError(message string, loginType string, eType string) *NotLoginError {
	return &NotLoginError{
		message:   message,
		LoginType: loginType,
		Type:      eType,
	}
}

// NewNotLoginInstance 新建登录错误实例
func NewNotLoginInstance(opts ...NotLoginErrorOption) *NotLoginError {
	e := new(NotLoginErrorOptions)
	for _, opt := range opts {
		opt(e)
	}
	msg := ""
	switch e.ErrType {
	case NotToken:
		msg = NotTokenMessage
	case InvalidToken:
		msg = InvalidTokenMessage
	case TokenTimeout:
		msg = TokenTimeoutMessage
	case BeReplaced:
		msg = BeReplacedMessage
	case KickOut:
		msg = KickOutMessage
	default:
		msg = DefaultMessage
	}
	if strings.TrimSpace(e.Token) != "" {
		msg += "：" + e.Token
	}
	return NewNotLoginError(msg, e.LoginType, e.ErrType)
}
