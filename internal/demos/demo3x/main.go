package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/yylego/simplejsonx"
)

func main() {
	// Sample JSON data with mixed-type price field
	// 示例 JSON 数据，包含混合类型的 price 字段
	data := []byte(`{
		"items": [
			{"name": "A", "price": 100.5},
			{"name": "B", "price": "-"},
			{"name": "C", "price": null}
		]
	}`)

	// Load the JSON data
	// 加载 JSON 数据
	object, err := simplejsonx.Load(data)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	// Get items list
	// 获取 items 列表
	items, err := simplejsonx.GetList(object, "items")
	if err != nil {
		log.Fatalf("Error getting items: %v", err)
	}

	// Define harvest option with mixed-type handling
	// 定义收获选项，处理混合类型
	option := &simplejsonx.HarvestOption[float64]{
		Nilable: true, // Accept null values // 接受 null 值
		Strconv: func(s string) (float64, error) {
			if s == "-" || s == "N/A" {
				return 0, nil // Treat "-" as zero // 将 "-" 视为零值
			}
			return strconv.ParseFloat(s, 64)
		},
	}

	// Harvest price from each item
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

		fmt.Printf("Item %s: price = %.2f\n", name, price)
	}
}
