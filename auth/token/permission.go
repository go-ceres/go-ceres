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

type Permission interface {
	// GetRouterInfoWithPath 获取指定路径的路由信息
	GetRouterInfoWithPath(path string, unescape bool) (*RouterInfo, error)
	// GetPermissionSlice 获取指定账号指定设备的权限列表
	GetPermissionSlice(loginId string, logicType string) ([]string, error)
	// GetRoleListSlice 获取指定
	GetRoleListSlice(loginId string, logicType string) ([]string, error)
}
