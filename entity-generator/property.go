package entity_generator

import "github.com/mabels/wueste/entity-generator/rusty"

type Type = string

const (
	OBJECT  Type = "object"
	STRING  Type = "string"
	NUMBER  Type = "number"
	INTEGER Type = "integer"
	BOOLEAN Type = "boolean"
	ARRAY   Type = "array"
)

type Property interface {
	Id() string
	Type() Type
	Description() rusty.Optional[string]
	Format() rusty.Optional[string]
	Optional() bool
	SetOptional()
}

type PropertyParam struct {
	Id          string
	Type        Type
	Description rusty.Optional[string]
	Format      rusty.Optional[string]
	Optional    bool
}

type property struct {
	param PropertyParam
}

func NewProperty(p PropertyParam) Property {
	return &property{
		param: p,
	}
}

func (p *property) Id() string {
	return p.param.Id
}

func (p *property) Format() rusty.Optional[string] {
	return p.param.Format
}

func (p *property) Type() Type {
	panic("implement me")
}

// Required implements PropertyString.
func (p *property) Optional() bool {
	return p.param.Optional
}

func (p *property) SetOptional() {
	p.param.Optional = true
}

// Description implements PropertyString.
func (p *property) Description() rusty.Optional[string] {
	return p.param.Description
}
