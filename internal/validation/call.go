package validation

import (
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

func CallHasAtMinimumNArguments(c syntax.FunctionCall, n int) *errors.PositionedErr {
	if len(c.Args) < n {
		return errors.Errorf(c.Position, "call should have at least `%d` arguments", n)
	}
	return nil
}

func CallHasNArguments(c syntax.FunctionCall, n int) *errors.PositionedErr {
	if len(c.Args) != n {
		return errors.Errorf(c.Position, "call should have `%d` arguments", n)
	}
	return nil
}

func CallHasNthArgumentOfType[T syntax.FunctionCallArgument](c syntax.FunctionCall, n int) *errors.PositionedErr {
	arg := c.Args[n]
	_, ok := arg.(T)
	if !ok {
		return errors.Errorf(arg.Position(), "wrong type")
	}

	return nil
}
