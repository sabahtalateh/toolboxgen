package syntax

type kind int

const (
	kUnknown kind = iota + 1
	kType
	kMap
	kChan
	kFuncType
	kStructType
	kInterfaceType
)
