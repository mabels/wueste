package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyBoolean interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]
	Default() rusty.Optional[bool] // match Type

	Ref() rusty.Optional[string]
	Runtime() *PropertyRuntime
}

type PropertyBooleanParam struct {
	// __loader SchemaLoader
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Default     rusty.Optional[bool]
	Ref         rusty.Optional[string]

	Runtime PropertyRuntime
	Ctx     PropertyCtx
	// Optional    bool
}

func (b *PropertyBooleanParam) FromJson(rt PropertyRuntime, js JSONProperty) *PropertyBooleanParam {
	// b.Id = getFromAttributeString(js, "$id")
	b.Type = "boolean"
	b.Runtime.Assign(rt)
	ensureAttributeId(js, func(id string) { b.Id = id })
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Default = getFromAttributeOptionalBoolean(js, "default")
	return b
}

func PropertyBooleanToJson(b PropertyBoolean) JSONProperty {
	jsp := NewJSONProperty()
	JSONsetId(jsp, b)
	JSONsetString(jsp, "type", b.Type())
	JSONsetOptionalString(jsp, "description", b.Description())
	JSONsetOptionalBoolean(jsp, "default", b.Default())
	return jsp
}

func (b *PropertyBooleanParam) Build() rusty.Result[Property] {
	return ConnectRuntime(NewPropertyBoolean(*b))
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

func (p *propertyBoolean) Runtime() *PropertyRuntime {
	return &p.param.Runtime
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
		return rusty.Some[bool](p.param.Default.Value())

	}
	return rusty.None[bool]()
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

func NewPropertyBoolean(p PropertyBooleanParam) rusty.Result[Property] {
	p.Type = BOOLEAN
	return rusty.Ok[Property](&propertyBoolean{
		param: p,
	})
}
