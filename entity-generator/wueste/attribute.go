package wueste

import (
	"errors"

	"github.com/mabels/wueste/entity-generator/rusty"
)

// why wueste german word for desert?
// why de + serialize = desert
// look for a name which is not used by any other package

type Attribute[T any] struct {
	mustSet bool
	value   T
}

func (a *Attribute[T]) IsValid() rusty.Optional[error] {
	if !a.mustSet {
		return rusty.Some[error](errors.New("Attribute not set"))
	}
	return rusty.None[error]()
}

func (a *Attribute[T]) Set(v T) {
	a.mustSet = true
	a.value = v
}

func (a *Attribute[T]) Get() T {
	if !a.mustSet {
		panic("Attribute not set")
	}
	return a.value
}

func (a *Attribute[T]) IsSet() bool {
	return a.mustSet
}

func MustAttribute[T any]() Attribute[T] {
	return Attribute[T]{mustSet: true}
}

func DefaultAttribute[T any](t T) Attribute[T] {
	return Attribute[T]{mustSet: false, value: t}
}

func OptionalAttribute[T any]() Attribute[T] {
	return Attribute[T]{mustSet: false}
}
