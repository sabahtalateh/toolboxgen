package errors

import (
	"fmt"
	"go/token"
)

type PositionedErr struct {
	pos token.Position
	err error
}

func Errorf(pos token.Position, format string, a ...any) *PositionedErr {
	return &PositionedErr{pos: pos, err: fmt.Errorf(format, a...)}
}

func Error(pos token.Position, err error) *PositionedErr {
	return &PositionedErr{pos: pos, err: err}
}

func (p *PositionedErr) Error() string {
	if p == nil {
		return ""
	}

	res := fmt.Sprintf("%s", p.err)
	if p.pos.IsValid() {
		res += fmt.Sprintf("\n\tat %s", p.pos)
	}

	return res
}

func Check(ff ...func() *PositionedErr) *PositionedErr {
	for _, f := range ff {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
