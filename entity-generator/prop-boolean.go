package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

type PropertyBoolean interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]
	Optional() bool
	SetOptional()
	Format() rusty.Optional[string]
	// Property
	Default() rusty.Optional[wueste.Literal[bool]] // match Type
}

type PropertyBooleanParam struct {
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Default     rusty.Optional[bool]
	Optional    bool
}

type propertyBoolean struct {
	// property
	// propertyLiteral[bool]
	param PropertyBooleanParam
}

// Description implements PropertyBoolean.
func (p *propertyBoolean) Description() rusty.Optional[string] {
	return p.param.Description
}

func (p *propertyBoolean) Format() rusty.Optional[string] {
	panic("implement me")
}

// Id implements PropertyBoolean.
func (p *propertyBoolean) Id() string {
	return p.param.Id
}

// Optional implements PropertyBoolean.
func (p *propertyBoolean) Optional() bool {
	return p.param.Optional
}

// SetOptional implements PropertyBoolean.
func (p *propertyBoolean) SetOptional() {
	p.param.Optional = true
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
