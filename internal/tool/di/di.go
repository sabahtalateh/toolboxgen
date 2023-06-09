package di

import (
	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
	"github.com/sabahtalateh/toolboxgen/internal/validation"
)

type Component struct {
	Name       string
	Type       tool.TypeRef
	With       map[string][]*With
	Parameters []tool.TypeRef
	WithError  bool
	Function   *tool.FuncRef
}

func (c *Component) Tool() {
}

func (c *Component) EnrichWithFunction(f *tool.FuncRef) *errors.PositionedErr {
	if err := errors.Check(
		func() *errors.PositionedErr { return validation.FuncDefHasAtMinimumNResults(f.Def, 1) },
		func() *errors.PositionedErr { return validation.FuncDefHasAtMaximumNResults(f.Def, 2) },
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
	c.Parameters = slices.Map(f.Parameters, func(el tool.FuncParam) tool.TypeRef { return el.Type })

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

type With struct {
	Key       string
	Type      tool.TypeRef
	WithError bool
	Function  *tool.FuncRef
}
