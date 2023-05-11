package discovery

import "fmt"

func validate(ff ...func() error) error {
	for _, f := range ff {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

type validators struct {
	component struct {
		typed           func(c call) error
		typeNotEmpty    func(c call) error
		name            func(c []call) error
		constructorArgs func(c call) error
		providerArgs    func(c call) error
	}
}

var v = validators{
	component: struct {
		typed           func(c call) error
		typeNotEmpty    func(c call) error
		name            func(c []call) error
		constructorArgs func(c call) error
		providerArgs    func(c call) error
	}{},
}

func init() {
	v.component.typed = func(c call) error {
		if !c.typed {
			return fmt.Errorf("`component.Component` should have type parameter\n\tat %s", c.position)
		}
		return nil
	}

	v.component.typeNotEmpty = func(c call) error {
		if c.typed && c.typeParameter.typ == "" {
			return fmt.Errorf("`component.Component` should have type parameter\n\tat %s", c.position)
		}
		return nil
	}

	v.component.name = func(calls []call) error {
		names := 0
		var nameCall *call
		for _, c := range calls {
			if c.funcName == "Name" {
				nameCall = &c
				names++
			}
			if names > 1 {
				return fmt.Errorf("second occurrence of `component.Name` not allowed\n\tat %s", c.position)
			}
		}

		if nameCall != nil {
			if len(nameCall.arguments) != 1 {
				return fmt.Errorf("`component.Name` should have at exactly one argument\n\tat %s", nameCall.position)
			}
			arg := nameCall.arguments[0]
			if arg.typ != stringType {
				return fmt.Errorf("`component.Name` should have string argument\n\tat %s", nameCall.position)
			}
		}

		return nil
	}

	v.component.constructorArgs = func(c call) error {
		if len(c.arguments) != 1 {
			return fmt.Errorf("`component.Component` should have exactly one argument\n\tat %s", c.position)
		}

		arg := c.arguments[0]
		if arg.typ != refType {
			return fmt.Errorf("`component.Component` argument should be of function type\n\tat %s", c.position)
		}

		return nil
	}
	v.component.providerArgs = v.component.constructorArgs
}
