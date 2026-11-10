package simplejsonx_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/require"
	"github.com/yylego/simplejsonx"
)

// TestExtract_Int tests Extract function with int type conversion
// TestExtract_Int 测试 Extract 函数的 int 类型转换
func TestExtract_Int(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"age": 18
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Extract[int](object, "age")
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 18, res)
}

// TestExtract_Mismatch tests Extract function when key is missing
// TestExtract_Mismatch 测试 Extract 函数在键缺失时的行为
func TestExtract_Mismatch(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": "yylego"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Extract[int](object, "age")
	require.Error(t, err)
	t.Log(err)
	t.Log(res)
	require.Equal(t, 0, res)
}

// TestExtract_NotObject tests Extract on a JSON list, which has no keys to look up
// TestExtract_NotObject 测试 Extract 遇到 JSON 数组，数组里没有键可查
func TestExtract_NotObject(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`[1, 2, 3]`))
	require.NoError(t, err)

	_, err = simplejsonx.Extract[int](object, "age")
	require.Error(t, err)
}

// TestSolicit_NonZeroInt tests Solicit function returns non-zero int value
// TestSolicit_NonZeroInt 测试 Solicit 函数返回非零 int 值
func TestSolicit_NonZeroInt(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"age": 18
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[int](object, "age")
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 18, res)
}

// TestSolicit_ZeroInt tests Solicit function rejects zero int value
// TestSolicit_ZeroInt 测试 Solicit 函数拒绝零值 int
func TestSolicit_ZeroInt(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"count": 0
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[int](object, "count")
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, 0, res)
}

// TestSolicit_EmptyString tests Solicit function rejects empty string value
// TestSolicit_EmptyString 测试 Solicit 函数拒绝空字符串值
func TestSolicit_EmptyString(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": ""
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[string](object, "name")
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, "", res)
}

// TestSolicit_NonZeroString tests Solicit function returns non-empty string value
// TestSolicit_NonZeroString ��试 Solicit 函数返回非空字符串值
func TestSolicit_NonZeroString(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": "yylego"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[string](object, "name")
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, "yylego", res)
}

// TestSolicit_FalseBool tests Solicit function rejects false boolean value
// TestSolicit_FalseBool 测试 Solicit 函数拒绝 false 布尔值
func TestSolicit_FalseBool(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"active": false
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[bool](object, "active")
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, false, res)
}

// TestSolicit_MissingKey tests Solicit function rejects missing key
// TestSolicit_MissingKey ����� Solicit 函数拒绝缺失的键
func TestSolicit_MissingKey(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": "yylego"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[int](object, "age")
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, 0, res)
}

// TestSolicit_NullValue tests Solicit function rejects null value
// TestSolicit_NullValue 测试 Solicit 函数拒绝 null 值
func TestSolicit_NullValue(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": null
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Solicit[string](object, "name")
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, "", res)
}

// TestInspect tests Inspect function with existing key
// TestInspect 测试 Inspect 函数处理存在的键
func TestInspect(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"age": 18
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Inspect[int](object, "age")
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 18, res)
}

// TestInspect_Mismatch tests Inspect function returning zero value when key is missing
// TestInspect_Mismatch 测试 Inspect 函数在键缺失时返回零值
func TestInspect_Mismatch(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": "yylego"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Inspect[int](object, "age")
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 0, res)
}

// TestResolve_Int tests Resolve function with int type
// TestResolve_Int 测试 Resolve 函数的 int 类型转换
func TestResolve_Int(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"age": 18
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Resolve[int](object.Get("age"))
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 18, res)
}

// TestResolve_Mismatch tests Resolve function with type mismatch scenarios
// TestResolve_Mismatch 测试 Resolve 函数的类型不匹配场景
func TestResolve_Mismatch(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": "yylego",
		"age": 18
	}`))
	require.NoError(t, err)
	t.Run("case-1", func(t *testing.T) {
		res, err := simplejsonx.Resolve[int](object.Get("age"))
		require.NoError(t, err)
		require.Equal(t, 18, res)
	})
	t.Run("case-2", func(t *testing.T) {
		res, err := simplejsonx.Resolve[string](object.Get("age"))
		require.Error(t, err)
		t.Log(err)
		require.Equal(t, "", res) // zero value of string
	})
	t.Run("case-3", func(t *testing.T) {
		res, err := simplejsonx.Resolve[string](object.Get("name"))
		require.NoError(t, err)
		require.Equal(t, "yylego", res)
	})
	t.Run("case-4", func(t *testing.T) {
		res, err := simplejsonx.Resolve[int](object.Get("name"))
		require.Error(t, err)
		t.Log(err)
		require.Equal(t, 0, res) // zero value
	})
}

// TestResolve_Int64 tests Resolve function with int64 type
// TestResolve_Int64 测试 Resolve 函数的 int64 类型转换
func TestResolve_Int64(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"height": 175,
		"weight": 80
	}`))
	require.NoError(t, err)

	t.Run("case-1", func(t *testing.T) {
		res, err := simplejsonx.Resolve[int64](object.Get("height"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, int64(175), res)
	})
	t.Run("case-2", func(t *testing.T) {
		res, err := simplejsonx.Resolve[int64](object.Get("weight"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, int64(80), res)
	})
}

// TestResolve_Float64 tests Resolve function with float64 type
// TestResolve_Float64 测试 Resolve 函数的 float64 类型转换
func TestResolve_Float64(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"size": 18.5
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Resolve[float64](object.Get("size"))
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 18.5, res)
}

// TestResolve_String tests Resolve function with string type
// TestResolve_String 测试 Resolve 函数的 string 类型转换
func TestResolve_String(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"name": "yylego",
		"age": 18,
		"like": "rice"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Resolve[string](object.Get("like"))
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, "rice", res)
}

// TestResolve_Uint64 tests Resolve function with uint64 type
// TestResolve_Uint64 测试 Resolve 函数的 uint64 类型转换
func TestResolve_Uint64(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"endurance": 30,
		"persistence": 60
	}`))
	require.NoError(t, err)

	t.Run("case-1", func(t *testing.T) {
		res, err := simplejsonx.Resolve[uint64](object.Get("endurance"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, uint64(30), res)
	})
	t.Run("case-2", func(t *testing.T) {
		res, err := simplejsonx.Resolve[uint64](object.Get("persistence"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, uint64(60), res)
	})
}

// TestResolve_Bool tests Resolve function with bool type
// TestResolve_Bool 测试 Resolve 函数的 bool 类型转换
func TestResolve_Bool(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"is_tall": true,
		"is_rich": true,
		"is_cool": true
	}`))
	require.NoError(t, err)

	t.Run("case-1", func(t *testing.T) {
		res, err := simplejsonx.Resolve[bool](object.Get("is_tall"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, true, res)
	})
	t.Run("case-2", func(t *testing.T) {
		res, err := simplejsonx.Resolve[bool](object.Get("is_rich"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, true, res)
	})
	t.Run("case-3", func(t *testing.T) {
		res, err := simplejsonx.Resolve[bool](object.Get("is_cool"))
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, true, res)
	})
}

// TestResolve_StringArray tests Resolve function with []string type
// TestResolve_StringArray 测试 Resolve 函数的 []string 类型转换
func TestResolve_StringArray(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`[
		"a",
		"b",
		"c"
	]`))
	require.NoError(t, err)

	res, err := simplejsonx.Resolve[[]string](object)
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, []string{"a", "b", "c"}, res)
}

// TestResolve_Array tests Resolve function with []any type
// TestResolve_Array 测试 Resolve 函数的 []any 类型转换
func TestResolve_Array(t *testing.T) {
	t.Run("case-1", func(t *testing.T) {
		object, err := simplejsonx.Load([]byte(`[
			1,
			"two",
			3.3
		]`))
		require.NoError(t, err)

		res, err := simplejsonx.Resolve[[]any](object)
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, []any{json.Number("1"), "two", json.Number("3.3")}, res)
	})
	t.Run("case-2", func(t *testing.T) {
		var value any
		decoder := json.NewDecoder(bytes.NewBufferString(`[
			1,
			"two",
			3.3
		]`))
		decoder.UseNumber()
		require.NoError(t, decoder.Decode(&value))

		object := simplejsonx.Wrap(value)

		res, err := simplejsonx.Resolve[[]any](object)
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, []any{json.Number("1"), "two", json.Number("3.3")}, res)
	})
	t.Run("case-3", func(t *testing.T) {
		var a any
		require.NoError(t, json.Unmarshal([]byte(`[
			1,
			"two",
			3.3
		]`), &a))

		object := simplejsonx.Wrap(a)

		res, err := simplejsonx.Resolve[[]any](object)
		require.NoError(t, err)
		t.Log(res)
		require.Equal(t, []any{float64(1), "two", 3.3}, res)
	})
}

// TestResolve_Map tests Resolve function with map[string]any type
// TestResolve_Map 测试 Resolve 函数的 map[string]any 类型转换
func TestResolve_Map(t *testing.T) {
	object := simplejsonx.Wrap(map[string]any{"key": "value"})

	res, err := simplejsonx.Resolve[map[string]any](object)
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, map[string]any{"key": "value"}, res)
}

// TestResolve_Bytes tests Resolve function with []byte type
// TestResolve_Bytes 测试 Resolve 函数的 []byte 类型转换
func TestResolve_Bytes(t *testing.T) {
	a := map[string]any{
		"value": "abc",
	}
	data, err := json.Marshal(a)
	require.NoError(t, err)
	t.Log(string(data))

	object, err := simplejsonx.Load(data)
	require.NoError(t, err)

	res, err := simplejsonx.Extract[[]byte](object, "value")
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, "abc", string(res))
}

// TestGetList tests GetList function retrieving JSON list as simplejson.Json slice
// TestGetList 测试 GetList 函数获取 JSON 数组并转换成 simplejson.Json 切片
func TestGetList(t *testing.T) {
	jsonData := `{
		"key": [
			{"name": "item1"},
			{"name": "item2"},
			{"name": "item3"}
		]
	}`

	object, err := simplejsonx.Load([]byte(jsonData))
	require.NoError(t, err)

	t.Run("case-1", func(t *testing.T) {
		objects, err := simplejsonx.GetList(object, "key")
		require.NoError(t, err)
		require.Len(t, objects, 3)

		var names = make([]string, 0, len(objects))
		for _, item := range objects {
			name, err := simplejsonx.Extract[string](item, "name")
			require.NoError(t, err)
			names = append(names, name)
		}
		require.Equal(t, []string{"item1", "item2", "item3"}, names)
	})

	t.Run("case-2", func(t *testing.T) {
		objects, err := simplejsonx.GetList(object, "invalidKey")
		require.Error(t, err)
		require.Len(t, objects, 0)
	})

	t.Run("not-a-list", func(t *testing.T) {
		item, err := simplejsonx.Load([]byte(`{"key": "text"}`))
		require.NoError(t, err)

		objects, err := simplejsonx.GetList(item, "key")
		require.Error(t, err)
		require.Len(t, objects, 0)
	})
}

// TestInquire tests Inquire function with tri-state result pattern
// TestInquire 测试 Inquire 函数的三态结果模式
func TestInquire(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"age": 18
	}`))
	require.NoError(t, err)

	// Test success scenario
	// 测试成功场景
	t.Run("case-1", func(t *testing.T) {
		res, exists, err := simplejsonx.Inquire[int](object, "age")
		require.NoError(t, err)
		require.True(t, exists, "key should exist")
		t.Log(res)
		require.Equal(t, 18, res)
	})

	// Test missing key scenario
	// 测试键不存在场景
	t.Run("case-2", func(t *testing.T) {
		res, exists, err := simplejsonx.Inquire[int](object, "name")
		require.NoError(t, err)
		require.False(t, exists, "key should not exist")
		t.Log(res)
		require.Equal(t, 0, res) // zero value / 零值
	})

	// Test resolution failure scenario
	// 测试解析失败场景
	t.Run("case-3", func(t *testing.T) {
		res, exists, err := simplejsonx.Inquire[string](object, "age")
		require.Error(t, err, "should fail to resolve int to string")
		require.False(t, exists, "should return false on resolve failure")
		t.Log(res)
		require.Equal(t, "", res) // zero value / 零值
	})
}

// TestAttempt tests Attempt function with success indication
// TestAttempt 测试 Attempt 函数的成功指示器模式
func TestAttempt(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"age": 18
	}`))
	require.NoError(t, err)

	// Test success scenario
	// 测试成功场景
	t.Run("case-1", func(t *testing.T) {
		res, ok := simplejsonx.Attempt[int](object, "age")
		require.True(t, ok, "should succeed")
		t.Log(res)
		require.Equal(t, 18, res)
	})

	// Test missing key scenario
	// 测试键不存在场景
	t.Run("case-2", func(t *testing.T) {
		res, ok := simplejsonx.Attempt[int](object, "name")
		require.False(t, ok, "should fail due to missing key")
		t.Log(res)
		require.Equal(t, 0, res) // zero value / 零值
	})

	// Test resolution failure scenario
	// 测试解析失败场景
	t.Run("case-3", func(t *testing.T) {
		res, ok := simplejsonx.Attempt[string](object, "age")
		require.False(t, ok, "should fail to resolve int to string")
		t.Log(res)
		require.Equal(t, "", res) // zero value / 零值
	})
}

// TestExplore tests Explore function navigating nested JSON via dot-separated path
// TestExplore 测试 Explore 函数通过点分隔路径导航嵌套 JSON
func TestExplore(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"user": {
			"name": "Alice"
		}
	}`))
	require.NoError(t, err)

	// Test success scenario with nested path
	// 测试成功场景（嵌套路径）
	t.Run("case-1", func(t *testing.T) {
		res, exists, err := simplejsonx.Explore[string](object, "user.name")
		require.NoError(t, err)
		require.True(t, exists, "path should exist")
		t.Log(res)
		require.Equal(t, "Alice", res)
	})

	// Test non-existent path scenario
	// 测试路径不存在场景
	t.Run("case-2", func(t *testing.T) {
		res, exists, err := simplejsonx.Explore[string](object, "user.address")
		require.NoError(t, err)
		require.False(t, exists, "path should not exist")
		t.Log(res)
		require.Equal(t, "", res) // zero value / 零值
	})

	// Test resolution failure scenario
	// 测试解析失败场景
	t.Run("case-3", func(t *testing.T) {
		res, exists, err := simplejsonx.Explore[int](object, "user.name")
		require.Error(t, err, "should fail to resolve string to int")
		require.False(t, exists, "should return false on resolve failure")
		t.Log(res)
		require.Equal(t, 0, res) // zero value / 零值
	})

	// Test empty path scenario
	// 测试空路径场景
	t.Run("case-4", func(t *testing.T) {
		res, exists, err := simplejsonx.Explore[int](object, "")
		require.Error(t, err, "should fail due to empty path")
		require.False(t, exists, "should return false for empty path")
		t.Log(res)
		require.Equal(t, 0, res) // zero value / 零值
	})
}

// TestExplore_GetJson tests Explore function returning *simplejson.Json type
// TestExplore_GetJson 测试 Explore 函数返回 *simplejson.Json 类型
func TestExplore_GetJson(t *testing.T) {
	// Test success scenario with nested path
	// 测试成功场景（嵌套路径）
	object, err := simplejsonx.Load([]byte(`{
		"user": {
			"name": "Alice",
			"events": [
				{"eventName": "eat"},
				{"eventName": "sleep"}
			]
		}
	}`))
	require.NoError(t, err)
	res, ok, err := simplejsonx.Explore[*simplejson.Json](object, "user")
	require.NoError(t, err)
	require.True(t, ok)
	t.Log(res)
	elements, err := simplejsonx.GetList(res, "events")
	require.NoError(t, err)
	t.Log(len(elements))
	require.Len(t, elements, 2)
	for _, item := range elements {
		t.Log(item)
	}
}

// TestExplore_GetList tests Explore function returning []*simplejson.Json type
// TestExplore_GetList 测试 Explore 函数返回 []*simplejson.Json 类型
func TestExplore_GetList(t *testing.T) {
	// Test success scenario with nested path
	// 测试成功场景（嵌套路径）
	object, err := simplejsonx.Load([]byte(`{
		"user": {
			"name": "Alice",
			"events": [
				{"eventName": "eat"},
				{"eventName": "sleep"}
			]
		}
	}`))
	require.NoError(t, err)
	elements, ok, err := simplejsonx.Explore[[]*simplejson.Json](object, "user.events")
	require.NoError(t, err)
	require.True(t, ok)
	t.Log(len(elements))
	require.Len(t, elements, 2)
	for _, item := range elements {
		t.Log(item)
	}
}

// TestHarvest_Float64 tests Harvest function with standard float64 value
// TestHarvest_Float64 测试 Harvest 函数处理正常的 float64 值
func TestHarvest_Float64(t *testing.T) {
	// Test normal float64 value
	// 测试正常的 float64 值
	object, err := simplejsonx.Load([]byte(`{
		"price": 123.45
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Harvest[float64](object, "price", &simplejsonx.HarvestOption[float64]{})
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, 123.45, res)
}

// TestHarvest_NullWithNullable tests Harvest function accepting null with Nullable=true
// TestHarvest_NullWithNullable 测试 Harvest 函数在 Nullable=true 时接受 null 值
func TestHarvest_NullWithNullable(t *testing.T) {
	// Test null value with Nullable=true
	// 测试 null 值且 Nullable=true
	object, err := simplejsonx.Load([]byte(`{
		"price": null
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Harvest[float64](object, "price", &simplejsonx.HarvestOption[float64]{
		Nullable: true,
	})
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, float64(0), res)
}

// TestHarvest_NullWithoutNullable tests Harvest function rejecting null with Nullable=false
// TestHarvest_NullWithoutNullable 测试 Harvest 函数在 Nullable=false 时拒绝 null 值
func TestHarvest_NullWithoutNullable(t *testing.T) {
	// Test null value with Nullable=false
	// 测试 null 值且 Nullable=false
	object, err := simplejsonx.Load([]byte(`{
		"price": null
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Harvest[float64](object, "price", &simplejsonx.HarvestOption[float64]{
		Nullable: false,
	})
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, float64(0), res)
}

// TestHarvest_StringWithStrconv tests Harvest function with custom string conversion
// TestHarvest_StringWithStrconv 测试 Harvest 函数使用自定义字符串转换
func TestHarvest_StringWithStrconv(t *testing.T) {
	// Test string value "-" with Strconv
	// 测试字符串值 "-" 使用 Strconv 转换
	object, err := simplejsonx.Load([]byte(`{
		"price": "-"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Harvest[float64](object, "price", &simplejsonx.HarvestOption[float64]{
		Strconv: func(s string) (float64, error) {
			if s == "-" || s == "N/A" {
				return 0, nil
			}
			return 0, errors.New("unexpected string value")
		},
	})
	require.NoError(t, err)
	t.Log(res)
	require.Equal(t, float64(0), res)
}

// TestHarvest_MixedTypes tests Harvest function with mixed-type JSON fields
// TestHarvest_MixedTypes 测试 Harvest 函数处理混合类型的 JSON 字段
func TestHarvest_MixedTypes(t *testing.T) {
	// Test mixed types scenario: some float64, some string "-"
	// 测试混合类型场景：部分是 float64，部分是字符串 "-"
	jsonData := `{
		"items": [
			{"name": "A", "price": 100.5},
			{"name": "B", "price": "-"},
			{"name": "C", "price": null}
		]
	}`
	object, err := simplejsonx.Load([]byte(jsonData))
	require.NoError(t, err)

	items, err := simplejsonx.GetList(object, "items")
	require.NoError(t, err)
	require.Len(t, items, 3)

	option := &simplejsonx.HarvestOption[float64]{
		Nullable: true,
		Strconv: func(s string) (float64, error) {
			if s == "-" || s == "N/A" {
				return 0, nil
			}
			return 0, errors.New("unexpected string value")
		},
	}

	// Item A: float64 value
	// 项目 A：float64 值
	t.Run("case-1", func(t *testing.T) {
		res, err := simplejsonx.Harvest[float64](items[0], "price", option)
		require.NoError(t, err)
		require.Equal(t, 100.5, res)
	})

	// Item B: string "-"
	// 项目 B：字符串 "-"
	t.Run("case-2", func(t *testing.T) {
		res, err := simplejsonx.Harvest[float64](items[1], "price", option)
		require.NoError(t, err)
		require.Equal(t, float64(0), res)
	})

	// Item C: null value
	// 项目 C：null 值
	t.Run("case-3", func(t *testing.T) {
		res, err := simplejsonx.Harvest[float64](items[2], "price", option)
		require.NoError(t, err)
		require.Equal(t, float64(0), res)
	})
}

// TestHarvest_KeyNotFound tests Harvest function when key is missing
// TestHarvest_KeyNotFound 测试 Harvest 函数在键缺失时的行为
func TestHarvest_KeyNotFound(t *testing.T) {
	// Test key not found scenario
	// 测试键不存在场景
	object, err := simplejsonx.Load([]byte(`{
		"name": "test"
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Harvest[float64](object, "price", &simplejsonx.HarvestOption[float64]{})
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, float64(0), res)
}

// TestHarvest_EmptyKey tests Harvest function with blank key input
// TestHarvest_EmptyKey 测试 Harvest 函数处理空键参数
func TestHarvest_EmptyKey(t *testing.T) {
	// Test empty key scenario
	// 测试空键场景
	object, err := simplejsonx.Load([]byte(`{
		"price": 100
	}`))
	require.NoError(t, err)

	res, err := simplejsonx.Harvest[float64](object, "", &simplejsonx.HarvestOption[float64]{})
	require.Error(t, err)
	t.Log(err)
	require.Equal(t, float64(0), res)
	_, err = simplejsonx.Harvest[int](nil, "", &simplejsonx.HarvestOption[int]{Nullable: true})
	require.Error(t, err)
}

// TestInspect_NotObject tests Inspect on a JSON list, which has no keys to look up
// TestInspect_NotObject 测试 Inspect 遇到 JSON 数组，数组没有键可查
func TestInspect_NotObject(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`[]`))
	require.NoError(t, err)
	_, err = simplejsonx.Inspect[int](object, "amount")
	require.Error(t, err)
}

// TestInquire_NotObject tests Inquire on a JSON list, which has no keys to look up
// TestInquire_NotObject 测试 Inquire 遇到 JSON 数组，数组没有键可查
func TestInquire_NotObject(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`[]`))
	require.NoError(t, err)
	_, _, err = simplejsonx.Inquire[int](object, "amount")
	require.Error(t, err)
}

// TestHarvest_NotObject tests Harvest on a JSON list, which the options do not accept
// TestHarvest_NotObject 测试 Harvest 遇到 JSON 数组，选项放宽不了这种情况
func TestHarvest_NotObject(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`[]`))
	require.NoError(t, err)
	_, err = simplejsonx.Harvest[int](object, "amount", &simplejsonx.HarvestOption[int]{AbsentOk: true, Nullable: true})
	require.Error(t, err)
}

// TestExplore_NotObject tests Explore when a path segment lands on a list / numeric value in place of a JSON object
// TestExplore_NotObject 测试 Explore 的某段路径落在数组或数字上，而不是 JSON 对象
func TestExplore_NotObject(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`[]`))
	require.NoError(t, err)
	_, _, err = simplejsonx.Explore[int](object, "amount")
	require.Error(t, err)
	object, err = simplejsonx.Load([]byte(`{"item":42}`))
	require.NoError(t, err)
	_, _, err = simplejsonx.Explore[int](object, "item.amount")
	require.Error(t, err)
	object, err = simplejsonx.Load([]byte(`{"item":{}}`))
	require.NoError(t, err)
	_, ok, err := simplejsonx.Explore[int](object, "item.amount")
	require.NoError(t, err)
	require.False(t, ok)
}

// TestHarvest_Absent tests Harvest with absent keys, JSON nulls and values that cannot be parsed
// TestHarvest_Absent 测试 Harvest 遇到缺失的键、JSON null 和解析不了的值时的表现
func TestHarvest_Absent(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{"null":null,"zero":0,"text":"12","bad":"oops"}`))
	require.NoError(t, err)
	opt := &simplejsonx.HarvestOption[int]{AbsentOk: true, Strconv: strconv.Atoi}
	value, err := simplejsonx.Harvest(object, "absent", opt)
	require.NoError(t, err)
	require.Zero(t, value)
	_, err = simplejsonx.Harvest(object, "null", opt)
	require.Error(t, err)
	_, err = simplejsonx.Harvest[int](nil, "absent", opt)
	require.Error(t, err)
	opt.Nullable = true
	value, err = simplejsonx.Harvest(object, "null", opt)
	require.NoError(t, err)
	require.Zero(t, value)
	value, err = simplejsonx.Harvest(object, "text", opt)
	require.NoError(t, err)
	require.Equal(t, 12, value)
	_, err = simplejsonx.Harvest(object, "bad", opt)
	require.Error(t, err)
	_, err = simplejsonx.Harvest(object, "", opt)
	require.Error(t, err)
	_, err = simplejsonx.Harvest[int](object, "absent", nil)
	require.Error(t, err)
	value, err = simplejsonx.Harvest[int](object, "zero", nil)
	require.NoError(t, err)
	require.Zero(t, value)
	_, err = simplejsonx.Harvest[int](nil, "absent", nil)
	require.Error(t, err)
	value, err = simplejsonx.Harvest[int](nil, "absent", &simplejsonx.HarvestOption[int]{Nullable: true})
	require.NoError(t, err)
	require.Zero(t, value)
	_, err = simplejsonx.Harvest[int](object, "absent", &simplejsonx.HarvestOption[int]{Nullable: true})
	require.Error(t, err)
}

// TestNumeric tests Numeric with JSON numbers, numeric strings, null, strings meaning no value and bad values
// TestNumeric 测试 Numeric 遇到数字、数字字符串、null、占位字符串和坏值时的表现
func TestNumeric(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{
		"views": 1200,
		"price": "9.5",
		"rating": null,
		"discount": "",
		"weight": "N/A",
		"title": "hello"
	}`))
	require.NoError(t, err)

	t.Run("numeric-value-and-numeric-string", func(t *testing.T) {
		views, err := simplejsonx.Numeric[int](object, "views", nil)
		require.NoError(t, err)
		require.Equal(t, 1200, views)

		price, err := simplejsonx.Numeric[float64](object, "price", nil)
		require.NoError(t, err)
		require.Equal(t, 9.5, price)
	})

	t.Run("strict-without-option", func(t *testing.T) {
		for _, key := range []string{"rating", "discount", "weight", "title", "absent"} {
			_, err := simplejsonx.Numeric[float64](object, key, nil)
			require.Error(t, err, key)
		}
		// "9.5" is numeric, yet not an int
		// "9.5" 是数字，但不是整数
		_, err := simplejsonx.Numeric[int](object, "price", nil)
		require.Error(t, err)
	})

	t.Run("option-accepts-what-it-names", func(t *testing.T) {
		option := &simplejsonx.NumericOption{AbsentOk: true, Nullable: true, Blanks: []string{"", "N/A"}}
		for _, key := range []string{"rating", "discount", "weight", "absent"} {
			value, err := simplejsonx.Numeric[float64](object, key, option)
			require.NoError(t, err, key)
			require.Zero(t, value, key)
		}
		// A bad value stays an error even with each option on
		// 选项全开时坏值照样报错
		_, err := simplejsonx.Numeric[float64](object, "title", option)
		require.Error(t, err)
	})

	t.Run("nil-object-follows-nullable", func(t *testing.T) {
		_, err := simplejsonx.Numeric[float64](nil, "views", nil)
		require.Error(t, err)

		value, err := simplejsonx.Numeric[float64](nil, "views", &simplejsonx.NumericOption{Nullable: true})
		require.NoError(t, err)
		require.Zero(t, value)
	})
}
