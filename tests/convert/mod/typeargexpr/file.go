package typeargexpr

type A[X any] struct {
	x *X
}

type B[Y, Z any] A[[]Z]
