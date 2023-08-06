package call

import "github.com/sabahtalateh/toolbox/di"

type A struct {
}

func NewA() {

}

func f(...any) any {
	return nil
}

var _ = f(map[string]int{})

func init() {
	di.Component(NewA)
}
