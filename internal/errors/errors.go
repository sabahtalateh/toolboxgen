package errors

import "go/token"

const (
	unexpectedType      = "unexpected type"
	builtinTypeExpected = "builtin type expected"
	structExpected      = "struct expected"
	interfaceExpected   = "interface expected"
	typeDefExpected     = "typedef expected"
	typeAliasExpected   = "type alias expected"
)

func UnexpectedType(pos token.Position) *PositionedErr {
	return Errorf(pos, unexpectedType)
}

func BuiltinTypeExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, builtinTypeExpected)
}

func StructExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, structExpected)
}

func InterfaceExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, interfaceExpected)
}

func TypeDefExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, typeDefExpected)
}

func TypeAliasExpectedErr(pos token.Position) *PositionedErr {
	return Errorf(pos, typeAliasExpected)
}
