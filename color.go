package debugkit

import "github.com/fatih/color"

var (
	colorType        = color.New(color.FgYellow).SprintFunc()
	colorField       = color.New(color.FgCyan).SprintFunc()
	colorString      = color.New(color.FgGreen).SprintFunc()
	colorNumber      = color.New(color.FgBlue).SprintFunc()
	colorBool        = color.New(color.FgMagenta).SprintFunc()
	colorNil         = color.New(color.FgRed).SprintFunc()
	colorPunctuation = color.New(color.FgWhite).SprintFunc()
)
