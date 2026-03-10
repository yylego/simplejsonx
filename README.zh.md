[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/simplejsonx/release.yml?branch=main&label=BUILD)](https://github.com/yylego/simplejsonx/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/simplejsonx)](https://pkg.go.dev/github.com/yylego/simplejsonx)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/simplejsonx/main.svg)](https://coveralls.io/github/yylego/simplejsonx?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.24--1.25-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yylego/simplejsonx.svg)](https://github.com/yylego/simplejsonx/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/simplejsonx)](https://goreportcard.com/report/github.com/yylego/simplejsonx)

# simplejsonx

`simplejsonx` 是个基于泛型的 JSON 解析库，依赖于 `github.com/bitly/go-simplejson`，通过增强类型安全性和灵活性来优化 JSON 处理，同时要求至少使用 Go 1.22 版本（需要支持泛型）。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

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
		Nilable: true, // 接受 null 值
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

| 函数                           | 说明                                  |
| ------------------------------ | ------------------------------------- |
| `Load(data []byte)`            | 解析 JSON 字节为 simplejson.Json 实例 |
| `Wrap(value interface{})`      | 包装 Go 值为 simplejson.Json 对象     |
| `List(elements []interface{})` | 转换切片为 simplejson.Json 对象切片   |

### 类型安全提取函数

| 函数                      | 说明                                       |
| ------------------------- | ------------------------------------------ |
| `Extract[T](object, key)` | 严格模式提取指定键的值（键缺失时返回错误） |
| `Inspect[T](object, key)` | 宽松模式提取指定键的值（键缺失时返回零值） |
| `Resolve[T](object)`      | 将 simplejson.Json 对象转换为目标类型 T    |
| `Inquire[T](object, key)` | 三态提取，返回 (value, exists, error)      |
| `Attempt[T](object, key)` | 简洁提取，返回 (value, ok)，不暴露错误详情 |

### 高级函数

| 函数                              | 说明                                                  |
| --------------------------------- | ----------------------------------------------------- |
| `GetList(object, key)`            | 获取指定键的 JSON 数组为 simplejson.Json 切片         |
| `Explore[T](object, path)`        | 通过点分隔路径导航嵌套 JSON（如 "user.profile.name"） |
| `Harvest[T](object, key, option)` | 提取混合类型字段，支持可配置的 null/字符串处理        |
| `Strconv[T](object)`              | 两阶段转换：JSON → 字符串 → 目标类型                  |

### 支持的类型

`Resolve` 和 `Extract` 支持以下类型：

- **基础类型**：`int`、`int64`、`uint64`、`float64`、`string`、`bool`
- **集合类型**：`[]string`、`[]interface{}`、`[]*simplejson.Json`
- **复杂类型**：`map[string]interface{}`、`[]byte`、`*simplejson.Json`

### HarvestOption 配置

`HarvestOption[T]` 用于配置混合类型字段的解析行为：

| 字段      | 类型                      | 说明                                |
| --------- | ------------------------- | ----------------------------------- |
| `Nilable` | `bool`                    | 接受 null 值不报错（返回零值）      |
| `Strconv` | `func(string) (T, error)` | 自定义字符串到目标类型 T 的转换函数 |

### 使用 Strconv 进行字符串转换

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

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/yylego/simplejsonx.svg?variant=adaptive)](https://starchart.cc/yylego/simplejsonx)
