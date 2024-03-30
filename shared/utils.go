package shared

func RemoveIf[T interface{}](l []T, p func(T) bool) []T {
	hits := make([]int, 0)
	for k, it := range l {
		if p(it) {
			hits = append(hits, k)
		}
	}
	shift := 0
	for _, k := range hits {
		k := k - shift
		l = append(l[:k], l[k+1:]...)
		shift++
	}
	return l
}
