package discovery

import (
	"errors"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func Discover(conf *Conf) error {
	d := discovery{mod: conf.Mod}
	if err := d.discoverTools(conf.RootDir); err != nil {
		return err
	}

	return nil
}

type discovery struct {
	mod   *mod.Module
	tools []Tool
}

func (d *discovery) discoverTools(dir string) error {
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip vendors
		if strings.HasPrefix(path, filepath.Join(d.mod.Dir, "vendor")) {
			return err
		}
		if !info.IsDir() {
			return err
		}

		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		for _, pkg := range pkgs {
			for filePath, file := range pkg.Files {
				currentPkgPath := filepath.Dir(filePath)
				if strings.HasPrefix(currentPkgPath, d.mod.Dir) {
					currentPkgPath = strings.Replace(currentPkgPath, d.mod.Dir, d.mod.Path, 1)
				} else {
					return errors.New("impossibru")
				}
				var tools []Tool
				tools, err = discoverToolsInFile(fset, file, currentPkgPath)
				if err != nil {
					return err
				}
				d.tools = append(d.tools, tools...)
			}
		}

		return nil
	})

	return err
}
