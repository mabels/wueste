package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type Type = string

const (
	OBJECT     Type = "object"
	OBJECTITEM Type = "objectitem"
	STRING     Type = "string"
	NUMBER     Type = "number"
	INTEGER    Type = "integer"
	BOOLEAN    Type = "boolean"
	ARRAY      Type = "array"
	ARRAYITEM  Type = "arrayitem"
)

type PropertyMeta interface {
	Parent() rusty.Optional[Property]
	SetParent(p Property) PropertyMeta
	FileName() rusty.Optional[string]
	SetFileName(fn string) PropertyMeta
	SetMeta(m Property) PropertyMeta
}

type propertyMeta struct {
	parent   rusty.Optional[Property]
	filename rusty.Optional[string]
}

// SetMeta implements PropertyMeta.
func (p *propertyMeta) SetMeta(m Property) PropertyMeta {
	if p.FileName().IsNone() {
		if m.Meta().FileName().IsSome() {
			p.SetFileName(m.Meta().FileName().Value())
		}
	}
	p.SetParent(m)
	return p
}

// FileName implements PropertyMeta.
func (m propertyMeta) FileName() rusty.Optional[string] {
	if m.filename.IsNone() {
		if m.parent.IsSome() {
			return m.parent.Value().Meta().FileName()
		}
	}
	return m.filename
}

// Parent implements PropertyMeta.
func (m propertyMeta) Parent() rusty.Optional[Property] {
	return m.parent
}

// SetFileName implements PropertyMeta.
func (m *propertyMeta) SetFileName(fn string) PropertyMeta {
	m.filename = rusty.Some(fn)
	return m
}

// SetParent implements PropertyMeta.
func (m *propertyMeta) SetParent(p Property) PropertyMeta {
	m.parent = rusty.Some(p)
	return m
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
	XProperties() map[string]interface{}
	Meta() PropertyMeta
}

type PropertyFormat interface {
	Format() rusty.Optional[string]
}

type PropertyBuilder struct {
	Id          string
	Type        Type
	Description rusty.Optional[string]
	XProperties map[string]interface{}
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
	if p.XProperties == nil {
		p.XProperties = make(map[string]interface{})
	}
	r := &property{
		param: p,
		meta:  NewPropertyMeta(),
	}
	return r
}

// func (p property) Clone() Property {
// 	return NewProperty(p.param)
// }

func (p *property) XProperties() map[string]interface{} {
	return p.param.XProperties
}

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
