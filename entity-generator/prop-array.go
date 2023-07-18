package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyArray interface {
	Property
	MinItems() rusty.Optional[int]
	MaxItems() rusty.Optional[int]
	Items() Property
	// UniqueItems() bool
	// Containts() Property
	// AdditionalItems() rusty.Optional[Property]
}

type PropertyArrayParam struct {
	PropertyParam
	MinItems rusty.Optional[int]
	MaxItems rusty.Optional[int]
	Items    Property
	// Default rusty.Optional[string]
	// Enum      []string
	// MinLength rusty.Optional[int]
	// MaxLength rusty.Optional[int]
	// Format    rusty.Optional[StringFormat]
}

type propertyArray struct {
	propertyLiteral[string]
	param PropertyArrayParam
}

// Items implements PropertyArray.
func (p *propertyArray) Items() Property {
	return p.param.Items
}

// MaxItems implements PropertyArray.
func (*propertyArray) MaxItems() rusty.Optional[int] {
	panic("unimplemented")
}

// MinItems implements PropertyArray.
func (*propertyArray) MinItems() rusty.Optional[int] {
	panic("unimplemented")
}

// Enum implements PropertyArray.
// func (p *propertyString) Enum() []string {
// 	return p.param.enum
// }

func (p *propertyArray) Type() Type {
	return ARRAY
}

// func (p *propertyArray) Default() rusty.Optional[wueste.Literal[string]] {
// 	if !p.param.Default.IsNone() {
// 		lit := wueste.StringLiteral(*p.param.Default.Value())
// 		return rusty.Some[wueste.Literal[string]](lit)

// 	}
// 	return rusty.None[wueste.Literal[string]]()
// }

// func (p *propertyString) MinLength() rusty.Optional[int] {
// 	return p.param.minLength
// }

// func (p *propertyString) MaxLength() rusty.Optional[int] {
// 	return p.param.maxLength
// }

// func (p *propertyString) Format() rusty.Optional[StringFormat] {
// 	return p.param.format
// }

func NewPropertyArray(p PropertyArrayParam) PropertyArray {
	return &propertyArray{
		param: p,
	}
}
