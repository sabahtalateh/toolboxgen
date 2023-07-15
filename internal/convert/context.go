package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/utils"
)

// Context represents state of convert process
// Package is currently parsing file package
// Imports is currently parsing file imports block
// Files is currently parsing package file set
// Pos is current parsing structure position
// Defined is actual type params values
type Context struct {
	Package string
	Imports []*ast.ImportSpec
	Files   *token.FileSet
	Pos     token.Pos
	Defined map[string]*types.TypeParam
}

func NewContext() Context {
	return Context{}
}

func (c Context) WithPackage(Package string) Context {
	return Context{
		Package: Package,
		Imports: c.Imports,
		Files:   c.Files,
		Pos:     c.Pos,
		Defined: c.Defined,
	}
}

func (c Context) WithImports(imports []*ast.ImportSpec) Context {
	return Context{
		Package: c.Package,
		Imports: imports,
		Files:   c.Files,
		Pos:     c.Pos,
		Defined: c.Defined,
	}
}

func (c Context) WithFiles(files *token.FileSet) Context {
	return Context{
		Package: c.Package,
		Imports: c.Imports,
		Files:   files,
		Pos:     c.Pos,
		Defined: c.Defined,
	}
}

func (c Context) WithPos(pos token.Pos) Context {
	return Context{
		Package: c.Package,
		Imports: c.Imports,
		Files:   c.Files,
		Pos:     pos,
		Defined: c.Defined,
	}
}

func (c Context) WithDefined(defined map[string]*types.TypeParam) Context {
	return Context{
		Package: c.Package,
		Imports: c.Imports,
		Files:   c.Files,
		Pos:     c.Pos,
		Defined: defined,
	}
}

func (c Context) Position(pos token.Pos) token.Position {
	return c.Files.Position(pos)
}

func (c Context) CurrentPosition() token.Position {
	return c.Position(c.Pos)
}

func (c Context) PackageFromAlias(packageAlias string) string {
	for _, spec := range c.Imports {
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

	return c.Package
}
