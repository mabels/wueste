package entity_generator

import (
	"fmt"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyArray interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]
	// Format() rusty.Optional[string]
	// Optional() bool
	// SetOptional()
	MinItems() rusty.Optional[int]
	MaxItems() rusty.Optional[int]
	Items() Property
	// UniqueItems() bool
	// Containts() Property
	// AdditionalItems() rusty.Optional[Property]
	Ref() rusty.Optional[string]
	// Runtime() *PropertyRuntime
	// Clone() Property

	// ToPropertyObject() rusty.Result[PropertyObject]
	Meta() PropertyMeta
}

type PropertyArrayBuilder struct {
	// __loader SchemaLoader
	Id          string
	Type        Type
	Ref         rusty.Optional[string]
	Description rusty.Optional[string]
	// Format      rusty.Optional[string]
	// Optional    bool
	MinItems rusty.Optional[int]
	MaxItems rusty.Optional[int]
	Items    rusty.Result[Property]
	// Errors   []error
	// Runtime PropertyRuntime
	// Ctx     PropertyCtx
	// Default rusty.Optional[string]
	// Enum      []string
	// MinLength rusty.Optional[int]
	// MaxLength rusty.Optional[int]
	// Format    rusty.Optional[StringFormat]
	_propertiesBuilder *PropertiesBuilder
}

func NewPropertyArrayBuilder(pb *PropertiesBuilder) *PropertyArrayBuilder {
	return &PropertyArrayBuilder{
		_propertiesBuilder: pb,
	}
}

func (b *PropertyArrayBuilder) FromJson(js JSONProperty) *PropertyArrayBuilder {
	b.Type = ARRAY
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.MaxItems = getFromAttributeOptionalInt(js, "maxItems")
	b.MinItems = getFromAttributeOptionalInt(js, "minItems")

	builder := NewPropertiesBuilder(b._propertiesBuilder.ctx)
	builder.filename = b._propertiesBuilder.filename
	b.Items = builder.FromJson(js.Get("items").(JSONProperty)).Build()
	return b
}

func PropertyArrayToJson(b PropertyArray) JSONProperty {
	jsp := NewJSONProperty()
	JSONsetId(jsp, b)
	JSONsetString(jsp, "type", b.Type())
	JSONsetOptionalString(jsp, "description", b.Description())
	JSONsetOptionalInt(jsp, "maxItems", b.MaxItems())
	JSONsetOptionalInt(jsp, "minItems", b.MinItems())
	jsp.Set("items", PropertyToJson(b.Items()))
	return jsp
}

func (b *PropertyArrayBuilder) Build() rusty.Result[Property] {
	if b.Items.IsOk() {
		pa := NewPropertyArray(*b)
		if b.Items.IsOk() {
			b.Items.Ok().Meta().SetParent(pa.Ok())
		}
		return pa
	} else {
		return rusty.Err[Property](fmt.Errorf("Array needs items:%v", b.Items.Err()))
	}
}

type propertyArray struct {
	param PropertyArrayBuilder
	meta  PropertyMeta
}

func (p *propertyArray) Meta() PropertyMeta {
	return p.meta
}

// func (p *propertyArray) Clone() Property {
// 	return NewPropertyArray(p.param).Ok()
// }

// func (p *propertyArray) ToPropertyObject() rusty.Result[PropertyObject] {
// 	return rusty.Err[PropertyObject](fmt.Errorf("not a PropertyObject"))
// }

// Description implements PropertyArray.
func (p *propertyArray) Description() rusty.Optional[string] {
	return p.param.Description
}

// Format implements PropertyArray.
// func (p *propertyArray) Format() rusty.Optional[string] {
// 	return p.param.Format
// }

// Id implements PropertyArray.
// func (p *propertyArray) Id() string {
// 	return p.param.Id
// }

// // Optional implements PropertyArray.
// func (p *propertyArray) Optional() bool {
// 	return p.param.Optional
// }

// // SetOptional implements PropertyArray.
// func (p *propertyArray) SetOptional() {
// 	p.param.Optional = true
// }

// Items implements PropertyArray.
func (p *propertyArray) Items() Property {
	return p.param.Items.Ok()
}

// MaxItems implements PropertyArray.
func (p *propertyArray) MaxItems() rusty.Optional[int] {
	return p.param.MaxItems
}

// MinItems implements PropertyArray.
func (p *propertyArray) MinItems() rusty.Optional[int] {
	return p.param.MinItems
}

// Enum implements PropertyArray.
// func (p *propertyString) Enum() []string {
// 	return p.param.enum
// }

func (p *propertyArray) Type() Type {
	return ARRAY
}

func (p *propertyArray) Id() string {
	return p.param.Id
}

func (p *propertyArray) Ref() rusty.Optional[string] {
	return p.param.Ref
}

// func (p *propertyArray) Runtime() *PropertyRuntime {
// 	return &p.param.Runtime
// }

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

func NewPropertyArray(p PropertyArrayBuilder) rusty.Result[Property] {
	p.Type = ARRAY
	pa := &propertyArray{
		param: p,
		meta:  NewPropertyMeta(),
	}
	if p.Id == "" {
		pa.param.Id = fmt.Sprintf("array-%p", pa)
	}
	return rusty.Ok[Property](pa)
}
