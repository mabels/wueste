package entity_generator

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyItem interface {
	Property() Property
	// SetOrder(order int)
	// Order() int
	Name() string
}
type PropertiesObject interface {
	Items() []PropertyItem
}

type PropertyObject interface {
	Property
	Type() Type
	Id() string
	Title() string
	Schema() string
	Ref() rusty.Optional[string]
	Description() rusty.Optional[string]
	FileName() string

	Properties() PropertiesObject
	Required() []string
	// Deref() map[string]PropertyItem
}

type SchemaBuilder struct {
	id          string
	fileName    string
	typ         Type
	title       string
	schema      string
	format      *string
	ref         *string
	description *string
	properties  PropertiesObject
	required    []string
	// derefs      map[string]PropertyItem
	__loader SchemaLoader
}

func NewSchemaBuilder(loader SchemaLoader) *SchemaBuilder {
	return &SchemaBuilder{
		// derefs:   map[string]PropertyItem{},
		__loader: loader,
	}
}

func (b *SchemaBuilder) ResolveRef(v *JSONProperty) (*JSONProperty, error) {
	if v.Ref != nil {
		ref := strings.TrimSpace(*v.Ref)
		if ref[0] == '#' {
			return nil, fmt.Errorf("local ref not supported")
		}
		if !strings.HasPrefix(ref, "file://") {
			return nil, fmt.Errorf("only file:// ref supported")
		}
		fname := ref[len("file://"):]
		if !strings.HasSuffix(fname, "/") {
			var err error
			fname, err = b.__loader.Abs(path.Join(path.Dir(b.fileName), fname))
			if err != nil {
				return nil, err
			}
		}
		pl, err := LoadSchema(fname, b.__loader)
		if err != nil {
			return nil, err
		}
		p := pl.(PropertyObject)
		// pref := b.__loader.SchemaRegistry().AddSchema(po)

		// p := pref.Property().(PropertyObject)
		myv := JSONProperty{}
		myv.FileName = fname
		myv.Id = p.Id()
		myv.Schema = p.Schema()
		myv.Title = p.Title()
		myv.Type = p.Type()
		myv.Description = rusty.OptionalToPtr(p.Description())
		myv.Properties = PropertiesToJson(p.Properties())
		myv.Required = p.Required()
		myv.Ref = rusty.OptionalToPtr(p.Ref())
		return &myv, nil
	} else if v.Type == "object" && v.Properties != nil {
		// register schema
		NewSchemaBuilder(b.__loader).JSON2PropertyObject(v.JSONSchema).Build()
	}
	return v, nil
}

func coerceString(v interface{}) rusty.Optional[string] {
	if v != nil {
		switch v.(type) {
		case string:
			return rusty.Some[string](v.(string))
		case bool:
			val := v.(bool)
			if val {
				return rusty.Some[string]("true")
			} else {
				return rusty.Some[string]("false")
			}
		case int, int8, int16, int32, int64, float32, float64:
			return rusty.Some[string](fmt.Sprintf("%v", v))
		default:
			panic(fmt.Errorf("unknown type %T", v))
		}
	}
	return rusty.None[string]()

}
func coerceBool(v interface{}) rusty.Optional[bool] {
	if v != nil {
		switch v.(type) {
		case string:
			val := strings.ToLower(v.(string))
			if val == "true" || val == "on" || val == "yes" {
				return rusty.Some[bool](true)
			}
			return rusty.Some[bool](false)
		case bool:
			val := v.(bool)
			return rusty.Some[bool](val)
		case int, int8, int16, int32, int64, float32, float64:
			val := coerceInt(v)
			if val.IsNone() {
				return rusty.None[bool]()
			}
			return rusty.Some[bool](*val.Value() != 0)
		default:
			panic(fmt.Errorf("unknown type %T", v))
		}

	}
	return rusty.None[bool]()
}

func coerceNumber[T int | int8 | int16 | int32 | int64 | float32 | float64](v interface{}) rusty.Optional[T] {
	if v != nil {
		switch v.(type) {
		case string:
			f, err := strconv.ParseFloat(v.(string), 8)
			if err != nil {
				return rusty.None[T]()
			}
			return rusty.Some[T](T(f))
		case int:
			return rusty.Some[T](T(v.(int)))
		case int8:
			return rusty.Some[T](T(v.(int8)))
		case int16:
			return rusty.Some[T](T(v.(int16)))
		case int32:
			return rusty.Some[T](T(v.(int32)))
		case int64:
			return rusty.Some[T](T(v.(int64)))
		case float32:
			return rusty.Some[T](T(v.(float32)))
		case float64:
			return rusty.Some[T](T(v.(float64)))
		default:
			panic(fmt.Errorf("unknown type %T", v))
		}
	}
	return rusty.None[T]()
}
func coerceInt(v interface{}) rusty.Optional[int] {
	return coerceNumber[int](v)
}
func coerceInt8(v interface{}) rusty.Optional[int8] {
	return coerceNumber[int8](v)
}
func coerceInt16(v interface{}) rusty.Optional[int16] {
	return coerceNumber[int16](v)
}
func coerceInt32(v interface{}) rusty.Optional[int32] {
	return coerceNumber[int32](v)
}
func coerceInt64(v interface{}) rusty.Optional[int64] {
	return coerceNumber[int64](v)
}
func coerceFloat32(v interface{}) rusty.Optional[float32] {
	return coerceNumber[float32](v)
}
func coerceFloat64(v interface{}) rusty.Optional[float64] {
	return coerceNumber[float64](v)
}

func addJSONProperty(v *JSONProperty, loader SchemaLoader) Property {
	switch v.Type {
	case "string":
		p := PropertyStringParam{
			PropertyParam: PropertyParam{
				Type: v.Type,
			},
		}
		p.Default = coerceString(v.Default)
		if v.Format != nil {
			p.Format = rusty.Some(*v.Format)
		}
		return NewPropertyString(p)
	case "number":
		if v.Format == nil {
			panic("number format is required")
		}
		switch *v.Format {
		case "float32":
			p := PropertyNumberParam[float32]{
				Format: rusty.Some("float32"),
				Type:   v.Type,
			}
			p.Default = coerceFloat32(v.Default)
			return NewPropertyNumber(p)
		case "float64":
			p := PropertyNumberParam[float64]{
				Format: rusty.Some("float64"),
				Type:   v.Type,
			}
			p.Default = coerceFloat64(v.Default)
			return NewPropertyNumber(p)
		default:
			panic("unknown format")
		}
	case "integer":
		switch *v.Format {
		case "int":
			p := PropertyIntegerParam[int]{
				Format: rusty.Some("int"),
				Type:   v.Type,
			}
			p.Default = coerceInt(v.Default)
			return NewPropertyInteger(p)
		case "int8":
			p := PropertyIntegerParam[int8]{
				Format: rusty.Some("int8"),
				Type:   v.Type,
			}
			p.Default = coerceInt8(v.Default)
			return NewPropertyInteger(p)
		case "int16":
			p := PropertyIntegerParam[int16]{
				Format: rusty.Some("int16"),
				Type:   v.Type,
			}
			p.Default = coerceInt16(v.Default)
			return NewPropertyInteger(p)
		case "int32":
			p := PropertyIntegerParam[int32]{
				Format: rusty.Some("int32"),
				Type:   v.Type,
			}
			p.Default = coerceInt32(v.Default)
			return NewPropertyInteger(p)
		case "int64":
			p := PropertyIntegerParam[int64]{
				Format: rusty.Some("int64"),
				Type:   v.Type,
			}
			p.Default = coerceInt64(v.Default)
			return NewPropertyInteger(p)
		default:
			panic("unknown format")
		}
	case "boolean":
		p := PropertyBooleanParam{
			Type: v.Type,
		}
		p.Default = coerceBool(v.Default)
		return NewPropertyBoolean(p)
	case "object":
		return NewPropertyObject(PropertyObjectParam{
			FileName:    v.FileName,
			Id:          v.Id,
			Title:       v.Title,
			Schema:      v.Schema,
			Description: rusty.OptionalFromPtr(v.Description),
			Properties:  NewPropertiesBuilder(loader).FromJson(v.Properties, v.Required).Build(),
			Required:    v.Required,
			Ref:         rusty.OptionalFromPtr(v.Ref),
		})
	default:
		panic("unknown type")
	}
}

func (b *SchemaBuilder) JSON2PropertyObject(js JSONSchema) *SchemaBuilder {
	x := b.Id(js.Id).
		FileName(js.FileName).
		Type(js.Type).
		Title(js.Title).
		Required(js.Required)
	if js.Description != nil {
		x.Description(*js.Description)
	}
	pb := NewPropertiesBuilder(b.__loader)
	for k, vin := range js.Properties {
		v, err := x.ResolveRef(&vin)
		if err != nil {
			panic(err)
		}
		pb.Add(NewPropertyItem(k, addJSONProperty(v, b.__loader)))
	}
	x.Properties(pb.Build())
	return x
}

func (b *SchemaBuilder) Id(id string) *SchemaBuilder {
	b.id = id
	return b
}

func (b *SchemaBuilder) FileName(fname string) *SchemaBuilder {
	b.fileName = fname
	return b
}

func (b *SchemaBuilder) Type(_type Type) *SchemaBuilder {
	b.typ = _type
	return b
}

func (b *SchemaBuilder) Schema(_schema string) *SchemaBuilder {
	b.schema = _schema
	return b
}

func (b *SchemaBuilder) Title(title string) *SchemaBuilder {
	b.title = title
	return b
}

func (b *SchemaBuilder) Description(description string) *SchemaBuilder {
	b.description = &description
	return b
}

func (b *SchemaBuilder) Properties(properties PropertiesObject) *SchemaBuilder {
	b.properties = properties
	return b
}

func (b *SchemaBuilder) Required(required []string) *SchemaBuilder {
	b.required = required
	return b
}

type propertyObject struct {
	id          string
	_type       Type
	fileName    string
	schema      string
	optional    bool
	title       string
	format      rusty.Optional[string]
	description rusty.Optional[string]
	ref         rusty.Optional[string]
	properties  PropertiesObject
	required    []string
	// deref       map[string]PropertyItem
}

// FileName implements PropertyObject.
func (p *propertyObject) FileName() string {
	return p.fileName
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
func (s *propertyObject) Properties() PropertiesObject {
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

func (s *propertyObject) Optional() bool {
	return s.optional
}

func (s *propertyObject) SetOptional() {
	s.optional = true
}

func (b *SchemaBuilder) Build() Property {
	requiredMap := make(map[string]bool)
	for _, v := range b.required {
		requiredMap[v] = true
	}
	for _, v := range b.properties.Items() {
		_, found := requiredMap[v.Name()]
		if !found {
			v.Property().SetOptional()
		}
	}
	// desc := rusty.None[string]()
	// if b.description != nil {
	// 	desc = rusty.Some[string](*b.description)
	// }
	ret := &propertyObject{
		id:          b.id,
		fileName:    b.fileName,
		format:      rusty.OptionalFromPtr(b.format),
		_type:       b.typ,
		title:       b.title,
		description: rusty.OptionalFromPtr(b.description),
		properties:  b.properties,
		required:    b.required,
		ref:         rusty.OptionalFromPtr(b.ref),
		// deref:       b.deref,
	}
	b.__loader.SchemaRegistry().AddSchema(ret)
	return ret
}

type PropertyObjectParam struct {
	FileName    string
	Id          string
	Title       string
	Schema      string
	Description rusty.Optional[string]
	Properties  PropertiesObject
	Required    []string
	Ref         rusty.Optional[string]
}

func NewPropertyObject(p PropertyObjectParam) PropertyObject {
	return &propertyObject{
		fileName:    p.FileName,
		id:          p.Id,
		_type:       OBJECT,
		title:       p.Title,
		schema:      p.Schema,
		description: p.Description,
		properties:  p.Properties,
		required:    p.Required,
		ref:         p.Ref,
		// deref:       map[string]PropertyItem{},
	}
}
