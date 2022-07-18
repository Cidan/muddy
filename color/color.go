package color

import "strings"

const (
	reset         = "\u001b[0m"
	black         = "\u001b[30m"
	red           = "\u001b[31m"
	green         = "\u001b[32m"
	yellow        = "\u001b[33m"
	blue          = "\u001b[34m"
	magenta       = "\u001b[35m"
	cyan          = "\u001b[36m"
	white         = "\u001b[37m"
	brightBlack   = "\u001b[30;1m"
	brightRed     = "\u001b[31;1m"
	brightGreen   = "\u001b[32;1m"
	brightYellow  = "\u001b[33;1m"
	brightBlue    = "\u001b[34;1m"
	brightMagenta = "\u001b[35;1m"
	brightCyan    = "\u001b[36;1m"
	brightWhite   = "\u001b[37;1m"
)

// Parse will convert color codes into ANSI colors. Color strings
// will always be terminated with a reset marker to avoid bleeding.
func Parse(input string) string {
	str := input
	str = strings.ReplaceAll(str, "{x", reset)
	str = strings.ReplaceAll(str, "{k", black)
	str = strings.ReplaceAll(str, "{r", red)
	str = strings.ReplaceAll(str, "{g", green)
	str = strings.ReplaceAll(str, "{y", yellow)
	str = strings.ReplaceAll(str, "{b", blue)
	str = strings.ReplaceAll(str, "{m", magenta)
	str = strings.ReplaceAll(str, "{c", cyan)
	str = strings.ReplaceAll(str, "{w", white)

	str = strings.ReplaceAll(str, "{K", brightBlack)
	str = strings.ReplaceAll(str, "{R", brightRed)
	str = strings.ReplaceAll(str, "{G", brightGreen)
	str = strings.ReplaceAll(str, "{Y", brightYellow)
	str = strings.ReplaceAll(str, "{B", brightBlue)
	str = strings.ReplaceAll(str, "{M", brightMagenta)
	str = strings.ReplaceAll(str, "{C", brightCyan)
	str = strings.ReplaceAll(str, "{W", brightWhite)
	return str + Reset()
}

// Strip will remove color codes from a string.
func Strip(input string) string {
	str := input
	str = strings.ReplaceAll(str, "{x", "")
	str = strings.ReplaceAll(str, "{k", "")
	str = strings.ReplaceAll(str, "{r", "")
	str = strings.ReplaceAll(str, "{g", "")
	str = strings.ReplaceAll(str, "{y", "")
	str = strings.ReplaceAll(str, "{b", "")
	str = strings.ReplaceAll(str, "{m", "")
	str = strings.ReplaceAll(str, "{c", "")
	str = strings.ReplaceAll(str, "{w", "")

	str = strings.ReplaceAll(str, "{K", "")
	str = strings.ReplaceAll(str, "{R", "")
	str = strings.ReplaceAll(str, "{G", "")
	str = strings.ReplaceAll(str, "{Y", "")
	str = strings.ReplaceAll(str, "{B", "")
	str = strings.ReplaceAll(str, "{M", "")
	str = strings.ReplaceAll(str, "{C", "")
	str = strings.ReplaceAll(str, "{W", "")
	return str
}

// Reset returns the color reset string.
func Reset() string {
	return reset
}
