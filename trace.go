package debugkit

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/fatih/color"
)

func Trace(v any) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	file = filepath.Base(file)

	colorHeader := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("\n%s\n", colorHeader(fmt.Sprintf("[%s:%d]", file, line)))

	val := reflect.ValueOf(v)
	visited := make(map[uintptr]bool)

	printValue(val, 0, visited)

	fmt.Println()
}

func TraceAll(values ...any) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	colorHeader := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("\n%s\n", colorHeader(fmt.Sprintf("[%s:%d]", file, line)))

	for i, v := range values {
		fmt.Printf("%s = ", colorField(fmt.Sprintf("arg%d", i)))

		val := reflect.ValueOf(v)
		visited := make(map[uintptr]bool)

		printValue(val, 0, visited)

		fmt.Println()
	}

	fmt.Println()
}
