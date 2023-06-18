package pkgdir

import (
	"fmt"
	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"golang.org/x/tools/go/packages"
	"path/filepath"
	"strings"
)

type PkgDir struct {
	mod *mod.Module
	std map[string]*packages.Package
}

func New(mod *mod.Module) (*PkgDir, error) {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		return nil, err
	}

	pd := &PkgDir{
		mod: mod,
		std: map[string]*packages.Package{},
	}
	for _, pkg := range pkgs {
		pd.std[pkg.ID] = pkg
	}

	return pd, nil
}

func (p *PkgDir) Dir(pakage string) (string, error) {
	if pkg, ok := p.std[pakage]; ok {
		return p.stdPackageDir(pkg)
	}

	if strings.HasPrefix(pakage, p.mod.Path) {
		suffix := strings.TrimPrefix(pakage, p.mod.Path)
		suffParts := strings.Split(suffix, "/")
		suffParts = slices.Filter(suffParts, func(el string) bool {
			return el != ""
		})
		d := filepath.Join(append([]string{p.mod.Dir}, suffParts...)...)
		return d, nil
	}
	return "", fmt.Errorf("failed to determine directory for `%s` package", pakage)
}

func (p *PkgDir) stdPackageDir(pkg *packages.Package) (string, error) {
	if len(pkg.GoFiles) > 0 {
		return filepath.Dir(pkg.GoFiles[0]), nil
	}

	for _, pp := range p.std {
		if len(pp.GoFiles) > 0 {
			file := pp.GoFiles[0]
			dir := filepath.Dir(file)
			dir = strings.TrimSuffix(dir, pp.ID)
			return filepath.Join(dir, pkg.ID), nil
		}
	}

	return "", fmt.Errorf("failed determine `%s` stdlib package dir", pkg.ID)
}
