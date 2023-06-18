package convert

import (
	"go/ast"
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
	"github.com/sabahtalateh/toolboxgen/internal/tool/di"
	"github.com/sabahtalateh/toolboxgen/internal/validation"
)

type component struct {
	converter *Converter
}

func (c *component) convert(
	calls []syntax.FunctionCall,
	currPkg string,
	imports []*ast.ImportSpec,
	files *token.FileSet,
) (comp *di.Component, err *errors.PositionedErr) {
	firstCall := calls[0]
	if err = errors.Validate(
		func() *errors.PositionedErr { return validation.CallHasAtMinimumNArguments(firstCall, 1) },
	); err != nil {
		return nil, err
	}

	comp = new(di.Component)
	arg := firstCall.Args[0]
	var function *tool.FunctionRef
	switch val := arg.(type) {
	case *syntax.Ref:
		function, err = c.converter.FuncRef(*val, currPkg, imports, files)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.FunctionExpectedErr(arg.Position())
	}

	if err = comp.EnrichWithFunction(function); err != nil {
		return nil, err
	}

	comp.Values, err = c.convertValues(firstCall.Args[1:], currPkg, imports, files)
	if err != nil {
		return nil, err
	}

	return comp, nil
}

func (c *component) convertNamed(
	calls []syntax.FunctionCall,
	currPkg string,
	imports []*ast.ImportSpec,
	files *token.FileSet,
) (comp *di.Component, err *errors.PositionedErr) {
	firstCall := calls[0]

	if err = errors.Validate(
		func() *errors.PositionedErr { return validation.CallHasAtMinimumNArguments(firstCall, 2) },
	); err != nil {
		return nil, err
	}

	comp = new(di.Component)
	firstArg := firstCall.Args[0]
	switch val := firstArg.(type) {
	case *syntax.String:
		comp.Name = val.Val
	default:
		return nil, errors.StringExpectedErr(firstArg.Position())
	}

	secondArg := firstCall.Args[1]
	var function *tool.FunctionRef
	switch val := secondArg.(type) {
	case *syntax.Ref:
		function, err = c.converter.FuncRef(*val, currPkg, imports, files)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.FunctionExpectedErr(secondArg.Position())
	}

	if err = comp.EnrichWithFunction(function); err != nil {
		return nil, err
	}

	comp.Values, err = c.convertValues(firstCall.Args[2:], currPkg, imports, files)
	if err != nil {
		return nil, err
	}

	return comp, nil
}

func (c *component) convertValues(
	args []syntax.FunctionCallArgument,
	currPkg string,
	imports []*ast.ImportSpec,
	files *token.FileSet,
) ([]*di.Value, *errors.PositionedErr) {
	var res []*di.Value

	for _, arg := range args {
		switch ref := arg.(type) {
		case *syntax.Ref:
			f, err := c.converter.FuncRef(*ref, currPkg, imports, files)
			if err != nil {
				return nil, err
			}

			v, err := di.ValueFromFunction(f)
			if err != nil {
				return nil, err
			}

			res = append(res, v)
		default:
			return nil, errors.FunctionExpectedErr(arg.Position())
		}
	}

	return res, nil
}
