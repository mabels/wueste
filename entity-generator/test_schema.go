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
	registry *SchemaRegistry
}

func NewTestSchemaLoader() *TestSchemaLoader {
	return &TestSchemaLoader{
		registry: NewSchemaRegistry(),
	}
}

// AddSchema implements SchemaLoader.
func (t *TestSchemaLoader) SchemaRegistry() *SchemaRegistry {
	return t.registry
}

func (l *TestSchemaLoader) Abs(path string) (string, error) {
	if strings.HasPrefix(path, "/abs/") {
		return path, nil
	}
	return filepath.Join("/abs/", path), nil
}

// ReadFile implements SchemaLoader.
func (l *TestSchemaLoader) ReadFile(path string) ([]byte, error) {
	switch filepath.Base(path) {
	case "base.schema.json":
		return JSONBase(), nil
	case "sub.schema.json":
		return JSONSub(), nil
	case "sub2.schema.json":
		return JSONSub2(), nil
	default:
		return nil, fmt.Errorf("Test file not found: %s", path)
	}
}

// Unmarshal implements SchemaLoader.
func (l *TestSchemaLoader) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

var BaseSchema = JSONSchema{
	FileName:    "base.schema.json",
	Id:          "http://example.com/base.schema.json",
	Schema:      "http://json-schema.org/draft-07/schema#",
	Title:       "Base",
	Type:        "object",
	Description: toPtrString("Base description"),
	Properties: map[string]JSONProperty{
		"foo": {
			Type: "string",
		},
		"sub": {
			Ref: toPtrString("file://sub.schema.json"),
		},
	},
	Required: []string{"foo", "sub"},
}

func JSONBase() []byte {
	out, _ := json.MarshalIndent(BaseSchema, "", "  ")
	return out
}

var SubSchema = JSONSchema{
	FileName:    "sub.schema.json",
	Id:          "http://example.com/sub.schema.json",
	Schema:      "http://json-schema.org/draft-07/schema#",
	Title:       "Sub",
	Type:        "object",
	Description: toPtrString("Sub description"),
	Properties: map[string]JSONProperty{
		"sub": {
			Type: "string",
		},
		"sub-down": {
			Ref: toPtrString("file://sub2.schema.json"),
		},
	},
	Required: []string{"bar", "sub-down"},
}

func JSONSub() []byte {
	out, _ := json.MarshalIndent(SubSchema, "", "  ")
	return out
}

var Sub2Schema = JSONSchema{
	FileName:    "sub2.schema.json",
	Id:          "http://example.com/sub2.schema.json",
	Schema:      "http://json-schema.org/draft-07/schema#",
	Title:       "Sub2",
	Type:        "object",
	Description: toPtrString("Sub2 description"),
	Properties: map[string]JSONProperty{
		"bar": {
			Type: "string",
		},
	},
	Required: []string{"bar"},
}

func JSONSub2() []byte {
	out, _ := json.MarshalIndent(Sub2Schema, "", "  ")
	return out
}

func SchemaSchema(sl SchemaLoader) Property {
	return NewSchemaBuilder(sl).
		Id("JsonSchema").
		Type("object").
		Title("JsonSchema").
		Description("JSON Schema").
		Properties(NewPropertiesBuilder(sl).
			Add(NewPropertyItem("$id", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("$schema", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("title", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("type", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("properties", NewPropertyObject(PropertyObjectParam{}))).
			Add(NewPropertyItem("required", NewPropertyArray(PropertyArrayParam{})))).
		Required([]string{"$id", "$schema", "title", "type", "properties"}).
		Build()
}

func toPtrString(s string) *string {
	return &s
}

func TestJSONSubSchema() JSONProperty {
	return JSONProperty{
		JSONSchema: JSONSchema{
			Id:          "Sub",
			Title:       "Sub",
			Description: toPtrString("Description"),
			Properties: map[string]JSONProperty{
				"Test": {
					Type: "string",
				},
				"opt-Test": {
					Type: "string",
				},
			},
			Required: []string{"Test"},
		},
		Type: "object",
	}
}

func TestJsonFlatSchema() JSONSchema {
	json := JSONSchema{
		Id:          "SimpleType",
		Schema:      "",
		Title:       "SimpleType",
		Type:        "object",
		Description: toPtrString("Jojo SimpleType"),
		Properties: map[string]JSONProperty{
			"string": {
				Type: "string",
			},
			"default-string": {
				Type:    "string",
				Default: "hallo",
			},
			"optional-string": {
				Type: "string",
			},
			"optional-default-string": {
				Type:    "string",
				Default: "hallo",
			},
			"createdAt": {
				Type:   "string",
				Format: toPtrString("date-time"),
			},
			"default-createdAt": {
				Type:    "string",
				Default: "2023-12-31T23:59:59Z",
				Format:  toPtrString("date-time"),
			},
			"optional-createdAt": {
				Type:   "string",
				Format: toPtrString("date-time"),
			},
			"optional-default-createdAt": {
				Type:    "string",
				Format:  toPtrString("date-time"),
				Default: "2023-12-31T23:59:59Z",
			},
			"float64": {
				Type:   "number",
				Format: toPtrString("float64"),
			},
			"default-float64": {
				Type:    "number",
				Format:  toPtrString("float64"),
				Default: float64(4711.4),
			},
			"optional-float32": {
				Type:   "number",
				Format: toPtrString("float32"),
			},
			"optional-default-float32": {
				Type:    "number",
				Default: float32(49.2),
				Format:  toPtrString("float32"),
			},
			"int64": {
				Type:   "integer",
				Format: toPtrString("int64"),
			},
			"default-int64": {
				Type:    "integer",
				Default: int64(64),
				Format:  toPtrString("int64"),
			},
			"optional-int32": {
				Type:   "integer",
				Format: toPtrString("int32"),
			},
			"optional-default-int32": {
				Type:    "integer",
				Default: int32(32),
				Format:  toPtrString("int32"),
			},
			"bool": {
				Type: "boolean",
			},
			"default-bool": {
				Type:    "boolean",
				Default: true,
			},
			"optional-bool": {
				Type: "boolean",
			},
			"optional-default-bool": {
				Type:    "boolean",
				Default: true,
			},
			"sub":     TestJSONSubSchema(),
			"opt-sub": TestJSONSubSchema(),
		},
		Required: []string{
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
			"sub",
		},
	}
	return json
}

func TestSubSchema(sl SchemaLoader) Property {
	return NewSchemaBuilder(sl).
		Id("Sub").
		Title("Sub").
		Type("object").
		Description("Description").
		Properties(NewPropertiesBuilder(sl).
			Add(NewPropertyItem("Test", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("opt-Test", NewPropertyString(PropertyStringParam{}))).
			Build()).
		Required([]string{"Test"}).
		Build()
}

func TestFlatSchema(sl SchemaLoader) Property {
	return NewSchemaBuilder(sl).
		Id("SimpleType").
		Type("object").
		Title("SimpleType").
		Description("Jojo SimpleType").
		Properties(NewPropertiesBuilder(sl).
			Add(NewPropertyItem("string", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("default-string", NewPropertyString(PropertyStringParam{Default: rusty.Some("hallo")}))).
			Add(NewPropertyItem("optional-string", NewPropertyString(PropertyStringParam{}))).
			Add(NewPropertyItem("optional-default-string", NewPropertyString(PropertyStringParam{Default: rusty.Some("hallo")}))).
			Add(NewPropertyItem("createdAt", NewPropertyString(PropertyStringParam{
				Format: rusty.Some("date-time"),
			}))).
			Add(NewPropertyItem("default-createdAt", NewPropertyString(PropertyStringParam{
				Default: rusty.Some("2023-12-31T23:59:59Z"),
				Format:  rusty.Some("date-time"),
			}))).
			Add(NewPropertyItem("optional-createdAt", NewPropertyString(PropertyStringParam{
				Format: rusty.Some("date-time"),
			}))).
			Add(NewPropertyItem("optional-default-createdAt", NewPropertyString(PropertyStringParam{
				Default: rusty.Some("2023-12-31T23:59:59Z"),
				Format:  rusty.Some("date-time"),
			}))).
			Add(NewPropertyItem("float64", NewPropertyNumber(PropertyNumberParam[float64]{}))).
			Add(NewPropertyItem("default-float64", NewPropertyNumber(PropertyNumberParam[float64]{Default: rusty.Some(4711.4)}))).
			Add(NewPropertyItem("optional-float32", NewPropertyNumber(PropertyNumberParam[float32]{}))).
			Add(NewPropertyItem("optional-default-float32", NewPropertyNumber(PropertyNumberParam[float32]{Default: rusty.Some(float32(49.2))}))).
			Add(NewPropertyItem("int64", NewPropertyInteger(PropertyIntegerParam[int64]{}))).
			Add(NewPropertyItem("default-int64", NewPropertyInteger(PropertyIntegerParam[int64]{Default: rusty.Some(int64(64))}))).
			Add(NewPropertyItem("optional-int32", NewPropertyInteger(PropertyIntegerParam[int32]{}))).
			Add(NewPropertyItem("optional-default-int32", NewPropertyInteger(PropertyIntegerParam[int32]{Default: rusty.Some(int32(32))}))).
			Add(NewPropertyItem("bool", NewPropertyBoolean(PropertyBooleanParam{}))).
			Add(NewPropertyItem("default-bool", NewPropertyBoolean(PropertyBooleanParam{Default: rusty.Some(true)}))).
			Add(NewPropertyItem("optional-bool", NewPropertyBoolean(PropertyBooleanParam{}))).
			Add(NewPropertyItem("optional-default-bool", NewPropertyBoolean(PropertyBooleanParam{Default: rusty.Some(true)}))).
			Add(NewPropertyItem("sub", TestSubSchema(sl))).
			Add(NewPropertyItem("opt-sub", TestSubSchema(sl))).
			Build()).
		Required([]string{
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
}

func TestSchema(sl SchemaLoader) Property {
	ps := NewPropertiesBuilder(sl)
	fls := TestFlatSchema(sl).(PropertyObject)
	for _, item := range fls.Properties().Items() {
		ps.Add(item)
	}
	return NewSchemaBuilder(sl).
		Id("NestedType").
		Type("object").
		Title("NestedType").
		Description("Jojo NestedType").
		Properties(ps.
			Add(NewPropertyItem("arrayarrayBool", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: NewPropertyBoolean(PropertyBooleanParam{})}),
					})})}))).
			Add(NewPropertyItem("opt-arrayarrayBool", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: NewPropertyBoolean(PropertyBooleanParam{})}),
					})})}))).
			Add(NewPropertyItem("arrayString", NewPropertyArray(PropertyArrayParam{Items: NewPropertyString(PropertyStringParam{})}))).
			Add(NewPropertyItem("opt-arrayString", NewPropertyArray(PropertyArrayParam{Items: NewPropertyString(PropertyStringParam{})}))).
			Add(NewPropertyItem("arrayNumber", NewPropertyArray(PropertyArrayParam{Items: NewPropertyNumber(PropertyNumberParam[float64]{})}))).
			Add(NewPropertyItem("opt-arrayNumber", NewPropertyArray(PropertyArrayParam{Items: NewPropertyNumber(PropertyNumberParam[float64]{})}))).
			Add(NewPropertyItem("arrayInteger", NewPropertyArray(PropertyArrayParam{Items: NewPropertyInteger(PropertyIntegerParam[int]{})}))).
			Add(NewPropertyItem("opt-arrayInteger", NewPropertyArray(PropertyArrayParam{Items: NewPropertyInteger(PropertyIntegerParam[int]{})}))).
			Add(NewPropertyItem("arrayBool", NewPropertyArray(PropertyArrayParam{Items: NewPropertyBoolean(PropertyBooleanParam{})}))).
			Add(NewPropertyItem("opt-arrayBool", NewPropertyArray(PropertyArrayParam{Items: NewPropertyBoolean(PropertyBooleanParam{})}))).
			Add(NewPropertyItem("arrayarrayFlatSchema", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: TestSubSchema(sl)}),
					})})}))).
			Add(NewPropertyItem("opt-arrayarrayFlatSchema", NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: NewPropertyArray(PropertyArrayParam{
							Items: TestSubSchema(sl)}),
					})})}))).
			Add(NewPropertyItem("sub-flat", TestSubSchema(sl))).
			Add(NewPropertyItem("opt-sub-flat", TestSubSchema(sl))).
			Add(NewPropertyItem("arraySubType", NewPropertyArray(PropertyArrayParam{Items: TestSubSchema(sl)}))).
			Add(NewPropertyItem("opt-arraySubType", NewPropertyArray(PropertyArrayParam{Items: TestSubSchema(sl)}))).
			Build()).
		Required(append([]string{
			"arrayarrayBool",
			"sub",
			"arrayString",
			"arrayNumber",
			"arrayInteger",
			"arrayBool",
			"arraySubType",
			"arrayarrayFlatSchema",
		}, fls.Required()...)).
		Build()
}

func WriteTestSchema(cfg *GeneratorConfig) string {
	bytes, _ := json.MarshalIndent(TestJsonFlatSchema(), "", "  ")
	schemaFile := path.Join(cfg.OutputDir, "simple_type.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)
	return schemaFile
}
