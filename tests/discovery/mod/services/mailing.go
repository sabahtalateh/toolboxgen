package services

import (
	"github.com/sabahtalateh/toolbox/di"
)

type Mailing struct{}

var _ = di.Component()

func NewMailing() *Mailing {
	return &Mailing{}
}
