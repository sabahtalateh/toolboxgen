// package mid is a short-word for intermediate
// includes functions to parse ast into more convenient types

package mid

import "github.com/sabahtalateh/toolboxgen/internal/errors"

func (c FunctionCall) Error() *errors.PositionedErr {
	// var err *errors.PositionedErr
	//
	// for _, x := range c.TypeParams {
	// 	if err = x.ParseError(); err != nil {
	// 		return err
	// 	}
	// }
	//
	// for _, x := range c.Args {
	// 	if err = x.ParseError(); err != nil {
	// 		return err
	// 	}
	// }

	return c.Err
}

func (c *FuncCallRef) Error() *errors.PositionedErr {
	return c.err
}

func (s *String) Error() *errors.PositionedErr {
	return s.err
}

func (i *Int) Error() *errors.PositionedErr {
	return i.err
}
