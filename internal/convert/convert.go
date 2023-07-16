package convert

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/life4/genesis/slices"

	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/internal/pkgdir"
	"github.com/sabahtalateh/toolboxgen/internal/utils"
)

type Converter struct {
	pkgDir  *pkgdir.PkgDir
	builtin *builtin.Builtin
}

func New(mod *mod.Module) (*Converter, error) {
	var (
		c   = new(Converter)
		err error
	)

	c.pkgDir, err = pkgdir.Init(mod)
	if err != nil {
		return nil, err
	}

	c.builtin = builtin.Init()
	return c, nil
}

func packagePath(ctx Context, alias string) (pakage string, err error) {
	if alias == "" {
		pakage = ctx.Package()
	} else {
		var imp *ast.ImportSpec
		for _, spec := range ctx.Imports() {
			if spec.Name != nil {
				if spec.Name.Name == alias {
					imp = spec
					break
				}
			} else {
				parts := strings.Split(utils.Unquote(spec.Path.Value), "/")
				if len(parts) == 0 {
					err = fmt.Errorf("empty import")
					break
				}
				lastPart := parts[len(parts)-1]
				if alias == lastPart {
					imp = spec
					break
				}
			}
		}

		if imp == nil {
			err = fmt.Errorf("package `%s` not resolved", alias)
		} else {
			pakage = imp.Path.Value
		}
	}

	return utils.Unquote(pakage), err
}

func (c *Converter) dir(mod *mod.Module, pakage string) (string, error) {
	if strings.HasPrefix(pakage, mod.Path) {
		suffix := strings.TrimPrefix(pakage, mod.Path)
		suffParts := strings.Split(suffix, "/")
		suffParts = slices.Filter(suffParts, func(el string) bool {
			return el != ""
		})
		d := filepath.Join(append([]string{mod.Dir}, suffParts...)...)
		return d, nil
	}
	return "", fmt.Errorf("failed to determine directory for `%s` package", pakage)
}
