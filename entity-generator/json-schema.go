package entity_generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type JSONProperty struct {
	JSONSchema
	Type    string      `json:"type,omitempty"`
	Format  *string     `json:"format,omitempty"`
	Default interface{} `json:"default,omitempty"`
	Minimum interface{} `json:"minimum,omitempty"`
	Maximum interface{} `json:"maximum,omitempty"`
	Ref     *string     `json:"$ref,omitempty"`
}

type JSONProperties map[string]JSONProperty

type JSONSchema struct {
	FileName    string         `json:"$fileName,omitempty"`
	Id          string         `json:"$id,omitempty"`
	Schema      string         `json:"$schema,omitempty"`
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description *string        `json:"description,omitempty"`
	Properties  JSONProperties `json:"properties,omitempty"`
	Required    []string       `json:"required,omitempty"`
	// Deref       map[string]JSONSchema `json:"deref"`
}

func PropertiesToJson(props PropertiesObject) JSONProperties {
	ret := JSONProperties{}
	for _, p := range props.Items() {

		jsp := JSONProperty{
			Type: p.Property().Type(),
		}
		switch p.Property().Type() {
		case "string":
			p := p.Property().(PropertyString)
			if p.Format().IsSome() {
				jsp.Format = p.Format().Value()
			}
			if p.Default().IsSome() {
				jsp.Default = p.Default().Value()
			}
		case "boolean":
			p := p.Property().(PropertyBoolean)
			if p.Default().IsSome() {
				jsp.Default = p.Default().Value()
			}
		case "integer":
			switch p.Property().(type) {
			case PropertyInteger[int]:
				p := p.Property().(PropertyInteger[int])
				jsp.Format = toPtrString("int")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			case PropertyInteger[int8]:
				p := p.Property().(PropertyInteger[int8])
				jsp.Format = toPtrString("int8")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			case PropertyInteger[int16]:
				p := p.Property().(PropertyInteger[int16])
				jsp.Format = toPtrString("int16")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			case PropertyInteger[int32]:
				p := p.Property().(PropertyInteger[int32])
				jsp.Format = toPtrString("int32")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			case PropertyInteger[int64]:
				p := p.Property().(PropertyInteger[int64])
				jsp.Format = toPtrString("int64")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			default:
				panic("integer unknown type: " + p.Property().Type())
			}
		case "number":
			switch p.Property().(type) {
			case PropertyNumber[float64]:
				p := p.Property().(*propertyNumber[float64])
				jsp.Format = toPtrString("float64")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			case PropertyNumber[float32]:
				p := p.Property().(*propertyNumber[float32])
				jsp.Format = toPtrString("float32")
				if p.Default().IsSome() {
					jsp.Default = p.Default().Value()
				}
			default:
				panic("number unknown type: " + p.Property().Type())
			}
		case "object":
			p := p.Property().(PropertyObject)
			jsp.Id = p.Id()
			jsp.Title = p.Title()
			jsp.Schema = p.Schema()
			jsp.FileName = p.FileName()
			jsp.Description = rusty.OptionalToPtr(p.Description())
			jsp.Properties = PropertiesToJson(p.Properties())
			jsp.Required = p.Required()
			jsp.Ref = rusty.OptionalToPtr(p.Ref())
		default:
			panic("unknown type: " + p.Property().Type())
		}
		ret[p.Name()] = jsp
	}
	return ret
}

func JSONFromProperty(prop PropertyObject) JSONSchema {
	ret := JSONSchema{
		Id:         prop.Id(),
		Type:       prop.Type(),
		Title:      prop.Title(),
		Properties: JSONProperties{},
	}
	if prop.Description().IsSome() {
		ret.Description = prop.Description().Value()
	}

	ret.Properties = PropertiesToJson(prop.Properties())

	ret.Required = prop.Required()

	return ret
}

type SchemaLoader interface {
	ReadFile(path string) ([]byte, error)
	Abs(path string) (string, error)
	Unmarshal(bytes []byte, v interface{}) error
	SchemaRegistry() *SchemaRegistry
}

func LoadSchemaFromBytes(file string, bytes []byte, loader SchemaLoader) (Property, error) {
	absFile, err := loader.Abs(file)
	if err != nil {
		return nil, err
	}
	jsonSchema := JSONSchema{}
	err = loader.Unmarshal(bytes, &jsonSchema)
	if err != nil {
		return nil, fmt.Errorf("error parsing json schema: %s:%v:%w", absFile, string(bytes), err)
	}
	jsonSchema.FileName = absFile
	prop := NewSchemaBuilder(loader).JSON2PropertyObject(jsonSchema).Build()
	return prop, nil
}

type SchemaRegistryItem interface {
	Written() bool
	PropertItem() PropertyItem
}

type schemaRegistryItem struct {
	written bool
	prop    PropertyItem
}

func (sri *schemaRegistryItem) Written() bool {
	return sri.written
}

func (sri *schemaRegistryItem) PropertItem() PropertyItem {
	return sri.prop
}

type SchemaRegistry struct {
	registry map[string]*schemaRegistryItem
}

func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		registry: map[string]*schemaRegistryItem{},
	}
}

func (sr *SchemaRegistry) AddSchema(prop PropertyObject) PropertyItem {
	sri, found := sr.registry[prop.Id()]
	if found {
		return sri.prop
	}
	pi := NewPropertyItem(prop.Id(), prop)
	sr.registry[prop.Id()] = &schemaRegistryItem{
		written: false,
		prop:    pi,
	}
	return pi
}

func (sr *SchemaRegistry) SetWritten(prop PropertyObject) bool {
	sri, found := sr.registry[prop.Id()]
	if !found {
		panic("schema not found in registry: " + prop.Id())
	}
	sri.written = true
	return sri.written
}

func (sr *SchemaRegistry) IsWritten(prop PropertyObject) bool {
	sri, found := sr.registry[prop.Id()]
	if !found {
		return false
	}
	return sri.written
}

func (sr *SchemaRegistry) Items() []SchemaRegistryItem {
	ret := []SchemaRegistryItem{}
	for _, v := range sr.registry {
		ret = append(ret, v)
	}
	return ret
}

type schemaLoader struct {
	registry *SchemaRegistry
}

// var DefaultSchemaLoader = &schemaLoader{}
func NewSchemaLoader() SchemaLoader {
	return &schemaLoader{
		registry: NewSchemaRegistry(),
	}
}

func (s schemaLoader) SchemaRegistry() *SchemaRegistry {
	return s.registry
}

// Abs implements SchemaLoader.
func (schemaLoader) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

// ReadFile implements SchemaLoader.
func (schemaLoader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Unmarshal implements SchemaLoader.
func (schemaLoader) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

func LoadSchema(path string, loader SchemaLoader) (Property, error) {
	bytes, err := loader.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadSchemaFromBytes(path, bytes, loader)
}
