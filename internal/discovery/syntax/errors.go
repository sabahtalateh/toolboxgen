package syntax

import "github.com/sabahtalateh/toolboxgen/internal/errors"

func (c FunctionCall) Error() *errors.PositionedErr {
	var err *errors.PositionedErr

	for _, x := range c.TypeParams {
		if err = x.Error(); err != nil {
			return err
		}
	}

	for _, x := range c.Args {
		if err = x.Error(); err != nil {
			return err
		}
	}

	return c.Err
}

func (c *Ref) Error() *errors.PositionedErr {
	return c.err
}

func (s *String) Error() *errors.PositionedErr {
	return s.err
}

func (i *Int) Error() *errors.PositionedErr {
	return i.err
}

func (f *FuncDef) Error() *errors.PositionedErr {
	var err *errors.PositionedErr

	if err = f.Receiver.Error(); err != nil {
		return err
	}

	for _, x := range f.TypeParams {
		if err = x.Error(); err != nil {
			return err
		}
	}

	for _, x := range f.Args {
		if err = x.Error(); err != nil {
			return err
		}
	}

	for _, x := range f.Results {
		if err = x.Error(); err != nil {
			return err
		}
	}

	return f.Err
}

func (a *TypeParam) Error() *errors.PositionedErr {
	return a.Err
}

func (t *TypeRef) Error() *errors.PositionedErr {
	var err *errors.PositionedErr

	for _, x := range t.TypeParams {
		if err = x.Error(); err != nil {
			return err
		}
	}

	return t.Err
}
