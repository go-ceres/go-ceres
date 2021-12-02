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

package token

import (
	"encoding/json"
	"fmt"
	"github.com/go-ceres/go-ceres/auth/token/entity"
	"github.com/go-ceres/go-ceres/auth/token/errors"
	"github.com/go-ceres/go-ceres/auth/token/manager"
	"github.com/go-ceres/go-ceres/auth/token/session"
	"github.com/go-ceres/go-ceres/auth/token/stp"
	"github.com/go-ceres/go-ceres/cache"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

const (
	TokenConnectorChat = " " // 前缀与token的连接字符

	NeverExpire    = -1 // 标注一个key永不过期
	NotValueExpire = -2 // 对一个不存在的key返回此值
)

// Token 逻辑处理
type Token struct {
	loginType string  // 当前逻辑处理的类型
	conf      *Config // 配置信息
}

// WithStore 设置持久化存储
func (t *Token) WithStore(store cache.Cache) *Token {
	manager.SetStorage(store)
	return t
}

// WithStpInterface 设置权限认证接口
func (t *Token) WithStpInterface(stp stp.Interface) *Token {
	manager.SetStp(stp)
	return t
}

// =====================获取token相关=====================

// GetTokenName 获取token名称
func (t *Token) GetTokenName() string {
	return t.conf.TokenName
}

// CreateTokenValue 创建TokenValue
func (t *Token) CreateTokenValue(loginId interface{}, device string, timeout int64) string {

	return t.createToken(loginId, device)
}

// GetTokenInfo 根据token Value名获取token信息
func (t *Token) GetTokenInfo(tokenValue string) *entity.TokenInfo {
	info := new(entity.TokenInfo)
	info.TokenName = t.GetTokenName()
	info.TokenValue = tokenValue
	info.IsLogin = t.IsLogin(tokenValue)
	info.LoginType = t.loginType
	info.TokenTimeout = t.GetTokenTimeout(tokenValue)
	info.SessionTimeout = t.GetSessionTimeoutByLoginId(t.GetLoginIdNotHandle(tokenValue))
	info.TokenSessionTimeout = t.GetTokenSessionTimeoutByTokenValue(tokenValue)
	info.TokenActivityTimeout = t.GetTokenActivityTimeoutByToken(tokenValue)
	info.LoginDevice = t.GetLoginDevice(tokenValue)
	return info
}

// ========================登录相关========================

// Login 登录，成功返回tokenValue，错误
func (t *Token) Login(id interface{}, opts ...LoginOption) (string, error) {
	if id == nil {
		return "", errors.NewError("账号id不能为空")
	}
	// 检查账号是否被封禁
	if t.IsDisable(id) {
		return "", errors.NewDisableLoginError(t.loginType, id, t.GetDisableTime(id))
	}

	options := DefaultOption(t.conf.Timeout)
	for _, opt := range opts {
		opt(options)
	}

	// 2.生成token
	var tokenValue string
	// 如果允许并发登录
	if t.conf.IsConcurrent {
		// 如果配置为共享token，则从session签名中获取token
		if t.conf.IsShare {
			tokenValue = t.GetTokenValueByLoginId(id, options.Device())
		}
	} else {
		// 该分支为不允许并发登录,则强制下线其他token
	}
	if tokenValue == "" {
		tokenValue = t.CreateTokenValue(id, options.Device(), options.Timeout())
	}

	// 3.获取 user-session
	userSession := t.GetSessionByLoginId(id, true)
	userSession.UpdateTimeout(options.Timeout())

	// 4.在user-session上记录签名
	userSession.AddTokenSign(&entity.TokenSign{Value: tokenValue, Divice: options.Device()})
	// 5.存储映射关系
	t.SaveTokenToMapping(tokenValue, id, options.Timeout())
	// 6.给指定token设置最后活动时间
	t.setLastActivityToNow(tokenValue)
	return tokenValue, nil
}

// NewToken 构建token
func NewToken(c *Config, logicType string) *Token {
	return &Token{
		conf:      c,
		loginType: logicType,
	}
}

// GetLogicType 获取当前逻辑处理的类型
func (t *Token) GetLogicType() string {
	return t.loginType
}

// SetLogicType 设置当前逻辑处理类型
func (t *Token) SetLogicType(logicType string) {
	t.loginType = logicType
}

// GetConfig 获取配置信息
func (t *Token) GetConfig() *Config {
	return t.conf
}

// IsLogin 判断该token是否登录
func (t *Token) IsLogin(tokenValue string) bool {
	return len(t.GetLoginIdNotHandle(tokenValue)) != 0
}

// ========================User-Session 相关=======================

// GetSessionBySessionId 根据sessionId获取用户session
func (t *Token) GetSessionBySessionId(sessionId string, isCreate bool) *session.Session {
	sess := new(session.Session)
	err := json.Unmarshal([]byte(manager.Storage().Get(sessionId)), sess)
	if err != nil && isCreate {
		sess = t.CreateSession(sessionId)
		marshal, _ := json.Marshal(sess)
		manager.Storage().Set(sessionId, string(marshal), t.conf.Timeout)
	}
	return sess
}

// GetSessionByLoginId 获取用户session
func (t *Token) GetSessionByLoginId(loginId interface{}, isCreate bool) *session.Session {
	return t.GetSessionBySessionId(t.splicingKeySession(loginId), isCreate)
}

// CreateSession 创建一个session
func (t *Token) CreateSession(sessionId string) *session.Session {
	return session.NewSession(sessionId)
}

// ========================Token-Session相关=======================

// GetTokenSessionByToken 获取指定token-session，如果session没有创建，可传递isCreate创建并返回
func (t *Token) GetTokenSessionByToken(tokenValue string, isCreate bool) *session.Session {
	return t.GetSessionBySessionId(t.splicingKeyTokenSession(tokenValue), isCreate)
}

// DeleteTokenSession 删除token-session
func (t *Token) DeleteTokenSession(tokenValue string) {
	manager.Storage().Del(t.splicingKeyTokenValue(tokenValue))
}

// ========================账号封禁=======================

// Disable 封禁指定账号
func (t *Token) Disable(id interface{}, disableTime int64) {
	manager.Storage().Set(t.splicingKeyDisable(id), errors.DisableLoginErrorValue, disableTime)
}

// IsDisable 查询账号是否被封禁
func (t *Token) IsDisable(id interface{}) bool {
	return manager.Storage().Get(t.splicingKeyDisable(id)) != ""
}

// GetDisableTime 获取封禁时间
func (t *Token) GetDisableTime(id interface{}) int64 {
	return manager.Storage().TTl(t.splicingKeyDisable(id))
}

// ========================根据loginid反查=======================

// GetTokenValueByLoginId 根据指定的账号和设备反查TokenValue
func (t *Token) GetTokenValueByLoginId(loginId interface{}, device string) string {
	tokenValueList := t.GetTokenValueListByLoginId(loginId, device)
	if len(tokenValueList) == 0 {
		return ""
	}
	return tokenValueList[len(tokenValueList)-1]
}

// GetTokenValueListByLoginId 根据指定id查
func (t *Token) GetTokenValueListByLoginId(loginId interface{}, device string) []string {
	userSession := t.GetSessionByLoginId(loginId, false)
	if userSession == nil {
		return []string{}
	}
	tokenSignList := userSession.GetTokenSignList()
	var ret []string
	for _, sign := range tokenSignList {
		if device == "" || strings.EqualFold(sign.Value, device) {
			ret = append(ret, sign.Value)
		}
	}
	return ret
}

// GetLoginDevice 根据token反查登录设备
func (t *Token) GetLoginDevice(tokenValue string) string {
	if len(tokenValue) == 0 {
		return ""
	}
	if !t.IsLogin(tokenValue) {
		return ""
	}
	sess := t.GetSessionByLoginId(t.GetLoginIdNotHandle(tokenValue), false)
	if sess == nil {
		return ""
	}

	tokenSignList := sess.GetTokenSignList()
	for _, sign := range tokenSignList {
		if strings.EqualFold(tokenValue, sign.Value) {
			return sign.Divice
		}
	}
	return ""
}

// ========================查询相关========================

// GetLoginId 获取登录用户id，如果没有登录则返回错误
func (t *Token) GetLoginId(tokenValue string) (interface{}, error) {
	if tokenValue == "" {
		return nil, errors.NewNotLoginInstance(errors.NotLoginErrorLoginType(t.loginType), errors.NotLoginErrorErrType(errors.NotToken))
	}
	loginId := t.GetLoginIdNotHandle(tokenValue)
	// 为空
	if len(loginId) == 0 {
		return nil, errors.NewNotLoginInstance(errors.NotLoginErrorLoginType(t.loginType), errors.NotLoginErrorErrType(errors.InvalidToken), errors.NotLoginErrorToken(tokenValue))
	}
	// 如果是已经过期
	switch loginId {
	case errors.TokenTimeout:
		return nil, errors.NewNotLoginInstance(errors.NotLoginErrorLoginType(t.loginType), errors.NotLoginErrorErrType(errors.TokenTimeout), errors.NotLoginErrorToken(tokenValue))
	case errors.BeReplaced:
		return nil, errors.NewNotLoginInstance(errors.NotLoginErrorLoginType(t.loginType), errors.NotLoginErrorErrType(errors.BeReplaced), errors.NotLoginErrorToken(tokenValue))
	case errors.KickOut:
		return nil, errors.NewNotLoginInstance(errors.NotLoginErrorLoginType(t.loginType), errors.NotLoginErrorErrType(errors.KickOut), errors.NotLoginErrorToken(tokenValue))
	}
	// 检查是否过期（ActivityTimeout）

	// 如果配置了自动续签, 则: 更新[最后操作时间]
	if t.conf.AutoRenew {
		//updateLastActivityToNow(tokenValue)
	}

	return loginId, nil
}

// GetLoginIdNotHandle 获取指定token的用户id
func (t *Token) GetLoginIdNotHandle(tokenValue string) string {
	return manager.Storage().Get(t.splicingKeyTokenValue(tokenValue))
}

// ========================返回对应的key========================

// splicingKeyTokenValue 获取token存储key
func (t *Token) splicingKeyTokenValue(tokenValue string) string {
	return t.conf.TokenName + ":" + t.loginType + ":token:" + tokenValue
}

// splicingKeySession 获取session存储key
func (t *Token) splicingKeySession(loginId interface{}) string {
	return t.conf.TokenName + ":" + t.loginType + ":session:" + fmt.Sprintf("%v", loginId)
}

// splicingKeyTokenSession 获取tokensession存储key
func (t *Token) splicingKeyTokenSession(tokenValue string) string {
	return t.conf.TokenName + ":" + t.loginType + ":token-session:" + tokenValue
}

// splicingKeyLastActivtyTime 获取token的最后操作时间存储key
func (t *Token) splicingKeyLastActivtyTime(tokenValue string) string {
	return t.conf.TokenName + ":" + t.loginType + ":last-activity:" + tokenValue
}

// splicingKeyDisable 获取账号封禁key
func (t *Token) splicingKeyDisable(loginId interface{}) string {
	return t.conf.TokenName + ":" + t.loginType + ":disable:" + fmt.Sprintf("%v", loginId)
}

// ========================临时有效期操作========================
func (t *Token) setLastActivityToNow(tokenValue string) {
	// 如果token为空 或者设置了永不过期，则直接返回
	if &tokenValue == nil || t.conf.ActivityTimeout == NeverExpire {
		return
	}
	manager.Storage().Set(t.splicingKeyLastActivtyTime(tokenValue), strconv.Itoa(int(time.Now().UnixMilli())), t.conf.Timeout)
}

// ========================过期相关查询========================

// GetTokenTimeout 获取token过期时间
func (t *Token) GetTokenTimeout(tokenValue string) int64 {
	return manager.Storage().TTl(t.splicingKeyTokenValue(tokenValue))
}

// GetSessionTimeoutByLoginId 根据用户id获取session的过期时间
func (t *Token) GetSessionTimeoutByLoginId(loginId interface{}) int64 {
	return manager.Storage().TTl(t.splicingKeySession(loginId))
}

// GetTokenSessionTimeoutByTokenValue 获取指定token的过期剩余时间
func (t *Token) GetTokenSessionTimeoutByTokenValue(tokenValue string) int64 {
	return manager.Storage().TTl(t.splicingKeyTokenSession(tokenValue))
}

// GetTokenActivityTimeoutByToken 获取临时过期时间
func (t *Token) GetTokenActivityTimeoutByToken(tokenValue string) int64 {
	if len(tokenValue) == 0 {
		return NotValueExpire
	}
	if t.conf.ActivityTimeout == NotValueExpire {
		return NeverExpire
	}
	lastActivtyTimeKey := t.splicingKeyLastActivtyTime(tokenValue)
	lastActivtyTimeStr := manager.Storage().Get(lastActivtyTimeKey)
	// 如果查询不到
	if len(lastActivtyTimeStr) == 0 {
		return NotValueExpire
	}
	// 计算时间差
	lastActivtyTime, _ := strconv.Atoi(lastActivtyTimeStr)
	apartSecond := time.Now().UnixMilli() - int64(lastActivtyTime)/1000
	timeout := t.conf.ActivityTimeout - apartSecond
	// 如果小于0，则代表已经过期
	if timeout < 0 {
		return NotValueExpire
	}
	return timeout
}

// ========================其他操作========================

func (t *Token) SaveTokenToMapping(tokenValue string, loginId interface{}, timeout int64) {
	manager.Storage().Set(t.splicingKeyTokenValue(tokenValue), fmt.Sprintf("%v", loginId), timeout)
}

// createToken 生成token
func (t *Token) createToken(id interface{}, device string) string {
	switch t.conf.TokenStyle {
	case TokenStyleUuid:
		return uuid.NewString()
	case TokenStyleSimpleUuid:
		return strings.ReplaceAll(uuid.NewString(), "-", "")
	case TOKEN_STYLE_RANDOM_32:
	}
	return uuid.NewString()
}
