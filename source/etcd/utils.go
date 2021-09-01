package etcd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/BurntSushi/toml"
	"github.com/go-ceres/go-ceres/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

var (
	Marshals = map[string]Marshal{
		"json": json.Marshal,
		"xml":  xml.Marshal,
		"yaml": yaml.Marshal,
		"yml":  yaml.Marshal,
		"toml": func(v interface{}) ([]byte, error) {
			b := bytes.NewBuffer(nil)
			defer b.Reset()
			err := toml.NewEncoder(b).Encode(v)
			if err != nil {
				return nil, err
			}
			return b.Bytes(), nil
		},
	}
	Unmarshals = map[string]Unmarshal{
		"json": json.Unmarshal,
		"xml":  xml.Unmarshal,
		"yaml": yaml.Unmarshal,
		"yml":  yaml.Unmarshal,
		"toml": toml.Unmarshal,
	}
)

type Marshal func(v interface{}) ([]byte, error)
type Unmarshal func(data []byte, v interface{}) error

// makeMapData 把[]byte根据配置转为map
func makeMapData(kv []*mvccpb.KeyValue, trimPrefix string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, value := range kv {
		data = modifyMapData(trimPrefix, data, value)
	}
	return data
}

// 获取解码方法
func getUnmarshal(path string) (Unmarshal, string) {
	UnmarshalStr := strings.TrimPrefix(filepath.Ext(path), ".")
	if fn, ok := Unmarshals[UnmarshalStr]; ok {
		return fn, UnmarshalStr
	}
	return func(data []byte, v interface{}) error {
		v = string(data)
		return nil
	}, "txt"
}

// modifyMapData 调整数据
func modifyMapData(trimPrefix string, data map[string]interface{}, kv *mvccpb.KeyValue) map[string]interface{} {
	// 删除前缀，例如：/ceres/config/etcd/default.json,操作后的为：etcd/default.json
	key := strings.TrimPrefix(strings.TrimPrefix(string(kv.Key), trimPrefix), "/")
	fn, str := getUnmarshal(key)
	// 判断是否有后缀,如果没有后缀，默认使用json解码
	if !strings.HasSuffix(key, "."+str) {
		fn = Unmarshals["json"]
	}
	// 去掉后缀
	key = strings.TrimSuffix(key, "."+str)
	// 分割为["etcd","default"]
	keys := strings.Split(key, "/")
	// 序列化数据
	var value interface{}
	err := fn(kv.Value, &value)
	if err != nil {
		logger.Error("解析etcd错误，错误信息为：", err)
	}
	if len(keys) > 0 && len(keys) == 1 {
		v, ok := value.(map[string]interface{})
		if ok {
			data = v
		} else {
			data = make(map[string]interface{})
			if str == "txt" {
				data[keys[0]] = value
			}
		}
		return data
	}

	tempData := data
	for i, k := range keys {
		// 先判断该key是否已经存在值
		kData, ok := data[k].(map[string]interface{})
		if !ok {
			kData = make(map[string]interface{})
			tempData[k] = kData
		}
		// 如果是最后一个key，则设置数据
		if len(keys)-1 == i {
			tempData[k] = value
		} else {
			tempData = kData
		}
	}
	return data
}
