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

import "testing"

func TestAddRouter(t *testing.T) {
	router := NewRouter()
	routes := [...]string{
		"/",
		"/cmd/:tool/",
		"/cmd/:tool/:sub",
		"/cmd/whoami",
		"/cmd/whoami/root",
		"/cmd/whoami/root/",
		"/src/*filepath",
		"/search/",
		"/search/:query",
		"/search/gin-gonic",
		"/search/google",
		"/user_:name",
		"/user_:name/about",
		"/files/:dir/*filepath",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/:user/public",
		"/info/:user/project/:project",
		"/info/:user/project/golang",
		"/aa/*xx",
		"/ab/*xx",
		"/:cc",
		"/c1/:dd/e",
		"/c1/:dd/e1",
		"/:cc/cc",
		"/:cc/:dd/ee",
		"/:cc/:dd/:ee/ff",
		"/:cc/:dd/:ee/:ff/gg",
		"/:cc/:dd/:ee/:ff/:gg/hh",
		"/get/test/abc/",
		"/get/:param/abc/",
		"/something/:paramname/thirdthing",
		"/something/secondthing/test",
		"/get/abc",
		"/get/:param",
		"/get/abc/123abc",
		"/get/abc/:param",
		"/get/abc/123abc/xxx8",
		"/get/abc/123abc/:param",
		"/get/abc/123abc/xxx8/1234",
		"/get/abc/123abc/xxx8/:param",
		"/get/abc/123abc/xxx8/1234/ffas",
		"/get/abc/123abc/xxx8/1234/:param",
		"/get/abc/123abc/xxx8/1234/kkdd/12c",
		"/get/abc/123abc/xxx8/1234/kkdd/:param",
		"/get/abc/:param/test",
		"/get/abc/123abd/:param",
		"/get/abc/123abddd/:param",
		"/get/abc/123/:param",
		"/get/abc/123abg/:param",
		"/get/abc/123abf/:param",
		"/get/abc/123abfff/:param",
	}

	for _, route := range routes {
		router.AddRouter(route, []string{route})
	}
	router.GetRouter("/cmd/test/", true)
}
