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

package session

import (
	"encoding/json"
	"github.com/go-ceres/go-ceres/auth/token/entity"
	"github.com/go-ceres/go-ceres/auth/token/manager"
	"github.com/go-ceres/go-ceres/utils/structure"
	"time"
)

const (
	// RoleList 在 Session 上存储角色时建议使用的key
	RoleList = "ROLE_LIST"
	// PermissionList 在 Session 上存储权限时建议使用的key
	PermissionList = "PERMISSION_LIST"
)

// Session user-session
type Session struct {
	Id         string                 `json:"id"`          // 会话id
	CreateTime int64                  `json:"create_time"` // 会话创建的时间
	DataMap    map[string]interface{} `json:"data_map"`    // 会话所存储的数据
	SignList   []*entity.TokenSign    `json:"sign_list"`   // 存储的登录信息
}

// NewSession 构建session
func NewSession(id string) *Session {
	ss := &Session{
		Id:         id,
		CreateTime: time.Now().UnixMilli(),
		DataMap:    map[string]interface{}{},
		SignList:   []*entity.TokenSign{},
	}
	return ss
}

// GetTokenSignList 获取sign列表
func (s *Session) GetTokenSignList() []*entity.TokenSign {
	var sess = new(Session)
	_ = structure.Copy(sess, s)
	return sess.SignList
}

// Update 更新session会话（从持久库）
func (s *Session) Update() {
	str, _ := json.Marshal(s)
	manager.Storage().Update(s.Id, string(str))
}

// AddTokenSign 添加签名
func (s *Session) AddTokenSign(tokenSign *entity.TokenSign) {
	for _, sign := range s.SignList {
		if tokenSign.Value == sign.Value {
			return
		}
	}
	s.SignList = append(s.SignList, tokenSign)
	s.Update()
}

// UpdateTimeout 设置当前session的过期时间
func (s *Session) UpdateTimeout(timeout int64) {
}
