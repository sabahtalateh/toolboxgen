package syntax

type kind int

const (
	kUnexpected kind = iota + 1
	kType
	kMap
	kChan
	kFuncType
	kStructType
	kInterfaceType
)
