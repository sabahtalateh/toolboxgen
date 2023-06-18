package errors

import "go/token"

const (
	inconsistentTypeParams = "inconsistent type params"
	stringExpected         = "string expected"
	errorExpected          = "error expected"
	functionExpected       = "function expected"
)

func InconsistentTypeParamsErr(pos token.Position) *PositionedErr {
	return Errorf(pos, inconsistentTypeParams)
}

func StringExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, stringExpected)
}

func ErrorExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, errorExpected)
}

func FunctionExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, functionExpected)
}

func BuiltinFunctionErr(pos token.Position, funcName string) *PositionedErr {
	return Errorf(pos, "`%s` builtin function not allowed", funcName)
}
