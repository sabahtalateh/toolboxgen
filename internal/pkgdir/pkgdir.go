package pkgdir

import (
	"errors"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
)

var (
	ErrDirNotFound    = errors.New("package dir not found")
	ErrStdDirNotFound = errors.New("standard library package dir not found")
)

type PkgDir struct {
	mod *mod.Module
	std map[string]*packages.Package
}

func Init(mod *mod.Module) (*PkgDir, error) {
	pd := &PkgDir{
		mod: mod,
		std: map[string]*packages.Package{},
	}

	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		return nil, err
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

	return "", errors.Join(ErrDirNotFound, errors.New(pakage))
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

	return "", errors.Join(ErrStdDirNotFound, errors.New(pkg.ID))
}
