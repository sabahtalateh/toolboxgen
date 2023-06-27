package errors

import "go/token"

const (
	inconsistentTypeParams = "inconsistent type params"
	stringLitExpected      = "string literal expected"
	errorExpected          = "error expected"
	functionRefExpected    = "function reference expected"
)

func InconsistentTypeParamsErr(pos token.Position) *PositionedErr {
	return Errorf(pos, inconsistentTypeParams)
}

func StringLitExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, stringLitExpected)
}

func ErrorExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, errorExpected)
}

func FunctionRefExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, functionRefExpected)
}

func BuiltinFunctionErr(pos token.Position, funcName string) *PositionedErr {
	return Errorf(pos, "`%s` builtin function not allowed", funcName)
}
