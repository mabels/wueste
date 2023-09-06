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

/*
{
	"$id": "https://schema.6265746f6f.tech/nft_ask_data.schema.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "NftAskData",
	"type": "object",
	"properties": {
		"id": {
			"type": "string"
		},
	}
}
*/

type TestSchemaLoader struct {
	// registry *SchemaRegistry
}

func NewTestRuntime() PropertyRuntime {
	return PropertyRuntime{}
}

func NewTestContext() PropertyCtx {
	return PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
}

// AddSchema implements SchemaLoader.
// func (t *TestSchemaLoader) SchemaRegistry() *SchemaRegistry {
// 	return t.registry
// }

func (l *TestSchemaLoader) Abs(path string) (string, error) {
	if strings.HasPrefix(path, "/abs/") {
		return path, nil
	}
	return filepath.Join("/abs/", path), nil
}

// ReadFile implements SchemaLoader.
func (l *TestSchemaLoader) ReadFile(path string) ([]byte, error) {
	switch path {
	case "/abs/unnamed_nested_object.schema.json":
		return JSONUnnamedNestedObject(), nil
	case "/abs/base.schema.json":
		return JSONBase(), nil
	case "/abs/sub.schema.json":
		return JSONSub(), nil
	case "/abs/wurst/sub2.schema.json":
		return JSONSub2(), nil
	case "/abs/wurst/sub3.schema.json":
		return JSONSub3(), nil
	default:
		return nil, fmt.Errorf("Test file not found: %s", path)
	}
}

// Unmarshal implements SchemaLoader.
func (l *TestSchemaLoader) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

func json2JSONProperty(inp string) JSONProperty {
	jp := NewJSONProperty()
	err := json.Unmarshal([]byte(inp), &jp)
	if err != nil {
		panic(fmt.Errorf("json2JSONProperty: %w:%v", err, inp))
	}
	return jp
}

var UnnamedNestedObject = json2JSONProperty(`{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"$id": "https://github.com/betooinc/architectures/schema/opensea/event/item_listed.json",
	"title": "OpenSeaEventItemListed",
	"type": "object",
	"properties": {
		"collection": {
			"type": "object",
			"properties": {
			"slug": { "type": "string" }
			}
		}
	}
}`)

func JSONUnnamedNestedObject() []byte {
	out, _ := json.MarshalIndent(UnnamedNestedObject, "", "  ")
	return out
}

var BaseSchema = json2JSONProperty(`{
	"fileName":    "base.schema.json",
	"$id":         "http://example.com/base.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Base",
	"type":        "object",
	"description": "Base description",
	"properties": {
		"foo": {
			"type": "string"
		},
		"sub": {
			"$ref": "file://sub.schema.json"
		}
	},
	"required": ["foo", "sub"]
}`)

func JSONBase() []byte {
	out, _ := json.MarshalIndent(BaseSchema, "", "  ")
	return out
}

func TestJsonSubSchema() JSONProperty {
	return json2JSONProperty(`{
	"fileName":    "sub.schema.json",
	"$id":         "http://example.com/sub.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Sub",
	"type":        "object",
	"description": "Sub description",
	"properties": {
		"sub": {
			"type": "string"
		},
		"sub-down": {
			"$ref": "file://wurst/sub2.schema.json"
		}
	},
	"required": ["bar", "sub-down"]
}`)

}

func JSONSub() []byte {
	out, _ := json.MarshalIndent(TestJsonSubSchema(), "", "  ")
	return out
}

var Sub2Schema = json2JSONProperty(`{
	"fileName":    "wurst/sub2.schema.json",
	"$id":         "http://example.com/sub2.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Result",
	"type":        "object",
	"description": "Sub2 description",
	"properties": {
		"bar": {
			"$ref": "file://./sub3.schema.json"
		}
	},
	"required": ["bar"]
}`)

func JSONSub2() []byte {
	out, _ := json.MarshalIndent(Sub2Schema, "", "  ")
	return out
}

var Sub3Schema = json2JSONProperty(`{
	"fileName":    "wurst/sub3.schema.json",
	"$id":         "http://example.com/sub3.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Sub3",
	"type":        "object",
	"description": "Sub3 description",
	"properties": {
		"bar": {
			"type": "string"
		}
	},
	"required": ["bar"]
}`)

func JSONSub3() []byte {
	out, _ := json.MarshalIndent(Sub3Schema, "", "  ")
	return out
}

func SchemaSchema(sl PropertyCtx) Property {
	return NewPropertiesBuilder(sl).BuildObject().
		id("JsonSchema").
		title("JsonSchema").
		description("JSON Schema").
		propertiesAdd(NewPropertyItem("$id", NewPropertyString(PropertyStringParam{}))).
		propertiesAdd(NewPropertyItem("$schema", NewPropertyString(PropertyStringParam{}))).
		propertiesAdd(NewPropertyItem("title", NewPropertyString(PropertyStringParam{}))).
		propertiesAdd(NewPropertyItem("type", NewPropertyString(PropertyStringParam{}))).
		propertiesAdd(NewPropertyItem("properties", NewPropertyObject(PropertyObjectParam{}))).
		propertiesAdd(NewPropertyItem("required", NewPropertyArray(PropertyArrayParam{}))).
		required([]string{"$id", "$schema", "title", "type", "properties"}).
		Build().Ok()
}

// func toPtrString(s string) *string {
// 	return &s
// }

func TestJSONSubSchema() JSONProperty {
	return json2JSONProperty(`{
		"$id":         "https://Sub",
		"title":       "Payload",
		"description": "Description",
		"properties": {
			"Test": {
				"type": "string"
			},
			"opt-Test": {
				"type": "string"
			}
		},
		"required": ["Test"],
		"type":     "object"
	}`)
}

func TestJsonFlatSchema() JSONProperty {
	json := json2JSONProperty(`{
		"type":        "object",
		"$id": "https://SimpleType",
		"title":       "SimpleType",
		"description": "Jojo SimpleType",
		"properties": {
			"string": {
				"type": "string"
			},
			"default-string": {
				"type":    "string",
				"default": "hallo"
			},
			"optional-string": {
				"type": "string"
			},
			"optional-default-string": {
				"type":    "string",
				"default": "hallo"
			},
			"createdAt": {
				"type":   "string",
				"format": "date-time"
			},
			"default-createdAt": {
				"type":    "string",
				"format":  "date-time",
				"default": "2023-12-31T23:59:59Z"
			},
			"optional-createdAt": {
				"type":   "string",
				"format": "date-time"
			},
			"optional-default-createdAt": {
				"type":    "string",
				"format":  "date-time",
				"default": "2023-12-31T23:59:59Z"
			},
			"float64": {
				"type": "number"
			},
			"default-float64": {
				"type":    "number",
				"format":  "float32",
				"default": 4711.4
			},
			"optional-float32": {
				"type":   "number",
				"format": "float32"
			},
			"optional-default-float32": {
				"type":    "number",
				"format":  "float32",
				"default": 49.2
			},
			"int64": {
				"type":   "integer",
				"format": "int64"
			},
			"default-int64": {
				"type":    "integer",
				"format":  "int64",
				"default": 64
			},
			"optional-int32": {
				"type":   "integer",
				"format": "int32"
			},
			"optional-default-int32": {
				"type":    "integer",
				"format":  "int32",
				"default": 32
			},
			"bool": {
				"type": "boolean"
			},
			"default-bool": {
				"type":    "boolean",
				"default": true
			},
			"optional-bool": {
				"type": "boolean"
			},
			"optional-default-bool": {
				"type":    "boolean",
				"default": true
			}
		},
		"required": [
			"string",
			"createdAt",
			"default-string",
			"default-createdAt",
			"float64",
			"default-float64",
			"int64",
			"default-int64",
			"uint64",
			"default-uint64",
			"default-bool",
			"bool",
			"sub"
		]
	}`)
	_prop := json.Get("properties")
	prop, ok := _prop.(JSONProperty)
	if !ok {
		panic("not ok")
	}
	prop.Set("sub", TestJSONSubSchema())
	prop.Set("opt-sub", TestJSONSubSchema())
	return json
}

func TestSubSchema(sl PropertyCtx, rt PropertyRuntime) rusty.Result[Property] {
	return sl.Registry.EnsureSchema("file://./sub.schema.json", rt, func(fname string, rt PropertyRuntime) rusty.Result[Property] {
		return NewPropertiesBuilder(sl).BuildObject().
			id("https://Sub").
			title("Payload").
			description("Description").
			propertiesAdd(NewPropertyItem("Test", NewPropertyString(PropertyStringParam{}))).
			propertiesAdd(NewPropertyItem("opt-Test", NewPropertyString(PropertyStringParam{}))).
			required([]string{"Test"}).
			Build()
	})
}

func TestFlatSchema(sl PropertyCtx, rt PropertyRuntime) rusty.Result[Property] {
	return sl.Registry.EnsureSchema("file://./simple_type.schema.json", rt, func(fname string, rt PropertyRuntime) rusty.Result[Property] {
		return NewPropertiesBuilder(sl).BuildObject().
			id("https://SimpleType").
			title("SimpleType").
			description("Jojo SimpleType").
			propertiesAdd(NewPropertyItem("string", NewPropertyString(PropertyStringParam{}))).
			propertiesAdd(NewPropertyItem("default-string", NewPropertyString(PropertyStringParam{Default: rusty.Some("hallo")}))).
			propertiesAdd(NewPropertyItem("optional-string", NewPropertyString(PropertyStringParam{}))).
			propertiesAdd(NewPropertyItem("optional-default-string", NewPropertyString(PropertyStringParam{Default: rusty.Some("hallo")}))).
			propertiesAdd(NewPropertyItem("createdAt", NewPropertyString(PropertyStringParam{
				Format: rusty.Some("date-time"),
			}))).
			propertiesAdd(NewPropertyItem("default-createdAt", NewPropertyString(PropertyStringParam{
				Default: rusty.Some("2023-12-31T23:59:59Z"),
				Format:  rusty.Some("date-time"),
			}))).
			propertiesAdd(NewPropertyItem("optional-createdAt", NewPropertyString(PropertyStringParam{
				Format: rusty.Some("date-time"),
			}))).
			propertiesAdd(NewPropertyItem("optional-default-createdAt", NewPropertyString(PropertyStringParam{
				Default: rusty.Some("2023-12-31T23:59:59Z"),
				Format:  rusty.Some("date-time"),
			}))).
			propertiesAdd(NewPropertyItem("float64", NewPropertyNumber(PropertyNumberParam{}))).
			propertiesAdd(NewPropertyItem("default-float64", NewPropertyNumber(PropertyNumberParam{Default: rusty.Some(4711.4), Format: rusty.Some("float32")}))).
			propertiesAdd(NewPropertyItem("optional-float32", NewPropertyNumber(PropertyNumberParam{Format: rusty.Some("float32")}))).
			propertiesAdd(NewPropertyItem("optional-default-float32", NewPropertyNumber(PropertyNumberParam{Default: rusty.Some(49.2)}))).
			propertiesAdd(NewPropertyItem("int64", NewPropertyInteger(PropertyIntegerParam{}))).
			propertiesAdd(NewPropertyItem("default-int64", NewPropertyInteger(PropertyIntegerParam{Default: rusty.Some(64)}))).
			propertiesAdd(NewPropertyItem("optional-int32", NewPropertyInteger(PropertyIntegerParam{Format: rusty.Some("int32")}))).
			propertiesAdd(NewPropertyItem("optional-default-int32", NewPropertyInteger(PropertyIntegerParam{Default: rusty.Some(32)}))).
			propertiesAdd(NewPropertyItem("bool", NewPropertyBoolean(PropertyBooleanParam{}))).
			propertiesAdd(NewPropertyItem("default-bool", NewPropertyBoolean(PropertyBooleanParam{Default: rusty.Some(true)}))).
			propertiesAdd(NewPropertyItem("optional-bool", NewPropertyBoolean(PropertyBooleanParam{}))).
			propertiesAdd(NewPropertyItem("optional-default-bool", NewPropertyBoolean(PropertyBooleanParam{Default: rusty.Some(true)}))).
			propertiesAdd(NewPropertyItem("sub", TestSubSchema(sl, rt))).
			propertiesAdd(NewPropertyItem("opt-sub", TestSubSchema(sl, rt))).
			required([]string{
				"string",
				"createdAt",
				"default-string",
				"default-createdAt",
				"float64",
				"default-float64",
				"int64",
				"default-int64",
				"uint64",
				"default-uint64",
				"default-bool",
				"bool",
				"sub"}).
			Build()
	})
}

func TestSchema(sl PropertyCtx, rts ...PropertyRuntime) Property {
	rt := PropertyRuntime{}
	if len(rts) > 0 {
		rt = rts[0]
	}
	return sl.Registry.EnsureSchema("file://./nested_type.schema.json", rt, func(fname string, rt PropertyRuntime) rusty.Result[Property] {
		ps := NewPropertiesBuilder(sl).BuildObject()
		fls := TestFlatSchema(sl, rt).Ok().Runtime().ToPropertyObject().Ok()
		for _, item := range fls.Items() {
			ps.propertiesAdd(rusty.Ok(item))
		}
		return ps.
			id("https://NestedType").
			title("NestedType").
			description("Jojo NestedType").
			propertiesAdd(NewPropertyItem("arrayarrayBool", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: NewPropertyBoolean(PropertyBooleanParam{})}),
					})})}))).
			propertiesAdd(NewPropertyItem("opt-arrayarrayBool", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: NewPropertyBoolean(PropertyBooleanParam{})}),
					})})}))).
			propertiesAdd(NewPropertyItem("arrayString", NewPropertyArray(PropertyArrayParam{Items: NewPropertyString(PropertyStringParam{})}))).
			propertiesAdd(NewPropertyItem("opt-arrayString", NewPropertyArray(PropertyArrayParam{Items: NewPropertyString(PropertyStringParam{})}))).
			propertiesAdd(NewPropertyItem("arrayNumber", NewPropertyArray(PropertyArrayParam{Items: NewPropertyNumber(PropertyNumberParam{})}))).
			propertiesAdd(NewPropertyItem("opt-arrayNumber", NewPropertyArray(PropertyArrayParam{Items: NewPropertyNumber(PropertyNumberParam{})}))).
			propertiesAdd(NewPropertyItem("arrayInteger", NewPropertyArray(PropertyArrayParam{Items: NewPropertyInteger(PropertyIntegerParam{})}))).
			propertiesAdd(NewPropertyItem("opt-arrayInteger", NewPropertyArray(PropertyArrayParam{Items: NewPropertyInteger(PropertyIntegerParam{})}))).
			propertiesAdd(NewPropertyItem("arrayBool", NewPropertyArray(PropertyArrayParam{Items: NewPropertyBoolean(PropertyBooleanParam{})}))).
			propertiesAdd(NewPropertyItem("opt-arrayBool", NewPropertyArray(PropertyArrayParam{Items: NewPropertyBoolean(PropertyBooleanParam{})}))).
			propertiesAdd(NewPropertyItem("arrayarrayFlatSchema", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: TestSubSchema(sl, rt)}),
					})})}))).
			propertiesAdd(NewPropertyItem("opt-arrayarrayFlatSchema", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: TestSubSchema(sl, rt)}),
					})})}))).
			propertiesAdd(NewPropertyItem("sub-flat", TestSubSchema(sl, rt))).
			propertiesAdd(NewPropertyItem("opt-sub-flat", TestSubSchema(sl, rt))).
			propertiesAdd(NewPropertyItem("arraySubType", NewPropertyArray(PropertyArrayParam{Items: TestSubSchema(sl, rt)}))).
			propertiesAdd(NewPropertyItem("opt-arraySubType", NewPropertyArray(PropertyArrayParam{Items: TestSubSchema(sl, rt)}))).
			required(append([]string{
				"arrayarrayBool",
				"sub-flat",
				"arrayString",
				"arrayNumber",
				"arrayInteger",
				"arrayBool",
				"arraySubType",
				"arrayarrayFlatSchema",
			}, fls.Required()...)).
			Build()
	}).Ok()
}

func WriteTestSchema(cfg *GeneratorConfig) string {
	jsonSchema := PropertyToJson(TestFlatSchema(NewTestContext(), PropertyRuntime{}).Ok())
	bytes, _ := json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile := path.Join(cfg.OutputDir, "simple_type.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	jsonSchema = PropertyToJson(TestSchema(NewTestContext(), PropertyRuntime{}))
	bytes, _ = json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile = path.Join(cfg.OutputDir, "nested_type.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)
	return schemaFile
}
