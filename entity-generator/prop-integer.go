package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

type PropertyInteger[T uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64] interface {
	PropertyLiteralType[T]
	Maximum() rusty.Optional[T]
	Minimum() rusty.Optional[T]

	// ExclusiveMinimum() rusty.Optional[int]
	// ExclusiveMaximum() rusty.Optional[int]
	// MultipleOf() rusty.Optional[int]
}

type PropertyIntegerParam[T uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64] struct {
	PropertyLiteralParam[T]
	Default rusty.Optional[T]
	Maximum rusty.Optional[T]
	Minimum rusty.Optional[T]
	// ExclusiveMinimum() rusty.Optional[int]
	// ExclusiveMaximum() rusty.Optional[int]
	// MultipleOf() rusty.Optional[int]
}

type propertyInteger[T uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64] struct {
	propertyLiteral[T]
	param PropertyIntegerParam[T]
}

func NewPropertyInteger[T uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64](p PropertyIntegerParam[T]) PropertyInteger[T] {
	return &propertyInteger[T]{
		param: p,
	}
}

func (p *propertyInteger[T]) Type() Type {
	return INTEGER
}

func (p *propertyInteger[T]) Default() rusty.Optional[wueste.Literal[T]] {
	if p.param.Default.IsSome() {
		lit := wueste.IntegerLiteral(*p.param.Default.Value())
		return rusty.Some[wueste.Literal[T]](lit)

	}
	return rusty.None[wueste.Literal[T]]()
}

func (p *propertyInteger[T]) Maximum() rusty.Optional[T] {
	return p.param.Maximum
}

func (p *propertyInteger[T]) Minimum() rusty.Optional[T] {
	return p.param.Minimum
}
