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

package ceres

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ceres/go-ceres/cmd"
	"strings"
)

// defaultBanner 默认的banner信息
const defaultBanner = `
  ____ _____ ____  _____ ____
 / ___| ____|  _ \| ____/ ___|
| |   |  _| | |_) |  _| \___ \
| |___| |___|  _ <| |___ ___) |
 \____|_____|_| \_|_____|____/
`

// version 版本信息
var version = `go-ceres@` + cmd.GetCeresVersion() + `    http://go-ceres.com/`

// customBanner 自定义 Banner 字符串
var customBanner = ""

// SetBanner 设置自定义 Banner 字符串
func SetBanner(banner string) {
	customBanner = banner
}

// printBanner 打印 Banner 到控制台
func printBanner(banner string) {

	// 确保 Banner 前面有空行
	if banner[0] != '\n' {
		fmt.Println()
	}

	maxLength := 0
	for _, s := range strings.Split(banner, "\n") {
		color.Cyan(s)
		if len(s) > maxLength {
			maxLength = len(s)
		}
	}

	// 确保 Banner 后面有空行
	if banner[len(banner)-1] != '\n' {
		fmt.Println()
	}

	var padding []byte
	if n := (maxLength - len(version)) / 2; n > 0 {
		padding = make([]byte, n)
		for i := range padding {
			padding[i] = ' '
		}
	}
	color.Cyan(string(padding) + version + "\n")
}
