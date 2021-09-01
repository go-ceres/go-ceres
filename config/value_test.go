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
	"log"
	"testing"
)

func TestJSONValue_Exists(t *testing.T) {
	datas := &JSONValues{}
	dataMap := make(map[string]interface{})
	if data, err := json.Marshal(dataMap); err != nil {

	} else {
		datas = NewJSONValues(data)
	}
	b := datas.Get("ceshi").IsEmpty()
	log.Print(b)
}
