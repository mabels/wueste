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

type SchemaLoader interface {
	ReadFile(path string) ([]byte, error)
	Abs(path string) (string, error)
	Unmarshal(bytes []byte, v interface{}) error
	// SchemaRegistry() *SchemaRegistry
	// LoadRef(refVal string) (Property, error)
}

type SchemaRegistryItem interface {
	Written() bool
	JSonFile() JSonFile
}

type JSonFile struct {
	FileName     string   `json:"filename"`
	JSONProperty JSONDict `json:"jsonProperty"`
}

type schemaRegistryItem struct {
	written bool
	// prop    Property
	jsonFile JSonFile
}

func (sri *schemaRegistryItem) Written() bool {
	return sri.written
}

func (sri *schemaRegistryItem) JSonFile() JSonFile {
	return sri.jsonFile
}

type SchemaRegistry struct {
	registry map[string]*schemaRegistryItem
	BaseDir  rusty.Optional[string]
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

func (sr *SchemaRegistry) EnsureJSONProperty(parentFname rusty.Optional[string], inRef string) rusty.Result[JSonFile] {
	ref := strings.TrimSpace(inRef)
	if ref[0] == '#' {
		return rusty.Err[JSonFile](fmt.Errorf("local ref not supported"))
	}
	if !strings.HasPrefix(ref, "file://") {
		return rusty.Err[JSonFile](fmt.Errorf("only file:// ref supported"))
	}
	fname := ref[len("file://"):]
	if !strings.HasSuffix(fname, "/") {
		dir := "./"
		if sr.BaseDir.IsSome() {
			dir = sr.BaseDir.Value()
		}
		if parentFname.IsSome() {
			dir = path.Dir(parentFname.Value())
		}
		fname = path.Join(dir, fname)
	}
	fname, err := sr.loader.Abs(fname)
	if err != nil {
		var err error = fmt.Errorf("only directory ref supported")
		return rusty.Err[JSonFile](err)
	}
	sri, found := sr.registry[fname]
	if found {
		return rusty.Ok[JSonFile](sri.jsonFile)
	}
	rjsonFile := loadSchema(fname, sr.loader)
	if rjsonFile.IsErr() {
		return rjsonFile
	}
	sr.registry[rjsonFile.Ok().FileName] = &schemaRegistryItem{
		jsonFile: rjsonFile.Ok(),
	}
	return rjsonFile
}

// func (sr *SchemaRegistry) EnsureSchema(key string, parentFname rusty.Optional[string], fn func(fname string) rusty.Result[Property]) rusty.Result[Property] {
// 	ref := strings.TrimSpace(key)
// 	if ref[0] == '#' {
// 		return rusty.Err[Property](fmt.Errorf("local ref not supported"))
// 	}
// 	if !strings.HasPrefix(ref, "file://") {
// 		return rusty.Err[Property](fmt.Errorf("only file:// ref supported"))
// 	}
// 	fname := ref[len("file://"):]
// 	if !strings.HasSuffix(fname, "/") {
// 		dir := "./"
// 		if sr.BaseDir.IsSome() {
// 			dir = sr.BaseDir.Value()
// 		}
// 		if parentFname.IsSome() {
// 			dir = path.Dir(parentFname.Value())
// 		}
// 		fname = path.Join(dir, fname)
// 	}
// 	fname, err := sr.loader.Abs(fname)
// 	if err != nil {
// 		var err error = fmt.Errorf("only directory ref supported")
// 		return rusty.Err[Property](err)
// 	}

// 	_ /*sri */, found := sr.registry[fname]
// 	if found {
// 		panic("schema already loaded: " + fname)
// 		// return rusty.Ok[Property](sri.prop.Clone())
// 	}
// 	// rt := ort.Clone()
// 	// rt.SetRef(key)
// 	// rt.SetFileName(fname)
// 	pi := fn(fname) // , *rt)
// 	if pi.IsErr() {
// 		return pi
// 	}
// 	// pip1 := pi.Ok().Runtime()
// 	// pip2 := pi.Ok().Runtime()
// 	// if pip1 != pip2 {
// 	// 	panic("runtime not equal")
// 	// }
// 	// pi.Ok().Runtime().Assign(*rt)
// 	item := &schemaRegistryItem{
// 		written: false,
// 		prop:    pi.Ok(),
// 	}
// 	sr.registry[key] = item
// 	sr.registry[fname] = item
// 	if pi.Ok().Id() != "" {
// 		sr.registry[pi.Ok().Id()] = item
// 	}
// 	return pi
// }

// func (sr *SchemaRegistry) SetWritten(prop Property) bool {
// 	sri, found := sr.registry[prop.Meta().FileName().Value()]
// 	if !found {
// 		panic("schema not found in registry: " + prop.Id())
// 	}
// 	sri.written = true
// 	return sri.written
// }

// func (sr *SchemaRegistry) IsWritten(prop JSonFile) bool {
// 	sri, found := sr.registry[prop.FileName]
// 	if !found {
// 		return false
// 	}
// 	return sri.written
// }

func (sr *SchemaRegistry) Items() []SchemaRegistryItem {
	ret := []SchemaRegistryItem{}
	for _, v := range sr.registry {
		ret = append(ret, v)
	}
	return ret
}

type schemaLoader struct {
	// registry *SchemaRegistry
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

func loadSchemaFromBytes(bytes []byte, loader SchemaLoader) rusty.Result[JSONDict] {
	jsonSchema := NewJSONDict()
	err := loader.Unmarshal(bytes, jsonSchema)
	if err != nil {
		return rusty.Err[JSONDict](fmt.Errorf("error parsing json schema: %v:%w", string(bytes), err))
	}
	return rusty.Ok(jsonSchema)
}

func loadSchema(fname string, loader SchemaLoader) rusty.Result[JSonFile] {
	fname, err := loader.Abs(fname)
	if err != nil {
		return rusty.Err[JSonFile](err)
	}
	bytes, err := loader.ReadFile(fname)
	if err != nil {
		return rusty.Err[JSonFile](err)
	}
	rjs := loadSchemaFromBytes(bytes, loader)
	if rjs.IsErr() {
		return rusty.Err[JSonFile](rjs.Err())
	}
	return rusty.Ok[JSonFile](JSonFile{
		FileName:     fname,
		JSONProperty: rjs.Ok(),
	})

}
