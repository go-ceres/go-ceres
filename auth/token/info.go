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

// Info token详细信息
type Info struct {
	TokenName            string `json:"token_name"`             // token名称
	TokenValue           string `json:"token_value"`            // token值
	IsLogin              bool   `json:"is_login"`               // 是否登录
	LogicType            string `json:"logic_type"`             // 逻辑类型
	TokenTimeout         int64  `json:"token_timeout"`          // 当前token剩余时间
	SessionTimeout       int64  `json:"session_timeout"`        // 当前用户session有效时间
	TokenSessionTimeout  int64  `json:"token_session_timeout"`  // 当前token的有效时间
	TokenActivityTimeout int64  `json:"token_activity_timeout"` // 无操作时token剩余时间
	LoginDevice          string `json:"login_device"`           // 登录设备
}
