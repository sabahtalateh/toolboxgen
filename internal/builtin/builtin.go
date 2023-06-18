package builtin

var builtinTypes = map[string]struct{}{
	"bool":       {},
	"uint8":      {},
	"uint16":     {},
	"uint32":     {},
	"uint64":     {},
	"int8":       {},
	"int16":      {},
	"int32":      {},
	"int64":      {},
	"float32":    {},
	"float64":    {},
	"complex64":  {},
	"complex128": {},
	"string":     {},
	"int":        {},
	"uint":       {},
	"uintptr":    {},
	"byte":       {},
	"rune":       {},
	"any":        {},
	"error":      {},
}

func Type(t string) bool {
	_, ok := builtinTypes[t]
	return ok
}

var builtinFuncs = map[string]struct{}{
	"append":  {},
	"copy":    {},
	"delete":  {},
	"len":     {},
	"cap":     {},
	"make":    {},
	"new":     {},
	"complex": {},
	"real":    {},
	"imag":    {},
	"close":   {},
	"panic":   {},
	"recover": {},
	"print":   {},
	"println": {},
}

func Func(t string) bool {
	_, ok := builtinFuncs[t]
	return ok
}
