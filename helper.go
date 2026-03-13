package debugkit

import "strings"

func indentStr(level int) string {
	return strings.Repeat("  ", level)
}
