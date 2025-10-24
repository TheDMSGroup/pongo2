# Raw Value Execution in Pongo2

This document describes the new raw value execution feature that allows templates to return their internal values without converting them to strings.

## Problem

Previously, the `Execute` method always returned a string, even when the template evaluated to other types like numbers, booleans, or complex objects. This forced unnecessary type conversions and lost type information.

## Solution

Three new methods have been added to work with raw values:

### 1. ExecuteValue

Returns the first evaluated value from a template without string conversion.

```go
tpl := pongo2.Must(pongo2.FromString("{{ value * 2 }}"))
result, err := tpl.ExecuteValue(pongo2.Context{"value": 21})
// result is int(42), not "42"
```

### 2. ExecuteRaw

Returns all evaluated values from a template as a slice.

```go
tpl := pongo2.Must(pongo2.FromString("{{ name }} {{ age }} {{ active }}"))
values, err := tpl.ExecuteRaw(pongo2.Context{
    "name":   "Alice",
    "age":    30,
    "active": true,
})
// values[0] is string("Alice")
// values[1] is int(30)  
// values[2] is bool(true)
```

### 3. Convenience Functions

Package-level functions for quick evaluation:

```go
// Single value
result, err := pongo2.ExecuteTemplateValue("{{ value * 2 }}", pongo2.Context{"value": 21})

// Multiple values
values, err := pongo2.ExecuteTemplateRaw("{{ a }} {{ b }}", pongo2.Context{"a": 1, "b": 2})
```

## TemplateSet Methods

The TemplateSet also provides corresponding methods:

```go
set := pongo2.NewSet("mySet", pongo2.DefaultLoader)

// From string
value, err := set.RenderTemplateStringValue("{{ value }}", ctx)

// From bytes
value, err := set.RenderTemplateBytesValue([]byte("{{ value }}"), ctx)

// From file
value, err := set.RenderTemplateFileValue("template.html", ctx)
```

## Use Cases

1. **Mathematical Calculations**: Get numeric results directly without parsing strings
2. **Boolean Logic**: Evaluate conditions and get actual boolean values
3. **Data Processing**: Process templates that generate structured data
4. **Type Preservation**: Maintain original types when passing data between systems
5. **Performance**: Avoid unnecessary string conversions for non-string data

## Compatibility

These new methods are fully backward compatible. The original `Execute` method continues to work exactly as before, returning strings.

## Implementation Details

- HTML nodes are skipped when collecting raw values (they don't have evaluatable content)
- Only nodes implementing the `IEvaluator` interface contribute values
- If no evaluatable nodes are found, `ExecuteValue` falls back to string execution
- The methods preserve the original types from the context and expressions