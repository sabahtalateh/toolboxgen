package discovery

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/convert"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func Discover(modRoot string) (tt *Tools, err error) {
	d := new(discovery)

	if d.mod, err = mod.LookupDir(modRoot, true); err != nil {
		return nil, err
	}

	if d.converter, err = convert.New(d.mod); err != nil {
		return nil, err
	}

	if err = d.discoverDir(modRoot); err != nil {
		return nil, err
	}

	return d.result, nil
}

type discovery struct {
	mod       *mod.Module
	converter *convert.Converter
	result    *Tools
}

func (d *discovery) discoverDir(dir string) error {
	var err error

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip discovery at vendors
		if strings.HasPrefix(path, filepath.Join(d.mod.Dir, "vendor")) {
			return err
		}
		if !info.IsDir() {
			return err
		}

		files := token.NewFileSet()
		pkgs, err := parser.ParseDir(files, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, pkg := range pkgs {
			for filePath, file := range pkg.Files {
				Package := filepath.Dir(filePath)
				if strings.HasPrefix(Package, d.mod.Dir) {
					Package = strings.Replace(Package, d.mod.Dir, d.mod.Path, 1)
				} else {
					return errors.New("impossibru")
				}

				err = d.discoverFile(file, Package, files)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (d *discovery) discoverFile(file *ast.File, Package string, files *token.FileSet) error {
	var (
		insideFunction *ast.FuncDecl // current top level function
		err            error
	)

	ast.Inspect(file, func(node ast.Node) bool {
		if err != nil {
			return false
		}

		// reset insideFunction
		if node != nil && insideFunction != nil && node.Pos() >= insideFunction.End() {
			insideFunction = nil
		}

		ctx := convert.NewContext().WithPackage(Package).WithImports(file.Imports).WithFiles(files)
		switch n := node.(type) {
		case *ast.FuncDecl:
			if insideFunction == nil {
				insideFunction = n
			}

			var fd *types.Function
			fd, err = d.converter.Function(ctx.WithPos(n.Pos()), n)
			if err != nil {
				return false
			}
			println(fd)
		case *ast.TypeSpec:
			var t types.Type
			t, err = d.converter.Type(ctx.WithPos(n.Pos()), n)
			if err != nil {
				return false
			}
			println(t)
		case *ast.CallExpr:
			d.converter.Call(ctx.WithPos(n.Pos()), n)

			if insideFunction == nil || isInit(insideFunction) {

				// calls := parse.ParseFuncCalls(n, ctx.files)
				// for _, call := range calls {
				// 	if err = call.error(); err != nil {
				// 		return false
				// 	}
				// }
				//
				// tool, err := d.converter.ToolBox(ctx.WithImports(file.imports), calls)
				// if err != nil {
				// 	return false
				// }
				//
				// println(tool)

				return false
			}
		}

		return true
	})

	return err
}

func isInit(f *ast.FuncDecl) bool {
	return f != nil &&
		f.Name.Name == "init" &&
		(f.Type.Params == nil || len(f.Type.Params.List) == 0) &&
		(f.Type.Results == nil || len(f.Type.Results.List) == 0)
}
