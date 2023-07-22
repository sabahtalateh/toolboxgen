package typeparamref

type A[X any] struct{}

type B[Y, Z any] A[Z]
