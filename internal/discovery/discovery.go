package discovery

import (
	"fmt"
	"github.com/kim89098/slice"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

var lookupPrefixes = []string{
	"github.com/sabahtalateh/toolbox",
}

type discovery struct {
	calls []callSequence
}

func (d *discovery) discoverCallsRec(dir string) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return err
		}
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		var inspectErr error

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				if !needLookup(file) {
					return err
				}

				ast.Inspect(file, func(node ast.Node) bool {
					switch n := node.(type) {
					case *ast.FuncDecl:
						if isInit(n) {
							for _, stmt := range n.Body.List {
								switch exprStmt := stmt.(type) {
								case *ast.ExprStmt:
									switch expr := exprStmt.X.(type) {
									case *ast.CallExpr:
										callSeq, analyzeErr := analyzeCallExpr(expr, fset, file.Imports)
										_, ok := slice.Find(lookupPrefixes, func(s string) bool {
											return strings.HasPrefix(callSeq.pakage, s)
										})
										if !ok {
											continue
										}
										if analyzeErr != nil {
											inspectErr = fmt.Errorf("%w\n\tat %s", analyzeErr.err, analyzeErr.position.String())
											return false
										}
										d.calls = append(d.calls, callSeq)
									}
								}
							}
						}
					}
					return true
				})
			}
		}

		if inspectErr != nil {
			return inspectErr
		}

		return nil
	})
}

func needLookup(file *ast.File) bool {
	imports := file.Imports
	for _, imp := range imports {
		for _, pref := range lookupPrefixes {
			if strings.HasPrefix(unquote(imp.Path.Value), pref) {
				return true
			}
		}
	}

	return false
}

func Discover(conf *Conf) error {
	d := discovery{}
	if err := d.discoverCallsRec(conf.RootDir); err != nil {
		return err
	}

	return nil
}
