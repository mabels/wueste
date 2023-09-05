package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyBoolean interface {
	// Id() string
	Type() Type
	Description() rusty.Optional[string]
	// Optional() bool
	// SetOptional()
	// Format() rusty.Optional[string]
	// Property
	Default() rusty.Optional[bool] // match Type
}

type PropertyBooleanParam struct {
	__loader    SchemaLoader
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Default     rusty.Optional[bool]
	// Optional    bool
}

func (b *PropertyBooleanParam) FromJson(js JSONProperty) *PropertyBooleanParam {
	b.Id = getFromAttributeString(js, "$id")
	b.Type = "boolean"
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Default = getFromAttributeOptionalBoolean(js, "default")
	return b
}

func PropertyBooleanToJson(b PropertyBoolean) JSONProperty {
	jsp := JSONProperty{}
	jsp.setString("type", b.Type())
	jsp.setOptionalString("description", b.Description())
	jsp.setOptionalBoolean("default", b.Default())
	return jsp
}

func (b *PropertyBooleanParam) Build() PropertyBoolean {
	return NewPropertyBoolean(*b)
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

func (p *propertyBoolean) Default() rusty.Optional[bool] {
	if p.param.Default.IsSome() {
		// lit := wueste.BoolLiteral(*p.param.Default.Value())
		return rusty.Some[bool](*p.param.Default.Value())

	}
	return rusty.None[bool]()
}

func (p *propertyBoolean) Type() Type {
	return BOOLEAN
}

func NewPropertyBoolean(p PropertyBooleanParam) PropertyBoolean {
	p.Type = BOOLEAN
	return &propertyBoolean{
		param: p,
	}
}
