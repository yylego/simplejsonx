// Package simplejsonx: Type-safe JSON parsing engine with Go generics support
// Built on github.com/bitly/go-simplejson with typed results and runtime validation
// Provides Extract, Solicit, Inspect, Resolve, Inquire, Attempt, Harvest, Numeric, Explore functions
// Each function offers a different degree of strictness in field validation
// Supports primitives, lists, maps, and nested JSON objects
// Features code-generated error handling variants in must/omit/soft sub-packages
//
// simplejsonx: 类型安全的 JSON 解析引擎，支持 Go 泛型
// 基于 github.com/bitly/go-simplejson 构建，提供确定类型的返回值与运行时校验
// 提供 Extract、Solicit、Inspect、Resolve、Inquire、Attempt、Harvest、Numeric、Explore 函数
// 每个函数提供不同的严格程度，精确控制字段校验
// 支持基础类型、数组、映射和嵌套 JSON 对象
// 具有代码生成的错误处理变体，包含 must/omit/soft 子包
package simplejsonx

import (
	"slices"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/pkg/errors"
	"github.com/yylego/simplejsonx/internal/utils"
)

// Extract retrieves and parses the value at the specified key into the target type
// The key must exist and the value must not be null; use Inspect when the key can be absent
// Supports automatic type conversion based on generic type parameters
//
// Extract 检索指定键的值并解析成目标类型
// 键必须存在、值不能是 null；键可以缺时用 Inspect
// 支持基于泛型类型参数的自动类型转换
func Extract[T any](object *simplejson.Json, key string) (T, error) {
	if object == nil {
		return utils.Zero[T](), errors.New("Extract: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), errors.New("Extract: param key is blank")
	}
	if _, err := object.Map(); err != nil {
		return utils.Zero[T](), errors.Wrapf(err, "Extract: key %q requires a JSON object", key)
	}
	// CheckGet tells a missing key apart from a null value, so the error names the cause as it is
	// CheckGet 能分清缺键和值为 null，报错时写对原因
	value, exist := object.CheckGet(key)
	if !exist {
		return utils.Zero[T](), errors.Errorf("Extract: key %q not found in JSON object", key)
	}
	return Resolve[T](value)
}

// Solicit retrieves and parses the value at the specified key with non-zero validation
// Goes beyond Extract: rejects zero values (0, "", false) in addition to missing/null/wrong-type
//
// Solicit 检索指定键的值并解析成目标类型，同时校验非零值
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
//
// Inspect 检索并解析指定键的值（当键存在时）
// 当键不存在时返回零值且不报错，仅在转换失败时返回错误
func Inspect[T any](object *simplejson.Json, key string) (T, error) {
	if object == nil {
		return utils.Zero[T](), errors.New("Inspect: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), errors.New("Inspect: param key is blank")
	}
	if _, err := object.Map(); err != nil {
		return utils.Zero[T](), errors.Wrapf(err, "Inspect: key %q requires a JSON object", key)
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
// Handles arrays ([]string, []any, []*simplejson.Json)
// Handles complex types (map[string]any, []byte, *simplejson.Json)
//
// Resolve 提取 JSON 值并转换成目标类型
// 支持使用 simplejson.Json 方法进行全面的类型转换
// 处理基础类型（int、int64、float64、string、uint64、bool）
// 处理数组类型（[]string、[]any、[]*simplejson.Json）
// 处理复杂类型（map[string]any、[]byte、*simplejson.Json）
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
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to int")
		}
		return any(res).(T), nil
	case int64:
		res, err := object.Int64()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to int64")
		}
		return any(res).(T), nil
	case float64:
		res, err := object.Float64()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to float64")
		}
		return any(res).(T), nil
	case string:
		res, err := object.String()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to string")
		}
		return any(res).(T), nil
	case uint64:
		res, err := object.Uint64()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to uint64")
		}
		return any(res).(T), nil
	case bool:
		res, err := object.Bool()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to bool")
		}
		return any(res).(T), nil
	case []string:
		res, err := object.StringArray()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to []string")
		}
		return any(res).(T), nil
	case []any:
		res, err := object.Array()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to []any")
		}
		return any(res).(T), nil
	case map[string]any:
		res, err := object.Map()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to map[string]any")
		}
		return any(res).(T), nil
	case []byte:
		res, err := object.Bytes()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to []byte")
		}
		return any(res).(T), nil
	case *simplejson.Json:
		return any(object).(T), nil
	case []*simplejson.Json:
		elements, err := object.Array()
		if err != nil {
			return zero, errors.WithMessage(err, "Resolve: unable to resolve JSON value to []*simplejson.Json")
		}
		return any(List(elements)).(T), nil
	default:
		return zero, errors.Errorf("Resolve: unsupported generic type %T, unable to resolve JSON value", zero)
	}
}

// GetList retrieves JSON list at the specified key as simplejson.Json objects slice
// The same as Extract[[]*simplejson.Json], with no type argument to write
//
// GetList 检索指定键的 JSON 数组并转换成 simplejson.Json 对象切片
// 与 Extract[[]*simplejson.Json] 相同，只是不用写类型参数
func GetList(object *simplejson.Json, key string) (objects []*simplejson.Json, err error) {
	objects, err = Extract[[]*simplejson.Json](object, key)
	if err != nil {
		return nil, errors.WithMessage(err, "GetList: unable to get JSON elements as list")
	}
	return objects, nil
}

// Inquire queries JSON object at the specified key with tri-state result pattern
// Returns (value, true, nil) on success and (zero, false, nil) on absent keys
// Errors return (zero, false, err); check err before interpreting the flag
// Distinguishes between missing keys and type conversion failures
//
// Inquire 查询 JSON 对象中指定键的值，使用三态结果模式
// 成功返回 (值, true, nil)，缺键返回 (零值, false, nil)
// 出错返回 (零值, false, err)；应先检查 err，再判断布尔值
// 区分缺失的键和类型转换失败的情况
func Inquire[T any](object *simplejson.Json, key string) (T, bool, error) {
	if object == nil {
		return utils.Zero[T](), false, errors.New("Inquire: param object is nil")
	}
	if key == "" {
		return utils.Zero[T](), false, errors.New("Inquire: param key is blank")
	}
	if _, err := object.Map(); err != nil {
		return utils.Zero[T](), false, errors.Wrapf(err, "Inquire: key %q requires a JSON object", key)
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
// Returns zero value and false on each failure, without the cause
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
// Uses the same result convention as Inquire; each path segment requires a JSON object
//
// Explore 通过点分隔路径表示法导航嵌套 JSON 结构
// 使用类似 "user.profile.name" 的路径遍历多层级结构
// 返回值约定与 Inquire 相同；每段路径的取值对象必须是 JSON 对象
func Explore[T any](object *simplejson.Json, path string) (T, bool, error) {
	if object == nil {
		return utils.Zero[T](), false, errors.New("Explore: param object is nil")
	}
	if path == "" {
		return utils.Zero[T](), false, errors.New("Explore: param path is blank")
	}
	value := object
	var exist bool
	for key := range strings.SplitSeq(path, ".") {
		if _, err := value.Map(); err != nil {
			return utils.Zero[T](), false, errors.Wrapf(err, "Explore: path %q segment %q requires a JSON object", path, key)
		}
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

// HarvestOption configures mixed-type field parsing settings
// Enables flexible handling when JSON fields may contain numbers, strings, or null
// Used with Harvest function to customize extraction and conversion rules
//
// HarvestOption 配置混合类型字段的解析行为
// 当 JSON 字段可能包含数字、字符串或 null 时启用灵活处理
// 与 Harvest 函数配合使用，用于自定义提取和转换规则
type HarvestOption[T any] struct {
	AbsentOk bool                    // Accept absent keys as zero values // 缺键时返回零值，不影响 null 与类型校验
	Nullable bool                    // Accept null values and nil objects as zero values // null 值或 nil 对象返回零值
	Strconv  func(string) (T, error) // Custom string conversion function // 自定义字符串转换函数
}

// Harvest extracts mixed-type fields with configurable null and string handling
// Designed to handle JSON fields that may contain numbers, strings like "-", or null
// Uses HarvestOption to customize settings: Nullable accepts null, AbsentOk accepts absent keys, Strconv converts strings
//
// Harvest 提取混合类型字段，支持可配置的 null 和字符串处理
// 设计用于处理可能包含数字、字符串（如 "-"）或 null 的 JSON 字段
// 使用 HarvestOption 自定义行为：Nullable 允许 null，AbsentOk 允许缺键，Strconv 转换字符串
func Harvest[T any](object *simplejson.Json, key string, option *HarvestOption[T]) (T, error) {
	if key == "" {
		return utils.Zero[T](), errors.New("Harvest: param key is blank")
	}
	if option == nil {
		option = &HarvestOption[T]{}
	}
	if object == nil {
		if option.Nullable {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.New("Harvest: param object is nil")
	}
	if _, err := object.Map(); err != nil {
		return utils.Zero[T](), errors.Wrapf(err, "Harvest: key %q requires a JSON object", key)
	}
	value, exist := object.CheckGet(key)
	if !exist {
		if option.AbsentOk {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.Errorf("Harvest: key %q not found in JSON object", key)
	}
	if value.Interface() == nil {
		if option.Nullable {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.Errorf("Harvest: key %q value is null in JSON object", key)
	}
	// A JSON string goes through the custom Strconv, the rest of the values go through Resolve
	// JSON 字符串交给自定义的 Strconv，其余的值交给 Resolve
	if text, isText := value.Interface().(string); isText && option.Strconv != nil {
		res, err := option.Strconv(text)
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

// NumericOption configures how Numeric treats absent keys, null values and strings meaning no value
// The zero option is the strict one: the key must exist, null is not accepted, and each string must be numeric
//
// NumericOption 配置 Numeric 怎么对待缺键、null 值和表示没有值的字符串
// 零值选项最严：键必须存在，不接受 null，字符串必须是数字
type NumericOption struct {
	AbsentOk bool     // Accept absent keys as zero values // 缺键时返回零值
	Nullable bool     // Accept null values and nil objects as zero values, as HarvestOption does // null 值或 nil 对象返回零值，与 HarvestOption 一致
	Blanks   []string // Strings meaning no value, such as "" and "N/A", read as zero values // 表示没有值的字符串，如 "" 和 "N/A"，读作零值
}

// Numeric reads numeric values that the JSON gives as JSON numbers / numeric strings, such as 9.5 and "9.5"
// Option decides on absent keys, null values and the strings meaning no value
//
// Numeric 读取一个数字，JSON 里它可以是数字也可以是数字字符串，如 9.5 或 "9.5"
// 选项决定缺键、null 值以及表示没有值的字符串怎么处理
func Numeric[T int | int64 | uint64 | float64](object *simplejson.Json, key string, option *NumericOption) (T, error) {
	if key == "" {
		return utils.Zero[T](), errors.New("Numeric: param key is blank")
	}
	if option == nil {
		option = &NumericOption{}
	}
	if object == nil {
		if option.Nullable {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.New("Numeric: param object is nil")
	}
	if _, err := object.Map(); err != nil {
		return utils.Zero[T](), errors.Wrapf(err, "Numeric: key %q requires a JSON object", key)
	}
	value, exist := object.CheckGet(key)
	if !exist {
		if option.AbsentOk {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.Errorf("Numeric: key %q not found in JSON object", key)
	}
	if value.Interface() == nil {
		if option.Nullable {
			return utils.Zero[T](), nil
		}
		return utils.Zero[T](), errors.Errorf("Numeric: key %q value is null in JSON object", key)
	}
	if text, isText := value.Interface().(string); isText {
		if slices.Contains(option.Blanks, text) {
			return utils.Zero[T](), nil
		}
		res, err := Strconv[T](value)
		if err != nil {
			return utils.Zero[T](), errors.WithMessagef(err, "Numeric: key %q string value %q is not numeric", key, text)
		}
		return res, nil
	}
	res, err := Resolve[T](value)
	if err != nil {
		return utils.Zero[T](), errors.WithMessagef(err, "Numeric: key %q value is not numeric", key)
	}
	return res, nil
}
