package typealiasexpr

type C struct {
	a string
}

type B C

type A = B
