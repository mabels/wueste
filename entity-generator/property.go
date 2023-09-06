package entity_generator

import (
	"fmt"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type Type = string

const (
	OBJECT  Type = "object"
	STRING  Type = "string"
	NUMBER  Type = "number"
	INTEGER Type = "integer"
	BOOLEAN Type = "boolean"
	ARRAY   Type = "array"
)

type PropertyCtx struct {
	Registry *SchemaRegistry
}

type PropertyRuntime struct {
	FileName rusty.Optional[string]
	Ref      rusty.Optional[string]
	BaseDir  rusty.Optional[string]
	Of       Property
}

// func NewRuntime(b SchemaLoader) PropertyMeta {
// 	return PropertyMeta{
// 		Loader: b,
// 	}
// }

func ConnectRuntime[T Property](p T) T {
	p.Runtime().Of = p
	// if p.Runtime().Registry == nil {
	// 	panic("loader not set")
	// }
	return p
}

func (p *PropertyRuntime) SetFileName(name string) {
	p.FileName = rusty.Some(name)
}

func (p *PropertyRuntime) SetRef(name string) {
	p.Ref = rusty.Some(name)
}

func (p *PropertyRuntime) Assign(b PropertyRuntime) *PropertyRuntime {
	// p.Registry = b.Registry
	if b.Ref.IsSome() {
		p.Ref = b.Ref
	}
	if p.FileName.IsNone() {
		p.FileName = b.FileName
	}
	return p
}

func (p *PropertyRuntime) Clone() *PropertyRuntime {
	return (&PropertyRuntime{}).Assign(*p)
}

func (p *PropertyRuntime) ToPropertyObject() rusty.Result[PropertyObject] {
	var pi Property = p.Of
	po, ok := pi.(PropertyObject)
	if !ok {
		return rusty.Err[PropertyObject](fmt.Errorf("not a PropertyObject"))
	}
	return rusty.Ok[PropertyObject](po)
}

type Property interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]

	Ref() rusty.Optional[string]

	Runtime() *PropertyRuntime
}

type PropertyParam struct {
	Id          string
	Type        Type
	Description rusty.Optional[string]

	Ref     rusty.Optional[string]
	Runtime PropertyRuntime
	// Format      rusty.Optional[string]
	// Optional    bool
}

type property struct {
	param PropertyParam
}

func NewProperty(p PropertyParam) Property {
	if p.Ref.IsSome() {
		p.Runtime.Ref = p.Ref
	}
	r := &property{
		param: p,
	}
	return r
}

func (p *property) Ref() rusty.Optional[string] {
	return p.param.Ref
}

func (p *property) Id() string {
	return p.param.Id
}

func (p *property) Runtime() *PropertyRuntime {
	return &p.param.Runtime
}

// func (p *property) Id() string {
// 	return p.param.Id
// }

// func (p *property) Format() rusty.Optional[string] {
// 	return p.param.Format
// }

func (p *property) Type() Type {
	panic("Property Type implement me")
}

// Required implements PropertyString.
// func (p *property) Optional() bool {
// 	return p.param.Optional
// }

// func (p *property) SetOptional() {
// 	p.param.Optional = true
// }

// Description implements PropertyString.
func (p *property) Description() rusty.Optional[string] {
	return p.param.Description
}
