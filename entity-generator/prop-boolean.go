package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

type PropertyBoolean interface {
	Property
	Default() rusty.Optional[wueste.Literal[bool]] // match Type
}

type PropertyBooleanParam struct {
	PropertyParam
	Default rusty.Optional[bool]
}

type propertyBoolean struct {
	propertyLiteral[bool]
	param PropertyBooleanParam
}

func (p *propertyBoolean) Default() rusty.Optional[wueste.Literal[bool]] {
	if p.param.Default.IsSome() {
		lit := wueste.BoolLiteral(*p.param.Default.Value())
		return rusty.Some[wueste.Literal[bool]](lit)

	}
	return rusty.None[wueste.Literal[bool]]()
}

func (p *propertyBoolean) Type() Type {
	return BOOLEAN
}

func NewPropertyBoolean(p PropertyBooleanParam) PropertyBoolean {
	return &propertyBoolean{
		param: p,
	}
}
