package entity_generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mabels/wueste/entity-generator/rusty"
)

// type JSONProperty struct {
// 	JSONSchema
// 	Type    string      `json:"type,omitempty"`
// 	Format  *string     `json:"format,omitempty"`
// 	Default interface{} `json:"default,omitempty"`
// 	Minimum interface{} `json:"minimum,omitempty"`
// 	Maximum interface{} `json:"maximum,omitempty"`
// 	Ref     *string     `json:"$ref,omitempty"`
// }

// type JSONPropertyString struct {
// 	Type    string      `json:"type,omitempty"`
// 	Format  *string     `json:"format,omitempty"`
// 	Default interface{} `json:"default,omitempty"`
// 	Ref     *string     `json:"$ref,omitempty"`
// }

// type JSONProperties map[string]JSONProperty

// type JSONSchema struct {
// 	FileName    string         `json:"$fileName,omitempty"`
// 	Id          string         `json:"$id,omitempty"`
// 	Schema      string         `json:"$schema,omitempty"`
// 	Title       string         `json:"title"`
// 	Type        string         `json:"type"`
// 	Description *string        `json:"description,omitempty"`
// 	Properties  JSONProperties `json:"properties,omitempty"`
// 	Required    []string       `json:"required,omitempty"`
// 	Ref         *string        `json:"$ref,omitempty"`
// 	// Deref       map[string]JSONSchema `json:"deref"`
// }

// type JSONPropertyObject struct {
// 	Type        string  `json:"type"`
// 	Description *string `json:"description,omitempty"`

// 	FileName   string         `json:"$fileName,omitempty"`
// 	Id         string         `json:"$id,omitempty"`
// 	Schema     string         `json:"$schema,omitempty"`
// 	Title      string         `json:"title"`
// 	Properties JSONProperties `json:"properties,omitempty"`
// 	Required   []string       `json:"required,omitempty"`
// 	// Deref       map[string]JSONSchema `json:"deref"`
// }

// func PropertiesToJson(props PropertiesObject) JSONProperties {
// 	ret := JSONProperties{}
// 	for _, p := range props.Items() {

// 		jsp := JSONProperty{
// 			Type: p.Property().Type(),
// 		}
// 		switch p.Property().Type() {
// 		case "string":
// 			p := p.Property().(PropertyString)
// 			if p.Format().IsSome() {
// 				jsp.Format = p.Format().Value()
// 			}
// 			if p.Default().IsSome() {
// 				jsp.Default = p.Default().Value()
// 			}
// 		case "boolean":
// 			p := p.Property().(PropertyBoolean)
// 			if p.Default().IsSome() {
// 				jsp.Default = p.Default().Value()
// 			}
// 		case "integer":
// 			switch p.Property().(type) {
// 			case PropertyInteger[int]:
// 				p := p.Property().(PropertyInteger[int])
// 				jsp.Format = toPtrString("int")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int8]:
// 				p := p.Property().(PropertyInteger[int8])
// 				jsp.Format = toPtrString("int8")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int16]:
// 				p := p.Property().(PropertyInteger[int16])
// 				jsp.Format = toPtrString("int16")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int32]:
// 				p := p.Property().(PropertyInteger[int32])
// 				jsp.Format = toPtrString("int32")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int64]:
// 				p := p.Property().(PropertyInteger[int64])
// 				jsp.Format = toPtrString("int64")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			default:
// 				panic("integer unknown type: " + p.Property().Type())
// 			}
// 		case "number":
// 			switch p.Property().(type) {
// 			case PropertyNumber[float64]:
// 				p := p.Property().(*propertyNumber[float64])
// 				jsp.Format = toPtrString("float64")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyNumber[float32]:
// 				p := p.Property().(*propertyNumber[float32])
// 				jsp.Format = toPtrString("float32")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			default:
// 				panic("number unknown type: " + p.Property().Type())
// 			}
// 		case "object":
// 			p := p.Property().(PropertyObject)
// 			jsp.Id = p.Id()
// 			jsp.Title = p.Title()
// 			jsp.Schema = p.Schema()
// 			jsp.FileName = p.FileName()
// 			jsp.Description = rusty.OptionalToPtr(p.Description())
// 			jsp.Properties = PropertiesToJson(p.Properties())
// 			jsp.Required = p.Required()
// 			jsp.Ref = rusty.OptionalToPtr(p.Ref())
// 		default:
// 			panic("unknown type: " + p.Property().Type())
// 		}
// 		ret[p.Name()] = jsp
// 	}
// 	return ret
// }

// func JSONFromProperty(iprop any) JSONProperty {
// 	ret := JSONProperty{}
// 	switch prop := iprop.(type) {
// 	case PropertyString:
// 	case PropertyArray:
// 	case PropertyBoolean:
// 	case PropertyInteger:
// 	case PropertyNumber:
// 	case PropertyObject:
// 	default:
// 		panic("unknown type: " + prop.(Property).Type())
// 	}

// / 		case "string":
// 			p := p.Property().(PropertyString)
// 			if p.Format().IsSome() {
// 				jsp.Format = p.Format().Value()
// 			}
// 			if p.Default().IsSome() {
// 				jsp.Default = p.Default().Value()
// 			}
// 		case "boolean":
// 			p := p.Property().(PropertyBoolean)
// 			if p.Default().IsSome() {
// 				jsp.Default = p.Default().Value()
// 			}
// 		case "integer":
// 			switch p.Property().(type) {
// 			case PropertyInteger[int]:
// 				p := p.Property().(PropertyInteger[int])
// 				jsp.Format = toPtrString("int")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int8]:
// 				p := p.Property().(PropertyInteger[int8])
// 				jsp.Format = toPtrString("int8")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int16]:
// 				p := p.Property().(PropertyInteger[int16])
// 				jsp.Format = toPtrString("int16")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int32]:
// 				p := p.Property().(PropertyInteger[int32])
// 				jsp.Format = toPtrString("int32")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyInteger[int64]:
// 				p := p.Property().(PropertyInteger[int64])
// 				jsp.Format = toPtrString("int64")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			default:
// 				panic("integer unknown type: " + p.Property().Type())
// 			}
// 		case "number":
// 			switch p.Property().(type) {
// 			case PropertyNumber[float64]:
// 				p := p.Property().(*propertyNumber[float64])
// 				jsp.Format = toPtrString("float64")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			case PropertyNumber[float32]:
// 				p := p.Property().(*propertyNumber[float32])
// 				jsp.Format = toPtrString("float32")
// 				if p.Default().IsSome() {
// 					jsp.Default = p.Default().Value()
// 				}
// 			default:
// 				panic("number unknown type: " + p.Property().Type())
// 			}
// 		case "object":
// 			p := p.Property().(PropertyObject)
// 			jsp.Id = p.Id()
// 			jsp.Title = p.Title()
// 			jsp.Schema = p.Schema()
// 			jsp.FileName = p.FileName()
// 			jsp.Description = rusty.OptionalToPtr(p.Description())
// 			jsp.Properties = PropertiesToJson(p.Properties())
// 			jsp.Required = p.Required()
// 			jsp.Ref = rusty.OptionalToPtr(p.Ref())

// return ret
// }

type SchemaLoader interface {
	ReadFile(path string) ([]byte, error)
	Abs(path string) (string, error)
	Unmarshal(bytes []byte, v interface{}) error
	// SchemaRegistry() *SchemaRegistry
	// LoadRef(refVal string) (Property, error)
}

func loadSchemaFromBytes(file string, bytes []byte, loader PropertyCtx, fn func(abs string, prop JSONProperty) rusty.Result[Property]) rusty.Result[Property] {
	absFile, err := loader.Registry.loader.Abs(file)
	if err != nil {
		return rusty.Err[Property](err)
	}
	jsonSchema := NewJSONProperty()
	err = loader.Registry.loader.Unmarshal(bytes, jsonSchema)
	if err != nil {
		return rusty.Err[Property](fmt.Errorf("error parsing json schema: %s:%v:%w", absFile, string(bytes), err))
	}
	return fn(absFile, jsonSchema)

}

type SchemaRegistryItem interface {
	Written() bool
	Property() Property
}

type schemaRegistryItem struct {
	written bool
	prop    Property
}

func (sri *schemaRegistryItem) Written() bool {
	return sri.written
}

func (sri *schemaRegistryItem) Property() Property {
	return sri.prop
}

type SchemaRegistry struct {
	registry map[string]*schemaRegistryItem
	loader   SchemaLoader
}

func NewSchemaRegistry(loaders ...SchemaLoader) *SchemaRegistry {
	var loader SchemaLoader
	if len(loaders) == 0 {
		loader = &schemaLoader{}
	} else {
		loader = loaders[0]
	}
	return &SchemaRegistry{
		loader:   loader,
		registry: map[string]*schemaRegistryItem{},
	}
}

func (sr *SchemaRegistry) EnsureSchema(key string, ort PropertyRuntime, fn func(fname string, rt PropertyRuntime) rusty.Result[Property]) rusty.Result[Property] {

	ref := strings.TrimSpace(key)
	if ref[0] == '#' {
		return rusty.Err[Property](fmt.Errorf("local ref not supported"))
	}
	if !strings.HasPrefix(ref, "file://") {
		return rusty.Err[Property](fmt.Errorf("only file:// ref supported"))
	}
	fname := ref[len("file://"):]
	if !strings.HasSuffix(fname, "/") {
		dir := "./"
		if ort.FileName.IsSome() {
			dir = path.Dir(ort.FileName.Value())
		}
		fname = path.Join(dir, fname)
		// if prop.Runtime().FileName.IsSome() {
	}
	fname, err := sr.loader.Abs(fname)
	if err != nil {
		var err error = fmt.Errorf("only directory ref supported")
		return rusty.Err[Property](err)
	}

	sri, found := sr.registry[fname]
	if found {
		return rusty.Ok[Property](sri.prop)
	}
	rt := ort.Clone()
	rt.SetRef(key)
	rt.SetFileName(fname)
	pi := fn(fname, *rt)
	if pi.IsErr() {
		return pi
	}
	pip1 := pi.Ok().Runtime()
	pip2 := pi.Ok().Runtime()
	if pip1 != pip2 {
		panic("runtime not equal")
	}
	pi.Ok().Runtime().Assign(*rt)
	item := &schemaRegistryItem{
		written: false,
		prop:    pi.Ok(),
	}
	sr.registry[key] = item
	sr.registry[fname] = item
	if pi.Ok().Id() != "" {
		sr.registry[pi.Ok().Id()] = item
	}
	return pi
}

func (sr *SchemaRegistry) SetWritten(prop PropertyObject) bool {
	sri, found := sr.registry[prop.Runtime().FileName.Value()]
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

// // var DefaultSchemaLoader = &schemaLoader{}
// func NewSchemaLoader() SchemaLoader {
// 	return &schemaLoader{
// 		registry: NewSchemaRegistry(),
// 	}
// }

// func (b schemaLoader) ResolveRef(ref rusty.Optional[string]) (any, error) {
// 	if !ref.IsNone() {
// 		ref := strings.TrimSpace(*ref.Value())
// 		if ref[0] == '#' {
// 			return nil, fmt.Errorf("local ref not supported")
// 		}
// 		if !strings.HasPrefix(ref, "file://") {
// 			return nil, fmt.Errorf("only file:// ref supported")
// 		}
// 		fname := ref[len("file://"):]
// 		if !strings.HasSuffix(fname, "/") {
// 			var err error
// 			fname, err = b.Abs(path.Join(path.Dir(b.fileName), fname))
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 		pl, err := LoadSchema(fname, b)
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
// }

// func (s schemaLoader) SchemaRegistry() *SchemaRegistry {
// 	return s.registry
// }

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

func loadSchema(path string, loader PropertyCtx, fn func(abs string, prop JSONProperty) rusty.Result[Property]) rusty.Result[Property] {
	bytes, err := loader.Registry.loader.ReadFile(path)
	if err != nil {
		return rusty.Err[Property](err)
	}
	return loadSchemaFromBytes(path, bytes, loader, fn)
}
