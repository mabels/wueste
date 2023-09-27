package entity_generator

import (
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

type PropertyMeta interface {
	Parent() rusty.Optional[Property]
	SetParent(p Property)
	FileName() rusty.Optional[string]
	SetFileName(fn string)
	SetMeta(m Property)
}

type propertyMeta struct {
	parent   rusty.Optional[Property]
	filename rusty.Optional[string]
}

// SetMeta implements PropertyMeta.
func (p *propertyMeta) SetMeta(m Property) {
	if p.FileName().IsNone() {
		p.SetFileName(m.Meta().FileName().Value())
	}
	p.SetParent(m)
}

// FileName implements PropertyMeta.
func (m propertyMeta) FileName() rusty.Optional[string] {
	return m.filename
}

// Parent implements PropertyMeta.
func (m propertyMeta) Parent() rusty.Optional[Property] {
	return m.parent
}

// SetFileName implements PropertyMeta.
func (m *propertyMeta) SetFileName(fn string) {
	m.filename = rusty.Some(fn)
}

// SetParent implements PropertyMeta.
func (m *propertyMeta) SetParent(p Property) {
	m.parent = rusty.Some(p)
}

func NewPropertyMeta() PropertyMeta {
	return &propertyMeta{}
}

//
// FileName rusty.Optional[string]
// parent   rusty.Optional[Property]
// }

type Property interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]
	Ref() rusty.Optional[string]
	Meta() PropertyMeta
}

type PropertyBuilder struct {
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Ref         rusty.Optional[string]
}

type property struct {
	param PropertyBuilder
	meta  PropertyMeta
}

func NewProperty(p PropertyBuilder) Property {
	// if p.Ref.IsSome() {
	// 	p.Runtime.Ref = p.Ref
	// }
	r := &property{
		param: p,
		meta:  NewPropertyMeta(),
	}
	return r
}

// func (p property) Clone() Property {
// 	return NewProperty(p.param)
// }

func (p *property) Ref() rusty.Optional[string] {
	return p.param.Ref
}

func (p *property) Id() string {
	return p.param.Id
}

func (p *property) Meta() PropertyMeta {
	return p.meta
}

// func (p *property) Runtime() *PropertyRuntime {
// 	return &p.param.Runtime
// }

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
