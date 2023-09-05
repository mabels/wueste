package entity_generator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

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

var BaseSchema = JSONProperty{
	"fileName":    "base.schema.json",
	"$id":         "http://example.com/base.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Base",
	"type":        "object",
	"description": "Base description",
	"properties": JSONProperty{
		"foo": JSONProperty{
			"type": "string",
		},
		"sub": JSONProperty{
			"$ref": "file://sub.schema.json",
		},
	},
	"required": []string{"foo", "sub"},
}

func JSONBase() []byte {
	out, _ := json.MarshalIndent(BaseSchema, "", "  ")
	return out
}

var SubSchema = JSONProperty{
	"fileName":    "sub.schema.json",
	"$id":         "http://example.com/sub.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Sub",
	"type":        "object",
	"description": "Sub description",
	"properties": JSONProperty{
		"sub": JSONProperty{
			"type": "string",
		},
		"sub-down": JSONProperty{
			"$ref": "file://sub2.schema.json",
		},
	},
	"required": []string{"bar", "sub-down"},
}

func JSONSub() []byte {
	out, _ := json.MarshalIndent(SubSchema, "", "  ")
	return out
}

var Sub2Schema = JSONProperty{
	"fileName":    "sub2.schema.json",
	"$id":         "http://example.com/sub2.schema.json",
	"$schema":     "http://json-schema.org/draft-07/schema#",
	"title":       "Sub2",
	"type":        "object",
	"description": "Sub2 description",
	"properties": JSONProperty{
		"bar": JSONProperty{
			"type": "string",
		},
	},
	"required": []string{"bar"},
}

func JSONSub2() []byte {
	out, _ := json.MarshalIndent(Sub2Schema, "", "  ")
	return out
}

func SchemaSchema(sl SchemaLoader) Property {
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
		Build()
}

func toPtrString(s string) *string {
	return &s
}

func TestJSONSubSchema() JSONProperty {
	return JSONProperty{
		"$id":         "https://Sub",
		"title":       "Sub",
		"description": "Description",
		"properties": JSONProperty{
			"Test": JSONProperty{
				"type": "string",
			},
			"opt-Test": JSONProperty{
				"type": "string",
			},
		},
		"required": []string{"Test"},
		"type":     "object",
	}
}

func TestJsonFlatSchema() JSONProperty {
	json := JSONProperty{
		"$id": "https://SimpleType",
		// "$schema":     "",
		"title":       "SimpleType",
		"type":        "object",
		"description": "Jojo SimpleType",
		"properties": JSONProperty{
			"string": JSONProperty{
				"type": "string",
			},
			"default-string": JSONProperty{
				"type":    "string",
				"default": "hallo",
			},
			"optional-string": JSONProperty{
				"type": "string",
			},
			"optional-default-string": JSONProperty{
				"type":    "string",
				"default": "hallo",
			},
			"createdAt": JSONProperty{
				"type":   "string",
				"format": "date-time",
			},
			"default-createdAt": JSONProperty{
				"type":    "string",
				"default": "2023-12-31T23:59:59Z",
				"format":  "date-time",
			},
			"optional-createdAt": JSONProperty{
				"type":   "string",
				"format": "date-time",
			},
			"optional-default-createdAt": JSONProperty{
				"type":    "string",
				"format":  "date-time",
				"default": "2023-12-31T23:59:59Z",
			},
			"float64": JSONProperty{
				"type": "number",
				// "format": "float64",
			},
			"default-float64": JSONProperty{
				"type":    "number",
				"format":  "float32",
				"default": 4711.4,
			},
			"optional-float32": JSONProperty{
				"type":   "number",
				"format": "float32",
			},
			"optional-default-float32": JSONProperty{
				"type":    "number",
				"default": 49.2,
				"format":  "float32",
			},
			"int64": JSONProperty{
				"type":   "integer",
				"format": "int64",
			},
			"default-int64": JSONProperty{
				"type":    "integer",
				"default": 64,
				"format":  "int64",
			},
			"optional-int32": JSONProperty{
				"type":   "integer",
				"format": "int32",
			},
			"optional-default-int32": JSONProperty{
				"type":    "integer",
				"default": 32,
				"format":  "int32",
			},
			"bool": JSONProperty{
				"type": "boolean",
			},
			"default-bool": JSONProperty{
				"type":    "boolean",
				"default": true,
			},
			"optional-bool": JSONProperty{
				"type": "boolean",
			},
			"optional-default-bool": JSONProperty{
				"type":    "boolean",
				"default": true,
			},
			"sub":     TestJSONSubSchema(),
			"opt-sub": TestJSONSubSchema(),
		},
		"required": []string{
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
	return NewPropertiesBuilder(sl).BuildObject().
		id("https://Sub").
		title("Sub").
		description("Description").
		propertiesAdd(NewPropertyItem("Test", NewPropertyString(PropertyStringParam{}))).
		propertiesAdd(NewPropertyItem("opt-Test", NewPropertyString(PropertyStringParam{}))).
		required([]string{"Test"}).
		Build()
}

func TestFlatSchema(sl SchemaLoader) Property {
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
		propertiesAdd(NewPropertyItem("sub", TestSubSchema(sl))).
		propertiesAdd(NewPropertyItem("opt-sub", TestSubSchema(sl))).
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
}

func TestSchema(sl SchemaLoader) Property {
	ps := NewPropertiesBuilder(sl).BuildObject()
	fls := TestFlatSchema(sl).(PropertyObject)
	for _, item := range fls.Items() {
		ps.propertiesAdd(item)
	}
	return NewPropertiesBuilder(sl).BuildObject().
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
						Items: TestSubSchema(sl)}),
				})})}))).
		propertiesAdd(NewPropertyItem("opt-arrayarrayFlatSchema", NewPropertyArray(PropertyArrayParam{
			Items: NewPropertyArray(PropertyArrayParam{
				Items: NewPropertyArray(PropertyArrayParam{
					Items: NewPropertyArray(PropertyArrayParam{
						Items: TestSubSchema(sl)}),
				})})}))).
		propertiesAdd(NewPropertyItem("sub-flat", TestSubSchema(sl))).
		propertiesAdd(NewPropertyItem("opt-sub-flat", TestSubSchema(sl))).
		propertiesAdd(NewPropertyItem("arraySubType", NewPropertyArray(PropertyArrayParam{Items: TestSubSchema(sl)}))).
		propertiesAdd(NewPropertyItem("opt-arraySubType", NewPropertyArray(PropertyArrayParam{Items: TestSubSchema(sl)}))).
		required(append([]string{
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

func TestXxx(tx *testing.T) {
	dec := json.NewDecoder(strings.NewReader(`{
		"doof": "doof",
		"type": "object",
		"properties": {
			"xxxx": {
				"type": "string"
			}
		},
	}`))

	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for dec.More() {
		var m struct {
			Type string `json:"type"`
		}
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", m.Type)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

}
