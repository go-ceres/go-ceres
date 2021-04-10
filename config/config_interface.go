//   Copyright 2021 Go-Ceres
//   Author https://github.com/go-ceres/go-ceres
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package config

import (
	"time"
)

type ChangeFunc func(v Values)

type Config interface {
	LoadSource(source Source) error
	Load(content []byte, format string) error
	Get(path string) Value
	Root() Values
	Set(Path string, data interface{}) error
	OnChange(change ChangeFunc)
	UnWatch()
	Watch()
	Write() error
}
type Values interface {
	Get(path string) Value
	Delete(path string)
	Set(path string, val interface{})
	Bytes() []byte
	Map() map[string]interface{}
	Scan(v interface{}) error
	String() string
}
type Value interface {
	IsEmpty() bool
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	StringSlice(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
	Bytes() []byte
}
