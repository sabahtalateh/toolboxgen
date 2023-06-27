package discovery

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/context"
	"github.com/sabahtalateh/toolboxgen/internal/convert"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
)

func Discover(rootDir string) (tt *Tools, err error) {
	d := new(discovery)

	d.mod, err = mod.LookupDir(rootDir, true)
	if err != nil {
		return nil, err
	}

	d.converter, err = convert.New(d.mod)
	if err != nil {
		return nil, err
	}

	if err = d.discoverDir(rootDir); err != nil {
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
		for _, pkg := range pkgs {
			for filePath, file := range pkg.Files {
				Package := filepath.Dir(filePath)
				if strings.HasPrefix(Package, d.mod.Dir) {
					Package = strings.Replace(Package, d.mod.Dir, d.mod.Path, 1)
				} else {
					return errors.New("impossibru")
				}

				err = d.discoverFile(context.New(Package, file.Imports, files), file)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

// TODO в функциях, которые принимают syntax.TypeRef и др проверять ошибки в этих типах - `.Error()`

func (d *discovery) discoverFile(ctx context.Context, file *ast.File) error {
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

		switch n := node.(type) {
		case *ast.FuncDecl:
			if insideFunction == nil {
				insideFunction = n
			}

			fd, pErr := d.converter.ConvertFuncDef(ctx.WithImports(file.Imports), n)
			if pErr != nil {
				err = pErr.Err()
				return false
			}

			println(fd)
		case *ast.CallExpr:
			if insideFunction == nil || isInit(insideFunction) {
				calls := syntax.ParseFuncCalls(n, ctx.Files)
				for _, call := range calls {
					if err = call.Error().Err(); err != nil {
						return false
					}
				}

				tool, pErr := d.converter.ToolBox(ctx.WithImports(file.Imports), calls)
				if pErr != nil {
					err = pErr.Err()
					return false
				}

				println(tool)

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
