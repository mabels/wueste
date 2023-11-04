package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyBoolean interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]
	Default() rusty.Optional[any] // match Type
	Format() rusty.Optional[string]
	Ref() rusty.Optional[string]
	Meta() PropertyMeta
}

type PropertyBooleanBuilder struct {
	Id          string
	Type        Type
	Format      rusty.Optional[string]
	Description rusty.Optional[string]
	Default     rusty.Optional[any]
	Ref         rusty.Optional[string]
}

func NewPropertyBooleanBuilder(pb *PropertiesBuilder) *PropertyBooleanBuilder {
	return &PropertyBooleanBuilder{}
}

func (b *PropertyBooleanBuilder) FromJson(js JSONDict) *PropertyBooleanBuilder {
	b.Type = "boolean"
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Default = getFromAttributeOptionalAny(js, "default")
	b.Format = getFromAttributeOptionalString(js, "format")
	return b
}

func PropertyBooleanToJson(b PropertyBoolean) JSONDict {
	jsp := NewJSONDict()
	JSONsetId(jsp, b)
	JSONsetString(jsp, "type", b.Type())
	JSONsetOptionalString(jsp, "description", b.Description())
	JSONsetOptionalAny(jsp, "default", b.Default())
	JSONsetOptionalString(jsp, "format", b.Format())
	return jsp
}

func (b *PropertyBooleanBuilder) Build() rusty.Result[Property] {
	// return ConnectRuntime(NewPropertyBoolean(*b))
	return NewPropertyBoolean(*b)
}

type propertyBoolean struct {
	param PropertyBooleanBuilder
	meta  PropertyMeta
}

func NewPropertyBoolean(p PropertyBooleanBuilder) rusty.Result[Property] {
	p.Type = BOOLEAN
	return rusty.Ok[Property](&propertyBoolean{
		param: p,
		meta:  NewPropertyMeta(),
	})
}

func (p *propertyBoolean) Meta() PropertyMeta {
	return p.meta
}

// func (p propertyBoolean) Clone() Property {
// 	return NewPropertyBoolean(p.param).Ok()
// }

// Description implements PropertyBoolean.
func (p *propertyBoolean) Description() rusty.Optional[string] {
	return p.param.Description
}

func (p *propertyBoolean) Format() rusty.Optional[string] {
	return p.param.Format
}

// func (p *propertyBoolean) Runtime() *PropertyRuntime {
// 	return &p.param.Runtime
// }

// func (p *propertyBoolean) Format() rusty.Optional[string] {
// 	panic("implement me")
// }

// // Id implements PropertyBoolean.
// func (p *propertyBoolean) Id() string {
// 	return p.param.Id
// }

// // Optional implements PropertyBoolean.
// func (p *propertyBoolean) Optional() bool {
// 	return p.param.Optional
// }

// // SetOptional implements PropertyBoolean.
// func (p *propertyBoolean) SetOptional() {
// 	p.param.Optional = true
// }

func (p *propertyBoolean) Default() rusty.Optional[any] {
	return p.param.Default
	// if p.param.Default.IsSome() {
	// 	// lit := wueste.BoolLiteral(*p.param.Default.Value())
	// 	return rusty.Some[bool](p.param.Default.Value())

	// }
	// return rusty.None[bool]()
}

func (p *propertyBoolean) Type() Type {
	return BOOLEAN
}

func (p *propertyBoolean) Id() string {
	return p.param.Id
}

func (p *propertyBoolean) Ref() rusty.Optional[string] {
	return p.param.Ref
}
