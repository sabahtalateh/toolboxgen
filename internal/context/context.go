package context

import (
	"go/ast"
	"go/token"
)

type Context struct {
	Package string
	Imports []*ast.ImportSpec
	Files   *token.FileSet
}

func New(Package string, imports []*ast.ImportSpec, files *token.FileSet) Context {
	return Context{Package: Package, Imports: imports, Files: files}
}

func (c Context) WithPackage(Package string) Context {
	return Context{
		Package: Package,
		Imports: c.Imports,
		Files:   c.Files,
	}
}

func (c Context) WithImports(imports []*ast.ImportSpec) Context {
	return Context{
		Package: c.Package,
		Imports: imports,
		Files:   c.Files,
	}
}

func (c Context) WithFiles(files *token.FileSet) Context {
	return Context{
		Package: c.Package,
		Imports: c.Imports,
		Files:   files,
	}
}

func (c Context) Position(pos token.Pos) token.Position {
	return c.Files.Position(pos)
}
