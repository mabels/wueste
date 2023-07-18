package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

type PropertyNumber[T float32 | float64] interface {
	PropertyLiteralType[T]
	Maximum() rusty.Optional[T]
	Minimum() rusty.Optional[T]

	// ExclusiveMinimum() rusty.Optional[int]
	// ExclusiveMaximum() rusty.Optional[int]
	// MultipleOf() rusty.Optional[int]
}

type PropertyNumberParam[T float32 | float64] struct {
	PropertyLiteralParam[T]
	Default rusty.Optional[T]
	Maximum rusty.Optional[T]
	Minimum rusty.Optional[T]
	// ExclusiveMinimum() rusty.Optional[int]
	// ExclusiveMaximum() rusty.Optional[int]
	// MultipleOf() rusty.Optional[int]
}

type propertyNumber[T float32 | float64] struct {
	propertyLiteral[T]
	param PropertyNumberParam[T]
}

func NewPropertyNumber[T float32 | float64](p PropertyNumberParam[T]) PropertyNumber[T] {
	return &propertyNumber[T]{
		param: p,
	}
}

func (p *propertyNumber[T]) Default() rusty.Optional[wueste.Literal[T]] {
	if p.param.Default.IsSome() {
		lit := wueste.NumberLiteral(*p.param.Default.Value())
		return rusty.Some[wueste.Literal[T]](lit)

	}
	return rusty.None[wueste.Literal[T]]()
}

func (p *propertyNumber[T]) Type() Type {
	return NUMBER
}

func (p *propertyNumber[T]) Maximum() rusty.Optional[T] {
	return p.param.Maximum
}

func (p *propertyNumber[T]) Minimum() rusty.Optional[T] {
	return p.param.Minimum
}
