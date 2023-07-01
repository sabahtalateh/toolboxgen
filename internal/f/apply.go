package f

// Apply
// Sequentially applies ff to x
func Apply[T any](x T, ff ...func(x T) T) T {
	for _, f := range ff {
		x = f(x)
	}

	return x
}
