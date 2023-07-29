package services

import (
	"github.com/sabahtalateh/toolbox"
)

type Mailing struct{}

var _ = toolbox.Component(NewMailing)

func NewMailing() *Mailing {
	return &Mailing{}
}
