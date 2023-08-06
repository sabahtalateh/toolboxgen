package syntax

import "errors"

func IsUnknownExprErr(err error) bool {
	return errors.Is(err, ErrUnexpectedExpr)
}
