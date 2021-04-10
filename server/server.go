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

package server

import "github.com/go-ceres/go-ceres/cmd"

// 服务接口
type Server interface {
	// 获取命令
	Command() *cmd.Command
	// t停止
	Stop() error
	// 启动
	Start() error
	// 服务信息
	Info() *ServiceInfo
}

// 服务信息
type ServiceInfo struct {
	Id 			string 		`json:"id"`				// 应用ID
	Name 		string 		`json:"name"`			// 应用名称
	Namespace	string 		`json:"namespace"`		// 应用空间
	Address 	string 		`json:"address"`		// 服务地址
	Version		string 		`json:"version"`		// 当前版本
}



