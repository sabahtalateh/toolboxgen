package validation

import (
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func FuncDefHasExactlyNParams(f *tool.FuncDef, n int) *errors.PositionedErr {
	if len(f.Parameters) != n {
		return errors.Errorf(f.Position, "function should have exactly `%d` parameters", n)
	}

	return nil
}

func FuncDefHasAtMaximumNParams(f *tool.FuncDef, n int) *errors.PositionedErr {
	if len(f.Parameters) > n {
		return errors.Errorf(f.Position, "function should have at maximum `%d` parameters", n)
	}

	return nil
}

func FuncDefHasAtMinimumNParams(f *tool.FuncDef, n int) *errors.PositionedErr {
	if len(f.Parameters) < n {
		return errors.Errorf(f.Position, "function should have at minimum `%d` parameters", n)
	}

	return nil
}

func FuncDefHasAtMaximumNResults(f *tool.FuncDef, n int) *errors.PositionedErr {
	if len(f.Results) > n {
		return errors.Errorf(f.Position, "function should have at maximum `%d` results", n)
	}

	return nil
}

func FuncDefHasAtMinimumNResults(f *tool.FuncDef, n int) *errors.PositionedErr {
	if len(f.Results) < n {
		return errors.Errorf(f.Position, "function should have at minimum `%d` results", n)
	}

	return nil
}
