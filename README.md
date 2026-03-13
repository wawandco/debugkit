# DebugKit

A beautiful and powerful debugging toolkit for Go that provides syntax-highlighted output for inspecting complex data structures.

## Features

- 🎨 **Colorized Output** - Syntax highlighting for different data types (structs, strings, numbers, booleans, etc.)
- 🔍 **Deep Inspection** - Recursively prints nested structs, slices, maps, and pointers
- 🔄 **Circular Reference Detection** - Safely handles circular references without infinite loops
- 📍 **Source Location Tracking** - Shows file and line number where debug calls are made
- 🚀 **Zero Configuration** - Just import and use

## Installation

```bash
go get github.com/YanDeLeon/debugkit
```

## Usage

### Basic Dumping

Use `Dump()` to pretty-print any Go value with syntax highlighting:

```go
package main

import "github.com/YanDeLeon/debugkit"

type Person struct {
    Name string
    Age  int
    Tags []string
}

func main() {
    p := Person{
        Name: "Alice",
        Age:  30,
        Tags: []string{"developer", "gopher"},
    }
    
    debugkit.Dump(p)
}
```

**Output:**
```
Person {
  Name: "Alice"
  Age: 30
  Tags: [
    "developer",
    "gopher",
  ]
}
```

### Tracing with Location

Use `Trace()` to print a value along with the file and line number:

```go
func calculateTotal(items []int) int {
    total := 0
    for _, item := range items {
        total += item
    }
    
    debugkit.Trace(total)  // Shows: [main.go:15]
    return total
}
```

### Tracing Multiple Values

Use `TraceAll()` to trace multiple values at once:

```go
func process(x, y int) {
    result := x + y
    debugkit.TraceAll(x, y, result)
    // Shows: [main.go:10]
    // arg0 = 5
    // arg1 = 10
    // arg2 = 15
}
```

## Supported Types

DebugKit handles all Go types with appropriate formatting:

- **Primitives**: `int`, `float`, `bool`, `string`
- **Composite**: `struct`, `slice`, `array`, `map`
- **Pointers**: Dereferences and marks `nil` pointers
- **Complex Structures**: Deeply nested structs with multiple levels
- **Circular References**: Detects and marks circular references as `<circular>`

## Color Scheme

- **Yellow**: Type names and headers
- **Cyan**: Field names and argument labels
- **Green**: String values
- **Blue**: Numeric values
- **Magenta**: Boolean values
- **Red**: `nil` values and circular references
- **White**: Punctuation (braces, brackets, colons)

## API Reference

### `Dump(v any)`

Pretty-prints any Go value with syntax highlighting.

**Parameters:**
- `v any` - The value to dump

**Example:**
```go
debugkit.Dump(myStruct)
```

### `Trace(v any)`

Prints a value along with the source location (file and line number) where it was called.

**Parameters:**
- `v any` - The value to trace

**Example:**
```go
debugkit.Trace(result)
```

### `TraceAll(values ...any)`

Traces multiple values with their source location, labeling each as `arg0`, `arg1`, etc.

**Parameters:**
- `values ...any` - Variadic list of values to trace

**Example:**
```go
debugkit.TraceAll(x, y, z)
```

## Examples

### Nested Structs

```go
type Address struct {
    City  string
    State string
}

type Person struct {
    Name    string
    Age     int
    Address Address
}

p := Person{
    Name: "Bob",
    Age:  25,
    Address: Address{
        City:  "SF",
        State: "CA",
    },
}

debugkit.Dump(p)
```

### Slices and Maps

```go
data := map[string][]int{
    "scores": {95, 87, 92},
    "ages":   {25, 30, 28},
}

debugkit.Dump(data)
```

### Pointers

```go
x := 42
ptr := &x

debugkit.Dump(ptr)  // Shows dereferenced value
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Credits

Built with:
- [fatih/color](https://github.com/fatih/color) - Terminal color output
