package convert

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/maps"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"github.com/sabahtalateh/toolboxgen/internal/utils"
)

// Context represents state of convert process
// pakage is currently parsing file package
// imports is currently parsing file imports block
// files is currently parsing package file set
// position is current parsing ast.Node position
// defined is defined type params
type Context struct {
	pakage   string
	imports  []*ast.ImportSpec
	files    *token.FileSet
	position token.Position
	defined  map[string]*types.TypeParam
}

func NewContext(
	Package string,
	imports []*ast.ImportSpec,
	files *token.FileSet,
	position token.Position,
	defined map[string]*types.TypeParam,
) Context {
	return Context{
		pakage:   Package,
		imports:  imports,
		files:    files,
		position: position,
		defined:  defined,
	}
}

func (c Context) WithPackage(Package string) Context {
	return Context{
		pakage:   Package,
		imports:  c.imports,
		files:    c.files,
		position: c.position,
		defined:  c.defined,
	}
}

func (c Context) Package() string {
	return c.pakage
}

func (c Context) WithImports(imports []*ast.ImportSpec) Context {
	return Context{
		pakage:   c.pakage,
		imports:  imports,
		files:    c.files,
		position: c.position,
		defined:  c.defined,
	}
}

func (c Context) Imports() []*ast.ImportSpec {
	return c.imports
}

func (c Context) WithFiles(files *token.FileSet) Context {
	return Context{
		pakage:   c.pakage,
		imports:  c.imports,
		files:    files,
		position: c.position,
		defined:  c.defined,
	}
}

func (c Context) Files() *token.FileSet {
	return c.files
}

func (c Context) WithPos(pos token.Pos) Context {
	return Context{
		pakage:   c.pakage,
		imports:  c.imports,
		files:    c.files,
		position: c.files.Position(pos),
		defined:  c.defined,
	}
}

func (c Context) WithPosition(position token.Position) Context {
	return Context{
		pakage:   c.pakage,
		imports:  c.imports,
		files:    c.files,
		position: position,
		defined:  c.defined,
	}
}

func (c Context) Position() token.Position {
	return c.position
}

func (c Context) NodePosition(n ast.Node) token.Position {
	return c.files.Position(n.Pos())
}

func (c Context) WithDefined(defined types.TypeParams) Context {
	definedMap := maps.FromSlice(defined, func(v *types.TypeParam) string { return v.Original })

	return Context{
		pakage:   c.pakage,
		imports:  c.imports,
		files:    c.files,
		position: c.position,
		defined:  definedMap,
	}
}

func (c Context) Defined(name string) (*types.TypeParam, bool) {
	x, ok := c.defined[name]
	return x, ok
}

func (c Context) ResolvePackage(packageAlias string) string {
	for _, spec := range c.Imports() {
		var alias string
		if spec.Name != nil {
			alias = spec.Name.Name
		} else {
			parts := strings.Split(utils.Unquote(spec.Path.Value), "/")
			alias = parts[len(parts)-1]
		}

		if alias == packageAlias {
			return utils.Unquote(spec.Path.Value)
		}
	}

	return c.pakage
}
