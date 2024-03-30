package shared

type Pair[S, T any] struct {
	First  S
	Second T
}

func NewPair[S, T any](first S, second T) Pair[S, T] {
	return Pair[S, T]{First: first, Second: second}
}
