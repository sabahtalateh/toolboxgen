package syntax

import (
	"bytes"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"go/ast"
	"go/printer"
	"go/token"
)

func code(n ast.Node, p token.Position) (string, *errors.PositionedErr) {
	bb := new(bytes.Buffer)
	if err := printer.Fprint(bb, token.NewFileSet(), n); err != nil {
		return "", errors.Errorf(p, "failed to get string representation")
	}
	return bb.String(), nil
}
