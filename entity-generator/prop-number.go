package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyNumber interface {
	Id() string
	Type() Type
	Ref() rusty.Optional[string]
	Description() rusty.Optional[string]
	Format() rusty.Optional[string]
	Default() rusty.Optional[float64] // match Type
	// Enum() []float64
	Maximum() rusty.Optional[float64]
	Minimum() rusty.Optional[float64]

	Runtime() *PropertyRuntime
}

type PropertyNumberParam struct {
	// __loader    SchemaLoader
	Id          string
	Ref         rusty.Optional[string]
	Type        Type
	Description rusty.Optional[string]
	Format      rusty.Optional[string]
	Default     rusty.Optional[float64]
	// Enum        []float64
	Maximum rusty.Optional[float64]
	Minimum rusty.Optional[float64]

	Runtime PropertyRuntime
	Ctx     PropertyCtx
}

func (b *PropertyNumberParam) FromJson(rt PropertyRuntime, js JSONProperty) *PropertyNumberParam {
	b.Type = "number"
	b.Runtime.Assign(rt)
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Format = getFromAttributeOptionalString(js, "format")
	b.Default = getFromAttributeOptionalFloat64(js, "default")
	b.Maximum = getFromAttributeOptionalFloat64(js, "maximum")
	b.Minimum = getFromAttributeOptionalFloat64(js, "minimum")
	b.Runtime.Assign(rt)
	return b
}

func PropertyNumberToJson(b PropertyNumber) JSONProperty {
	jsp := NewJSONProperty()
	JSONsetId(jsp, b)
	JSONsetString(jsp, "type", b.Type())
	JSONsetOptionalString(jsp, "format", b.Format())
	JSONsetOptionalString(jsp, "description", b.Description())
	JSONsetOptionalFloat64(jsp, "default", b.Default())
	JSONsetOptionalFloat64(jsp, "maximum", b.Maximum())
	JSONsetOptionalFloat64(jsp, "minimum", b.Minimum())
	return jsp
}

func (b *PropertyNumberParam) Build() PropertyNumber {
	return ConnectRuntime(NewPropertyNumber(*b))
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

func (p *propertyNumber) Id() string {
	return p.param.Id
}

func (p *propertyNumber) Runtime() *PropertyRuntime {
	return &p.param.Runtime
}

func (p *propertyNumber) Ref() rusty.Optional[string] {
	return p.param.Ref
}
func (p *propertyNumber) Enum() []float64 {
	panic("PropNumber Enum implement me")
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
		return rusty.Some[float64](p.param.Default.Value())

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
