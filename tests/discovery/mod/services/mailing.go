package services

import (
	"github.com/sabahtalateh/toolbox"
	in2 "mod/services/in"
)

const hello = "hello"

type Config struct {
	F1 string
}

func (c Config) Bub() Config {
	return c
}

func C() Config {
	return Config{}
}

func From[T any]() T {
	var t T
	return t
}

type Mailing struct {
	a struct {
		b string
	}
}

var m Mailing

var z = []any{"1"}

var _ = toolbox.Component(
	*func() Mailing {
		var in in2.In
		return Mailing{}
	}().a,
	Config{F1: hello}.F1().F2,
	a.b.c,
	func() {}(),
	a.b().c.d(),
	[]**[]func(A[string], ...int),
	[]*map[string]a.B[a.B, int],
	a.B[string, int],

	Config{F1: hello}.F1().F2,
	a.b.c(func(a ...string) int {
		x := 1
		y := 2
		return x + y
	}),
	// z...,
	// a.b().c.d(),
	// с(),
	// a.b().c.d(),

	// map[string]string{"1": "2"},
	// []string{"Sat", "Sun"},
)

// With("c", Config{}.F1).
// With("c", From[Config]().Bub().F1)

func NewMailing() *Mailing {
	return &Mailing{}
}
