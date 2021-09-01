//  Copyright 2020 Go-Ceres
//  Author https://github.com/go-ceres/go-ceres
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package logger

import "go.uber.org/zap"

// FieldAid 应用id
func FieldAid(aid string) zap.Field {
	return zap.String("aid", aid)
}

// FieldMod 日志发生的模块
func FieldMod(mod string) zap.Field {
	return zap.String("mod", mod)
}

// FieldTid 日志发生的链路追踪id
func FieldTid(tid string) zap.Field {
	return zap.String("tid", tid)
}

// FieldHost 日志发生的主机名
func FieldHost(host string) zap.Field {
	return zap.String("host", host)
}

// FieldIp 日志发生的主机ip
func FieldIp(ip string) zap.Field {
	return zap.String("ip", ip)
}

// FieldCode 日志发生的执行代码
func FieldCode(code interface{}) zap.Field {
	return zap.Any("code", code)
}

// FieldErr 错误信息
func FieldErr(err error) zap.Field {
	return zap.Any("err", err)
}

// FieldValue 错误发生相关的值
func FieldValue(value interface{}) zap.Field {
	return zap.Any("value", value)
}

// FieldString 二级封装zap.string
func FieldString(key, val string) zap.Field {
	return zap.String(key, val)
}

// FieldAny 二级封装zap.any
func FieldAny(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
