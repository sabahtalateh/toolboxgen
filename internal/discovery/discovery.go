package discovery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func Discover(conf *Conf) error {
	d := discovery{
		vendorDir: filepath.Join(conf.Mod.Dir, "vendor"),
	}
	if err := d.discoverCalls(conf.RootDir); err != nil {
		return err
	}

	d.typesFromCalls()

	for _, c := range d.calls {
		println(c)
	}

	return nil
}

var lookupPrefixes = []string{
	"github.com/sabahtalateh/toolbox",
}

type discovery struct {
	calls     []callSequence
	vendorDir string
}

func (d *discovery) discoverCalls(dir string) error {
	var calls []callSequence
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(path, d.vendorDir) {
			return err
		}
		if !info.IsDir() {
			return err
		}

		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				if !needLookup(file) {
					continue
				}

				var callsInFile []callSequence
				callsInFile, err = extractCalls(fset, file)
				if err != nil {
					return err
				}
				calls = append(calls, callsInFile...)
			}
		}

		return nil
	})

	d.calls = append(d.calls, calls...)
	return err
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
