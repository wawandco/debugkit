package debugkit

import (
	"fmt"
	"reflect"
)

func Dump(v any) {
	val := reflect.ValueOf(v)
	visited := make(map[uintptr]bool)

	indent := 1

	printValue(val, indent, visited)
	fmt.Println()
}
