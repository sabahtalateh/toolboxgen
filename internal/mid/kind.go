package mid

type kind int

const (
	kUnknown kind = iota + 1

	kSelector
	kCall

	kType
	kMap
	kChan
	kFuncType
	kStructType
	kInterfaceType
)
