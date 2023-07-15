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

func FromSlice[K comparable, V any](s []V, keyFn func(V) K) map[K]V {
	res := map[K]V{}
	for _, el := range s {
		res[keyFn(el)] = el
	}
	return res
}
