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
	includeDirs []string
	// registry *SchemaRegistry
}

func NewTestSchemaLoader() *TestSchemaLoader {
	return &TestSchemaLoader{
		includeDirs: []string{"/abs", "./"},
	}
}

func NewTestContext() PropertyCtx {
	return PropertyCtx{
		Registry: NewSchemaRegistry(NewTestSchemaLoader()),
	}
}

func (sl TestSchemaLoader) Clone(prependIncDirs ...string) SchemaLoader {
	return &TestSchemaLoader{
		includeDirs: append(prependIncDirs, sl.includeDirs...),
	}
}

func (t *TestSchemaLoader) IncludeDirs() []string {
	return t.includeDirs
}

func (l *TestSchemaLoader) Abs(path string) (string, error) {
	if strings.HasPrefix(path, "/abs/") {
		return path, nil
	}
	return filepath.Join(l.includeDirs[0], path), nil
}

// ReadFile implements SchemaLoader.
func (l *TestSchemaLoader) ReadFile(path string) ([]byte, error) {
	switch path {
	case "/abs/unnamed_nested_object.schema.json":
		return JSONUnnamedNestedObject(), nil
	case "/abs/payload.schema.json":
		jf := TestJSONPayloadSchema()
		bytes, _ := json.MarshalIndent(jf.JSONProperty, "", "  ")
		return bytes, nil
	case "/abs/base.schema.json":
		return JSONBase(), nil
	case "/abs/sub.schema.json":
		return JSONSub(), nil
	case "/abs/wurst/sub2.schema.json":
		return JSONSub2(), nil
	case "/abs/wurst/sub3.schema.json":
		return JSONSub3(), nil
	case "/abs/simple_type.schema.json":
		jf := TestJsonFlatSchema()
		bytes, _ := json.MarshalIndent(jf.JSONProperty, "", "  ")
		return bytes, nil
	case "/abs/nested_type.schema.json":
		jf := TestJSONSchema()
		bytes, _ := json.MarshalIndent(jf.JSONProperty, "", "  ")
		return bytes, nil
	default:
		return nil, fmt.Errorf("Test file not found: %s", path)
	}
}

// Unmarshal implements SchemaLoader.
func (l *TestSchemaLoader) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

func json2JSonFile(inp string) JSonFile {
	jsonFile := NewJSONDict()
	err := json.Unmarshal([]byte(inp), &jsonFile)
	if err != nil {
		panic(fmt.Errorf("json2JSonFile: %w:%v", err, inp))
	}
	return JSonFile{
		FileName:     jsonFile.Get("filename").(string),
		JSONProperty: jsonFile.Get("jsonProperty").(JSONDict),
	}
}

func UnnamedNestedObject() JSonFile {
	return json2JSonFile(`{
	"filename":    "unnamed_nested_object.schema.json",
	"jsonProperty": {
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
	}
}`)
}

func JSONUnnamedNestedObject() []byte {
	out, _ := json.MarshalIndent(UnnamedNestedObject, "", "  ")
	return out
}

func BaseSchema() JSonFile {
	return json2JSonFile(`{
	"filename":    "base.schema.json",
	"jsonProperty": {
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
	}
}`)
}

func JSONBase() []byte {
	out, _ := json.MarshalIndent(BaseSchema().JSONProperty, "", "  ")
	return out
}

func TestJsonSubSchema() JSonFile {
	return json2JSonFile(`{
	"filename":    "sub.schema.json",
	"jsonProperty": {
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
	}
}`)

}

func JSONSub() []byte {
	out, _ := json.MarshalIndent(TestJsonSubSchema().JSONProperty, "", "  ")
	return out
}

func Sub2Schema() JSonFile {
	return json2JSonFile(`{
	"filename":    "wurst/sub2.schema.json",
	"jsonProperty": {
		"$id":         "http://example.com/sub2.schema.json",
		"$schema":     "http://json-schema.org/draft-07/schema#",
		"title":       "Result",
		"type":        "object",
		"description": "Sub2 description",
		"properties": {
			"bar": {
				"$ref": "file://./sub3.schema.json"
			},
			"maxSheep": {
				"$id": "Sheep",
				"title": "Sheep",
				"type": "object",
				"properties": {
					"flat": {
						"$ref": "file://./sub3.schema.json"
					},
					"nested": {
						"type": "array",
						"items": {
							"$ref": "file://./sub3.schema.json"
						}
					}
				}
			}
		},
		"required": ["bar"]
	}
}`)
}

func JSONSub2() []byte {
	out, _ := json.MarshalIndent(Sub2Schema().JSONProperty, "", "  ")
	return out
}

func Sub3Schema() JSonFile {
	return json2JSonFile(`{
	"filename":    "wurst/sub3.schema.json",
	"jsonProperty": {
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
	}
}`)
}

func JSONSub3() []byte {
	out, _ := json.MarshalIndent(Sub3Schema().JSONProperty, "", "  ")
	return out
}

// func SchemaSchema(sl PropertyCtx) Property {
// 	return NewPropertiesBuilder(sl).BuildObject().
// 		id("JsonSchema").
// 		title("JsonSchema").
// 		description("JSON Schema").
// 		propertiesAdd(NewPropertyItem("$id", NewPropertyString(PropertyStringParam{}))).
// 		propertiesAdd(NewPropertyItem("$schema", NewPropertyString(PropertyStringParam{}))).
// 		propertiesAdd(NewPropertyItem("title", NewPropertyString(PropertyStringParam{}))).
// 		propertiesAdd(NewPropertyItem("type", NewPropertyString(PropertyStringParam{}))).
// 		propertiesAdd(NewPropertyItem("properties", NewPropertyObject(PropertyObjectParam{}))).
// 		propertiesAdd(NewPropertyItem("required", NewPropertyArray(PropertyArrayParam{}))).
// 		required([]string{"$id", "$schema", "title", "type", "properties"}).
// 		Build().Ok()
// }

// func toPtrString(s string) *string {
// 	return &s
// }

func TestJSONPayloadSchema() JSonFile {
	return json2JSonFile(`{
		"filename":	"payload.schema.json",
		"jsonProperty": {
			"type":     "object",
			"$id":         "https://IPayload",
			"title":       "IPayload",
			"description": "Description",
			"properties": {
				"Test": {
					"type": "string"
				},
				"opt-Test": {
					"type": "string"
				},
				"Open": {
					"type": "object"
				},
				"opt-Open": {
					"type": "object"
				}
			},
			"required": ["Test", "Open"]
		}
	}`)
}

func TestJsonFlatSchema() JSonFile {
	json := json2JSonFile(`{
		"filename":    "simple_type.schema.json",
		"jsonProperty": {
			"type":        "object",
			"$id": "https://SimpleType",
			"title":       "SimpleType",
			"description": "Jojo SimpleType",
			"properties": {
				"string": {
					"type": "string",
					"description": "string description",
					"x-groups": ["string", "key", "primary-key"]
				},
				"default-string": {
					"type":    "string",
					"default": "hallo",
					"x-groups": ["key"]
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
					"type": "number",
					"x-groups": ["number", "key", "primary-key"]
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
					"format": "int64",
					"x-groups": ["integer", "key"]
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
					"type": "boolean",
					"x-groups": ["boolean", "key"]
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
		}
	}`)
	// _prop := json.JSONProperty.Get("properties").(JSONProperty)
	// prop, ok := _prop.(JSONProperty)
	// if !ok {
	// 	panic("not ok")
	// }
	json.JSONProperty.Get("properties").(JSONDict).Set("sub", TestJSONPayloadSchema().JSONProperty)
	json.JSONProperty.Get("properties").(JSONDict).Set("opt-sub", TestJSONPayloadSchema().JSONProperty)
	return json
}

func TestPayloadSchema(sl PropertyCtx) rusty.Result[Property] {
	prop := NewJSONDict()
	prop.Set("$ref", "file://payload.schema.json")
	return NewPropertiesBuilder(sl).FromJson(prop).Build()

	// return NewPropertiesBuilder(sl).SetFileName("/abs/payload.schema.json").assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
	// 	return NewPropertyObjectBuilder(b).
	// 		id("https://Sub").
	// 		title("Payload").
	// 		description("Description").
	// 		propertiesAdd(NewPropertyItem("Test", NewPropertyString(PropertyStringBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("opt-Test", NewPropertyString(PropertyStringBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("Open", NewPropertyObject(PropertyObjectBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("opt-Open", NewPropertyObject(PropertyObjectBuilder{}))).
	// 		required([]string{"Test", "Open"}).
	// 		Build()
	// }).Build()
}

func TestFlatSchema(sl PropertyCtx) rusty.Result[Property] {
	prop := NewJSONDict()
	prop.Set("$ref", "file://simple_type.schema.json")
	return NewPropertiesBuilder(sl).FromJson(prop).Build()

	// return NewPropertiesBuilder(sl).SetFileName("/abs/simple_type.schema.json").assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
	// 	return NewPropertyObjectBuilder(b).
	// 		id("https://SimpleType").
	// 		title("SimpleType").
	// 		description("Jojo SimpleType").
	// 		propertiesAdd(NewPropertyItem("string", NewPropertyString(PropertyStringBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("default-string", NewPropertyString(PropertyStringBuilder{Default: rusty.Some("hallo")}))).
	// 		propertiesAdd(NewPropertyItem("optional-string", NewPropertyString(PropertyStringBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("optional-default-string", NewPropertyString(PropertyStringBuilder{Default: rusty.Some("hallo")}))).
	// 		propertiesAdd(NewPropertyItem("createdAt", NewPropertyString(PropertyStringBuilder{
	// 			Format: rusty.Some("date-time"),
	// 		}))).
	// 		propertiesAdd(NewPropertyItem("default-createdAt", NewPropertyString(PropertyStringBuilder{
	// 			Default: rusty.Some("2023-12-31T23:59:59Z"),
	// 			Format:  rusty.Some("date-time"),
	// 		}))).
	// 		propertiesAdd(NewPropertyItem("optional-createdAt", NewPropertyString(PropertyStringBuilder{
	// 			Format: rusty.Some("date-time"),
	// 		}))).
	// 		propertiesAdd(NewPropertyItem("optional-default-createdAt", NewPropertyString(PropertyStringBuilder{
	// 			Default: rusty.Some("2023-12-31T23:59:59Z"),
	// 			Format:  rusty.Some("date-time"),
	// 		}))).
	// 		propertiesAdd(NewPropertyItem("float64", NewPropertyNumber(PropertyNumberBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("default-float64", NewPropertyNumber(PropertyNumberBuilder{Default: rusty.Some(4711.4), Format: rusty.Some("float32")}))).
	// 		propertiesAdd(NewPropertyItem("optional-float32", NewPropertyNumber(PropertyNumberBuilder{Format: rusty.Some("float32")}))).
	// 		propertiesAdd(NewPropertyItem("optional-default-float32", NewPropertyNumber(PropertyNumberBuilder{Default: rusty.Some(49.2)}))).
	// 		propertiesAdd(NewPropertyItem("int64", NewPropertyInteger(PropertyIntegerBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("default-int64", NewPropertyInteger(PropertyIntegerBuilder{Default: rusty.Some(64)}))).
	// 		propertiesAdd(NewPropertyItem("optional-int32", NewPropertyInteger(PropertyIntegerBuilder{Format: rusty.Some("int32")}))).
	// 		propertiesAdd(NewPropertyItem("optional-default-int32", NewPropertyInteger(PropertyIntegerBuilder{Default: rusty.Some(32)}))).
	// 		propertiesAdd(NewPropertyItem("bool", NewPropertyBoolean(PropertyBooleanBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("default-bool", NewPropertyBoolean(PropertyBooleanBuilder{Default: rusty.Some(true)}))).
	// 		propertiesAdd(NewPropertyItem("optional-bool", NewPropertyBoolean(PropertyBooleanBuilder{}))).
	// 		propertiesAdd(NewPropertyItem("optional-default-bool", NewPropertyBoolean(PropertyBooleanBuilder{Default: rusty.Some(true)}))).
	// 		propertiesAdd(NewPropertyItem("sub", TestPayloadSchema(sl))).
	// 		propertiesAdd(NewPropertyItem("opt-sub", TestPayloadSchema(sl))).
	// 		required([]string{
	// 			"string",
	// 			"createdAt",
	// 			"default-string",
	// 			"default-createdAt",
	// 			"float64",
	// 			"default-float64",
	// 			"int64",
	// 			"default-int64",
	// 			"uint64",
	// 			"default-uint64",
	// 			"default-bool",
	// 			"bool",
	// 			"sub"}).
	// 		Build()
	// }).Build()
}

func TestJSONSchema() JSonFile {
	dict := NewJSONDict()
	jsonStr := `{
		"$id": "https://NestedType",
		"title": "NestedType",
		"description": "Jojo NestedType",
		"type": "object",
		"properties": {
			"arrayarrayBool": {
				"type": "array",
				"items": {
					"type": "array",
					"items": {
						"type": "array",
						"items": {
							"type": "array",
							"items": {
								"type": "boolean"
							}
						}
					}
				}
			},
			"opt-arrayarrayBool": {
				"type": "array",
				"items": {
					"type": "array",
					"items": {
						"type": "array",
						"items": {
							"type": "array",
							"items": {
								"type": "boolean"
							}
						}
					}
				}
			},
			"arrayString": {
				"type": "array",
				"items": {
					"type": "string"
				}
			},
			"opt-arrayString": {
				"type": "array",
				"items": {
					"type": "string"
				}
			},
			"arrayNumber": {
				"type": "array",
				"items": {
					"type": "number"
				}
			},
			"opt-arrayNumber": {
				"type": "array",
				"items": {
					"type": "number"
				}
			},
			"arrayInteger": {
				"type": "array",
				"items": {
					"type": "integer"
				}
			},
			"opt-arrayInteger": {
				"type": "array",
				"items": {
					"type": "integer"
				}
			},
			"arrayBool": {
				"type": "array",
				"items": {
					"type": "boolean"
				}
			},
			"opt-arrayBool": {
				"type": "array",
				"items": {
					"type": "boolean"
				}
			},
			"arrayarrayFlatSchema": {
				"type": "array",
				"items": {
					"type": "array",
					"items": {
						"type": "array",
						"items": {
							"type": "array",
							"items": {}
						}
					}
				}
			},
			"opt-arrayarrayFlatSchema": {
				"type": "array",
				"items": {
					"type": "array",
					"items": {
						"type": "array",
						"items": {
							"type": "array",
							"items": {}
						}
					}
				}
			},
			"arraySubType": {
				"type": "array",
				"items": {}
			},
			"opt-arraySubType": {
				"type": "array",
				"items": {}
			}
		},
		"required": [
			"arrayarrayBool",
			"sub-flat",
			"arrayString",
			"arrayNumber",
			"arrayInteger",
			"arrayBool",
			"arraySubType",
			"arrayarrayFlatSchema"
		]
	}`
	err := json.Unmarshal([]byte(jsonStr), dict)
	if err != nil {
		panic(err)
	}
	payloadFile := TestJSONPayloadSchema()

	props := dict.Get("properties").(JSONDict)

	props.Get("arraySubType").(JSONDict).Set("items", payloadFile.JSONProperty)
	props.Get("opt-arraySubType").(JSONDict).Set("items", payloadFile.JSONProperty)

	props.Set("sub-flat", payloadFile.JSONProperty)
	props.Set("opt-sub-flat", payloadFile.JSONProperty)

	props.Get("arrayarrayFlatSchema").(JSONDict).
		Get("items").(JSONDict).
		Get("items").(JSONDict).
		Get("items").(JSONDict).
		Set("items", payloadFile.JSONProperty)

	props.Get("opt-arrayarrayFlatSchema").(JSONDict).
		Get("items").(JSONDict).
		Get("items").(JSONDict).
		Get("items").(JSONDict).
		Set("items", payloadFile.JSONProperty)

	flat := TestJsonFlatSchema()
	flatProp := flat.JSONProperty.Get("properties").(JSONDict)

	for _, item := range flatProp.Keys() {
		val := flatProp.Get(item)
		props.Set(item, val)
	}
	dict.Set("required", append(
		dict.Get("required").([]interface{}),
		flat.JSONProperty.Get("required").([]interface{})...))

	return JSonFile{
		FileName:     "nested_type.schema.json",
		JSONProperty: dict,
	}
}

func TestSchema(sl PropertyCtx) Property {
	prop := NewJSONDict()
	prop.Set("$ref", "file://./nested_type.schema.json")
	return NewPropertiesBuilder(sl).FromJson(prop).Build().Ok()

	// pb := NewPropertiesBuilder(sl).SetFileName("/abs/nested_type.schema.json")
	// ps := NewPropertyObjectBuilder(pb)
	// fls := TestFlatSchema(sl).Ok().(PropertyObject)
	// for _, item := range fls.Items() {
	// 	ps.propertiesAdd(rusty.Ok(item))
	// }
	// return ps.
	// 	id("https://NestedType").
	// 	title("NestedType").
	// 	description("Jojo NestedType").
	// 	propertiesAdd(NewPropertyItem("arrayarrayBool", NewPropertyArray(PropertyArrayBuilder{
	// 		Items: NewPropertyArray(PropertyArrayBuilder{
	// 			Items: NewPropertyArray(PropertyArrayBuilder{
	// 				Items: NewPropertyArray(PropertyArrayBuilder{
	// 					Items: NewPropertyBoolean(PropertyBooleanBuilder{})}),
	// 			})})}))).
	// 	propertiesAdd(NewPropertyItem("opt-arrayarrayBool", NewPropertyArray(PropertyArrayBuilder{
	// 		Items: NewPropertyArray(PropertyArrayBuilder{
	// 			Items: NewPropertyArray(PropertyArrayBuilder{
	// 				Items: NewPropertyArray(PropertyArrayBuilder{
	// 					Items: NewPropertyBoolean(PropertyBooleanBuilder{})}),
	// 			})})}))).
	// 	propertiesAdd(NewPropertyItem("arrayString", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyString(PropertyStringBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("opt-arrayString", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyString(PropertyStringBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("arrayNumber", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyNumber(PropertyNumberBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("opt-arrayNumber", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyNumber(PropertyNumberBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("arrayInteger", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyInteger(PropertyIntegerBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("opt-arrayInteger", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyInteger(PropertyIntegerBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("arrayBool", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyBoolean(PropertyBooleanBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("opt-arrayBool", NewPropertyArray(PropertyArrayBuilder{Items: NewPropertyBoolean(PropertyBooleanBuilder{})}))).
	// 	propertiesAdd(NewPropertyItem("arrayarrayFlatSchema", NewPropertyArray(PropertyArrayBuilder{
	// 		Items: NewPropertyArray(PropertyArrayBuilder{
	// 			Items: NewPropertyArray(PropertyArrayBuilder{
	// 				Items: NewPropertyArray(PropertyArrayBuilder{
	// 					Items: TestPayloadSchema(sl)}),
	// 			})})}))).
	// 	propertiesAdd(NewPropertyItem("opt-arrayarrayFlatSchema", NewPropertyArray(PropertyArrayBuilder{
	// 		Items: NewPropertyArray(PropertyArrayBuilder{
	// 			Items: NewPropertyArray(PropertyArrayBuilder{
	// 				Items: NewPropertyArray(PropertyArrayBuilder{
	// 					Items: TestPayloadSchema(sl)}),
	// 			})})}))).
	// 	propertiesAdd(NewPropertyItem("sub-flat", TestPayloadSchema(sl))).
	// 	propertiesAdd(NewPropertyItem("opt-sub-flat", TestPayloadSchema(sl))).
	// 	propertiesAdd(NewPropertyItem("arraySubType", NewPropertyArray(PropertyArrayBuilder{Items: TestPayloadSchema(sl)}))).
	// 	propertiesAdd(NewPropertyItem("opt-arraySubType", NewPropertyArray(PropertyArrayBuilder{Items: TestPayloadSchema(sl)}))).
	// 	required(append([]string{
	// 		"arrayarrayBool",
	// 		"sub-flat",
	// 		"arrayString",
	// 		"arrayNumber",
	// 		"arrayInteger",
	// 		"arrayBool",
	// 		"arraySubType",
	// 		"arrayarrayFlatSchema",
	// 	}, fls.Required()...)).
	// 	Build().Ok()
}

func WriteTestSchema(cfg *GeneratorConfig) string {
	jsonSchema := PropertyToJson(TestFlatSchema(NewTestContext()).Ok())
	bytes, _ := json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile := path.Join(cfg.OutputDir, "simple_type.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	jsonSchema = PropertyToJson(TestSchema(NewTestContext()))
	bytes, _ = json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile = path.Join(cfg.OutputDir, "nested_type.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	jsonSchema = BaseSchema().JSONProperty
	bytes, _ = json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile = path.Join(cfg.OutputDir, "base.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	jsonSchema = TestJsonSubSchema().JSONProperty
	bytes, _ = json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile = path.Join(cfg.OutputDir, "sub.schema.json")
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	jsonSchema = Sub2Schema().JSONProperty
	bytes, _ = json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile = path.Join(cfg.OutputDir, "wurst/sub2.schema.json")
	os.MkdirAll(path.Dir(schemaFile), 0755)
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	jsonSchema = Sub3Schema().JSONProperty
	bytes, _ = json.MarshalIndent(jsonSchema, "", "  ")
	schemaFile = path.Join(cfg.OutputDir, "wurst/sub3.schema.json")
	os.MkdirAll(path.Dir(schemaFile), 0755)
	os.WriteFile(schemaFile, bytes, 0644)
	fmt.Println("Wrote schema to -> ", schemaFile)

	return schemaFile
}
