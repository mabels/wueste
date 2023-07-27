package result

type Result[T any] interface {
	IsOk() bool
	IsErr() bool
	Unwrap() T
	UnwrapErr() error
}

type ResultOK[T any] struct {
	t T
}

func (r ResultOK[T]) IsOk() bool {
	return true
}
func (r ResultOK[T]) IsErr() bool {
	return false
}
func (r ResultOK[T]) UnwrapErr() error {
	panic("Result is Ok")
}
func (r ResultOK[T]) Unwrap() T {
	return r.t
}

type ResultError[T any] struct {
	t error
}

func (r ResultError[T]) IsOk() bool {
	return false
}
func (r ResultError[T]) IsErr() bool {
	return false
}

func (r ResultError[T]) Unwrap() T {
	panic("Result is Err")
}

func (r ResultError[T]) UnwrapErr() error {
	return r.t
}

func Ok[T any](t T) Result[T] {
	return ResultOK[T]{
		t: t,
	}
}

func Err[T any](t error) Result[T] {
	return ResultError[T]{t: t}
}
