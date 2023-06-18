package convert

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

const (
	toolBoxPackage      = "github.com/sabahtalateh/toolbox"
	diPackage           = toolBoxPackage + "/di"
	diComponentFunction = "Component"
)

func (c *Converter) ToolBox(
	calls []syntax.FunctionCall,
	imports []*ast.ImportSpec,
	currPkg string,
	files *token.FileSet,
) (tool.Tool, *errors.PositionedErr) {
	if len(calls) == 0 {
		return nil, nil
	}

	firstCall := calls[0]
	pkg, err := c.packagePath(firstCall.PkgAlias, currPkg, imports)
	if err != nil {
		return nil, errors.Errorf(firstCall.Position, "unsupported package `%s`", pkg)
	}

	if !strings.HasPrefix(pkg, toolBoxPackage) {
		return nil, nil
	}

	switch pkg {
	case diPackage:
		switch firstCall.FuncName {
		case diComponentFunction:
			return c.component.convert(calls, currPkg, imports, files)
		default:
			return nil, errors.Errorf(firstCall.Position, "unsupported di function `%s`", firstCall.Path())
		}
	default:
		return nil, nil
	}
}
