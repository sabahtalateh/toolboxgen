package code

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/life4/genesis/slices"

	"github.com/sabahtalateh/toolboxgen/internal/f"
)

func OfNode(n ast.Node) string {
	switch field := n.(type) {
	case *ast.Field:
		if len(field.Names) == 0 {
			return OfNode(field.Type)
		}
		names := slices.Map(field.Names, func(el *ast.Ident) string { return OfNode(el) })
		return fmt.Sprintf("%s %s", strings.Join(names, ", "), OfNode(field.Type))
	case *ast.BlockStmt:
		// to handle empty function bodies within builtin.go
		if field == nil {
			return ""
		}
	}

	bb := new(bytes.Buffer)
	if err := printer.Fprint(bb, token.NewFileSet(), n); err != nil {
		return fmt.Sprintf("<< %s >>", err)
	}
	res := bb.String()
	switch nn := n.(type) {
	case *ast.TypeSpec:
		suffix := ""
		_, isStruct := nn.Type.(*ast.StructType)
		if isStruct {
			suffix = "struct {...}"
		}
		_, isInterface := nn.Type.(*ast.InterfaceType)
		if isInterface {
			suffix = "interface {...}"
		}

		// don't trim body for typedef and type alias
		if isStruct || isInterface {
			body := OfNode(nn.Type)
			res = f.Apply(
				strings.TrimSuffix(res, body),
				strings.TrimSpace,
				func(x string) string {
					if suffix == "" {
						return x
					}
					return fmt.Sprintf("%s %s", x, suffix)
				},
			)
		}
	case *ast.FuncDecl:
		body := OfNode(nn.Body)
		res = f.Apply(strings.TrimSuffix(res, body), strings.TrimSpace)
	}

	return res
}
