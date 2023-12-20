package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type StringFormat = string

const (
	DATE_TIME StringFormat = "date-time"
	TIME      StringFormat = "time"
	DATE      StringFormat = "date"
)

type PropertyString interface {
	// Property
	Id() string
	Type() string
	Description() rusty.Optional[string]
	Default() rusty.Optional[string] // match Type
	Format() rusty.Optional[StringFormat]
	Ref() rusty.Optional[string]
	XProperties() map[string]interface{}
	Meta() PropertyMeta

	// Runtime() *PropertyRuntime
	// Clone() Property
	// MinLength() rusty.Optional[int]
	// MaxLength() rusty.Optional[int]
	// Pattern() rusty.Optional[string]
	// CententEncoding() rusty.Optional[string]
	// ContentMediaType() rusty.Optional[string]
}

type PropertyStringBuilder struct {
	// PropertyParam
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Default     rusty.Optional[string]
	Ref         rusty.Optional[string]
	XProperties map[string]interface{}
	// Enum      []string
	// MinLength rusty.Optional[int]
	// MaxLength rusty.Optional[int]

	Format rusty.Optional[StringFormat]
	// Runtime PropertyRuntime
	// Ctx     PropertyCtx
}

func NewPropertyStringBuilder(pb *PropertiesBuilder) *PropertyStringBuilder {
	return &PropertyStringBuilder{
		XProperties: make(map[string]interface{}),
	}
}

// func (b *PropertyStringBuilder) FromProperty(prop Property) *PropertyStringBuilder {
// 	ps, found := prop.(PropertyString)
// 	if !found {
// 		panic("not a PropertyString")
// 	}

// 	return b
// }

func (b *PropertyStringBuilder) FromJson(js JSONDict) *PropertyStringBuilder {
	b.Type = STRING
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Format = getFromAttributeOptionalString(js, "format")
	b.Default = getFromAttributeOptionalString(js, "default")
	b.XProperties = getFromAttributeXProperties(js)
	return b
}

func PropertyStringToJson(b PropertyString) JSONDict {
	jsp := NewJSONDict()
	JSONsetId(jsp, b)
	JSONsetString(jsp, "type", b.Type())
	JSONsetOptionalString(jsp, "description", b.Description())
	JSONsetOptionalString(jsp, "format", b.Format())
	JSONsetOptionalString(jsp, "default", b.Default())
	JSONsetXProperties(jsp, b.XProperties())
	return jsp
}

func (b *PropertyStringBuilder) Build() rusty.Result[Property] {
	return NewPropertyString(*b)
}

type propertyString struct {
	// propertyLiteral[string]
	param PropertyStringBuilder
	meta  PropertyMeta
}

func (p *propertyString) Meta() PropertyMeta {
	return p.meta
}

// func (p propertyString) Clone() Property {
// 	return NewPropertyString(p.param).Ok()
// }

// func (p *propertyString) Runtime() *PropertyRuntime {
// 	return &p.param.Runtime
// }

// Description implements PropertyString.
func (p *propertyString) Description() rusty.Optional[string] {
	return p.param.Description
}

// Id implements PropertyString.
// func (p *propertyString) Id() string {
// 	return p.param.Id
// }

// Optional implements PropertyString.
// func (p *propertyString) Optional() bool {
// 	return p.param.Optional
// }

// // SetOptional implements PropertyString.
// func (p *propertyString) SetOptional() {
// 	p.param.Optional = true
// }

// Enum implements PropertyString.
// func (p *propertyString) Enum() []string {
// 	return p.param.enum
// }

func (p *propertyString) Type() Type {
	return STRING
}

func (p *propertyString) Default() rusty.Optional[string] {
	if !p.param.Default.IsNone() {
		// lit := wueste.StringLiteral(*p.param.Default.Value())
		return rusty.Some[string](p.param.Default.Value())

	}
	return rusty.None[string]()
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

func (p *propertyString) XProperties() map[string]interface{} {
	return p.param.XProperties
}

func (p *propertyString) Ref() rusty.Optional[string] {
	return p.param.Ref
}

func (p *propertyString) Id() string {
	return p.param.Id
}

func NewPropertyString(p PropertyStringBuilder) rusty.Result[Property] {
	p.Type = STRING
	return rusty.Ok[Property](&propertyString{
		param: p,
		meta:  NewPropertyMeta(),
	})
}
