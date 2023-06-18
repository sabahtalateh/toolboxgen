package syntax

import "github.com/sabahtalateh/toolboxgen/internal/errors"

func (c FunctionCall) FirstError() *errors.PositionedErr {
	var err *errors.PositionedErr

	for _, x := range c.TypeParams {
		if err = x.FirstError(); err != nil {
			return err
		}
	}

	for _, x := range c.Args {
		if err = x.FirstError(); err != nil {
			return err
		}
	}

	return c.Error
}

func (c *Ref) FirstError() *errors.PositionedErr {
	return c.err
}

func (s *String) FirstError() *errors.PositionedErr {
	return s.err
}

func (i *Int) FirstError() *errors.PositionedErr {
	return i.err
}

func (f *FuncDef) FirstError() *errors.PositionedErr {
	var err *errors.PositionedErr

	if err = f.Receiver.FirstError(); err != nil {
		return err
	}

	for _, x := range f.TypeParams {
		if err = x.FirstError(); err != nil {
			return err
		}
	}

	for _, x := range f.Args {
		if err = x.FirstError(); err != nil {
			return err
		}
	}

	for _, x := range f.Results {
		if err = x.FirstError(); err != nil {
			return err
		}
	}

	return f.Error
}

func (a *TypeParam) FirstError() *errors.PositionedErr {
	return a.Error
}

func (t *TypeRef) FirstError() *errors.PositionedErr {
	var err *errors.PositionedErr

	for _, x := range t.TypeParams {
		if err = x.FirstError(); err != nil {
			return err
		}
	}

	return t.Error
}
