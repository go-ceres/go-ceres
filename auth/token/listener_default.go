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

import "fmt"

// defaultListener 默认的监听器
type defaultListener struct {
}

func (d defaultListener) DoLogin(loginType string, loginId string, options loginOptions) {
	fmt.Printf("登录逻辑：%s的用户：%s在设备：%s上登录成功", loginType, loginId, options.device)
}

func (d defaultListener) DoLogout(loginType string, loginId string, tokenValue string) {
	fmt.Printf("登录逻辑：%s的用户：%s使用的token：%s退出成功", loginType, loginId, tokenValue)
}

func (d defaultListener) DoKickout(loginType string, loginId string, tokenValue string) {
	fmt.Printf("登录逻辑：%s的用户：%s使用的token：%s被踢下线", loginType, loginId, tokenValue)
}

func (d defaultListener) DoReplaced(loginType string, loginId string, tokenValue string) {
	fmt.Printf("登录逻辑：%s的用户：%s使用的token：%s被顶下线", loginType, loginId, tokenValue)
}

func (d defaultListener) DoDisable(loginType string, loginId string, disableTime int64) {
	fmt.Printf("登录逻辑：%s的用户：%s被禁止登录，禁止时间：%d秒", loginType, loginId, disableTime)
}

func (d defaultListener) DoUntieDisable(loginType string, loginId string) {
	fmt.Printf("登录逻辑：%s的用户：%s解除封禁", loginType, loginId)
}

func (d defaultListener) DoCreateSession(id string) {
	fmt.Printf("用户：%s创建session成功", id)
}

func (d defaultListener) DoLogoutSession(id string) {
	fmt.Printf("用户：%s退出成功", id)
}

// newDefaultListener 创建默认的监听器
func newDefaultListener() Listener {
	return &defaultListener{}
}
