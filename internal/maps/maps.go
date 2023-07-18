package maps

func Values[KeyT comparable, ValT any](m map[KeyT]ValT) []ValT {
	var res []ValT

	for _, v := range m {
		res = append(res, v)
	}

	return res
}

func ValuesMap[KeyT, ValT comparable](m map[KeyT]ValT) map[ValT]struct{} {
	res := map[ValT]struct{}{}

	for _, v := range m {
		res[v] = struct{}{}
	}

	return res
}

func FromSlice[K comparable, V, S any](slice []S, kvFn func(S) (K, V)) map[K]V {
	res := map[K]V{}
	for _, el := range slice {
		k, v := kvFn(el)
		res[k] = v
	}
	return res
}
