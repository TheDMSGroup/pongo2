//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
)

func main() {
	// Test 1: Simple value extraction
	tpl1, err := pongo2.FromString("{{ value * 2 }}")
	if err != nil {
		panic(err)
	}

	result1, err := tpl1.ExecuteValue(pongo2.Context{"value": 21})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Test 1 - ExecuteValue: type=%T, value=%v\n", result1, result1)

	// Compare with traditional Execute
	strResult1, err := tpl1.Execute(pongo2.Context{"value": 21})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Test 1 - Execute: type=%T, value=%v\n", strResult1, strResult1)

	// Test 2: Multiple values
	tpl2, err := pongo2.FromString("{{ name }} {{ age }} {{ active }}")
	if err != nil {
		panic(err)
	}

	values, err := tpl2.ExecuteRaw(pongo2.Context{
		"name":   "Alice",
		"age":    30,
		"active": true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("\nTest 2 - ExecuteRaw (multiple values):")
	for i, v := range values {
		fmt.Printf("  Value %d: type=%T, value=%v\n", i, v, v)
	}

	// Test 3: Complex expression
	tpl3, err := pongo2.FromString("{{ items|length }}")
	if err != nil {
		panic(err)
	}

	result3, err := tpl3.ExecuteValue(pongo2.Context{
		"items": []string{"a", "b", "c"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nTest 3 - Filter result: type=%T, value=%v\n", result3, result3)

	// Test 4: Boolean value
	tpl4, err := pongo2.FromString("{{ value > 10 }}")
	if err != nil {
		panic(err)
	}

	result4, err := tpl4.ExecuteValue(pongo2.Context{"value": 15})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nTest 4 - Boolean: type=%T, value=%v\n", result4, result4)
}
