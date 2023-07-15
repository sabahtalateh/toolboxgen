package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

const (
	toolBoxPackage      = "github.com/sabahtalateh/toolbox"
	diPackage           = toolBoxPackage + "/di"
	diComponentFunction = "Component"
)

func (c *Converter) ToolBox(ctx Context, calls []mid.FunctionCall) (tool.Tool, *errors.PositionedErr) {
	if len(calls) == 0 {
		return nil, nil
	}

	firstCall := calls[0]
	pkg, err := packagePath(ctx, firstCall.PkgAlias)
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
			return nil, nil
			// return c.component.convert(ctx, calls)
		default:
			return nil, errors.Errorf(firstCall.Position, "unsupported di function %s", firstCall.Path())
		}
	default:
		return nil, nil
	}
}
