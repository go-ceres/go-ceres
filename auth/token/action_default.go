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
	s2 "github.com/go-ceres/go-ceres/utils/strings"
	"github.com/google/uuid"
	"strings"
)

// defaultAction 默认的token操作
type defaultAction struct {
	logic *Logic
}

func (d *defaultAction) createToken(loginId string, logicType string, device string) string {
	style := d.logic.config.TokenStyle
	switch style {
	case TOKEN_STYLE_UUID:
		return uuid.New().String()
	case TOKEN_STYLE_SIMPLE_UUID:
		return strings.ReplaceAll(uuid.New().String(), "_", "")
	case TOKEN_STYLE_RANDOM_32:
		return s2.RandStr(32)
	case TOKEN_STYLE_RANDOM_64:
		return s2.RandStr(64)
	case TOKEN_STYLE_JWT:
		token, _ := makeToken(loginId, logicType, device)
		return token
	default:
		return uuid.New().String()
	}
}

func newDefaultAction(logic *Logic) tokenAction {
	return &defaultAction{
		logic: logic,
	}
}
