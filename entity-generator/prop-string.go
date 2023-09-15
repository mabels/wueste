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

	Runtime() *PropertyRuntime
	// MinLength() rusty.Optional[int]
	// MaxLength() rusty.Optional[int]
	// Pattern() rusty.Optional[string]
	// CententEncoding() rusty.Optional[string]
	// ContentMediaType() rusty.Optional[string]
}

type PropertyStringParam struct {
	// PropertyParam
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Default     rusty.Optional[string]
	Ref         rusty.Optional[string]
	// Enum      []string
	// MinLength rusty.Optional[int]
	// MaxLength rusty.Optional[int]
	Format  rusty.Optional[StringFormat]
	Runtime PropertyRuntime
	Ctx     PropertyCtx
}

func (b *PropertyStringParam) FromJson(js JSONProperty) *PropertyStringParam {
	b.Type = STRING
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Format = getFromAttributeOptionalString(js, "format")
	b.Default = getFromAttributeOptionalString(js, "default")
	return b
}

func PropertyStringToJson(b PropertyString) JSONProperty {
	jsp := NewJSONProperty()
	JSONsetId(jsp, b)
	JSONsetString(jsp, "type", b.Type())
	JSONsetOptionalString(jsp, "description", b.Description())
	JSONsetOptionalString(jsp, "format", b.Format())
	JSONsetOptionalString(jsp, "default", b.Default())
	return jsp
}

func (b *PropertyStringParam) Build() rusty.Result[Property] {
	return ConnectRuntime(NewPropertyString(*b))
}

type propertyString struct {
	// propertyLiteral[string]
	param PropertyStringParam
}

func (p *propertyString) Runtime() *PropertyRuntime {
	return &p.param.Runtime
}

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

func (p *propertyString) Ref() rusty.Optional[string] {
	return p.Runtime().Ref
}

func (p *propertyString) Id() string {
	return p.param.Id
}

func NewPropertyString(p PropertyStringParam) rusty.Result[Property] {
	p.Type = STRING
	return rusty.Ok[Property](&propertyString{
		param: p,
	})
}
