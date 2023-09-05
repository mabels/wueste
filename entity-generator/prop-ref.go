package entity_generator

import (
	"fmt"
	"path"
	"strings"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyRef interface {
	Property
	Ref() rusty.Optional[string]
}

func (b *SchemaBuilder) ResolveRef(v JSONProperty) (*JSONProperty, error) {
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
