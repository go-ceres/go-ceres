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
	// CodeAddServerErrorNoSetup 添加服务错误（在添加服务前没有注册）
	CodeAddServerErrorNoSetup          = 4000
	CodeRegisterServerErrorNoNode      = 4001
	CodeWatcherServiceErrorCanceled    = 4002
	CodeWatcherServiceErrorPassFor     = 4003
	CodeGetServiceErrorNotFound        = 4004
	CodeWatchServiceErrorNoServiceName = 4005

	// CodeAddScheduleToMaximum 添加定时任务错误（添加数量超过定义数量）
	CodeAddScheduleToMaximum = 5000

	CodeCallGrpcUnaryTimeout = 6000
)
