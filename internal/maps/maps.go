package maps

func FromSlice[K comparable, V, S any](slice []S, kvFn func(S) (K, V)) map[K]V {
	res := map[K]V{}
	for _, el := range slice {
		k, v := kvFn(el)
		res[k] = v
	}
	return res
}
