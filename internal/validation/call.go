package validation

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

func FunctionCallHasAtMinimumNArguments(c mid.FunctionCall, n int) *errors.PositionedErr {
	if len(c.Args) < n {
		return errors.Errorf(c.Position, "call should have at minimum `%d` arguments", n)
	}
	return nil
}

func FunctionCallHasExactlyNArguments(c mid.FunctionCall, n int) *errors.PositionedErr {
	if len(c.Args) != n {
		return errors.Errorf(c.Position, "function call should have exactly `%d` arguments", n)
	}
	return nil
}
