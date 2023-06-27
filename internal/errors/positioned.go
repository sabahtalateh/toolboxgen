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

func (p *PositionedErr) Err() error {
	if p == nil {
		return nil
	}
	return fmt.Errorf("%w\n\tat %s", p.err, p.pos)
}

func Check(ff ...func() *PositionedErr) *PositionedErr {
	for _, f := range ff {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
