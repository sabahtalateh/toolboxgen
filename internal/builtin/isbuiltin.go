package builtin

import (
	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"go/ast"
	"go/parser"
	"go/token"
)

type Builtin struct {
	Types map[string]*types.Builtin
	// TODO parse function with standard mechanism
	Functions map[string]Function
}

type Type struct {
	Code string
}

type Function struct {
	Code string
}

func Init() *Builtin {
	r := &Builtin{
		Types:     map[string]*types.Builtin{},
		Functions: map[string]Function{},
	}
	r.init()

	return r
}

func (b *Builtin) init() {
	files := token.NewFileSet()
	f, err := parser.ParseFile(files, "", builtingo, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.TypeSpec:
			b.Types[n.Name.Name] = &types.Builtin{TypeName: n.Name.Name, Code: code.OfNode(n)}
			return false
		case *ast.FuncDecl:
			b.Functions[n.Name.Name] = Function{Code: code.OfNode(n)}
			return false
		}
		return true
	})
}
