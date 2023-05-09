package discovery

import (
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
	if err := d.extractToolsFromInits(conf.RootDir); err != nil {
		return err
	}

	return nil
}

type discovery struct {
	tools     []Tool
	vendorDir string
}

func (d *discovery) extractToolsFromInits(dir string) error {
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
				var tools []Tool
				tools, err = extractToolsFromFileInitBlock(fset, file)
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
