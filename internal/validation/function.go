package validation

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func FunctionHasNoParams(f *tool.FunctionRef) *errors.PositionedErr {
	return FunctionHasMaximumNParams(f, 0)
}

func FunctionHasMaximumNParams(f *tool.FunctionRef, n int) *errors.PositionedErr {
	if len(f.Parameters) > n {
		return errors.Errorf(f.Position, "function should have at maximum `%d` parameters", n)
	}

	return nil
}

func FunctionHasMinimumNParams(f *tool.FunctionRef, n int) *errors.PositionedErr {
	if len(f.Parameters) < n {
		return errors.Errorf(f.Position, "function should have at minimum `%d` parameters", n)
	}

	return nil
}

func FunctionHasMaximumNResults(f *tool.FunctionRef, n int) *errors.PositionedErr {
	if len(f.Results) > n {
		return errors.Errorf(f.Position, "function should have at maximum `%d` results", n)
	}

	return nil
}

func FunctionHasMinimumNResults(f *tool.FunctionRef, n int) *errors.PositionedErr {
	if len(f.Results) < n {
		return errors.Errorf(f.Position, "function should have at minimum `%d` results", n)
	}

	return nil
}
