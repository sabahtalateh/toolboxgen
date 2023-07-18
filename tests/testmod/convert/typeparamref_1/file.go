package typeparamref_1

type A[X any] struct{}

type B[Y, Z any] A[Z]
