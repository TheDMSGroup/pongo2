package pongo2_test

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
)

func ExampleTemplate_ExecuteValue() {
	// Create a template that evaluates to a number
	tpl := pongo2.Must(pongo2.FromString("{{ value * 2 }}"))

	// Execute and get the raw value
	result, err := tpl.ExecuteValue(pongo2.Context{"value": 21})
	if err != nil {
		panic(err)
	}

	// The result is the actual number, not a string
	fmt.Printf("Result type: %T, value: %v\n", result, result)
	// Output: Result type: int, value: 42
}

func ExampleTemplate_ExecuteRaw() {
	// Create a template with multiple expressions
	tpl := pongo2.Must(pongo2.FromString("{{ name }} {{ age }} {{ active }}"))

	// Execute and get all raw values
	values, err := tpl.ExecuteRaw(pongo2.Context{
		"name":   "Alice",
		"age":    30,
		"active": true,
	})
	if err != nil {
		panic(err)
	}

	// Each value retains its original type
	for i, v := range values {
		fmt.Printf("Value %d: type=%T, value=%v\n", i, v, v)
	}
	// Output:
	// Value 0: type=string, value=Alice
	// Value 1: type=int, value=30
	// Value 2: type=bool, value=true
}

func ExampleTemplate_Execute_comparison() {
	tpl := pongo2.Must(pongo2.FromString("{{ value }}"))
	ctx := pongo2.Context{"value": 42}

	// Traditional Execute - returns string
	strResult, _ := tpl.Execute(ctx)
	fmt.Printf("Execute: type=%T, value=%v\n", strResult, strResult)

	// New ExecuteValue - returns raw value
	rawResult, _ := tpl.ExecuteValue(ctx)
	fmt.Printf("ExecuteValue: type=%T, value=%v\n", rawResult, rawResult)

	// Output:
	// Execute: type=string, value=42
	// ExecuteValue: type=int, value=42
}
