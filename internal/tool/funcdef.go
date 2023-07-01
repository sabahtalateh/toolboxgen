package tool

import (
	"go/token"
)

type FuncParam struct {
	Name     string
	Type     TypeRef
	Position token.Position
}

type Receiver struct {
	Presented  bool
	Type       TypeRef
	TypeParams []TypeParam
	Position   token.Position
}

type FuncDef struct {
	Code       string
	Package    string
	FuncName   string
	Receiver   Receiver
	TypeParams []TypeParam
	Parameters []FuncParam
	Results    []TypeRef
	Position   token.Position
}
