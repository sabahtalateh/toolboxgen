package di

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
	"github.com/sabahtalateh/toolboxgen/internal/validation"
)

type Component struct {
	Name       string
	Type       tool.TypeRef
	Values     []*Value
	Parameters []tool.TypeRef
	WithError  bool
	Function   *tool.FunctionRef
}

func (c *Component) Tool() {
}

func (c *Component) EnrichWithFunction(f *tool.FunctionRef) *errors.PositionedErr {
	if err := errors.Validate(
		func() *errors.PositionedErr { return validation.FunctionHasMinimumNResults(f, 1) },
		func() *errors.PositionedErr { return validation.FunctionHasMaximumNResults(f, 2) },
	); err != nil {
		return err
	}

	firstResult := f.Results[0]
	switch res := firstResult.(type) {
	case *tool.BuiltinRef:
		if res.IsError() {
			return errors.Errorf(res.Position, "first return value can not be error")
		} else {
			return errors.Errorf(res.Position, "builtin type can not be a component")
		}
	}
	c.Type = firstResult
	c.Parameters = f.Parameters

	if len(f.Results) == 2 {
		secondResult := f.Results[1]
		switch res := secondResult.(type) {
		case *tool.BuiltinRef:
			if !res.IsError() {
				return errors.ErrorExpectedErr(res.Position)
			}
		default:
			return errors.Errorf(tool.Position(secondResult), "second return value should be error")
		}
		c.WithError = true
	}

	c.Function = f

	return nil
}

type Value struct {
	Type      tool.TypeRef
	WithError bool
	Function  *tool.FunctionRef
}

func ValueFromFunction(f *tool.FunctionRef) (*Value, *errors.PositionedErr) {
	if err := errors.Validate(
		func() *errors.PositionedErr { return validation.FunctionHasNoParams(f) },
		func() *errors.PositionedErr { return validation.FunctionHasMinimumNResults(f, 1) },
		func() *errors.PositionedErr { return validation.FunctionHasMaximumNResults(f, 2) },
	); err != nil {
		return nil, err
	}

	v := &Value{
		Function: f,
		Type:     f.Results[0],
	}

	if len(f.Results) == 2 {
		switch res := f.Results[1].(type) {
		case *tool.BuiltinRef:
			if !res.IsError() {
				return nil, errors.ErrorExpectedErr(res.Position)
			}
			v.WithError = true
		default:
			return nil, errors.ErrorExpectedErr(tool.Position(res))
		}
	}

	return v, nil
}
