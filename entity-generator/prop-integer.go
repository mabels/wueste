package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyInteger interface {
	// Id() string
	Type() Type
	Description() rusty.Optional[string]
	Format() rusty.Optional[string]
	// Optional() bool
	// SetOptional()
	Default() rusty.Optional[int] // match Type
	// Enum() []T
	Maximum() rusty.Optional[int]
	Minimum() rusty.Optional[int]

	// ExclusiveMinimum() rusty.Optional[int]
	// ExclusiveMaximum() rusty.Optional[int]
	// MultipleOf() rusty.Optional[int]
}

type PropertyIntegerParam struct {
	__loader    SchemaLoader
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Optional    bool
	Format      rusty.Optional[string]
	Default     rusty.Optional[int]
	// Enum        []T
	// Default rusty.Optional[T]
	Maximum rusty.Optional[int]
	Minimum rusty.Optional[int]
	// ExclusiveMinimum() rusty.Optional[int]
	// ExclusiveMaximum() rusty.Optional[int]
	// MultipleOf() rusty.Optional[int]
}

func (b *PropertyIntegerParam) FromJson(js JSONProperty) *PropertyIntegerParam {
	b.Id = getFromAttributeString(js, "$id")
	b.Type = "integer"
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Format = getFromAttributeOptionalString(js, "format")
	b.Default = getFromAttributeOptionalInt(js, "default")
	b.Maximum = getFromAttributeOptionalInt(js, "maximum")
	b.Minimum = getFromAttributeOptionalInt(js, "minimum")
	return b
}

func PropertyIntegerToJson(b PropertyInteger) JSONProperty {
	jsp := JSONProperty{}
	jsp.setString("type", b.Type())
	jsp.setOptionalString("format", b.Format())
	jsp.setOptionalString("description", b.Description())
	jsp.setOptionalInt("default", b.Default())
	jsp.setOptionalInt("maximum", b.Maximum())
	jsp.setOptionalInt("minimum", b.Minimum())
	return jsp
}

func (b *PropertyIntegerParam) Build() PropertyInteger {
	return NewPropertyInteger(*b)
}

type propertyInteger struct {
	param PropertyIntegerParam
}

func NewPropertyInteger(p PropertyIntegerParam) PropertyInteger {
	p.Type = INTEGER
	return &propertyInteger{
		param: p,
	}
}

func (p *propertyInteger) Description() rusty.Optional[string] {
	return p.param.Description
}

// Format implements PropertyBoolean.
func (p *propertyInteger) Format() rusty.Optional[string] {
	return p.param.Format
}

// Id implements PropertyBoolean.
// func (p *propertyInteger) Id() string {
// 	return p.param.Id
// }

func (p *propertyInteger) Type() Type {
	return INTEGER
}

func (p *propertyInteger) Default() rusty.Optional[int] {
	if p.param.Default.IsSome() {
		// lit := wueste.IntegerLiteral(*p.param.Default.Value())
		return rusty.Some[int](*p.param.Default.Value())

	}
	return rusty.None[int]()
}

func (p *propertyInteger) Maximum() rusty.Optional[int] {
	return p.param.Maximum
}

func (p *propertyInteger) Minimum() rusty.Optional[int] {
	return p.param.Minimum
}
