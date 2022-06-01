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

import "time"

// session 会话结构体
type session struct {
	Id         string                 `json:"id"`          // 当前session的id
	CreateTime int64                  `json:"create_time"` // 当前session的创建时间
	DataMap    map[string]interface{} `json:"data_map"`
	SignList   []*sign                `json:"sign_list"` // 存储的登录信息
	logic      *Logic
}

// NewSession 根据sessionId创建session
func NewSession(id string, logic *Logic) *session {
	ss := &session{
		Id:         id,
		CreateTime: time.Now().UnixMilli(),
		DataMap:    make(map[string]interface{}),
		SignList:   make([]*sign, 0),
		logic:      logic,
	}
	return ss
}

// GetSign 获取指定tokenValue的签名
func (s *session) GetSign(tokenValue string) *sign {
	for i := 0; i < len(s.SignList); i++ {
		if s.SignList[i].Value == tokenValue {
			return s.SignList[i]
		}
	}
	return nil
}

// GetTimeout 获取此session的剩余存活时间
func (s *session) GetTimeout() int64 {
	return s.logic.storage.TTl(s.Id)
}

// RemoveSign 移除指定签名
func (s *session) RemoveSign(tokenValue string) {
	for i := 0; i < len(s.SignList); i++ {
		if s.SignList[i].Value == tokenValue {
			s.SignList = append(s.SignList[:i], s.SignList[i+1:]...)
		}
	}
	s.Update()
}

// AddTokenSign 在user-session上记录签名
func (s *session) AddTokenSign(tokenValue string, device string) {
	// 如果已经存在于列表中，则无需再次添加
	for _, s2 := range s.SignList {
		if s2.Value == tokenValue {
			return
		}
	}
	// 添加并更新
	s.SignList = append(s.SignList, &sign{
		Value:  tokenValue,
		Device: device,
	})
	s.Update()
}

// LogoutByTokenSignCountToZero 当session上面的token签名为0时，注销用户级session
func (s *session) LogoutByTokenSignCountToZero() {
	if len(s.SignList) == 0 {
		s.LoginOut()
	}
}

// LoginOut 注销会话
func (s *session) LoginOut() {
	// 删除storage
	s.logic.storage.Del(s.Id)
}

// Update 更新持久库
func (s *session) Update() {
	s.logic.storage.UpdateObject(s.Id, s)
}

// UpdateMinTimeout 修改此Session的最小剩余存活时间 (只有在Session的过期时间低于指定的minTimeout时才会进行修改)
//形参:
//	minTimeout – 过期时间 (单位: 秒)
func (s *session) UpdateMinTimeout(minTimeout int64) {
	if s.GetTimeout() < minTimeout {
		s.logic.storage.UpdateObjectTTl(s.Id, minTimeout)
	}
}
