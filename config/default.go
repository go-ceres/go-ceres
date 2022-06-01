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
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

type config struct {
	mu        sync.RWMutex
	store     *JSONValues
	source    Source
	onChanges []ChangeFunc
}

func (c *config) Get(path string) Value {
	return c.store.Get(path)
}

func (c *config) Root() Values {
	return c.store
}

func (c *config) Set(path string, data interface{}) error {
	JSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var v interface{}
	err = json.Unmarshal(JSON, &v)
	if err != nil {
		return err
	}
	m, ok := v.(map[string]interface{})
	if ok {
		err := c.traverse(m, []string{path}, func(p string, value interface{}) error {
			c.store.Set(p, value)
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		c.store.Set(path, data)
	}
	c.notifyChange()
	return nil
}

func (c *config) OnChange(change ChangeFunc) {
	c.onChanges = append(c.onChanges, change)
}

func (c *config) LoadSource(source Source) error {
	c.source = source
	dataSet, redErr := source.Read()
	if redErr != nil {
		return redErr
	}
	if err := c.Load(dataSet.Data, dataSet.Format); err != nil {
		return err
	}
	return nil
}

func (c *config) Load(content []byte, unmarshal string) error {
	decode, ok := Unmarshals[unmarshal]
	if !ok {
		return errors.New("load err: No unmarshal method found")
	}
	dataMap := make(map[string]interface{})
	if err := decode(content, &dataMap); err != nil {
		return err
	}
	if data, err := json.Marshal(dataMap); err != nil {
		return err
	} else {
		c.store = NewJSONValues(data)
	}
	return nil
}

func (c *config) Watch() {
	c.source.Watch()
	go func() {
		for range c.source.IsChanged() {
			if dataSet, err := c.source.Read(); err == nil {
				_ = c.Load(dataSet.Data, dataSet.Format)
				c.notifyChange()
			}
		}
	}()
}

func (c *config) UnWatch() {
	c.source.UnWatch()
}

func (c *config) Write() error {
	return nil
}

func (c *config) traverse(m map[string]interface{}, paths []string, callback func(path string, value interface{}) error) error {
	for k, v := range m {
		val, ok := v.(map[string]interface{})
		if !ok {
			err := callback(strings.Join(append(paths, k), "."), v)
			if err != nil {
				return err
			}
			continue
		}
		err := c.traverse(val, append(paths, k), callback)
		if err != nil {
			return err
		}
	}
	return nil
}

// 通知监听，配置文件已经被改过了
func (c *config) notifyChange() {
	for _, change := range c.onChanges {
		change(c.store)
	}
}

// NewConfig 创建一个新的config管理器
func NewConfig() Config {
	conf := config{
		store: NewJSONValues([]byte("")),
	}
	return &conf
}
