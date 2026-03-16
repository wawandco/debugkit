package debugkit

import (
	"fmt"
	"reflect"
)

func printValue(v reflect.Value, indent int, visited map[uintptr]bool) {
	if !v.IsValid() {
		fmt.Print(colorNil("nil"))
		return
	}

	switch v.Kind() {

	case reflect.Pointer:
		handlePointer(v, indent, visited)

	case reflect.Struct:
		printStruct(v, indent, visited)

	case reflect.Map:
		printMap(v, indent, visited)

	case reflect.Slice, reflect.Array:
		printSlice(v, indent, visited)

	case reflect.String:
		fmt.Printf("%s", colorString(fmt.Sprintf("%q", v.String())))

	case reflect.Bool:
		fmt.Print(colorBool(fmt.Sprintf("%v", v.Bool())))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Print(colorNumber(fmt.Sprintf("%v", v.Int())))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Print(colorNumber(fmt.Sprintf("%v", v.Uint())))

	case reflect.Float32, reflect.Float64:
		fmt.Print(colorNumber(fmt.Sprintf("%v", v.Float())))

	default:
		if v.CanInterface() {
			fmt.Printf("%v", v.Interface())
		} else {
			fmt.Print(colorNil("<unexported>"))
		}
	}
}

func printStruct(v reflect.Value, indent int, visited map[uintptr]bool) {
	t := v.Type()

	fmt.Printf("%s %s\n", colorType(t.Name()), colorPunctuation("{"))

	for i := 0; i < v.NumField(); i++ {

		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("%s%s: ", indentStr(indent+1), colorField(field.Name))

		printValue(value, indent+1, visited)

		fmt.Println()
	}

	fmt.Printf("%s%s", indentStr(indent), colorPunctuation("}"))
}

func printSlice(v reflect.Value, indent int, visited map[uintptr]bool) {
	fmt.Println(colorPunctuation("["))

	for i := 0; i < v.Len(); i++ {

		fmt.Printf("%s", indentStr(indent+1))

		printValue(v.Index(i), indent+1, visited)

		fmt.Println(colorPunctuation(","))
	}

	fmt.Printf("%s%s", indentStr(indent), colorPunctuation("]"))
}

func printMap(v reflect.Value, indent int, visited map[uintptr]bool) {
	fmt.Println(colorPunctuation("{"))

	for _, key := range v.MapKeys() {
		var keyStr string
		if key.CanInterface() {
			keyStr = fmt.Sprintf("%v", key.Interface())
		} else {
			keyStr = key.String()
		}

		fmt.Printf("%s%s: ", indentStr(indent+1), colorField(keyStr))

		printValue(v.MapIndex(key), indent+1, visited)

		fmt.Println()
	}

	fmt.Printf("%s%s", indentStr(indent), colorPunctuation("}"))
}

func handlePointer(v reflect.Value, indent int, visited map[uintptr]bool) {
	if v.IsNil() {
		fmt.Print(colorNil("nil"))
		return
	}

	ptr := v.Pointer()

	if visited[ptr] {
		fmt.Print(colorNil("<circular>"))
		return
	}

	visited[ptr] = true

	printValue(v.Elem(), indent, visited)
}
