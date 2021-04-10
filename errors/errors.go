/*
 * //  Copyright 2021 Go-Ceres
 * //  Author https://github.com/go-ceres/go-ceres
 * //
 * //  Licensed under the Apache License, Version 2.0 (the "License");
 * //  you may not use this file except in compliance with the License.
 * //  You may obtain a copy of the License at
 * //
 * //      http://www.apache.org/licenses/LICENSE-2.0
 * //
 * //  Unless required by applicable law or agreed to in writing, software
 * //  distributed under the License is distributed on an "AS IS" BASIS,
 * //  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * //  See the License for the specific language governing permissions and
 * //  limitations under the License.
 */

package errors

import "encoding/json"

// 错误码，低于10000则为系统错误

type Error struct {
	Code 		int 			`json:"code"`					// 错误码
	Msg			string 			`json:"msg"`					// 错误信息
	Data 		interface{} 	`json:"data"`					// 数据
	Tid			string 			`json:"tid,omitempty"`			// 链路追踪id，可忽略
	Aid			string 			`json:"aid,omitempty"`			// 应用id，可忽略
	Mod			string 			`json:"mod,omitempty"`			// 出错模块，可忽略
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Msg
}

// String 打印所有信息
func (e *Error) String() string {
	d, _ := json.Marshal(e)
	return string(d)
}

// New 创建一个基础错误信息
func New(code int,msg string) *Error {
	return &Error{
		Code: code,
		Msg: msg,
		Data: nil,
	}
}

// FromError 从错误获取
func FromError(err error) *Error {
	e,ok:=err.(*Error)
	if ok {
		return e
	}
	return New(99999,err.Error())
}

// WithMsg	设置错误信息
func (e *Error) WithMsg(msg string) *Error {
	e.Msg = msg
	return e
}

// WithData 设置数据
func (e *Error) WithData(data interface{}) *Error {
	e.Data = data
	return e
}

// WithTid 设置trace id
func (e *Error) WithTid(tid string) *Error {
	e.Tid = tid
	return e
}

// WithAid 设置应用ID
func (e *Error) WithAid(aid string) *Error {
	e.Aid = aid
	return e
}

// WithMod 设置模块
func (e *Error) WithMod(mod string) *Error {
	e.Mod = mod
	return e
}


