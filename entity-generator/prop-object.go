package entity_generator

import (
	"fmt"
	"sort"

	"github.com/iancoleman/orderedmap"
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyItem interface {
	Property() Property
	Optional() bool
	// SetOrder(order int)
	// Order() int
	Name() string
}

type properties struct {
	oprops orderedmap.OrderedMap
}

func newProperties() *properties {
	return &properties{
		oprops: *orderedmap.New(),
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
	Runtime() *PropertyRuntime
	// Ctx() PropertyCtx
	// Deref() map[string]PropertyItem
}

// type SchemaBuilder struct {
// 	id          string
// 	fileName    string
// 	typ         Type
// 	title       string
// 	schema      string
// 	format      *string
// 	ref         *string
// 	description *string
// 	properties  PropertiesObject
// 	required    []string
// 	// derefs      map[string]PropertyItem
// 	__loader SchemaLoader
// }

// func NewSchemaBuilder(loader SchemaLoader) *SchemaBuilder {
// 	return &SchemaBuilder{
// 		// derefs:   map[string]PropertyItem{},
// 		__loader: loader,
// 	}
// }

// func (b *SchemaBuilder) ResolveRef(v *JSONProperty) (*JSONProperty, error) {
// 	if v.Ref != nil {
// 		ref := strings.TrimSpace(*v.Ref)
// 		if ref[0] == '#' {
// 			return nil, fmt.Errorf("local ref not supported")
// 		}
// 		if !strings.HasPrefix(ref, "file://") {
// 			return nil, fmt.Errorf("only file:// ref supported")
// 		}
// 		fname := ref[len("file://"):]
// 		if !strings.HasSuffix(fname, "/") {
// 			var err error
// 			fname, err = b.__loader.Abs(path.Join(path.Dir(b.fileName), fname))
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 		pl, err := LoadSchema(fname, b.__loader)
// 		if err != nil {
// 			return nil, err
// 		}
// 		p := pl.(PropertyObject)
// 		// pref := b.__loader.SchemaRegistry().AddSchema(po)

// 		// p := pref.Property().(PropertyObject)
// 		myv := JSONProperty{}
// 		myv.FileName = fname
// 		myv.Id = p.Id()
// 		myv.Schema = p.Schema()
// 		myv.Title = p.Title()
// 		myv.Type = p.Type()
// 		myv.Description = rusty.OptionalToPtr(p.Description())
// 		myv.Properties = PropertiesToJson(p.Properties())
// 		myv.Required = p.Required()
// 		myv.Ref = rusty.OptionalToPtr(p.Ref())
// 		return &myv, nil
// 	} else if v.Type == "object" && v.Properties != nil {
// 		// register schema
// 		NewSchemaBuilder(b.__loader).JSON2PropertyObject(v.JSONSchema).Build()
// 	}
// 	return v, nil
// }

// func addJSONProperty(v *JSONProperty, loader SchemaLoader) Property {
// 	switch v.Type {
// 	case "string":
// 		p := PropertyStringParam{
// 			PropertyParam: PropertyParam{
// 				Type: v.Type,
// 			},
// 		}
// 		p.Default = coerceString(v.Default)
// 		if v.Format != nil {
// 			p.Format = rusty.Some(*v.Format)
// 		}
// 		return NewPropertyString(p)
// 	case "number":
// 		if v.Format == nil {
// 			panic("number format is required")
// 		}
// 		switch *v.Format {
// 		case "float32":
// 			p := PropertyNumberParam[float32]{
// 				Format: rusty.Some("float32"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceFloat32(v.Default)
// 			return NewPropertyNumber(p)
// 		case "float64":
// 			p := PropertyNumberParam[float64]{
// 				Format: rusty.Some("float64"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceFloat64(v.Default)
// 			return NewPropertyNumber(p)
// 		default:
// 			panic("unknown format")
// 		}
// 	case "integer":
// 		switch *v.Format {
// 		case "int":
// 			p := PropertyIntegerParam[int]{
// 				Format: rusty.Some("int"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceInt(v.Default)
// 			return NewPropertyInteger(p)
// 		case "int8":
// 			p := PropertyIntegerParam[int8]{
// 				Format: rusty.Some("int8"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceInt8(v.Default)
// 			return NewPropertyInteger(p)
// 		case "int16":
// 			p := PropertyIntegerParam[int16]{
// 				Format: rusty.Some("int16"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceInt16(v.Default)
// 			return NewPropertyInteger(p)
// 		case "int32":
// 			p := PropertyIntegerParam[int32]{
// 				Format: rusty.Some("int32"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceInt32(v.Default)
// 			return NewPropertyInteger(p)
// 		case "int64":
// 			p := PropertyIntegerParam[int64]{
// 				Format: rusty.Some("int64"),
// 				Type:   v.Type,
// 			}
// 			p.Default = coerceInt64(v.Default)
// 			return NewPropertyInteger(p)
// 		default:
// 			panic("unknown format")
// 		}
// 	case "boolean":
// 		p := PropertyBooleanParam{
// 			Type: v.Type,
// 		}
// 		p.Default = coerceBool(v.Default)
// 		return NewPropertyBoolean(p)
// 	case "object":
// 		return NewPropertyObject(PropertyObjectParam{
// 			FileName:    v.FileName,
// 			Id:          v.Id,
// 			Title:       v.Title,
// 			Schema:      v.Schema,
// 			Description: rusty.OptionalFromPtr(v.Description),
// 			Properties:  NewPropertiesBuilder(loader).FromJson(v.Properties, v.Required).Build(),
// 			Required:    v.Required,
// 			Ref:         rusty.OptionalFromPtr(v.Ref),
// 		})
// 	default:
// 		panic("unknown type")
// 	}
// }

// func (b *SchemaBuilder) JSON2PropertyObject(js JSONSchema) *SchemaBuilder {
// 	x := b.Id(js.Id).
// 		FileName(js.FileName).
// 		Type(js.Type).
// 		Title(js.Title).
// 		Required(js.Required)
// 	if js.Description != nil {
// 		x.Description(*js.Description)
// 	}
// 	pb := NewPropertiesBuilder(b.__loader)
// 	for k, vin := range js.Properties {
// 		v, err := x.ResolveRef(&vin)
// 		if err != nil {
// 			panic(err)
// 		}
// 		pb.Add(NewPropertyItem(k, addJSONProperty(v, b.__loader)))
// 	}
// 	x.Properties(pb.Build())
// 	return x
// }

// func (b *SchemaBuilder) Id(id string) *SchemaBuilder {
// 	b.id = id
// 	return b
// }

// func (b *SchemaBuilder) FileName(fname string) *SchemaBuilder {
// 	b.fileName = fname
// 	return b
// }

// func (b *SchemaBuilder) Type(_type Type) *SchemaBuilder {
// 	b.typ = _type
// 	return b
// }

// func (b *SchemaBuilder) Schema(_schema string) *SchemaBuilder {
// 	b.schema = _schema
// 	return b
// }

// func (b *SchemaBuilder) Title(title string) *SchemaBuilder {
// 	b.title = title
// 	return b
// }

// func (b *SchemaBuilder) Description(description string) *SchemaBuilder {
// 	b.description = &description
// 	return b
// }

// func (b *SchemaBuilder) Properties(properties PropertiesObject) *SchemaBuilder {
// 	b.properties = properties
// 	return b
// }

// func (b *SchemaBuilder) Required(required []string) *SchemaBuilder {
// 	b.required = required
// 	return b
// }

type propertyObject struct {
	id    string
	_type Type
	// fileName string
	schema string
	// optional    bool
	title       string
	format      rusty.Optional[string]
	description rusty.Optional[string]
	ref         rusty.Optional[string]
	properties  *properties
	required    []string
	_runtime    PropertyRuntime
	_ctx        PropertyCtx
	// deref       map[string]PropertyItem
}

// FileName implements PropertyObject.
// func (p *propertyObject) FileName() string {
// 	return p.fileName
// }

func (s *propertyObject) Runtime() *PropertyRuntime {
	return &s._runtime
}

// Deref implements PropertyObject.
// func (s *propertyObject) Deref() map[string]PropertyItem {
// 	return s.deref
// }

// Description implements Schema.
func (s *propertyObject) Description() rusty.Optional[string] {
	return s.description
}

// Id implements Schema.
func (s *propertyObject) Id() string {
	return s.id
}

func (s *propertyObject) Schema() string {
	return s.schema
}

func (s *propertyObject) Format() rusty.Optional[string] {
	return s.format
}

func (s *propertyObject) Ref() rusty.Optional[string] {
	return s.ref
}

// Properties implements Schema.
func (s *propertyObject) Properties() *properties {
	return s.properties
}

// Required implements Schema.
func (s *propertyObject) Required() []string {
	return s.required
}

// Title implements Schema.
func (s *propertyObject) Title() string {
	return s.title
}

// Type implements Schema.
func (s *propertyObject) Type() string {
	return s._type
}

func (s *propertyObject) PropertyByName(name string) rusty.Result[PropertyItem] {
	v, found := s.properties.Lookup(name)
	if !found {
		return rusty.Err[PropertyItem](fmt.Errorf("property not found:%s", name))
	}
	return NewPropertyItem(name, rusty.Ok(v), isOptional(name, s.required))
}

func (s *propertyObject) Items() []PropertyItem {
	if s.properties == nil {
		return []PropertyItem{}
	}
	items := make([]PropertyItem, 0, s.properties.Len())
	for _, k := range s.properties.Keys() {
		n := s.PropertyByName(k)
		if n.IsErr() {
			panic(n.Err())
		}
		items = append(items, n.Ok())
	}
	return items
}

// func (b *SchemaBuilder) Build() Property {
// 	requiredMap := make(map[string]bool)
// 	for _, v := range b.required {
// 		requiredMap[v] = true
// 	}
// 	for _, v := range b.properties.Items() {
// 		_, found := requiredMap[v.Name()]
// 		if !found {
// 			panic("setOptional not found")
// 			// v.Property().SetOptional()
// 		}
// 	}
// 	// desc := rusty.None[string]()
// 	// if b.description != nil {
// 	// 	desc = rusty.Some[string](*b.description)
// 	// }
// 	ret := &propertyObject{
// 		id:          b.id,
// 		fileName:    b.fileName,
// 		format:      rusty.OptionalFromPtr(b.format),
// 		_type:       b.typ,
// 		title:       b.title,
// 		description: rusty.OptionalFromPtr(b.description),
// 		properties:  b.properties,
// 		required:    b.required,
// 		ref:         rusty.OptionalFromPtr(b.ref),
// 		// deref:       b.deref,
// 	}
// 	b.__loader.SchemaRegistry().AddSchema(ret)
// 	return ret
// }

type PropertyObjectParam struct {
	items       []PropertyItem
	Type        Type
	Description rusty.Optional[string]

	Id         string
	Title      string
	Schema     string
	Properties *properties // PropertiesObject
	Required   []string
	Ref        rusty.Optional[string]

	Runtime PropertyRuntime
	Ctx     PropertyCtx
	Errors  []error
}

func (b *PropertyObjectParam) fileName(fnam string) *PropertyObjectParam {
	b.Runtime.FileName = rusty.Some(fnam)
	return b
}

func (b *PropertyObjectParam) propertiesAdd(pin rusty.Result[PropertyItem]) *PropertyObjectParam {
	if pin.IsErr() {
		b.Errors = append(b.Errors, pin.Err())
		return b
	}
	ConnectRuntime(rusty.Ok(pin.Ok().Property()))
	property := pin.Ok()
	// property.SetOrder(len(b.items))
	property.Property().Runtime().Assign(b.Runtime)
	if b.Properties == nil {
		b.Properties = newProperties()
	}
	b.Properties.Set(property.Name(), property.Property())
	if b.Required == nil {
		b.Required = []string{}
	}
	if !property.Optional() {
		b.Required = append(b.Required, property.Name())
	}
	return b
}

// func (b *PropertyObjectParam) fileName(fnam string) *PropertyObjectParam {
// 	b.FileName = fnam
// 	return b
// }

func (p *PropertyObjectParam) Items() []PropertyItem {
	sort.Slice(p.items, func(i, j int) bool {
		return p.items[i].Name() < p.items[j].Name()
	})
	return p.items
}

func (p *PropertyObjectParam) description(id string) *PropertyObjectParam {
	p.Description = rusty.Some(id)
	return p
}

// func (p *PropertyObjectParam) properties(property map[string]any) *PropertyObjectParam {
// 	p.Properties = property
// 	return p
// }

func (p *PropertyObjectParam) id(id string) *PropertyObjectParam {
	p.Id = id
	return p
}

func (p *PropertyObjectParam) title(title string) *PropertyObjectParam {
	p.Title = title
	return p
}

// func (p *PropertyObjectParam) schema(schema string) *PropertyObjectParam {
// 	p.Schema = schema
// 	return p
// }

func (p *PropertyObjectParam) required(required []string) *PropertyObjectParam {
	p.Required = required
	return p
}

// func (p *PropertyObjectParam) ref(ref string) *PropertyObjectParam {
// 	p.Ref = rusty.Some(ref)
// 	return p
// }

// func toJSONProperty(js JSONProperty) map[string]JSONProperty {
// 	ret := map[string]JSONProperty{}
// 	for k, v := range js {
// 		x, found := v.(JSONProperty)
// 		if found {
// 			ret[k] = x
// 			continue
// 		}
// 	}
// 	return ret
// 	// 	x, found = v.(map[string]interface{})
// 	// 	if found {
// 	// 		ret[k] = x
// 	// 		continue
// 	// 	}
// 	// 	panic("toJSONProperty:unknown type")
// 	// }
// 	// return ret
// }

// func toMapStringJSONProperty(js any) JSONProperty {
// 	isJSONProperty, found := js.(JSONProperty)
// 	if found {
// 		return toJSONProperty(isJSONProperty)
// 	}
// 	// isMapStringInterface, found := js.(map[string]interface{})
// 	// if found {
// 	// 	// return toJSONProperty(isMapStringInterface)
// 	// }
// 	panic("unknown type")
// }

func (b *PropertyObjectParam) FromJson(js JSONProperty) *PropertyObjectParam {
	b.Type = OBJECT
	b.Id = getFromAttributeString(js, "$id")
	b.Title = getFromAttributeString(js, "title")
	b.Schema = getFromAttributeString(js, "$schema")
	b.Description = getFromAttributeOptionalString(js, "description")
	b.Properties = newProperties()
	_properties, found := js.Lookup("properties")
	if found {
		properties, found := _properties.(JSONProperty)
		if !found {
			panic(fmt.Errorf("properties[%s] is not JSONProperty", b.Id))
		}
		for _, k := range properties.Keys() {
			_v := properties.Get(k)
			v, found := _v.(JSONProperty)
			if !found {
				b.Errors = append(b.Errors, fmt.Errorf("properties[%s->%s] is not JSONProperty", b.Id, k))
				continue
			}
			r := NewPropertiesBuilder(b.Ctx).FromJson(b.Runtime, v).Build()
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
				panic(fmt.Errorf("required[%s] is not []string", b.Id))
			}
		} else {
			out := make([]string, 0, len(stringArray))
			for _, v := range stringArray {
				out = append(out, coerceString(v).Value())
			}
			b.Required = out
		}
	}
	b.Ref = getFromAttributeOptionalString(js, "$ref")
	return b
}

func PropertyObjectToJson(b PropertyObject) JSONProperty {
	jsp := NewJSONProperty()
	JSONsetString(jsp, "type", b.Type())
	// if b.Runtime().FileName.IsSome() {
	// 	JSONsetString("fileName", *b.Runtime().FileName.Value())
	// }
	JSONsetString(jsp, "$id", b.Id())
	JSONsetString(jsp, "title", b.Title())
	if b.Schema() != "" {
		JSONsetString(jsp, "$schema", b.Schema())
	}
	JSONsetOptionalString(jsp, "description", b.Description())
	props := NewJSONProperty()
	items := b.Items()
	for _, v := range items {
		props.Set(v.Name(), PropertyToJson(v.Property()))
	}
	jsp.Set("properties", props)
	if len(b.Required()) > 0 {
		jsp.Set("required", b.Required())
	}
	// JSONsetOptionalString("$ref", b.Ref())
	return jsp
}

func (p *PropertyObjectParam) Build() rusty.Result[Property] {
	if len(p.Errors) > 0 {
		str := ""
		for _, v := range p.Errors {
			str += v.Error() + "\n"
		}
		return rusty.Err[Property](fmt.Errorf(str))
	}
	return ConnectRuntime(NewPropertyObject(*p))
}

func NewPropertyObject(p PropertyObjectParam) rusty.Result[Property] {
	if !(p.Properties == nil || p.Properties.Len() == 0) && p.Id == "" {
		return rusty.Err[Property](fmt.Errorf("PropertyObject Id is required"))
	}
	r := &propertyObject{
		id:          p.Id,
		_type:       OBJECT,
		title:       p.Title,
		schema:      p.Schema,
		description: p.Description,
		properties:  p.Properties,
		required:    p.Required,
		ref:         p.Ref,
		_runtime:    p.Runtime,
		_ctx:        p.Ctx,
		// deref:       map[string]PropertyItem{},
	}
	return rusty.Ok[Property](r)
}

type propertyItem struct {
	name     string
	optional bool
	property Property
	// order    int
}

// Description implements PropertyItem.
func (pi *propertyItem) Name() string {
	return pi.name
}

// Optional implements PropertyItem.
func (pi *propertyItem) Optional() bool {
	return pi.optional
}

// func (pi *propertyItem) Order() int {
// 	return pi.order
// }

// func (pi *propertyItem) SetOrder(order int) {
// 	pi.order = order
// }

func (pi *propertyItem) Property() Property {
	return pi.property
}

func NewPropertyItem(name string, property rusty.Result[Property], optionals ...bool) rusty.Result[PropertyItem] {
	if property.IsErr() {
		return rusty.Err[PropertyItem](property.Err())
	}
	optional := true
	if len(optionals) > 0 {
		optional = optionals[0]
	}
	return rusty.Ok[PropertyItem](&propertyItem{
		name:     name,
		optional: optional,
		property: property.Ok(),
	})
}
