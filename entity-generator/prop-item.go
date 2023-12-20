package entity_generator

import "github.com/mabels/wueste/entity-generator/rusty"

type PropertyItem interface {
	Property() Property
	Optional() bool
	Idx() int
	Name() string

	Id() string
	Type() Type
	Description() rusty.Optional[string]
	XProperties() map[string]interface{}
	Ref() rusty.Optional[string]
	Meta() PropertyMeta
}

type propertyItem struct {
	typ      Type
	name     string
	optional bool
	property Property
	idx      int
	// order    int
}

func (pi *propertyItem) XProperties() map[string]interface{} {
	return nil
}

func (pi *propertyItem) Description() rusty.Optional[string] {
	panic("propertyItem:Description: implement me")
}
func (pi *propertyItem) Id() string {
	panic("propertyItem:Id: implement me")
}
func (pi *propertyItem) Type() Type {
	return pi.typ
}
func (pi *propertyItem) Ref() rusty.Optional[string] {
	panic("propertyItem:Ref: implement me")
}
func (pi *propertyItem) Meta() PropertyMeta {
	panic("propertyItem:Meta: implement me")
}

func (pi *propertyItem) Idx() int {
	return pi.idx
}
func (pi *propertyItem) Name() string {
	return pi.name
}

func (pi *propertyItem) Optional() bool {
	return pi.optional
}

func (pi *propertyItem) Property() Property {
	return pi.property
}

func NewPropertyObjectItem(name string, property rusty.Result[Property], idx int, optionals ...bool) rusty.Result[PropertyItem] {
	if property.IsErr() {
		return rusty.Err[PropertyItem](property.Err())
	}
	optional := true
	if len(optionals) > 0 {
		optional = optionals[0]
	}
	return rusty.Ok[PropertyItem](&propertyItem{
		typ:      OBJECTITEM,
		idx:      idx,
		name:     name,
		optional: optional,
		property: property.Ok(),
	})
}

func NewPropertyArrayItem(name string, property rusty.Result[Property], optionals ...bool) rusty.Result[PropertyItem] {
	if property.IsErr() {
		return rusty.Err[PropertyItem](property.Err())
	}
	optional := true
	if len(optionals) > 0 {
		optional = optionals[0]
	}
	return rusty.Ok[PropertyItem](&propertyItem{
		typ:      ARRAYITEM,
		name:     name,
		optional: optional,
		property: property.Ok(),
	})
}
