package typealiasref

type C struct {
	a string
}

type B C

type A = B
