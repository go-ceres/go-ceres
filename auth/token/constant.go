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

// TokenStyle token样式
type TokenStyle string

// NoLoginErrorValue 没有登录的错误
type NoLoginErrorValue int8

const (
	TOKEN_STYLE_UUID        TokenStyle = "uuid"        // uuid样式
	TOKEN_STYLE_SIMPLE_UUID TokenStyle = "simple-uuid" // uuid不带下划线
	TOKEN_STYLE_RANDOM_32   TokenStyle = "random-32"   // 随机32位字符串
	TOKEN_STYLE_RANDOM_64   TokenStyle = "random-64"   // 随机64位字符串
	TOKEN_STYLE_JWT         TokenStyle = "jwt"         // jwt格式
)

const (
	defaultLoginDevice = "device_default" // 默认设备
)

const (
	disableValue = "disable"
)
