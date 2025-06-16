package result

type Result[T any] struct {
	Result *T
	Err    error
}

func Failed[T any](err error) Result[T] {
	return Result[T]{Err: err}
}

func Success[T any](result T) Result[T] {
	return Result[T]{Err: nil, Result: &result}
}
