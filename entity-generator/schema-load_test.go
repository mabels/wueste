package entity_generator

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type JSONPropertyItems struct {
	Name string
	Prop JSONDict
}

func TestFlatJsonAndProp(t *testing.T) {
	jsobj := TestJsonFlatSchema()
	prop := NewPropertiesBuilder(NewTestContext()).
		FromJson(jsobj.JSONProperty).Build().Ok().(PropertyObject)
	pjs := PropertyToJson(prop)
	// assert.Equal(t, jsobj, pjs)

	jsonJsObj, err := json.MarshalIndent(jsobj.JSONProperty, "", "  ")
	assert.NoError(t, err)
	jsonPjs, err := json.MarshalIndent(pjs, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, string(jsonJsObj), string(jsonPjs))

	// assert.Equal(t, prop, ref)
	// pjsProps := toSorted(pjs.Properties)
	// pjs.Properties = nil
	// jsp := JSONFromProperty(jsprop)
	// jspProps := toSorted(jsp.Properties)
	// jsp.Properties = nil
	// assert.Equal(t, pjs, jsp)
	// assert.Equal(t, len(pjsProps), len(jspProps))
	// for i, _ := range pjsProps {
	// 	assert.Equal(t, pjsProps[i], jspProps[i], "Property %d:%s", i, pjsProps[i].Name)
	// }
}

func TestPayloadOpen(t *testing.T) {
	ctx := NewTestContext()
	v := TestPayloadSchema(ctx)
	assert.True(t, v.IsOk())
	assert.Equal(t, v.Ok().Meta().FileName().Value(), "/abs/payload.schema.json")
	assert.True(t, v.Ok().Meta().Parent().IsNone())

	vo := v.Ok().(PropertyObject)
	test, found := vo.Properties().Lookup("Test")
	assert.True(t, found)
	// assert.Equal(t, open.Id(), "")
	assert.Equal(t, test.Meta().FileName().Value(), "/abs/payload.schema.json")

	vo = v.Ok().(*propertyObject)
	open, found := vo.Properties().Lookup("Open")
	assert.True(t, found)
	// assert.Equal(t, open.Id(), "")
	assert.Equal(t, open.Meta().FileName().Value(), "/abs/payload.schema.json")
	assert.Equal(t, open.Meta().Parent().Value(), v.Ok())
}

func TestFileNames(t *testing.T) {

	ctx := NewTestContext()
	prop := NewJSONDict()
	prop.Set("$ref", "file://./payload.schema.json")
	rsub := NewPropertiesBuilder(ctx).FromJson(prop).Build()
	assert.True(t, rsub.IsOk())
	sub := rsub.Ok().(PropertyObject)
	assert.Equal(t, sub.Meta().FileName().Value(), "/abs/payload.schema.json")
	_, found := ctx.Registry.registry[sub.Meta().FileName().Value()]
	assert.True(t, found)

	prop.Set("$ref", "file://./simple_type.schema.json")
	base := NewPropertiesBuilder(ctx).FromJson(prop).Build().Ok().(PropertyObject)
	assert.Equal(t, base.Meta().FileName().Value(), "/abs/simple_type.schema.json")
	_, found = ctx.Registry.registry[base.Meta().FileName().Value()]
	assert.True(t, found)
	baseSub, _ := base.Properties().Lookup("opt-sub")
	assert.Equal(t, baseSub.Meta().FileName().Value(), "/abs/simple_type.schema.json")

	prop.Set("$ref", "file://./nested_type.schema.json")
	nested := NewPropertiesBuilder(ctx).FromJson(prop).Build().Ok().(PropertyObject)
	assert.Equal(t, nested.Meta().FileName().Value(), "/abs/nested_type.schema.json")
	_, found = ctx.Registry.registry[nested.Meta().FileName().Value()]
	assert.True(t, found)

}

func TestNestedJsonAndProp(t *testing.T) {
	ref := TestSchema(NewTestContext())
	refJs := PropertyToJson(ref)
	ret := NewPropertiesBuilder(NewTestContext()).FromJson(refJs).Build().Ok().(PropertyObject)
	retJs := PropertyToJson(ret)
	// ref := TestFlatSchema(NewTestSchemaLoader()).(PropertyObject)

	jsonRefJs, err := json.MarshalIndent(refJs, "", "  ")
	assert.NoError(t, err)
	jsonRetJs, err := json.MarshalIndent(retJs, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, string(jsonRefJs), string(jsonRetJs))

	assert.Equal(t, refJs, retJs)
	// assert.Equal(t, prop, ref)
	// pjsProps := toSorted(pjs.Properties)
	// pjs.Properties = nil
	// jsp := JSONFromProperty(jsprop)
	// jspProps := toSorted(jsp.Properties)
	// jsp.Properties = nil
	// assert.Equal(t, pjs, jsp)
	// assert.Equal(t, len(pjsProps), len(jspProps))
	// for i, _ := range pjsProps {
	// 	assert.Equal(t, pjsProps[i], jspProps[i], "Property %d:%s", i, pjsProps[i].Name)
	// }
}

func TestSchemaLoadImpl(t *testing.T) {
	baseFname, _ := filepath.Abs("./base.schema.json")
	baseDir := filepath.Dir(baseFname)
	sl := NewSchemaLoaderImpl("./wueste", "./ts")
	fname, err := sl.Abs("schema-load_test.go")
	assert.NoError(t, err)
	assert.Equal(t, fname, filepath.Join(baseDir, "schema-load_test.go"))
	_, err = sl.Abs("xchema-load_test.go")
	assert.Error(t, err)
	fname, err = sl.Abs("ts-generator_test.go")
	assert.NoError(t, err)
	assert.Equal(t, fname, filepath.Join(baseDir, "ts", "ts-generator_test.go"))
}

// type JsonProperty struct {
// 	Type        string  `json:"type"`
// 	Description *string `json:"description,omitempty"`
// 	FullType    interface{}
// }

// func (jp *JsonProperty) UnmarshalJSON(data []byte) error {
// 	var my struct {
// 		Type string `json:"type"`
// 	}
// 	err := json.Unmarshal(data, &my)
// 	if err != nil {
// 		return err
// 	}
// 	switch my.Type {
// 	case "object":
// 		jp.FullType = &JSONPropertyObject{}
// 	case "string":
// 		jp.FullType = &JSONPropertyString{}
// 	default:
// 		return fmt.Errorf("unknown type %s", my.Type)
// 	}
// 	err = json.Unmarshal(data, jp.FullType)
// 	return err
// }

// type SchemaString interface {
// 	Schema
// 	Default() *string
// }

// type SchemaObject interface {
// 	Schema
// 	Properties() map[string]Schema
// }

// type _SchemaBase struct {
// 	Type        string  `json:"type"`
// 	Description *string `json:"description,omitempty"`
// }

// type _Schema struct {
// 	_SchemaBase
// 	_fullType interface{}
// }
// type _SchemaObject struct {
// 	_SchemaBase
// 	Properties map[string]_Schema `json:"properties"`
// 	// _fullType  interface{}
// }
// type _SchemaString struct {
// 	_SchemaBase
// 	Default *string `json:"default,omitempty"`
// 	// _fullType interface{}
// }

// func (st *SchameLoader[int]) UnmarshalJSON(data []byte) error {
// var my struct {
// 	Type string `json:"type"`
// }
// err := json.Unmarshal(data, &my)
// if err != nil {
// 	return err
// }
// switch my.Type {
// case "object":
// 	st._fullType = &_SchemaObject{}
// case "string":
// 	st._fullType = &_SchemaString{}
// default:
// 	return fmt.Errorf("unknown type %s", my.Type)
// }
// err = json.Unmarshal(data, st._fullType)
// return err
// }

// type JsonSchemaBuilder struct {
// 	typ           string
// 	objectBuilder *JsonObjectBuilder
// }

// func NewJsonSchemaBuilder() *JsonSchemaBuilder {
// 	return &JsonSchemaBuilder{}
// }

// func (b *JsonSchemaBuilder) ObjectType() *JsonObjectBuilder {
// 	b.typ = "object"
// 	b.objectBuilder = &JsonObjectBuilder{}
// 	return b.objectBuilder
// }

// func (b *JsonSchemaBuilder) StringType() *JsonStringBuilder {
// 	b.typ = "string"
// 	return &JsonStringBuilder{}
// }

// func (b *JsonSchemaBuilder) BooleanType() *JsonBooleanBuilder {
// 	b.typ = "boolean"
// 	return &JsonBooleanBuilder{}
// }

// func (b *JsonSchemaBuilder) IntegerType() *JsonIntegerBuilder {
// 	b.typ = "integer"
// 	return &JsonIntegerBuilder{}
// }

// func (b *JsonSchemaBuilder) NumberType() *JsonNumberBuilder {
// 	b.typ = "number"
// 	return &JsonNumberBuilder{}
// }

// func (b *JsonSchemaBuilder) ArrayType() *JsonArrayBuilder {
// 	b.typ = "array"
// 	return &JsonArrayBuilder{}
// }

// type JsonObjectBuilder struct {
// 	param PropertyObjectParam
// }

// func (b *JsonObjectBuilder) FromObject(obj map[string]any) *JsonObjectBuilder {

// 	// FileName    string
// 	// Id          string
// 	// Title       string
// 	// Schema      string
// 	// Description rusty.Optional[string]
// 	// Properties  PropertiesObject
// 	// Required    []string
// 	// Ref         rusty.Optional[string]
// 	return b
// }

// type JsonStringBuilder struct {
// 	param PropertyStringParam
// }

// type JsonBooleanBuilder struct {
// 	param PropertyBooleanParam
// }

// type JsonNumberBuilder struct {
// 	param PropertyNumberParam
// }

// type JsonIntegerBuilder struct {
// 	param PropertyIntegerParam
// }

// type JsonArrayBuilder struct {
// 	param PropertyArrayParam
// }

// func TestXxx(tx *testing.T) {

// 	builder := NewJsonSchemaBuilder()
// 	builder.ObjectType()

// dec := []byte(`{
// 	"doof": "doof",
// 	"type": "object",
// 	"properties": {
// 		"xxxx": {
// 			"type": "string"
// 		}
// 	}
// }`)
// var my interface{}
// err := json.Unmarshal(dec, &my)
// t := reflect.TypeOf(my)
// switch t.Kind() {
// case reflect.Map:
// 	_map := my.(map[string]interface{})
// 	typ, found := _map["type"].(string)
// 	if found {
// 		fmt.Printf(">>>map>>>:%v", typ)
// 	} else {
// 		fmt.Printf("no type")
// 	}
// 	switch typ {
// 	case "object":
// 		jsobj := JSONPropertyObject{}
// 		rjsobj := reflect.TypeOf(jsobj)
// 		for i := 0; i < rjsobj.NumField(); i++ {
// 			field := rjsobj.Field(i)
// 			jsonTag, found := field.Tag.Lookup("json")
// 			var key string
// 			if found {
// 				key = strings.Split(jsonTag, ",")[0]
// 			} else {
// 				key = field.Name
// 			}
// 			val, found := _map[key]
// 			if !found {
// 				continue
// 			}
// 			switch val := val.(type) {
// 			case string:
// 				jsval := reflect.ValueOf(jsobj).Field(i)
// 				jsval.SetString(val)
// 			default:
// 				panic("xxxxx")
// 			}
// 			fmt.Println(jsobj)
// 		}
// 	default:
// 		panic("unknown type:" + typ)
// 	}
// 	// for fieldName, _ := range my {
// 	// field, found := t.FieldByName(fieldName)
// 	// if !found {
// 	// continue
// 	// }
// 	// fmt.Printf("\nField: Me.%s\n", fieldName, field.Tag.Get("json"))

// 	// }
// 	// for _, fieldName := range []string{"Firstname", "Lastname"} {
// 	// field, found := t.FieldByName(fieldName)
// 	// if !found {
// 	// 	continue
// 	// }
// 	// fmt.Printf("\nField: Me.%s\n", fieldName)
// 	// fmt.Printf("\tWhole tag value : %q\n", field.Tag)
// 	// fmt.Printf("\tValue of 'mytag': %q\n", field.Tag.Get("mytag"))
// case reflect.Array:
// 	panic("unknown type")
// default:
// 	panic("unknown type")
// }

// // for _, fieldName := range []string{"Firstname", "Lastname"} {
// // 	field, found := t.FieldByName(fieldName)
// // 	if !found {
// // 		continue
// // 	}
// // 	fmt.Printf("\nField: Me.%s\n", fieldName)
// // 	fmt.Printf("\tWhole tag value : %q\n", field.Tag)
// // 	fmt.Printf("\tValue of 'mytag': %q\n", field.Tag.Get("mytag"))
// // }
// // switch my.()
// assert.NoError(tx, err)
// }
