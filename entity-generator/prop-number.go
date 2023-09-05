package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyNumber interface {
	Type() Type
	Description() rusty.Optional[string]
	Format() rusty.Optional[string]
	Default() rusty.Optional[float64] // match Type
	// Enum() []float64
	Maximum() rusty.Optional[float64]
	Minimum() rusty.Optional[float64]
}

type PropertyNumberParam struct {
	__loader    SchemaLoader
	Type        Type
	Description rusty.Optional[string]
	Format      rusty.Optional[string]
	Default     rusty.Optional[float64]
	// Enum        []float64
	Maximum rusty.Optional[float64]
	Minimum rusty.Optional[float64]
}

func (b *PropertyNumberParam) FromJson(js JSONProperty) *PropertyNumberParam {
	b.Type = "number"
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Format = getFromAttributeOptionalString(js, "format")
	b.Default = getFromAttributeOptionalFloat64(js, "default")
	b.Maximum = getFromAttributeOptionalFloat64(js, "maximum")
	b.Minimum = getFromAttributeOptionalFloat64(js, "minimum")
	return b
}

func PropertyNumberToJson(b PropertyNumber) JSONProperty {
	jsp := JSONProperty{}
	jsp.setString("type", b.Type())
	jsp.setOptionalString("format", b.Format())
	jsp.setOptionalString("description", b.Description())
	jsp.setOptionalFloat64("default", b.Default())
	jsp.setOptionalFloat64("maximum", b.Maximum())
	jsp.setOptionalFloat64("minimum", b.Minimum())
	return jsp
}

func (b *PropertyNumberParam) Build() PropertyNumber {
	return NewPropertyNumber(*b)
}

type propertyNumber struct {
	param PropertyNumberParam
}

func NewPropertyNumber(p PropertyNumberParam) PropertyNumber {
	p.Type = NUMBER
	return &propertyNumber{
		param: p,
	}
}

func (p *propertyNumber) Enum() []float64 {
	panic("implement me")
}
func (p *propertyNumber) Description() rusty.Optional[string] {
	return p.param.Description
}

// // Format implements PropertyBoolean.
func (p *propertyNumber) Format() rusty.Optional[string] {
	return p.param.Format
}

func (p *propertyNumber) Default() rusty.Optional[float64] {
	if p.param.Default.IsSome() {
		// lit := wueste.NumberLiteral()
		return rusty.Some[float64](*p.param.Default.Value())

	}
	return rusty.None[float64]()
}

func (p *propertyNumber) Type() Type {
	return NUMBER
}

func (p *propertyNumber) Maximum() rusty.Optional[float64] {
	return p.param.Maximum
}

func (p *propertyNumber) Minimum() rusty.Optional[float64] {
	return p.param.Minimum
}
