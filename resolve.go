// Package simplejsonx: Type-safe JSON parsing engine with Go generics support
// Built on github.com/bitly/go-simplejson with enhanced compile-time type-safe checking
// Provides Extract, Inspect, Resolve, Inquire, Attempt, Harvest, Explore functions
// Each function offers a different level of strictness for field validation
// Supports primitives, lists, maps, and nested JSON objects with zero-cost type checking
// Features code-generated error handling variants in must/omit/soft sub-packages
//
// simplejsonx: 类型安全的 JSON 解析引擎，支持 Go 泛型
// 基于 github.com/bitly/go-simplejson 构建，具有增强的编译时类型安全性
// 提供 Extract、Inspect、Resolve、Inquire、Attempt、Harvest、Explore 函数
// 每个函数提供不同的严格程度，精确控制字段校验
// 支持基础类型、数组、映射和嵌套 JSON 对象，具有零成本类型检查
// 具有代码生成的错误处理变体，包含 must/omit/soft 子包
package simplejsonx

import (
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/pkg/errors"
	"github.com/yylego/simplejsonx/internal/utils"
)

// Extract retrieves and parses the value at the specified key into the target type
// Returns errors when the key is missing, when type conversion fails
// Supports automatic type conversion based on generic type parameters
//
// Extract 检索指定键的值并解析成目标类型
// 当键缺失或类型转换失败时返回错误
// 支持基于泛型类型参数的自动类型转换
func Extract[T any](object *simplejson.Json, key string) (T, error) {
	if object == nil {
		return utils.Zero[T](), errors.New("Extract: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), errors.New("Extract: param key is blank")
	}
	return Resolve[T](object.Get(key))
}

// Solicit retrieves and parses the value at the specified key with non-zero validation
// Returns errors when the key is missing, when type conversion fails, or when the value is zero
// Goes beyond Extract: rejects zero values (0, "", false) in addition to missing/null/wrong-type
//
// Solicit 检索指定键的值并解析成目标类型，同时校验非零值
// 当键缺失、类型转换失败、或值为零值时返回错误
// 比 Extract 更严格：除了缺失/null/类型不对之外，零值（0、""、false）也会报错
func Solicit[T comparable](object *simplejson.Json, key string) (T, error) {
	res, err := Extract[T](object, key)
	if err != nil {
		return res, errors.WithMessage(err, "Solicit: unable to extract value from JSON object")
	}
	if res == utils.Zero[T]() {
		return res, errors.Errorf("Solicit: key %q value is zero, expected non-zero", key)
	}
	return res, nil
}

// Inspect retrieves and parses the value at the specified key when present
// Returns zero value without errors when key is absent, errors on conversion failures
// Provides lenient access patterns that accept fields in JSON data
//
// Inspect 检索并解析指定键的值（当键存在时）
// 当键不存在时返回零值且不报错，仅在转换失败时返回错误
// 提供宽松的访问模式，允许 JSON 结构中存在可选字段
func Inspect[T any](object *simplejson.Json, key string) (T, error) {
	if object == nil {
		return utils.Zero[T](), errors.New("Inspect: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), errors.New("Inspect: param key is blank")
	}
	value, exist := object.CheckGet(key)
	if !exist {
		return utils.Zero[T](), nil
	}
	return Resolve[T](value)
}

// Resolve extracts and converts JSON value into the target type
// Supports comprehensive type conversion via simplejson.Json methods
// Handles primitives (int, int64, float64, string, uint64, bool)
// Handles arrays ([]string, []interface{}, []*simplejson.Json)
// Handles complex types (map[string]interface{}, []byte, *simplejson.Json)
//
// Resolve 提取 JSON 值并转换成目标类型
// 支持使用 simplejson.Json 方法进行全面的类型转换
// 处理基础类型（int、int64、float64、string、uint64、bool）
// 处理数组类型（[]string、[]interface{}、[]*simplejson.Json）
// 处理复杂类型（map[string]interface{}、[]byte、*simplejson.Json）
func Resolve[T any](object *simplejson.Json) (T, error) {
	if object == nil {
		return utils.Zero[T](), errors.New("Resolve: param object is nil")
	}
	if object.Interface() == nil {
		return utils.Zero[T](), errors.New("Resolve: JSON value is null, cannot convert to target type")
	}
	switch zero := utils.Zero[T](); any(zero).(type) {
	case int:
		res, err := object.Int()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to int")
		}
		return any(res).(T), nil
	case int64:
		res, err := object.Int64()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to int64")
		}
		return any(res).(T), nil
	case float64:
		res, err := object.Float64()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to float64")
		}
		return any(res).(T), nil
	case string:
		res, err := object.String()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to string")
		}
		return any(res).(T), nil
	case uint64:
		res, err := object.Uint64()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to uint64")
		}
		return any(res).(T), nil
	case bool:
		res, err := object.Bool()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to bool")
		}
		return any(res).(T), nil
	case []string:
		res, err := object.StringArray()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to []string")
		}
		return any(res).(T), nil
	case []interface{}:
		res, err := object.Array()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to []interface{}")
		}
		return any(res).(T), nil
	case map[string]interface{}:
		res, err := object.Map()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to map[string]interface{}")
		}
		return any(res).(T), nil
	case []byte:
		res, err := object.Bytes()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to []byte")
		}
		return any(res).(T), nil
	case *simplejson.Json:
		return any(object).(T), nil
	case []*simplejson.Json:
		elements, err := object.Array()
		if err != nil {
			return zero, errors.WithMessage(err, "unable to resolve JSON value to []interface{}")
		}
		return any(List(elements)).(T), nil
	default:
		return zero, errors.Errorf("unsupported generic type: %T. unable to resolve JSON value.", zero)
	}
}

// GetList retrieves JSON list at the specified key as simplejson.Json objects slice
// Converts each list element into distinct simplejson.Json objects
// Returns errors when key is missing, when value is not a list
//
// GetList 检索指定键的 JSON 数组并转换成 simplejson.Json 对象切片
// 将每个数组元素转换成独立的 simplejson.Json 包装对象
// 当键缺失或值不是数组时返回错误
func GetList(object *simplejson.Json, key string) (objects []*simplejson.Json, err error) {
	if object == nil {
		return nil, errors.New("GetList: param object is nil")
	}
	if key == "" {
		return nil, errors.New("GetList: param key is blank")
	}
	elements, err := object.Get(key).Array()
	if err != nil {
		return objects, errors.WithMessage(err, "GetList: unable to get JSON elements as list")
	}
	return List(elements), nil
}

// Inquire queries JSON object at the specified key with tri-state result pattern
// Returns parsed value, existence boolean, and possible conversion errors
// Distinguishes between missing keys and type conversion failures
//
// Inquire 查询 JSON 对象中指定键的值，使用三态结果模式
// 返回解析后的值、存在性布尔值和可能的转换错误
// 区分缺失的键和类型转换失败的情况
func Inquire[T any](object *simplejson.Json, key string) (T, bool, error) {
	if object == nil {
		return utils.Zero[T](), false, errors.New("Inquire: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), false, errors.New("Inquire: param key is blank")
	}
	value, exist := object.CheckGet(key)
	if !exist {
		return utils.Zero[T](), false, nil
	}
	res, err := Resolve[T](value)
	if err != nil {
		return utils.Zero[T](), false, errors.WithMessage(err, "Inquire: unable to resolve JSON value to target type")
	}
	return res, true, nil
}

// Attempt retrieves and parses value at the specified key with success indication
// Returns parsed value and boolean indicating complete operation success
// Without notification returns zero value and false on failures without exposing errors
//
// Attempt 检索并解析指定键的值，附带成功指示器
// 返回解析后的值和表示整体操作成功的布尔值
// 在任何失败情况下静默返回零值和 false，不暴露错误
func Attempt[T any](object *simplejson.Json, key string) (T, bool) {
	if object == nil {
		return utils.Zero[T](), false
	}
	if key == "" {
		return utils.Zero[T](), false
	}
	value, exist := object.CheckGet(key)
	if !exist {
		return utils.Zero[T](), false
	}
	res, err := Resolve[T](value)
	if err != nil {
		return utils.Zero[T](), false
	}
	return res, true
}

// Explore navigates nested JSON structure via dot-separated path notation
// Traverses multiple levels using path like "user.profile.name"
// Returns parsed value, existence boolean, and possible conversion errors
//
// Explore 通过点分隔路径表示法导航嵌套 JSON 结构
// 使用类似 "user.profile.name" 的路径遍历多层级结构
// 返回解析后的值、存在性布尔值和可能的转换错误
func Explore[T any](object *simplejson.Json, path string) (T, bool, error) {
	if object == nil {
		return utils.Zero[T](), false, errors.New("Explore: param object is nil")
	}
	if path == "" {
		return utils.Zero[T](), false, errors.New("Explore: param path is blank")
	}
	value := object
	var exist bool
	for _, key := range strings.Split(path, ".") {
		value, exist = value.CheckGet(key)
		if !exist {
			return utils.Zero[T](), false, nil
		}
	}
	res, err := Resolve[T](value)
	if err != nil {
		return utils.Zero[T](), false, errors.WithMessage(err, "Explore: unable to resolve JSON value to target type")
	}
	return res, true, nil
}

// HarvestOption configures mixed-type field parsing behavior
// Enables flexible handling when JSON fields may contain numbers, strings, or null
// Used with Harvest function to customize extraction and conversion rules
//
// HarvestOption 配置混合类型字段的解析行为
// 当 JSON 字段可能包含数字、字符串或 null 时启用灵活处理
// 与 Harvest 函数配合使用，用于自定义提取和转换规则
type HarvestOption[T any] struct {
	Nilable bool                    // Accept null values without error // 接受 null 值不报错
	Strconv func(string) (T, error) // Custom string conversion function // 自定义字符串转换函数
}

// Harvest extracts mixed-type fields with configurable null and string handling
// Designed to handle JSON fields that may contain numbers, strings like "-", or null
// Uses HarvestOption to customize behavior: Nilable allows null, Strconv converts strings
//
// Harvest 提取混合类型字段，支持可配置的 null 和字符串处理
// 设计用于处理可能包含数字、字符串（如 "-"）或 null 的 JSON 字段
// 使用 HarvestOption 自定义行为：Nilable 允许 null，Strconv 转换字符串
func Harvest[T any](object *simplejson.Json, key string, option *HarvestOption[T]) (T, error) {
	if object == nil {
		if option.Nilable {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.New("Harvest: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), errors.New("Harvest: param key is blank")
	}
	value, exist := object.CheckGet(key)
	if !exist {
		return utils.Zero[T](), errors.Errorf("Harvest: key %q not found in JSON object", key)
	}
	if value == nil || value.Interface() == nil {
		if option.Nilable {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.Errorf("Harvest: key %q value is null in JSON object", key)
	}
	if option.Strconv != nil {
		stringValue, err := value.String()
		if err != nil {
			res, err := Resolve[T](value)
			if err != nil {
				return utils.Zero[T](), errors.WithMessage(err, "Harvest: unable to resolve JSON value to target type")
			}
			return res, nil
		}
		res, err := option.Strconv(stringValue)
		if err != nil {
			return utils.Zero[T](), errors.WithMessage(err, "Harvest: unable to convert string value using custom Strconv function")
		}
		return res, nil
	}
	res, err := Resolve[T](value)
	if err != nil {
		return utils.Zero[T](), errors.WithMessage(err, "Harvest: unable to resolve JSON value to target type")
	}
	return res, nil
}
