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
	"encoding/json"
	"github.com/coocood/freecache"
)

// 默认的持久化存储
type defaultStorage struct {
	cache *freecache.Cache
}

func (d *defaultStorage) Has(key string) bool {
	get, err := d.cache.Get([]byte(key))
	if err != nil {
		return false
	}
	return len(get) != 0
}

func (d *defaultStorage) SetObject(key string, value interface{}, timeout int64) bool {
	marshal, err := json.Marshal(value)
	if err != nil {
		return false
	}
	err = d.cache.Set([]byte(key), marshal, int(timeout))
	if err != nil {
		return false
	}
	return true
}

func (d *defaultStorage) GetObject(key string, obj interface{}) bool {
	get, err := d.cache.Get([]byte(key))
	if err != nil {
		return false
	}
	err = json.Unmarshal(get, obj)
	if err != nil {
		return false
	}
	return true
}

func (d *defaultStorage) Get(key string, def ...string) string {
	res := func(v ...string) string {
		if len(v) > 0 {
			return v[0]
		} else {
			return ""
		}
	}(def...)
	get, err := d.cache.Get([]byte(key))
	if err != nil {
		return res
	}
	return string(get)
}

func (d *defaultStorage) Set(key string, value string, timeout int64) bool {
	err := d.cache.Set([]byte(key), []byte(value), int(timeout))
	if err != nil {
		return false
	}
	return true
}

func (d *defaultStorage) Update(key string, value string) {
	ttl, err := d.cache.TTL([]byte(key))
	if err != nil {
		return
	}
	err = d.cache.Set([]byte(key), []byte(value), int(ttl))
	if err != nil {
		return
	}
}

func (d *defaultStorage) UpdateObject(key string, value interface{}) bool {
	ttl, err := d.cache.TTL([]byte(key))
	if err != nil {
		return false
	}
	marshal, err := json.Marshal(value)
	if err != nil {
		return false
	}
	err = d.cache.Set([]byte(key), marshal, int(ttl))
	if err != nil {
		return false
	}
	return true
}

func (d *defaultStorage) UpdateObjectTTl(key string, timeout int64) {
	err := d.cache.Touch([]byte(key), int(timeout))
	if err != nil {
		return
	}
}

func (d *defaultStorage) Del(key string) bool {
	return d.cache.Del([]byte(key))
}

func (d *defaultStorage) TTl(key string) int64 {
	ttl, err := d.cache.TTL([]byte(key))
	if err != nil {
		return 0
	}
	return int64(ttl)
}

func (d *defaultStorage) Clear() bool {
	d.cache.Clear()
	return true
}

// NewDefaultStorage 创建一个默认的存储器
func NewDefaultStorage(size int) Storage {
	if size == 0 {
		size = 100 * 1024 * 1024
	}
	cache := freecache.NewCache(size)
	// 返回默认存储器
	return &defaultStorage{
		cache: cache,
	}
}
