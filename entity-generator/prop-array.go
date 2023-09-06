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
	Runtime() *PropertyRuntime

	ToPropertyObject() rusty.Result[PropertyObject]
}

type PropertyArrayParam struct {
	// __loader SchemaLoader
	Id          string
	Type        Type
	Ref         rusty.Optional[string]
	Description rusty.Optional[string]
	// Format      rusty.Optional[string]
	// Optional    bool
	MinItems rusty.Optional[int]
	MaxItems rusty.Optional[int]
	Items    Property
	Runtime  PropertyRuntime
	Ctx      PropertyCtx
	// Default rusty.Optional[string]
	// Enum      []string
	// MinLength rusty.Optional[int]
	// MaxLength rusty.Optional[int]
	// Format    rusty.Optional[StringFormat]
}

func (b *PropertyArrayParam) FromJson(rt PropertyRuntime, js JSONProperty) *PropertyArrayParam {
	b.Type = ARRAY
	b.Runtime.Assign(rt)
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.MaxItems = getFromAttributeOptionalInt(js, "maxItems")
	b.MinItems = getFromAttributeOptionalInt(js, "minItems")
	b.Items = NewPropertiesBuilder(b.Ctx).FromJson(rt, js.Get("items").(JSONProperty)).Build()
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

func (b *PropertyArrayParam) Build() PropertyArray {
	return ConnectRuntime(NewPropertyArray(*b))
}

type propertyArray struct {
	param PropertyArrayParam
}

func (p *propertyArray) ToPropertyObject() rusty.Result[PropertyObject] {
	return rusty.Err[PropertyObject](fmt.Errorf("not a PropertyObject"))
}

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
	return p.param.Items
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

func (p *propertyArray) Runtime() *PropertyRuntime {
	return &p.param.Runtime
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
	p.Type = ARRAY
	return &propertyArray{
		param: p,
	}
}
