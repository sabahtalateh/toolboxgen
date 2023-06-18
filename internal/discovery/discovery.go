package discovery

import (
	"errors"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/convert"
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
				currentPkgPath := filepath.Dir(filePath)
				if strings.HasPrefix(currentPkgPath, d.mod.Dir) {
					currentPkgPath = strings.Replace(currentPkgPath, d.mod.Dir, d.mod.Path, 1)
				} else {
					return errors.New("impossibru")
				}

				err = d.discoverFile(files, file, currentPkgPath)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (d *discovery) discoverFile(files *token.FileSet, file *ast.File, currPkg string) error {
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
			fd := syntax.ParseFuncDef(n, files)
			if err = fd.Error().Err(); err != nil {
				return false
			}

			if insideFunction == nil {
				insideFunction = n
			}

			return false

		// if n.Recv != nil {
		// tt, pErr := syntax.ParseTypeRef(n.Recv.List[0].TypeRef, files)
		// if pErr != nil {
		// 	err = pErr.Err()
		// 	return false
		// }
		//
		// tpp := make([]tool.TypeParam, 0)
		// for _, param := range tt.TypeParams {
		// 	tpp = append(tpp, tool.TypeParam{Name: param.TypeName, Position: param.Position})
		// }
		//
		// t, pErr := d.converter.ConvertTypeRef(tt, file.Imports, currPkg, files, tpp)
		// if pErr != nil {
		// 	err = pErr.Err()
		// }
		// println(&t)
		// println(n.Recv.List[0].TypeRef)
		// }

		case *ast.CallExpr:
			if insideFunction == nil || isInit(insideFunction) {
				calls := syntax.ParseFuncCalls(n, files)
				for _, call := range calls {
					if err = call.Error().Err(); err != nil {
						return false
					}
				}

				return false
			}

			// if isInit(n) {
			// 	for _, stmt := range n.Body.List {
			// 		switch exprStmt := stmt.(type) {
			// 		case *ast.ExprStmt:
			// 			switch expr := exprStmt.X.(type) {
			// 			case *ast.CallExpr:
			// 				calls, pErr := syntax.ParseFuncCalls(expr, files)
			// 				if pErr != nil {
			// 					err = pErr.Err()
			// 					return false
			// 				}
			//
			// 				t, cErr := d.converter.TryConvertCallsToTool(calls, file.Imports, currPkg, files)
			// 				if cErr != nil {
			// 					err = cErr.Err()
			// 					return false
			// 				}
			//
			// 				if t != nil {
			// 					tools = append(tools, t)
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			// insideFunction = ""
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
