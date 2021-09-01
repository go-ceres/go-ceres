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

package file

import (
	"github.com/go-ceres/go-ceres/logger/writer"
	"github.com/mitchellh/mapstructure"
)

func init() {
	writer.Register("file", new(Build))
}

type Build struct {
}

// Build 构造器
func (b *Build) Build(conf interface{}) writer.Writer {
	config := NewDefaultRotateConfig()
	err := mapstructure.Decode(conf, config)
	if err != nil {
		return nil
	}
	return config.Build()
}
