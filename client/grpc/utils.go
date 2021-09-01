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

package grpc

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
)

func split2(s, sep string) (string, string, bool) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return "", "", false
	}
	return spl[0], spl[1], true
}

// parseTarget 解析sting
func parseTarget(target string) (ret resolver.Target) {
	var ok bool
	ret.Scheme, ret.Endpoint, ok = split2(target, "://")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	ret.Authority, ret.Endpoint, ok = split2(ret.Endpoint, "/")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	return ret
}

// stringTarget 转换string
func stringTarget(ret resolver.Target) string {
	if ret.Scheme == "" {
		return ret.Endpoint
	}
	return fmt.Sprintf("%s://%s/%s", ret.Scheme, ret.Authority, ret.Endpoint)
}
