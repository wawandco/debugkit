package debugkit

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestMain(m *testing.M) {
	color.NoColor = true
	os.Exit(m.Run())
}

// captureOutput redirects stdout and captures whatever is printed during fn().
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// --- indentStr tests ---

func TestIndentStr(t *testing.T) {
	tests := []struct {
		level int
		want  string
	}{
		{0, ""},
		{1, "  "},
		{2, "    "},
		{3, "      "},
	}
	for _, tt := range tests {
		got := indentStr(tt.level)
		if got != tt.want {
			t.Errorf("indentStr(%d) = %q, want %q", tt.level, got, tt.want)
		}
	}
}

// --- Dump tests ---

func TestDumpString(t *testing.T) {
	out := captureOutput(func() { Dump("hello") })
	want := "\"hello\"\n"
	if out != want {
		t.Errorf("Dump(string) = %q, want %q", out, want)
	}
}

func TestDumpInt(t *testing.T) {
	out := captureOutput(func() { Dump(42) })
	want := "42\n"
	if out != want {
		t.Errorf("Dump(int) = %q, want %q", out, want)
	}
}

func TestDumpFloat(t *testing.T) {
	out := captureOutput(func() { Dump(3.14) })
	want := "3.14\n"
	if out != want {
		t.Errorf("Dump(float) = %q, want %q", out, want)
	}
}

func TestDumpBool(t *testing.T) {
	out := captureOutput(func() { Dump(true) })
	want := "true\n"
	if out != want {
		t.Errorf("Dump(bool) = %q, want %q", out, want)
	}
}

func TestDumpNil(t *testing.T) {
	out := captureOutput(func() { Dump(nil) })
	want := "nil\n"
	if out != want {
		t.Errorf("Dump(nil) = %q, want %q", out, want)
	}
}

func TestDumpStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Alice", Age: 30}
	out := captureOutput(func() { Dump(p) })

	if !strings.Contains(out, "Person {") {
		t.Errorf("expected struct header, got %q", out)
	}
	if !strings.Contains(out, "Name: \"Alice\"") {
		t.Errorf("expected Name field, got %q", out)
	}
	if !strings.Contains(out, "Age: 30") {
		t.Errorf("expected Age field, got %q", out)
	}
}

func TestDumpNestedStruct(t *testing.T) {
	type Address struct {
		City string
	}
	type Person struct {
		Name    string
		Address Address
	}
	p := Person{Name: "Bob", Address: Address{City: "NYC"}}
	out := captureOutput(func() { Dump(p) })

	if !strings.Contains(out, "Person {") {
		t.Errorf("expected outer struct, got %q", out)
	}
	if !strings.Contains(out, "Address {") {
		t.Errorf("expected nested struct, got %q", out)
	}
	if !strings.Contains(out, "City: \"NYC\"") {
		t.Errorf("expected City field, got %q", out)
	}
}

func TestDumpSlice(t *testing.T) {
	out := captureOutput(func() { Dump([]int{1, 2, 3}) })

	if !strings.Contains(out, "[") {
		t.Errorf("expected opening bracket, got %q", out)
	}
	if !strings.Contains(out, "1,") {
		t.Errorf("expected element 1, got %q", out)
	}
	if !strings.Contains(out, "2,") {
		t.Errorf("expected element 2, got %q", out)
	}
	if !strings.Contains(out, "3,") {
		t.Errorf("expected element 3, got %q", out)
	}
	if !strings.Contains(out, "]") {
		t.Errorf("expected closing bracket, got %q", out)
	}
}

func TestDumpEmptySlice(t *testing.T) {
	out := captureOutput(func() { Dump([]int{}) })

	if !strings.Contains(out, "[") || !strings.Contains(out, "]") {
		t.Errorf("expected brackets for empty slice, got %q", out)
	}
}

func TestDumpArray(t *testing.T) {
	out := captureOutput(func() { Dump([2]string{"a", "b"}) })

	if !strings.Contains(out, `"a"`) {
		t.Errorf("expected element a, got %q", out)
	}
	if !strings.Contains(out, `"b"`) {
		t.Errorf("expected element b, got %q", out)
	}
}

func TestDumpMap(t *testing.T) {
	// Use a single-entry map to avoid key ordering issues.
	out := captureOutput(func() { Dump(map[string]int{"x": 1}) })

	if !strings.Contains(out, "{") {
		t.Errorf("expected opening brace, got %q", out)
	}
	if !strings.Contains(out, "x: 1") {
		t.Errorf("expected key-value pair, got %q", out)
	}
	if !strings.Contains(out, "}") {
		t.Errorf("expected closing brace, got %q", out)
	}
}

func TestDumpEmptyMap(t *testing.T) {
	out := captureOutput(func() { Dump(map[string]int{}) })

	if !strings.Contains(out, "{") || !strings.Contains(out, "}") {
		t.Errorf("expected braces for empty map, got %q", out)
	}
}

func TestDumpPointer(t *testing.T) {
	s := "hello"
	out := captureOutput(func() { Dump(&s) })
	want := "\"hello\"\n"
	if out != want {
		t.Errorf("Dump(&string) = %q, want %q", out, want)
	}
}

func TestDumpNilPointer(t *testing.T) {
	var p *int
	out := captureOutput(func() { Dump(p) })
	want := "nil\n"
	if out != want {
		t.Errorf("Dump(nil pointer) = %q, want %q", out, want)
	}
}

func TestDumpCircularReference(t *testing.T) {
	type Node struct {
		Value int
		Next  *Node
	}
	a := &Node{Value: 1}
	b := &Node{Value: 2}
	a.Next = b
	b.Next = a // circular

	out := captureOutput(func() { Dump(a) })

	if !strings.Contains(out, "<circular>") {
		t.Errorf("expected <circular> for circular reference, got %q", out)
	}
	if !strings.Contains(out, "Value: 1") {
		t.Errorf("expected Value: 1, got %q", out)
	}
	if !strings.Contains(out, "Value: 2") {
		t.Errorf("expected Value: 2, got %q", out)
	}
}

func TestDumpSliceOfStructs(t *testing.T) {
	type Item struct {
		ID int
	}
	items := []Item{{ID: 1}, {ID: 2}}
	out := captureOutput(func() { Dump(items) })

	if !strings.Contains(out, "Item {") {
		t.Errorf("expected struct inside slice, got %q", out)
	}
	if !strings.Contains(out, "ID: 1") {
		t.Errorf("expected ID: 1, got %q", out)
	}
	if !strings.Contains(out, "ID: 2") {
		t.Errorf("expected ID: 2, got %q", out)
	}
}

// --- printValue tests for edge cases ---

func TestPrintValueInvalid(t *testing.T) {
	out := captureOutput(func() {
		visited := make(map[uintptr]bool)
		printValue(reflect.Value{}, 0, visited)
	})
	if out != "nil" {
		t.Errorf("printValue(invalid) = %q, want %q", out, "nil")
	}
}

func TestHandlePointerVisitedFlag(t *testing.T) {
	x := 42
	p := &x

	out := captureOutput(func() {
		visited := make(map[uintptr]bool)
		v := reflect.ValueOf(p)
		// First call should print the value and mark as visited.
		handlePointer(v, 0, visited)
		fmt.Print("|")
		// Second call should detect circular.
		handlePointer(v, 0, visited)
	})

	parts := strings.Split(out, "|")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d: %q", len(parts), out)
	}
	if parts[0] != "42" {
		t.Errorf("first call = %q, want %q", parts[0], "42")
	}
	if parts[1] != "<circular>" {
		t.Errorf("second call = %q, want %q", parts[1], "<circular>")
	}
}

// --- Trace tests ---

func TestTraceInt(t *testing.T) {
	out := captureOutput(func() { Trace(99) })

	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if !strings.Contains(out, "99") {
		t.Errorf("expected value 99, got %q", out)
	}
}

func TestTraceString(t *testing.T) {
	out := captureOutput(func() { Trace("hello") })

	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if !strings.Contains(out, `"hello"`) {
		t.Errorf("expected quoted string, got %q", out)
	}
}

func TestTraceNil(t *testing.T) {
	out := captureOutput(func() { Trace(nil) })

	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if !strings.Contains(out, "nil") {
		t.Errorf("expected nil, got %q", out)
	}
}

func TestTraceStruct(t *testing.T) {
	type Item struct {
		Name string
	}
	out := captureOutput(func() { Trace(Item{Name: "x"}) })

	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if !strings.Contains(out, "Item {") {
		t.Errorf("expected struct output, got %q", out)
	}
	if !strings.Contains(out, `Name: "x"`) {
		t.Errorf("expected Name field, got %q", out)
	}
}

func TestTraceHeaderFormat(t *testing.T) {
	out := captureOutput(func() { Trace(1) })

	// Should start with newline, then [file:line]
	if !strings.HasPrefix(out, "\n[") {
		t.Errorf("expected output to start with newline and bracket, got %q", out)
	}
	if !strings.Contains(out, "]") {
		t.Errorf("expected closing bracket in header, got %q", out)
	}
}

// --- TraceAll tests ---

func TestTraceAllSingleValue(t *testing.T) {
	out := captureOutput(func() { TraceAll(42) })

	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if !strings.Contains(out, "arg0 = 42") {
		t.Errorf("expected arg0 = 42, got %q", out)
	}
}

func TestTraceAllMultipleValues(t *testing.T) {
	out := captureOutput(func() { TraceAll("a", 2, true) })

	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if !strings.Contains(out, `arg0 = "a"`) {
		t.Errorf("expected arg0, got %q", out)
	}
	if !strings.Contains(out, "arg1 = 2") {
		t.Errorf("expected arg1, got %q", out)
	}
	if !strings.Contains(out, "arg2 = true") {
		t.Errorf("expected arg2, got %q", out)
	}
}

func TestTraceAllNoValues(t *testing.T) {
	out := captureOutput(func() { TraceAll() })

	// Should still print the header, just no arg lines.
	if !strings.Contains(out, "debug_test.go:") {
		t.Errorf("expected file:line header, got %q", out)
	}
	if strings.Contains(out, "arg0") {
		t.Errorf("expected no args, got %q", out)
	}
}

func TestTraceAllWithNil(t *testing.T) {
	out := captureOutput(func() { TraceAll(nil, 5) })

	if !strings.Contains(out, "arg0 = nil") {
		t.Errorf("expected arg0 = nil, got %q", out)
	}
	if !strings.Contains(out, "arg1 = 5") {
		t.Errorf("expected arg1 = 5, got %q", out)
	}
}

func TestTraceAllWithStruct(t *testing.T) {
	type Coord struct {
		X int
		Y int
	}
	out := captureOutput(func() { TraceAll(Coord{X: 1, Y: 2}, "label") })

	if !strings.Contains(out, "Coord {") {
		t.Errorf("expected struct in output, got %q", out)
	}
	if !strings.Contains(out, "X: 1") {
		t.Errorf("expected X field, got %q", out)
	}
	if !strings.Contains(out, `arg1 = "label"`) {
		t.Errorf("expected arg1, got %q", out)
	}
}

func TestDumpStructWithUnexportedFields(t *testing.T) {
	type secret struct {
		hidden  string
		count   int
		flag    bool
		ratio   float64
		Visible string
	}

	s := secret{
		hidden:  "private",
		count:   42,
		flag:    true,
		ratio:   3.14,
		Visible: "public",
	}

	// Should not panic on unexported fields.
	out := captureOutput(func() { Dump(s) })

	if !strings.Contains(out, "secret {") {
		t.Errorf("expected struct header, got %q", out)
	}
	if !strings.Contains(out, `Visible: "public"`) {
		t.Errorf("expected exported field, got %q", out)
	}
	if !strings.Contains(out, "hidden:") {
		t.Errorf("expected unexported field name, got %q", out)
	}
	if !strings.Contains(out, "count:") {
		t.Errorf("expected unexported int field name, got %q", out)
	}
	if !strings.Contains(out, "flag:") {
		t.Errorf("expected unexported bool field name, got %q", out)
	}
	if !strings.Contains(out, "ratio:") {
		t.Errorf("expected unexported float field name, got %q", out)
	}
}

func TestDumpTimeValue(t *testing.T) {
	ts := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	out := captureOutput(func() { Dump(ts) })

	if !strings.Contains(out, "2025-06-15") {
		t.Errorf("expected human-readable time, got %q", out)
	}
	// Should NOT show internal fields
	if strings.Contains(out, "wall:") || strings.Contains(out, "ext:") {
		t.Errorf("should not show internal time fields, got %q", out)
	}
}

func TestDumpStructWithTimeField(t *testing.T) {
	type Event struct {
		Name string
		At   time.Time
	}
	e := Event{Name: "launch", At: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	out := captureOutput(func() { Dump(e) })

	if !strings.Contains(out, "Event {") {
		t.Errorf("expected struct header, got %q", out)
	}
	if !strings.Contains(out, `Name: "launch"`) {
		t.Errorf("expected Name field, got %q", out)
	}
	if !strings.Contains(out, "2025-01-01") {
		t.Errorf("expected human-readable time in At field, got %q", out)
	}
}
