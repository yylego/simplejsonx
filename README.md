[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/simplejsonx/release.yml?branch=main&label=BUILD)](https://github.com/yylego/simplejsonx/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/simplejsonx)](https://pkg.go.dev/github.com/yylego/simplejsonx)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/simplejsonx/main.svg)](https://coveralls.io/github/yylego/simplejsonx?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.24%2C%201.25-lightgrey.svg)](https://github.com/yylego/simplejsonx)
[![GitHub Release](https://img.shields.io/github/release/yylego/simplejsonx.svg)](https://github.com/yylego/simplejsonx/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/simplejsonx)](https://goreportcard.com/report/github.com/yylego/simplejsonx)

# simplejsonx

`simplejsonx` is a generic-based JSON parsing package that depends on `github.com/bitly/go-simplejson`, enabling type-safe JSON processing with Go generics.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

### The Problem with `go-simplejson`

The standard `go-simplejson` API forces you to provide a default value on every access:

```go
name := object.Get("name").MustString("")   // what if "name" doesn't exist?
age := object.Get("age").MustInt(0)         // is 0 the real value, or a silent fallback?
```

The `Must` functions are designed to **force you to provide a default value** — even when you have no idea what the default should be, you must provide one. When the field is missing, null, or has the wrong type, `MustString("")` / `MustInt(0)` **return that default without notice** — you have no idea if `0` means "the actual value is zero" or "the field doesn't exist" or "the API renamed it last week." Exceptions are masked behind a mandatory default value.

In the Go ecosystem, the mainstream understanding of `Must` is "**must succeed, or panic**" — as seen in standard functions like `template.Must()` and `regexp.MustCompile()`. The `go-simplejson` `Must` means the opposite — "attempt to get a value, use the default if it fails" — more like "best attempt" than "must succeed." The `simplejsonx` sub-package `simplejsonm` follows the Go mainstream convention: `Must` means assertion — on failure, panic.

`simplejsonx` takes the opposite stance: **the data source serves you.** Instead of providing defaults and expecting correctness, you **assert** what each field must be, and the code tells you the moment something is wrong:

```go
name, err := simplejsonx.Extract[string](object, "name")  // missing/null/wrong type → err
age, err := simplejsonx.Extract[int](object, "age")        // no default needed, no silent swallowing
```

The distinction: `go-simplejson` asks "what should I return when things go wrong?" while `simplejsonx` asks "**should things be going wrong at all?**"

> **Use cases**: `simplejsonx` is designed for scenarios where the data structure from the upstream is **uncertain or unstable** — the same field may have different types under different business conditions (e.g., sometimes a numeric value, sometimes a string like `"-"`, sometimes null), fields may be added, removed, or renamed across versions, and comprehensive documentation is often absent. In such cases, you need to extract and validate each field one at a time, with active detection of data changes instead of silent fallback to defaults. In contrast, `go-simplejson` is more suitable for scenarios where the data structure is known and fixed, and the expected default values are well understood.

**Function mapping**:

| `go-simplejson` | Behavior | `simplejsonx` replacement | Behavior |
|----------------|----------|--------------------------|----------|
| `.Get("k").MustString("")` | missing/null/wrong type → return `""` | `Extract[string](obj, "k")` | missing/null/wrong type → error |
| `.Get("k").MustInt(0)` | missing/null/wrong type → return `0` | `Extract[int](obj, "k")` | missing/null/wrong type → error |
| `.Get("k").MustFloat64(0)` | missing/null/wrong type → return `0` | `Extract[float64](obj, "k")` | missing/null/wrong type → error |
| `.Get("k").MustBool(false)` | missing/null/wrong type → return `false` | `Extract[bool](obj, "k")` | missing/null/wrong type → error |
| `.Get("k").MustInt(0)` (value must not be 0) | cannot distinguish "value is 0" from "field broken" | `Solicit[int](obj, "k")` | missing/null/wrong type/zero → error |
| `.Get("k").MustString("")` (key may not exist) | cannot distinguish "missing" from "empty string" | `Inspect[string](obj, "k")` | missing returns zero, present but abnormal errors |
| `.Get("k").MustFloat64(0)` (value may be null) | cannot distinguish "null" from "value is 0" | `Harvest[float64](obj, "k", &HarvestOption[float64]{Nilable: true})` | null returns zero, wrong type errors |

### Design Approach

**Premise**: The upstream API returns correct data. It may return null, zero values, or different types depending on business scenarios — but these are legitimate responses. Errors come from **our parsing side**: wrong field names, wrong type assumptions, or not keeping up with API changes.

**Core approach — Squeeze Theorem**: Like the math squeeze theorem that pins down a limit from both sides, `simplejsonx` lets you **squeeze each field from both the strict and lenient sides** to land on the API's exact contract edge. Required fields get strict validation (`Extract`), fields that accept null get controlled relaxation (`Harvest` with `Nilable:true`) — no more, no less. The result is a parsing gate that fits the API's true contract like a glove.

**Anchor validation**: Although the data structure from upstream is uncertain, it is not unknowable — through observation of actual data, you can identify **deterministic rules** that serve as anchors: some fields exist in every response, some fields have a fixed type, some fields are non-null under specific conditions. Locking these rules down with strict validation (`Extract`) is like driving anchor stakes into an uncertain data stream. The moment the upstream API changes one of these anchor fields, the code errors or crashes at once — intercepting the change at the parsing gate and preventing corrupted data from flowing into downstream business logic.

**What's at stake**: The parsing gate is the **first type assertion** where outside data enters the system — performing admission validation and contract checking on each field. Letting everything in without challenge (using `MustString("")` everywhere) skips assertions, accepting broken data and passing it downstream where it causes untraceable damage. When the upstream API changes structure, renames fields, or adjusts data types, the code must **crash at once** so you detect and fix the change within minutes, not discover it days afterward through corrupted business data.

**In short**: Be strict as the default. Relax just where the data source demands it. Never accept without verification.

> ⚠️ **No batch relaxation**: When a field fails validation, relax **just that one field** — do not assume "this field is null, so others must be too" and relax them as a group. Each relaxation must be backed by actual data evidence, not by assumption. Batch relaxation destroys anchor points, reducing the entire parsing gate's validation strength to zero.

### Recommended Workflow

1. **Start with `Extract` (strictest) for each field** when writing new parsing code
2. **Run tests with actual data** — observe which fields panic
3. **When panic occurs, diagnose the cause**:
   - Wrong field name → fix the name, keep using `Extract`
   - API unexpected change → investigate the API, do not relax the code
   - The field is confirmed null in certain business scenarios → proceed to step 4
4. **Relax just the specific field that is confirmed as not required or accepts null**, and comment the reason:
   - Field may not exist → switch to `Inspect`
   - Field may not exist, and you need to know if it was present → switch to `Inquire`
   - Field exists but value may be null → switch to `Harvest` with `Nilable:true`
   - Field is best-effort, you don't need the exact cause of absence → switch to `Attempt`
5. **Relax one field at a time** — re-run tests, confirm each one. Field A being null does not mean field B is also null — each field must have independent actual data evidence, no reasoning based on assumption

### Function Strictness Modes

🎯 **Type-safe JSON extraction with Go generics** — compile-time type checking, no runtime type assertions needed

Each function handles data scenarios with a different strictness:

| Function | Key Missing | Value = null | Wrong Type | Correct Type, Zero Value | Correct Type, Non-Zero |
|----------|------------|-------------|-----------|------------------------|----------------------|
| `Extract` | ❌ error | ❌ error | ❌ error | ✅ return zero | ✅ return value |
| `Solicit` | ❌ error | ❌ error | ❌ error | ❌ error | ✅ return value |
| `Inspect` | ✅ return zero | ❌ error | ❌ error | ✅ return zero | ✅ return value |
| `Resolve` | N/A | ❌ error | ❌ error | ✅ return zero | ✅ return value |
| `Inquire` | ✅ (zero, false) | ❌ error | ❌ error | ✅ (zero, true) | ✅ (val, true) |
| `Attempt` | ✅ (zero, false) | ✅ (zero, false) | ✅ (zero, false) | ✅ (zero, true) | ✅ (val, true) |
| `Harvest(Nilable:true)` | ❌ error | ✅ return zero | ❌ error | ✅ return zero | ✅ return value |
| `Harvest(Nilable:false)` | ❌ error | ❌ error | ❌ error | ✅ return zero | ✅ return value |

🔀 `Harvest` with `HarvestOption` also handles fields that might be numbers, strings like "-", or null

🔗 `Explore` traverses nested JSON via dot-separated paths like `"user.profile.name"`

🔄 `Strconv` converts string-encoded values (e.g., `"42"` → `int 42`) via two-stage process

### Sub-Package Error Strategies

Same functions, different error strategies:

| Package | On Error | Return Value | Use When |
|---------|----------|-------------|----------|
| `simplejsonx` (base) | return error | `(T, error)` | You want to handle each error |
| `simplejsonx/sure/simplejsonm` | panic | `T` | Broken data detected, should crash at once |
| `simplejsonx/sure/simplejsono` | ignore without noise | `T` | You accept and discard errors on purpose |
| `simplejsonx/sure/simplejsons` | log warning | `T` | You want to know about errors but not crash |

---

## Installation

```bash
go get github.com/yylego/simplejsonx
```

## Usage Example

### 1. Basic JSON Parsing with Error Handling

This example demonstrates how to load JSON data and extract typed fields with error handling.

```go
package main

import (
	"fmt"
	"log"

	"github.com/yylego/simplejsonx"
)

func main() {
	// Sample JSON data
	data := []byte(`{
		"name": "yylego",
		"age": 18,
		"is_rich": true
	}`)

	// Load the JSON data
	object, err := simplejsonx.Load(data)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	// Extract fields
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

	// Output the extracted values
	fmt.Println("name:", name, "age:", age, "rich:", isRich)  // Output: name: yylego age: 18 rich: true
}
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### 2. Must-Style API (Panic on Failure)

This example shows the must-style API that panics on failure, best suited when failures are unexpected.

```go
package main

import (
	"fmt"

	"github.com/yylego/simplejsonx/sure/simplejsonm"
)

func main() {
	// Sample JSON data
	data := []byte(`{
		"name": "yylego",
		"age": 18,
		"is_rich": true
	}`)

	// Must Load the JSON data (panics on failure)
	object := simplejsonm.Load(data)

	// Must Extract fields (panics on failure)
	name := simplejsonm.Extract[string](object, "name")
	age := simplejsonm.Extract[int](object, "age")
	isRich := simplejsonm.Extract[bool](object, "is_rich")

	// Output the extracted values
	fmt.Println("name:", name, "age:", age, "rich:", isRich)  // Output: name: yylego age: 18 rich: true
}
```

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

### 3. Handling Mixed-Type Fields with Harvest

This example shows how to use `Harvest` to handle fields that might be numbers, strings like "-", and null.

```go
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yylego/simplejsonx"
)

func main() {
	// Sample JSON data with mixed-type price field
	data := []byte(`{
		"items": [
			{"name": "A", "price": 100.5},
			{"name": "B", "price": "-"},
			{"name": "C", "price": null}
		]
	}`)

	// Load the JSON data
	object, err := simplejsonx.Load(data)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	// Get items list
	items, err := simplejsonx.GetList(object, "items")
	if err != nil {
		log.Fatalf("Error getting items: %v", err)
	}

	// Define harvest option with mixed-type handling
	option := &simplejsonx.HarvestOption[float64]{
		Nilable: true, // Accept null values
		Strconv: func(s string) (float64, error) {
			if s == "-" || s == "N/A" {
				return 0, nil // Treat "-" as zero
			}
			return strconv.ParseFloat(s, 64)
		},
	}

	// Harvest price from each item
	for _, item := range items {
		name, err := simplejsonx.Extract[string](item, "name")
		if err != nil {
			log.Fatalf("Error extracting name: %v", err)
		}

		price, err := simplejsonx.Harvest[float64](item, "price", option)
		if err != nil {
			log.Fatalf("Error harvesting price: %v", err)
		}

		fmt.Printf("Item %s: price = %.2f\n", name, price)  // Output: A=100.50, B=0.00, C=0.00
	}
}
```

⬆️ **Source:** [Source](internal/demos/demo3x/main.go)

## Examples

### Extracting JSON Fields

**Basic field extraction:**

```go
res, err := simplejsonx.Extract[int](object, "age")
if err != nil {
    log.Fatalf("Error extracting 'age': %v", err)
}
fmt.Println("Age:", res)  // Output: 18
```

### Resolving JSON Values

**Using Resolve to get typed values:**

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
fmt.Println("Height:", height)  // Output: 175
```

### Using Inspect to Handle Missing Fields

**Inspect returns zero value when field is missing:**

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
fmt.Println("Name:", name)  // Output: yylego

address, err := simplejsonx.Inspect[string](object, "address")
if err != nil {
    log.Fatalf("Error inspecting 'address': %v", err)
}
fmt.Println("Address:", address)  // Output: Blank string (zero value)
```

## API Reference

### Core Functions

| Function                       | Description                                    |
| ------------------------------ | ---------------------------------------------- |
| `Load(data []byte)`            | Parse JSON bytes into simplejson.Json instance |
| `Wrap(value interface{})`      | Wrap Go value as simplejson.Json object        |
| `List(elements []interface{})` | Convert slice to simplejson.Json objects slice |

### Type-Safe Extraction Functions

| Function | Params | Returns | Description |
|----------|--------|---------|-------------|
| `Extract[T]` | `(object, key)` | `(T, error)` | Strict: key must exist, value must not be null |
| `Solicit[T]` | `(object, key)` | `(T, error)` | Strictest: key must exist, value must not be null / zero (T must be comparable) |
| `Inspect[T]` | `(object, key)` | `(T, error)` | Lenient: returns zero if key missing, errors if null/wrong type |
| `Resolve[T]` | `(object)` | `(T, error)` | Direct: convert Json object to target type T |
| `Inquire[T]` | `(object, key)` | `(T, bool, error)` | Tri-state: returns value + exists flag + error |
| `Attempt[T]` | `(object, key)` | `(T, bool)` | Try: returns (value, true) on success, (zero, false) on any failure |
| `Harvest[T]` | `(object, key, option)` | `(T, error)` | Configurable: use HarvestOption to control null handling and string conversion |

### Advanced Functions

| Function                          | Description                                                             |
| --------------------------------- | ----------------------------------------------------------------------- |
| `GetList(object, key)`            | Get JSON list at key as simplejson.Json slice                           |
| `Explore[T](object, path)`        | Navigate nested JSON via dot-separated path (e.g., "user.profile.name") |
| `Harvest[T](object, key, option)` | Extract mixed-type fields with configurable null/string handling        |
| `Strconv[T](object)`              | Two-stage conversion: JSON → string → target type                       |

### Supported Types

`Resolve` and `Extract` support these types:

- **Primitives**: `int`, `int64`, `uint64`, `float64`, `string`, `bool`
- **Collections**: `[]string`, `[]interface{}`, `[]*simplejson.Json`
- **Complex**: `map[string]interface{}`, `[]byte`, `*simplejson.Json`

### HarvestOption Configuration

`HarvestOption[T]` configures mixed-type field parsing:

| Field     | Type                      | Description                                     |
| --------- | ------------------------- | ----------------------------------------------- |
| `Nilable` | `bool`                    | Accept null values without error (returns zero) |
| `Strconv` | `func(string) (T, error)` | Custom string-to-T conversion function          |

### Using Strconv for String-Based Conversion

**Convert string values to target types via two-stage process:**

```go
object, err := simplejsonx.Load([]byte(`{
	"count": "42",
	"rate": "3.14"
}`))
if err != nil {
	log.Fatalf("Error loading JSON: %v", err)
}

// String "42" → int 42
count, err := simplejsonx.Strconv[int](object.Get("count"))
if err != nil {
	log.Fatalf("Error converting count: %v", err)
}
fmt.Println("Count:", count)  // Output: 42

// String "3.14" → float64 3.14
rate, err := simplejsonx.Strconv[float64](object.Get("rate"))
if err != nil {
	log.Fatalf("Error converting rate: %v", err)
}
fmt.Println("Rate:", rate)  // Output: 3.14
```

### Error Handling Variants (Sub-packages)

The package provides three error handling modes via sub-packages:

| Package                        | Mode | Behavior                                    |
| ------------------------------ | ---- | ------------------------------------------- |
| `simplejsonx/sure/simplejsonm` | Must | Panics on error (use when failure is fatal) |
| `simplejsonx/sure/simplejsono` | Omit | Returns zero value, ignores errors          |
| `simplejsonx/sure/simplejsons` | Soft | Logs errors, returns zero value             |

**Example using Must mode:**

```go
import "github.com/yylego/simplejsonx/sure/simplejsonm"

// Panics if JSON parsing or extraction fails
object := simplejsonm.Load(data)
name := simplejsonm.Extract[string](object, "name")
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/yylego/simplejsonx.svg?variant=adaptive)](https://starchart.cc/yylego/simplejsonx)
