package strings

import (
	"strings"
)

func Unquote(x string) string {
	x = strings.Trim(x, "\"")
	x = strings.Trim(x, "`")
	x = strings.Trim(x, "'")

	return x
}
