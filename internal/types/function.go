package types

import (
	"go/token"
)

type (
	Function struct {
		Declared   string
		Package    string
		FuncName   string
		Receiver   *Field
		TypeParams TypeParams
		Parameters []*Field
		Results    []*Field
		Position   token.Position
	}
)
