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
		register struct {
			typed        func(c call) error
			typeNotEmpty func(c call) error
		}
	}
}

var v = validators{
	component: struct {
		register struct {
			typed        func(c call) error
			typeNotEmpty func(c call) error
		}
	}{},
}

func init() {
	v.component.register.typed = func(c call) error {
		if !c.typed {
			return fmt.Errorf("`component.Register` should have type parameter\n\tat %s", c.position)
		}
		return nil
	}

	v.component.register.typeNotEmpty = func(c call) error {
		if c.typed && c.typeParameter.typ == "" {
			return fmt.Errorf("`component.Register` should have type parameter\n\tat %s", c.position)
		}
		return nil
	}
}
