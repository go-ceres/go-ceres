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

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type jwtClaims struct {
	LoginId              string `json:"login_id"`   // 登录的用户
	LogicType            string `json:"logic_type"` // 登录逻辑
	Device               string `json:"device"`     // 登录设备
	jwt.RegisteredClaims        // 注意!这是jwt-go的v4版本新增的，原先是jwt.StandardClaims
}

var jwtSecret = []byte("lqowicnamzuwsdawegasdwweghcanm") // 定义secret，后面会用到

// makeToken 这里传入的是手机号，因为我项目登陆用的是手机号和密码
func makeToken(loginId string, logicType string, device string) (string, error) {
	claim := jwtClaims{
		LoginId:   loginId,
		LogicType: logicType,
		Device:    device,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	return token.SignedString(jwtSecret)
}
