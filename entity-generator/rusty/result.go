package rusty

import "fmt"

type Result[T any] interface {
	IsOk() bool
	IsErr() bool
	Err() error
	Ok() T
}

type ResultOK[T any] struct {
	t T
}

func (r ResultOK[T]) IsOk() bool {
	return true
}
func (r ResultOK[T]) IsErr() bool {
	return !r.IsOk()
}

func (r ResultOK[T]) Err() error {
	panic("Result is Ok")
}
func (r ResultOK[T]) Ok() T {
	return r.t
}

type ResultError[T any] struct {
	t error
}

func (r ResultError[T]) IsOk() bool {
	return false
}
func (r ResultError[T]) IsErr() bool {
	return !r.IsOk()
}

func (r ResultError[T]) Ok() T {
	panic(fmt.Errorf("Result is Err:%v", r.t.Error()))
}

func (r ResultError[T]) Err() error {
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
