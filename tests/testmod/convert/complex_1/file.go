package complex_1

type A[A, B, C, D any] struct {
	d *[]D
}

type B[AA, BB, CC comparable] A[*interface{ A() string }, []float32, *CC, []**[]BB]
