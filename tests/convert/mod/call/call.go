package call

import "github.com/sabahtalateh/toolbox/di"

type A struct {
}

func NewA() {

}

var _ = di.Component(NewA)

func init() {
	di.Component(NewA)
}
