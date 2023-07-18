package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

type StringFormat = string

const (
	DATE_TIME StringFormat = "date-time"
	TIME      StringFormat = "time"
	DATE      StringFormat = "date"
)

type PropertyString interface {
	Property
	Type() string
	Default() rusty.Optional[wueste.Literal[string]] // match Type
	Format() rusty.Optional[StringFormat]
	// MinLength() rusty.Optional[int]
	// MaxLength() rusty.Optional[int]
	// Pattern() rusty.Optional[string]
	// CententEncoding() rusty.Optional[string]
	// ContentMediaType() rusty.Optional[string]
}

type PropertyStringParam struct {
	PropertyParam
	Default rusty.Optional[string]
	// Enum      []string
	// MinLength rusty.Optional[int]
	// MaxLength rusty.Optional[int]
	Format rusty.Optional[StringFormat]
}

type propertyString struct {
	// propertyLiteral[string]
	param PropertyStringParam
}

// Description implements PropertyString.
func (p *propertyString) Description() rusty.Optional[string] {
	return p.param.Description
}

// Id implements PropertyString.
func (p *propertyString) Id() string {
	return p.param.Id
}

// Optional implements PropertyString.
func (p *propertyString) Optional() bool {
	return p.param.Optional
}

// SetOptional implements PropertyString.
func (p *propertyString) SetOptional() {
	p.param.Optional = true
}

// Enum implements PropertyString.
// func (p *propertyString) Enum() []string {
// 	return p.param.enum
// }

func (p *propertyString) Type() Type {
	return STRING
}

func (p *propertyString) Default() rusty.Optional[wueste.Literal[string]] {
	if !p.param.Default.IsNone() {
		lit := wueste.StringLiteral(*p.param.Default.Value())
		return rusty.Some[wueste.Literal[string]](lit)

	}
	return rusty.None[wueste.Literal[string]]()
}

// func (p *propertyString) MinLength() rusty.Optional[int] {
// 	return p.param.minLength
// }

// func (p *propertyString) MaxLength() rusty.Optional[int] {
// 	return p.param.maxLength
// }

func (p *propertyString) Format() rusty.Optional[StringFormat] {
	return p.param.Format
}

func NewPropertyString(p PropertyStringParam) PropertyString {
	return &propertyString{
		param: p,
	}
}
