package entity_generator

import (
	"fmt"

	"github.com/iancoleman/orderedmap"
	"github.com/mabels/wueste/entity-generator/rusty"
)

type properties struct {
	oprops orderedmap.OrderedMap
}

func newProperties() *properties {
	return &properties{
		oprops: orderedmap.New(),
	}
}

func (p *properties) Set(name string, property Property) {
	p.oprops.Set(name, property)
}

func (p *properties) Lookup(key string) (Property, bool) {
	v, found := p.oprops.Get(key)
	if !found {
		return nil, false
	}
	return v.(Property), true
}
func (p *properties) Len() int {
	return len(p.oprops.Keys())
}

func (p *properties) Keys() []string {
	return p.oprops.Keys()
}

type PropertyObject interface {
	Type() Type
	Id() string
	Title() string
	Schema() string
	Description() rusty.Optional[string]

	Properties() *properties
	Items() []PropertyItem
	PropertyByName(name string) rusty.Result[PropertyItem]
	Required() []string

	Ref() rusty.Optional[string]
	Meta() PropertyMeta
	// Runtime() *PropertyRuntime
	// Clone() Property
	// Ctx() PropertyCtx
	// Deref() map[string]PropertyItem
}

type propertyObject struct {
	param PropertyObjectBuilder
	meta  PropertyMeta
	// id    string
	// _type Type
	// // fileName string
	// schema string
	// // optional    bool
	// title       string
	// format      rusty.Optional[string]
	// description rusty.Optional[string]
	// ref         rusty.Optional[string]
	// properties  *properties
	// required    []string
	// _runtime    PropertyRuntime
	// _ctx        PropertyCtx
	// deref       map[string]PropertyItem
}

func (p *propertyObject) Meta() PropertyMeta {
	return p.meta
}

// FileName implements PropertyObject.
// func (p *propertyObject) FileName() string {
// 	return p.fileName
// }

// func (s propertyObject) Clone() Property {
// 	return NewPropertyObject(s.param).Ok()
// }

// func (s *propertyObject) Runtime() *PropertyRuntime {
// 	return &s.param.Runtime
// }

// Deref implements PropertyObject.
// func (s *propertyObject) Deref() map[string]PropertyItem {
// 	return s.deref
// }

// Description implements Schema.
func (s *propertyObject) Description() rusty.Optional[string] {
	return s.param.Description
}

// Id implements Schema.
func (s *propertyObject) Id() string {
	return s.param.Id
}

func (s *propertyObject) Schema() string {
	return s.param.Schema
}

// func (s *propertyObject) Format() rusty.Optional[string] {
// 	return s.param.Descri
// }

func (s *propertyObject) Ref() rusty.Optional[string] {
	return s.param.Ref
}

// Properties implements Schema.
func (s *propertyObject) Properties() *properties {
	return s.param.Properties
}

// Required implements Schema.
func (s *propertyObject) Required() []string {
	return s.param.Required
}

// Title implements Schema.
func (s *propertyObject) Title() string {
	return s.param.Title
}

// Type implements Schema.
func (s *propertyObject) Type() string {
	return s.param.Type
}

func (s *propertyObject) PropertyByName(name string) rusty.Result[PropertyItem] {
	v, found := s.param.Properties.Lookup(name)
	if !found {
		return rusty.Err[PropertyItem](fmt.Errorf("property not found:%s", name))
	}
	return NewPropertyObjectItem(name, rusty.Ok(v), -1, isOptional(name, s.param.Required))
}

func (s *propertyObject) Items() []PropertyItem {
	if s.param.Properties == nil {
		return []PropertyItem{}
	}
	items := make([]PropertyItem, 0, s.param.Properties.Len())
	for _, k := range s.param.Properties.Keys() {
		n := s.PropertyByName(k)
		if n.IsErr() {
			panic(n.Err())
		}
		items = append(items, n.Ok())
	}
	return items
}

type PropertyObjectBuilder struct {
	// items       []PropertyItem
	Type        Type
	Description rusty.Optional[string]

	Id         string
	Title      string
	Schema     string
	Properties *properties // PropertiesObject
	Required   []string
	Ref        rusty.Optional[string]

	// Runtime PropertyRuntime
	// Ctx     PropertyCtx
	Errors             []error
	_propertiesBuilder *PropertiesBuilder
}

func NewPropertyObjectBuilder(pb *PropertiesBuilder) *PropertyObjectBuilder {
	return &PropertyObjectBuilder{
		Properties:         newProperties(),
		_propertiesBuilder: pb,
	}
}

// func (b *PropertyObjectParam) fileName(fnam string) *PropertyObjectParam {
// 	b.Runtime.FileName = rusty.Some(fnam)
// 	return b
// }

// func (b *PropertyObjectBuilder) propertiesAdd(pin rusty.Result[PropertyItem]) *PropertyObjectBuilder {
// 	if pin.IsErr() {
// 		b.Errors = append(b.Errors, pin.Err())
// 		return b
// 	}
// 	// ConnectRuntime(rusty.Ok(pin.Ok().Property()))
// 	property := pin.Ok()
// 	// property.Property().Runtime().Assign(b.Runtime)
// 	if b.Properties == nil {
// 		b.Properties = newProperties()
// 	}
// 	b.Properties.Set(property.Name(), property.Property())
// 	if b.Required == nil {
// 		b.Required = []string{}
// 	}
// 	if !property.Optional() {
// 		b.Required = append(b.Required, property.Name())
// 	}
// 	return b
// }

// func (b *PropertyObjectParam) fileName(fnam string) *PropertyObjectParam {
// 	b.FileName = fnam
// 	return b
// }

// func (p *PropertyObjectParam) Items() []PropertyItem {
// 	sort.Slice(p.items, func(i, j int) bool {
// 		return p.items[i].Name() < p.items[j].Name()
// 	})
// 	return p.items
// }

// func (p *PropertyObjectBuilder) description(id string) *PropertyObjectBuilder {
// 	p.Description = rusty.Some(id)
// 	return p
// }

// func (p *PropertyObjectParam) properties(property map[string]any) *PropertyObjectParam {
// 	p.Properties = property
// 	return p
// }

// func (p *PropertyObjectBuilder) id(id string) *PropertyObjectBuilder {
// 	p.Id = id
// 	return p
// }

// func (p *PropertyObjectBuilder) title(title string) *PropertyObjectBuilder {
// 	p.Title = title
// 	return p
// }

// func (p *PropertyObjectParam) schema(schema string) *PropertyObjectParam {
// 	p.Schema = schema
// 	return p
// }

// func (p *PropertyObjectBuilder) required(required []string) *PropertyObjectBuilder {
// 	p.Required = required
// 	return p
// }

// func (b *PropertyObjectBuilder) FromProperty(prop Property) *PropertyObjectBuilder {
// 	po, found := prop.(PropertyObject)
// 	if !found {
// 		b.Errors = append(b.Errors, fmt.Errorf("FromProperty: not a PropertyObject"))
// 		return b
// 	}
// 	for _, v := range po.Items() {
// 		builder := NewPropertiesBuilder(b._propertiesBuilder.ctx)
// 		bProp := builder.FromProperty(v.Property()).Build()
// 		if bProp.IsErr() {
// 			b.Errors = append(b.Errors, bProp.Err())
// 			continue
// 		}
// 		po.Properties().Set(v.Name(), bProp.Ok())
// 	}
// 	return b
// }

func (b *PropertyObjectBuilder) FromJson(js JSONDict) *PropertyObjectBuilder {
	b.Type = OBJECT
	b.Id = getFromAttributeString(js, "$id")
	b.Title = getFromAttributeString(js, "title")
	b.Schema = getFromAttributeString(js, "$schema")
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Properties = newProperties()
	_properties, found := js.Lookup("properties")
	if found {
		properties, found := _properties.(JSONDict)
		if !found {
			b.Errors = append(b.Errors, fmt.Errorf("properties[%s] is not JSONProperty", b.Id))
			return b
		}
		for _, k := range properties.Keys() {
			_v := properties.Get(k)
			v, found := _v.(JSONDict)
			if !found {
				b.Errors = append(b.Errors, fmt.Errorf("properties[%s->%s] is not JSONProperty", b.Id, k))
				continue
			}
			builder := NewPropertiesBuilder(b._propertiesBuilder.ctx)
			builder.parentFileName = b._propertiesBuilder.filename
			r := builder.FromJson(v).Build()
			if r.IsErr() {
				b.Errors = append(b.Errors, r.Err())
			} else {
				b.Properties.Set(k, r.Ok())
			}
		}
	}
	required, found := js.Lookup("required")
	if found {
		stringArray, found := required.([]any)
		if !found {
			stringArray, found := required.([]string)
			if found {
				b.Required = stringArray
			} else {
				b.Errors = append(b.Errors, fmt.Errorf("required[%s] is not []string", b.Id))
				return b
			}
		} else {
			out := make([]string, 0, len(stringArray))
			for _, v := range stringArray {
				vs := coerceString(v)
				if vs.IsNone() {
					b.Errors = append(b.Errors, fmt.Errorf("required[%v] is not string", v))
					return b
				}
				out = append(out, coerceString(v).Value())
			}
			b.Required = out
		}
	}
	b.Ref = getFromAttributeOptionalString(js, "$ref")
	return b
}

func PropertyObjectToJson(b PropertyObject) JSONDict {
	jsp := NewJSONDict()
	JSONsetString(jsp, "type", b.Type())
	// if b.Runtime().FileName.IsSome() {
	// 	JSONsetString("fileName", *b.Runtime().FileName.Value())
	// }
	JSONsetId(jsp, b)
	if b.Title() != "" {
		JSONsetString(jsp, "title", b.Title())
	}
	if b.Schema() != "" {
		JSONsetString(jsp, "$schema", b.Schema())
	}
	JSONsetOptionalString(jsp, "description", b.Description())
	props := NewJSONDict()
	items := b.Items()
	for _, v := range items {
		props.Set(v.Name(), PropertyToJson(v.Property()))
	}
	if props.Len() > 0 {
		jsp.Set("properties", props)
	}
	if len(b.Required()) > 0 {
		jsp.Set("required", b.Required())
	}
	// JSONsetOptionalString("$ref", b.Ref())
	return jsp
}

func (p *PropertyObjectBuilder) Build() rusty.Result[Property] {
	if len(p.Errors) > 0 {
		str := ""
		for _, v := range p.Errors {
			str += v.Error() + "\n"
		}
		return rusty.Err[Property](fmt.Errorf(str))
	}
	po := NewPropertyObject(*p)
	if p._propertiesBuilder.filename.IsSome() {
		po.Ok().Meta().SetFileName(p._propertiesBuilder.filename.Value())
	}
	for _, v := range po.Ok().(PropertyObject).Items() {
		v.Property().Meta().SetMeta(po.Ok())
	}
	return po
}

func NewPropertyObject(p PropertyObjectBuilder) rusty.Result[Property] {
	if !(p.Properties == nil || p.Properties.Len() == 0) && p.Id == "" {
		return rusty.Err[Property](fmt.Errorf("PropertyObject Id is required"))
	}
	p.Type = OBJECT
	r := &propertyObject{
		param: p,
		meta:  NewPropertyMeta(),
		// id:          p.Id,
		// _type:       OBJECT,
		// title:       p.Title,
		// schema:      p.Schema,
		// description: p.Description,
		// properties:  p.Properties,
		// required:    p.Required,
		// ref:         p.Ref,
		// _runtime:    p.Runtime,
		// _ctx:        p.Ctx,
		// deref:       map[string]PropertyItem{},
	}
	return rusty.Ok[Property](r)
}
