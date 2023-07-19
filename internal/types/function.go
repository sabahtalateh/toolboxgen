package types

import (
	"go/token"
)

type (
	Function struct {
		Package    string
		FuncName   string
		Receiver   *Field
		TypeParams TypeParams
		Parameters Fields
		Results    Fields
		Position   token.Position
		Declared   string
	}
)
