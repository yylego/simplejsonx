package simplejsonx

import (
	"bytes"
	"encoding/json"

	"github.com/bitly/go-simplejson"
	"github.com/pkg/errors"
)

// Load creates simplejson.Json instance from raw JSON bytes
// Parses byte data into structured JSON object representation
//
// Load 从原始 JSON 字节创建 simplejson.Json 实例
// 将字节数据解析成结构化的 JSON 对象表示
func Load(data []byte) (object *simplejson.Json, err error) {
	object = simplejson.New()
	// Validate the complete document first, then invoke Json.UnmarshalJSON, which keeps json.Number precision.
	// 先校验整个文档，再调用 Json.UnmarshalJSON，保留其 json.Number 精度。
	err = json.Unmarshal(data, object)
	if err != nil {
		return simplejson.New(), errors.WithMessage(err, "Load: unable to parse JSON data")
	}
	return object, nil
}

// LoadOption names the text that wraps the JSON, such as the "callback(" and ");" around a JSONP response
// The code that calls it says what to trim, so Load does not have to guess the shape
//
// LoadOption 指明包在 JSON 外面的文本，如 JSONP 响应外面的 "callback(" 和 ");"
// 由调用方说清要去掉什么，Load 不用去猜外壳的样子
type LoadOption struct {
	TrimPrefix string // The prefix the data must have, trimmed once before parsing // 数据必须带的前缀，解析前去掉一次
	TrimSuffix string // The suffix the data must have, trimmed once before parsing // 数据必须带的后缀，解析前去掉一次
}

// LoadWithOption trims the surrounding whitespace and the text the option names, then parses the rest as JSON
//
// LoadWithOption 去掉首尾空白和选项指明的文本，再把剩下的按 JSON 解析
func LoadWithOption(data []byte, option *LoadOption) (object *simplejson.Json, err error) {
	if option == nil {
		option = &LoadOption{}
	}
	data = bytes.TrimSpace(data)
	if option.TrimPrefix != "" {
		rest, found := bytes.CutPrefix(data, []byte(option.TrimPrefix))
		if !found {
			return simplejson.New(), errors.Errorf("LoadWithOption: data does not have the prefix %q", option.TrimPrefix)
		}
		data = rest
	}
	if option.TrimSuffix != "" {
		rest, found := bytes.CutSuffix(data, []byte(option.TrimSuffix))
		if !found {
			return simplejson.New(), errors.Errorf("LoadWithOption: data does not have the suffix %q", option.TrimSuffix)
		}
		data = rest
	}
	return Load(data)
}

// Wrap creates simplejson.Json instance wrapping given Go values
// Encapsulates Go data structures into JSON object representation
// Enables uniform JSON operations on native Go values
//
// Wrap 创建包装任意 Go 值的 simplejson.Json 实例
// 将任何 Go 数据结构封装成 JSON 对象表示
// 在原生 Go 值上实现统一的 JSON 操作
func Wrap(value any) (object *simplejson.Json) {
	object = simplejson.New()
	object.SetPath([]string{}, value)
	return object
}

// List converts slice of any values into simplejson.Json objects slice
// Wraps each element as distinct JSON object instances
// Provides batch conversion enabling uniform processing of heterogeneous data
//
// List 将 any 值切片转换成 simplejson.Json 对象切片
// 独立包装每个元素，创建独立的 JSON 对象实例
// 提供批量转换，实现对异构数据的统一处理
func List(elements []any) (objects []*simplejson.Json) {
	objects = make([]*simplejson.Json, 0, len(elements))
	for _, elem := range elements {
		objects = append(objects, Wrap(elem))
	}
	return objects
}
