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
	Declared string
}

type Function struct {
	Declared string
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
	f, err := parser.ParseFile(files, "", builtin[latest], 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.TypeSpec:
			b.Types[n.Name.Name] = &types.Builtin{Declared: code.OfNode(n), TypeName: n.Name.Name}
			return false
		case *ast.FuncDecl:
			b.Functions[n.Name.Name] = Function{Declared: code.OfNode(n)}
			return false
		}
		return true
	})
}
