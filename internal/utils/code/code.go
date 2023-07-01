package code

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/f"
)

func OfNodeE(n ast.Node, p token.Position) (string, *errors.PositionedErr) {
	bb := new(bytes.Buffer)
	if err := printer.Fprint(bb, token.NewFileSet(), n); err != nil {
		return "", errors.Errorf(p, "failed to get string representation")
	}
	res := bb.String()
	switch nn := n.(type) {
	case *ast.TypeSpec:
		body, err := OfNodeE(nn.Type, p)
		if err == nil {
			res = f.Apply(strings.TrimSuffix(res, body), strings.TrimSpace)
		}
	case *ast.FuncDecl:
		body, err := OfNodeE(nn.Body, p)
		if err == nil {
			res = f.Apply(strings.TrimSuffix(res, body), strings.TrimSpace)
		}
	}

	return res, nil
}

func OfNode(n ast.Node) string {
	r, _ := OfNodeE(n, token.Position{})
	return r
}
