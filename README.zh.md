<!-- TEMPLATE (ZH) BEGIN: BADGES -->

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/simplejsonx/release.yml?branch=main&label=BUILD)](https://github.com/yylego/simplejsonx/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/simplejsonx)](https://pkg.go.dev/github.com/yylego/simplejsonx)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/simplejsonx/main.svg)](https://coveralls.io/github/yylego/simplejsonx?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.26%2B-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/simplejsonx.svg)](https://github.com/yylego/simplejsonx/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/simplejsonx)](https://goreportcard.com/report/github.com/yylego/simplejsonx)

<!-- TEMPLATE (ZH) CLOSE: BADGES -->

# simplejsonx

`simplejsonx` 是个基于泛型的 JSON 解析库，依赖于 `github.com/bitly/go-simplejson`，通过增强类型安全性和灵活性来优化 JSON 处理。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) CLOSE: LANGUAGE NAVIGATION -->

## 核心特性

### `go-simplejson` 的问题

`go-simplejson` 的 `Must` 方法支持可选的默认值：

```go
name := object.Get("name").MustString("")   // 如果 "name" 不存在呢？
age := object.Get("age").MustInt(0)         // 0 是真实值，还是静默的兜底？
```

`Must` 方法在转换失败时返回指定的默认值；不传默认值时使用类型的零值。这会让真实的零值与缺字段、null、类型错误难以区分。原库的 `String()`、`Int()` 等方法本身会返回错误；`simplejsonx` 在此基础上增加泛型提取和字段校验。

值得注意的是，Go 社区对 `Must` 的主流理解是"**必须成功，否则 panic**"——如 `template.Must()`、`regexp.MustCompile()` 等标准库函数均遵循此语义。而 `go-simplejson` 的 `Must` 含义与此相反——"尽量获取，做不到就用默认值"，更接近"尽力而为"而非"必须成功"。`simplejsonx` 的子包 `simplejsonm` 遵循 Go 主流语义：`Must` 即断言，失败即 panic。

`simplejsonx` 采取相反的立场：**数据源为你服务。** 无需提供默认值并期望一切正常，而是**断言**每个字段必须是什么，代码会在出问题的那一刻立刻告诉你：

```go
name, err := simplejsonx.Extract[string](object, "name")  // 缺失/null/类型不对 → err
age, err := simplejsonx.Extract[int](object, "age")        // 无需默认值，不会静默掩盖异常
```

核心区别：`go-simplejson` 问的是"出错时我该返回什么？"而 `simplejsonx` 问的是"**这里该出错吗？**"

> **适用场景**：`simplejsonx` 的主要使用场景是对方返回的数据结构**不确定或可变**——同一个字段在不同业务条件下可能返回不同类型（如有时是数字，有时是字符串 `"-"`，有时是 null），字段可能随版本更新增减或重命名，且通常没有完善的开放文档。在这种情况下，需要逐字段提取并校验类型，主动检测数据变化，而非静默兜底。相比之下，`go-simplejson` 更适合数据结构已知且稳定的场景，开发者清楚每个字段的含义，可以放心使用默认值。

**函数对照**：

| `go-simplejson` 写法                         | 行为                              | `simplejsonx` 替代                                                    | 行为                             |
| -------------------------------------------- | --------------------------------- | --------------------------------------------------------------------- | -------------------------------- |
| `.Get("k").MustString("")`                   | 缺失/null/类型不对 → 返回 `""`    | `Extract[string](obj, "k")`                                           | 缺失/null/类型不对 → 报错        |
| `.Get("k").MustInt(0)`                       | 缺失/null/类型不对 → 返回 `0`     | `Extract[int](obj, "k")`                                              | 缺失/null/类型不对 → 报错        |
| `.Get("k").MustFloat64(0)`                   | 缺失/null/类型不对 → 返回 `0`     | `Extract[float64](obj, "k")`                                          | 缺失/null/类型不对 → 报错        |
| `.Get("k").MustBool(false)`                  | 缺失/null/类型不对 → 返回 `false` | `Extract[bool](obj, "k")`                                             | 缺失/null/类型不对 → 报错        |
| `.Get("k").MustInt(0)`（值不能为 0）         | 无法区分"值为0"和"字段损坏"       | `Solicit[int](obj, "k")`                                              | 缺失/null/类型不对/零值 → 报错   |
| `.Get("k").MustString("")`（字段可能不存在） | 无法区分"不存在"和"值为空串"      | `Inspect[string](obj, "k")`                                           | 不存在 → 零值，存在但异常 → 报错 |
| `.Get("k").MustFloat64(0)`（值可能为 null）  | 无法区分"null"和"值为0"           | `Harvest[float64](obj, "k", &HarvestOption[float64]{Nullable: true})` | null → 零值，类型不对 → 报错     |

### 设计思路

**前提**：官方接口返回的数据始终正确。不同业务场景下可能返回 null、零值、甚至不同类型——这些均为合法响应。错误只来自**解析端**：字段名拼写有误、类型假设不当、或未能跟上接口改版。

**核心方法——两侧逼近**：如同数学中从两侧逼近极限值，`simplejsonx` 支持**从严格和宽松两侧逼近每个字段**，精准卡在官方数据的真实边界上。必填字段使用严格校验（`Extract`），允许为空的字段使用受控放宽（`Harvest` + `Nullable:true`）——不多不少，刚好贴合接口的真实契约。

**桩点校验**：虽然对方返回的数据结构不确定，但并非完全不可知——通过观察真实数据，总能找到一些**确定性规则**作为桩点：某些字段必然存在、某些字段必然为特定类型、某些字段在特定条件下必然非空。将这些规则用严格校验（`Extract`）固定下来，相当于在不确定的数据流中打下桩。一旦对方接口在某天变更了这些桩点字段，代码立即报错或崩溃——将变更拦截在解析层，防止错误数据流入下游业务。

**设计意义**：解析层是外部数据进入系统的**第一道类型断言**——对每个字段做准入校验和契约检查。若不加甄别地放行一切（到处使用 `MustString("")`），等于跳过所有断言，默默接受损坏数据，将问题传至下游，造成无法追踪的破坏。当官方接口变更结构、重命名字段、调整数据类型时，代码应**立即崩溃**，使开发者在数分钟内感知并修复——而非数日后才通过被污染的业务数据发现异常。

**核心原则**：默认从严，仅在数据源确实需要处放宽。绝不盲目接受。

> ⚠️ **禁止批量放宽**：当某个字段校验失败时，只放开**那一个字段**——绝不因为"这个字段是 null，所以别的大概也会是 null"而批量放宽。每个字段的放宽都必须有真实数据作为依据，不能靠类比推断。批量放宽等于自废所有桩点，整个解析层的校验能力将归零。

### 推荐工作流

1. **新写解析代码时，所有字段一律使用 `Extract`（最严格）**
2. **执行测试，使用真实数据验证**——观察哪些字段 panic
3. **遇到 panic 先判断原因**：
   - 字段名有误 → 修正字段名，继续使用 `Extract`
   - 接口异常 → 排查接口问题，不要放宽代码
   - 该字段在某些业务场景下确实为 null → 进入第 4 步
4. **仅对确认为可选/可空的字段做最小放宽**，并注释说明原因：
   - 字段可能不存在 → 改用 `Inspect`
   - 字段可能不存在，且需要知道是否存在 → 改用 `Inquire`
   - 字段存在但值可能为 null → 改用 `Harvest` + `Nullable:true`
   - 字段是尽力获取的，不需要知道具体缺失原因 → 改用 `Attempt`
5. **绝不批量放宽**——每次仅放开一个字段，重新执行测试，逐个确认。发现字段 A 为 null 不代表字段 B 也为 null，每个字段必须有独立的真实数据证据，禁止类比推断

### 函数严格程度

🎯 **基于泛型的 JSON 提取** — 返回类型在编译期确定，JSON 类型与支持的转换在运行时检查

各个函数对不同数据场景有不同的严格程度：

| 函数                      | 字段不存在       | 字段存在值为 null | 字段存在类型不对 | 字段存在类型正确值为零值 | 字段存在类型正确值非零 |
| ------------------------- | ---------------- | ----------------- | ---------------- | ------------------------ | ---------------------- |
| `Extract`                 | ❌ 报错          | ❌ 报错           | ❌ 报错          | ✅ 返回零值              | ✅ 返回值              |
| `Solicit`                 | ❌ 报错          | ❌ 报错           | ❌ 报错          | ❌ 报错                  | ✅ 返回值              |
| `Inspect`                 | ✅ 返回零值      | ❌ 报错           | ❌ 报错          | ✅ 返回零值              | ✅ 返回值              |
| `Resolve`                 | 不适用           | ❌ 报错           | ❌ 报错          | ✅ 返回零值              | ✅ 返回值              |
| `Inquire`                 | ✅ (零值, false) | ❌ 报错           | ❌ 报错          | ✅ (零值, true)          | ✅ (值, true)          |
| `Attempt`                 | ✅ (零值, false) | ✅ (零值, false)  | ✅ (零值, false) | ✅ (零值, true)          | ✅ (值, true)          |
| `Harvest(Nullable:true)`  | ❌ 报错          | ✅ 返回零值       | ❌ 报错          | ✅ 返回零值              | ✅ 返回值              |
| `Harvest(Nullable:false)` | ❌ 报错          | ❌ 报错           | ❌ 报错          | ✅ 返回零值              | ✅ 返回值              |

🔀 `Harvest` 配合 `HarvestOption` 还可处理可能是数字、字符串"-"或 null 的混合类型字段

表格以 JSON 对象及默认的 `AbsentOk:false` 为前提。`AbsentOk:true` 允许缺字段，但不接受错误的对象结构。`Inquire` 和 `Explore` 出错时也返回 false，应先检查错误，再判断布尔值。

🔗 `Explore` 通过点分隔路径遍历嵌套 JSON，如 `"user.profile.name"`

🔄 `Strconv` 将字符串编码的值（如 `"42"` → `int 42`）通过两阶段过程转换

### 子包错误处理策略

同样的函数，不同的错误应对方式：

| 包                             | 出错时       | 返回值       | 适用场景                     |
| ------------------------------ | ------------ | ------------ | ---------------------------- |
| `simplejsonx`（基础包）        | 返回 error   | `(T, error)` | 需要自行处理每个错误         |
| `simplejsonx/sure/simplejsonm` | panic        | `T`          | 出错表明数据损坏，应立即终止 |
| `simplejsonx/sure/simplejsono` | 静默忽略     | `T`          | 有意接受并丢弃错误           |
| `simplejsonx/sure/simplejsons` | 打印警告日志 | `T`          | 需要感知错误但不终止运行     |

---

## 安装

```bash
go get github.com/yylego/simplejsonx
```

## 使用示例

### 1. 基础 JSON 解析与错误处理

此示例展示如何加载 JSON 数据并使用错误处理提取类型化字段。

```go
package main

import (
	"fmt"
	"log"

	"github.com/yylego/simplejsonx"
)

func main() {
	// 示例 JSON 数据
	data := []byte(`{
		"name": "yylego",
		"age": 18,
		"is_rich": true
	}`)

	// 加载 JSON 数据
	object, err := simplejsonx.Load(data)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	// 提取字段
	name, err := simplejsonx.Extract[string](object, "name")
	if err != nil {
		log.Fatalf("Error extracting 'name': %v", err)
	}

	age, err := simplejsonx.Extract[int](object, "age")
	if err != nil {
		log.Fatalf("Error extracting 'age': %v", err)
	}

	isRich, err := simplejsonx.Extract[bool](object, "is_rich")
	if err != nil {
		log.Fatalf("Error extracting 'is_rich': %v", err)
	}

	// 输出提取的值
	fmt.Println("name:", name, "age:", age, "rich:", isRich)  // 输出: name: yylego age: 18 rich: true
}
```

⬆️ **源码：** [源码](internal/demos/demo1x/main.go)

### 2. Must 风格 API（错误时 Panic）

此示例展示 must 风格 API，在遇到错误时会 panic，适用于不期望出错的场景。

```go
package main

import (
	"fmt"

	"github.com/yylego/simplejsonx/sure/simplejsonm"
)

func main() {
	// 示例 JSON 数据
	data := []byte(`{
		"name": "yylego",
		"age": 18,
		"is_rich": true
	}`)

	// 必须成功加载 JSON 数据（失败会 panic）
	object := simplejsonm.Load(data)

	// 必须成功提取字段（失败会 panic）
	name := simplejsonm.Extract[string](object, "name")
	age := simplejsonm.Extract[int](object, "age")
	isRich := simplejsonm.Extract[bool](object, "is_rich")

	// 输出提取的值
	fmt.Println("name:", name, "age:", age, "rich:", isRich)  // 输出: name: yylego age: 18 rich: true
}
```

⬆️ **源码：** [源码](internal/demos/demo2x/main.go)

### 3. 使用 Harvest 处理混合类型字段

此示例展示如何使用 `Harvest` 处理可能是数字、字符串"-"或 null 的字段。

```go
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yylego/simplejsonx"
)

func main() {
	// 示例 JSON 数据，包含混合类型的 price 字段
	data := []byte(`{
		"items": [
			{"name": "A", "price": 100.5},
			{"name": "B", "price": "-"},
			{"name": "C", "price": null}
		]
	}`)

	// 加载 JSON 数据
	object, err := simplejsonx.Load(data)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	// 获取 items 列表
	items, err := simplejsonx.GetList(object, "items")
	if err != nil {
		log.Fatalf("Error getting items: %v", err)
	}

	// 定义收获选项，处理混合类型
	option := &simplejsonx.HarvestOption[float64]{
		Nullable: true, // 接受 null 值
		Strconv: func(s string) (float64, error) {
			if s == "-" || s == "N/A" {
				return 0, nil // 将 "-" 视为零值
			}
			return strconv.ParseFloat(s, 64)
		},
	}

	// 从每个项目中收获价格
	for _, item := range items {
		name, err := simplejsonx.Extract[string](item, "name")
		if err != nil {
			log.Fatalf("Error extracting name: %v", err)
		}

		price, err := simplejsonx.Harvest[float64](item, "price", option)
		if err != nil {
			log.Fatalf("Error harvesting price: %v", err)
		}

		fmt.Printf("Item %s: price = %.2f\n", name, price)  // 输出: A=100.50, B=0.00, C=0.00
	}
}
```

⬆️ **源码：** [源码](internal/demos/demo3x/main.go)

## 功能示例

### 提取 JSON 字段

**基础字段提取：**

```go
res, err := simplejsonx.Extract[int](object, "age")
if err != nil {
    log.Fatalf("Error extracting 'age': %v", err)
}
fmt.Println("Age:", res)  // 输出: 18
```

### 解析 JSON 值

**使用 Resolve 获取类型化值：**

```go
object, err := simplejsonx.Load([]byte(`{
	"height": 175,
	"weight": 80
}`))
if err != nil {
	log.Fatalf("Error loading JSON: %v", err)
}

height, err := simplejsonx.Resolve[int64](object.Get("height"))
if err != nil {
	log.Fatalf("Error resolving 'height': %v", err)
}
fmt.Println("Height:", height)  // 输出: 175
```

### 使用 Inspect 处理缺失字段

**Inspect 在字段缺失时返回零值：**

```go
object, err := simplejsonx.Load([]byte(`{
	"name": "yylego",
	"age": 18
}`))
if err != nil {
	log.Fatalf("Error loading JSON: %v", err)
}

name, err := simplejsonx.Inspect[string](object, "name")
if err != nil {
    log.Fatalf("Error inspecting 'name': %v", err)
}
fmt.Println("Name:", name)  // 输出: yylego

address, err := simplejsonx.Inspect[string](object, "address")
if err != nil {
    log.Fatalf("Error inspecting 'address': %v", err)
}
fmt.Println("Address:", address)  // 输出: 空字符串（零值）
```

## API 参考

### 核心函数

| 函数                           | 说明                                                 |
| ------------------------------ | ---------------------------------------------------- |
| `Load(data []byte)`            | 解析 JSON 字节为 simplejson.Json 实例                |
| `LoadWithOption(data, option)` | 去掉指明的前缀 / 后缀（如 JSONP 回调）再按 JSON 解析 |
| `Wrap(value any)`              | 包装 Go 值为 simplejson.Json 对象                    |
| `List(elements []any)`         | 转换切片为 simplejson.Json 对象切片                  |

### 类型安全提取函数

| 函数         | 参数                    | 返回值             | 说明                                                                |
| ------------ | ----------------------- | ------------------ | ------------------------------------------------------------------- |
| `Extract[T]` | `(object, key)`         | `(T, error)`       | 严格提取：key 必须存在，值不能为 null                               |
| `Solicit[T]` | `(object, key)`         | `(T, error)`       | 最严格索取：key 必须存在，值不能为 null / 零值（T 须为 comparable） |
| `Inspect[T]` | `(object, key)`         | `(T, error)`       | 宽松提取：key 不存在时返回零值不报错，存在但 null/类型不对则报错    |
| `Resolve[T]` | `(object)`              | `(T, error)`       | 直接转换：将 Json 对象转为目标类型 T                                |
| `Inquire[T]` | `(object, key)`         | `(T, bool, error)` | 三态提取：返回值 + 是否存在 + 错误                                  |
| `Attempt[T]` | `(object, key)`         | `(T, bool)`        | 尝试提取：成功返回 (值, true)，任何失败返回 (零值, false)           |
| `Harvest[T]` | `(object, key, option)` | `(T, error)`       | 可配置提取：通过 HarvestOption 控制 null 处理和字符串转换           |

### 高级函数

| 函数                              | 说明                                                  |
| --------------------------------- | ----------------------------------------------------- |
| `GetList(object, key)`            | 获取指定键的 JSON 数组为 simplejson.Json 切片         |
| `Explore[T](object, path)`        | 通过点分隔路径导航嵌套 JSON（如 "user.profile.name"） |
| `Harvest[T](object, key, option)` | 提取混合类型字段，支持可配置的 null/字符串处理        |
| `Numeric[T](object, key, option)` | 读取 JSON 数字或数字字符串，选项决定哪些情况读作零值  |
| `Strconv[T](object)`              | 两阶段转换：JSON → 字符串 → 目标类型                  |

### 支持的类型

`Resolve` 和 `Extract` 支持以下类型：

- **基础类型**：`int`、`int64`、`uint64`、`float64`、`string`、`bool`
- **集合类型**：`[]string`、`[]any`、`[]*simplejson.Json`
- **复杂类型**：`map[string]any`、`[]byte`、`*simplejson.Json`

`type Amount int` 这样的定义类型不会按 `int` 解析；不支持的类型在运行时报错。`Wrap` 保留传入的 Go 值，不执行 JSON 编码。数值转换沿用 `go-simplejson`，包装后的浮点数转整数可能截断，整数转换可能溢出；`[]string` 转换会将 null 元素变成空字符串。`[]byte` 读取 JSON 字符串的字节，不做 Base64 解码。`Load` 要求输入为单个完整 JSON 文档，并通过 `json.Number` 保留数值。

### HarvestOption 配置

`HarvestOption[T]` 用于配置混合类型字段的解析行为：

| 字段       | 类型                      | 说明                                |
| ---------- | ------------------------- | ----------------------------------- |
| `Nullable` | `bool`                    | null 值或 nil 对象返回零值          |
| `AbsentOk` | `bool`                    | 缺键时返回零值，保留类型检查        |
| `Strconv`  | `func(string) (T, error)` | 自定义字符串到目标类型 T 的转换函数 |

`AbsentOk` 与 `Nullable` 相互独立，同时启用才同时接受缺键和 null。option 为 nil 时使用严格默认配置；空键名和转换失败仍返回错误。`Nullable` 接受 Go 的 nil 对象；非 nil 对象按字段名取值时，其内容必须是 JSON 对象。

### NumericOption 配置

`NumericOption` 用于配置 `Numeric`，它读取以 JSON 数字或数字字符串给出的数字，如 `9.5` 和 `"9.5"`：

| 字段       | 类型       | 说明                                                   |
| ---------- | ---------- | ------------------------------------------------------ |
| `AbsentOk` | `bool`     | 缺键时返回零值                                         |
| `Nullable` | `bool`     | null 值或 nil 对象返回零值，与 `HarvestOption` 一致    |
| `Blanks`   | `[]string` | 表示没有值的字符串，如 `""` 和 `"N/A"`，遇到时返回零值 |

不传选项时最严：键必须存在，不接受 null，字符串必须是数字。

```go
price := simplejsonm.Numeric[float64](object, "price", &simplejsonx.NumericOption{Nullable: true, Blanks: []string{"", "N/A"}})
```

### 使用 Strconv — 字符串转换

**通过两阶段过程将字符串值转换为目标类型：**

```go
object, err := simplejsonx.Load([]byte(`{
	"count": "42",
	"rate": "3.14"
}`))
if err != nil {
	log.Fatalf("Error loading JSON: %v", err)
}

// 字符串 "42" → int 42
count, err := simplejsonx.Strconv[int](object.Get("count"))
if err != nil {
	log.Fatalf("Error converting count: %v", err)
}
fmt.Println("Count:", count)  // 输出: 42

// 字符串 "3.14" → float64 3.14
rate, err := simplejsonx.Strconv[float64](object.Get("rate"))
if err != nil {
	log.Fatalf("Error converting rate: %v", err)
}
fmt.Println("Rate:", rate)  // 输出: 3.14
```

### 错误处理变体（子包）

本包通过子包提供三种错误处理模式：

| 包路径                         | 模式 | 行为                                 |
| ------------------------------ | ---- | ------------------------------------ |
| `simplejsonx/sure/simplejsonm` | Must | 出错时 panic（用于不容许失败的场景） |
| `simplejsonx/sure/simplejsono` | Omit | 返回零值，忽略错误                   |
| `simplejsonx/sure/simplejsons` | Soft | 记录错误日志，返回零值               |

**Must 模式使用示例：**

```go
import "github.com/yylego/simplejsonx/sure/simplejsonm"

// 如果 JSON 解析或提取失败则 panic
object := simplejsonm.Load(data)
name := simplejsonm.Extract[string](object, "name")
```

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) CLOSE: STANDARD PROJECT FOOTER -->
