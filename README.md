[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/simplejsonx/release.yml?branch=main&label=BUILD)](https://github.com/yylego/simplejsonx/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/simplejsonx)](https://pkg.go.dev/github.com/yylego/simplejsonx)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/simplejsonx/main.svg)](https://coveralls.io/github/yylego/simplejsonx?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.24--1.25-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yylego/simplejsonx.svg)](https://github.com/yylego/simplejsonx/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/simplejsonx)](https://goreportcard.com/report/github.com/yylego/simplejsonx)

# simplejsonx

`simplejsonx` is a generic-based JSON parsing package that depends on `github.com/bitly/go-simplejson`, enabling type-safe JSON processing. It requires at least Go 1.22 (which supports generics).

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

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

| Function                  | Description                                                                |
| ------------------------- | -------------------------------------------------------------------------- |
| `Extract[T](object, key)` | Extract value at key with strict validation (returns error if key missing) |
| `Inspect[T](object, key)` | Extract value at key with lenient mode (returns zero value if key missing) |
| `Resolve[T](object)`      | Convert simplejson.Json object to target type T                            |
| `Inquire[T](object, key)` | Tri-state extraction returning (value, exists, error)                      |
| `Attempt[T](object, key)` | Simple extraction returning (value, ok) without error details              |

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
