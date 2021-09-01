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
	// MsgAddServerErrorNoSetup 添加服务错误（在添加服务前没有注册）
	MsgAddServerErrorNoSetup          = "Please call setup before add service"
	MsgRegisterServerErrorNoNode      = "Require at least one node"
	MsgWatcherServiceErrorCanceled    = "could not get next,result is Canceled"
	MsgWatcherServiceErrorPassFor     = "could not get next,pass for"
	MsgGetServiceErrorNotFound        = "service not found"
	MsgWatchServiceErrorNoServiceName = "missing service parameter"
	// MsgAddScheduleToMaximum schedule 统一错误信息
	MsgAddScheduleToMaximum = "Maximum number of tasks added exceeded"

	MsgCallGrpcUnaryTimeout = "call grpc unary timeout"
)
