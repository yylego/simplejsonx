package simplejsonx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/simplejsonx"
)

// TestLoad_ValidJSON tests Load function with valid JSON data
// TestLoad_ValidJSON 测试 Load 函数解析有效 JSON 数据
func TestLoad_ValidJSON(t *testing.T) {
	data := []byte(`{
		"key": "value"
	}`)

	object, err := simplejsonx.Load(data)
	require.NoError(t, err)

	res, err := object.Get("key").String()
	require.NoError(t, err)
	require.Equal(t, "value", res)
}

// TestLoad_InvalidJSON tests Load function with invalid JSON data
// TestLoad_InvalidJSON 测试 Load 函数解析无效 JSON 数据
func TestLoad_InvalidJSON(t *testing.T) {
	data := []byte(`invalid message`)

	object, err := simplejsonx.Load(data)
	require.Error(t, err)
	require.NotNil(t, object)
}

// TestLoad_EmptyInput tests Load function with empty input data
// TestLoad_EmptyInput 测试 Load 函数处理空输入数据
func TestLoad_EmptyInput(t *testing.T) {
	data := []byte(``)

	object, err := simplejsonx.Load(data)
	require.Error(t, err)
	require.NotNil(t, object)
}

// TestWrap tests Wrap function with map data structure
// TestWrap 测试 Wrap 函数包装 map 数据结构
func TestWrap(t *testing.T) {
	value := map[string]interface{}{
		"key": "abc",
	}

	object := simplejsonx.Wrap(value)
	require.NotNil(t, object)

	res, err := object.Get("key").String()
	require.NoError(t, err)
	require.Equal(t, "abc", res)
}

// TestWrap_WithPrimitiveValue tests Wrap function with primitive values
// TestWrap_WithPrimitiveValue 测试 Wrap 函数包装基础类型值
func TestWrap_WithPrimitiveValue(t *testing.T) {
	value := 88

	object := simplejsonx.Wrap(value)
	require.NotNil(t, object)

	res, err := object.Int()
	require.NoError(t, err)
	require.Equal(t, 88, res)
}

// TestWrap_InvalidNone tests Wrap function with nil input value
// TestWrap_InvalidNone 测试 Wrap 函数处理 nil 输入值
func TestWrap_InvalidNone(t *testing.T) {
	object := simplejsonx.Wrap(nil)
	require.NotNil(t, object)

	t.Run("case-1", func(t *testing.T) {
		res, err := object.Int()
		require.Error(t, err)
		require.Equal(t, 0, res)
	})

	t.Run("case-2", func(t *testing.T) {
		res, err := object.Get("abc").String()
		require.Error(t, err)
		require.Equal(t, "", res)
	})
}

// TestList tests List function converting JSON elements to object slice
// TestList 测试 List 函数将 JSON 元素转换成对象切片
func TestList(t *testing.T) {
	data := []byte(`{
		"infos": [
			{"code":1, "name":"a"},
			{"code":2, "name":"b"},
			{"code":3, "name":"c"}
		],
		"ranks": ["x", "y", "z"]
	}`)
	object, err := simplejsonx.Load(data)
	require.NoError(t, err)
	infos, err := object.Get("infos").Array()
	require.NoError(t, err)
	elements := simplejsonx.List(infos)

	var resMap = map[string]int{}
	for _, elem := range elements {
		code, err := elem.Get("code").Int()
		require.NoError(t, err)
		name, err := elem.Get("name").String()
		require.NoError(t, err)
		t.Log(code, name)
		resMap[name] = code
	}
	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, resMap)
}
