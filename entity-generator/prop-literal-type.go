package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

type PropertyLiteralType[T string | uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64] interface {
	Property
	Default() rusty.Optional[wueste.Literal[T]] // match Type
	Enum() []T                                  // match Type
}

type PropertyLiteralParam[T string | uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64] struct {
	PropertyParam
	Format  string
	Default rusty.Optional[T]
	Enum    []T
}

type propertyLiteral[T string | uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64] struct {
	property
	param PropertyLiteralParam[T]
}

func NewPropertyLiteral[T string | uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64](p PropertyLiteralParam[T]) PropertyLiteralType[T] {
	return &propertyLiteral[T]{
		param: p,
	}
}

func (p *propertyLiteral[T]) Default() rusty.Optional[wueste.Literal[T]] {
	return rusty.None[wueste.Literal[T]]()
}

func (p *propertyLiteral[T]) Enum() []T {
	return p.param.Enum
}
