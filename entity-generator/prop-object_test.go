package entity_generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveRef(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}

	prop := NewJSONProperty()
	prop.Set("$ref", "file://./base.schema.json")

	po := NewPropertiesBuilder(ctx).FromJson(prop).Build().Ok().(PropertyObject)
	assert.True(t, po.Meta().Parent().IsNone())
	assert.Equal(t, po.Meta().FileName().Value(), "/abs/base.schema.json")
	assert.Equal(t, po.Ref().Value(), "file://./base.schema.json")

	foo := po.PropertyByName("foo")
	assert.True(t, foo.IsOk())
	assert.Equal(t, foo.Ok().Property().Meta().FileName().Value(), "/abs/base.schema.json")
	assert.Equal(t, foo.Ok().Property().Type(), "string")

	sub := po.PropertyByName("sub")
	assert.True(t, sub.IsOk())
	assert.Equal(t, sub.Ok().Property().(PropertyObject).Ref().Value(), "file://sub.schema.json")
	assert.Equal(t, sub.Ok().Property().Meta().FileName().Value(), "/abs/sub.schema.json")
	assert.Equal(t, sub.Ok().Property().Type(), "object")
}

func TestParentSimple(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	prop := NewJSONProperty()
	prop.Set("$ref", "file://./base.schema.json")
	po := NewPropertiesBuilder(ctx).FromJson(prop).Build().Ok().(PropertyObject)
	assert.True(t, po.Meta().Parent().IsNone())
	foo := po.PropertyByName("foo").Ok().Property()
	assert.Equal(t, foo.Meta().Parent().Value().Id(), po.Id())
	sub := po.PropertyByName("sub").Ok().Property().(PropertyObject)
	assert.True(t, sub.Meta().Parent().IsSome())
	assert.Equal(t, sub.Meta().Parent().Value(), po)

	sub_sub := sub.PropertyByName("sub").Ok().Property()
	assert.Equal(t, sub_sub.Meta().Parent().Value().Id(), sub.Id())

	sub_down := sub.PropertyByName("sub-down").Ok().Property().(PropertyObject)
	assert.Equal(t, sub_down.Meta().Parent().Value().Id(), sub.Id())

	sub_down_bar := sub_down.PropertyByName("bar").Ok().Property()
	assert.Equal(t, sub_down_bar.Meta().Parent().Value(), sub_down)

}

func TestParentSimpleArray(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	prop := NewJSONProperty()
	prop.Set("$ref", "file://./nested_types.schema.json")
	so := NewPropertiesBuilder(ctx).FromJson(prop).Build().Ok().(PropertyObject)
	str := so.PropertyByName("string").Ok().Property()
	assert.Equal(t, str.Meta().Parent().Value().Id(), so.Id())

	arr0 := so.PropertyByName("arrayarrayFlatSchema").Ok().Property().(PropertyArray)
	assert.Equal(t, arr0.Meta().Parent().Value().Id(), so.Id())

	arr1 := arr0.Items().(PropertyArray)
	assert.Equal(t, arr1.Meta().Parent().Value().Id(), arr0.Id())

	arr2 := arr1.Items().(PropertyArray)
	assert.Equal(t, arr2.Meta().Parent().Value().Id(), arr1.Id())

	arr3 := arr2.Items().(PropertyArray)
	assert.Equal(t, arr3.Meta().Parent().Value().Id(), arr2.Id())

	obj := arr3.Items()
	assert.Equal(t, obj.Id(), "https://Sub")
	assert.Equal(t, obj.Meta().Parent().Value().Id(), arr3.Id())

}

func TestDeref(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	prop := NewJSONProperty()
	prop.Set("$ref", "file://./base.schema.json")
	rpo := NewPropertiesBuilder(ctx).FromJson(prop).Build()
	if rpo.IsErr() {
		t.Fatal(rpo.Err())
	}
	pobj := rpo.Ok().(PropertyObject)
	assert.Equal(t, pobj.Meta().FileName().Value(), "/abs/base.schema.json")
	assert.Equal(t, pobj.Ref().Value(), "file://./base.schema.json")
	items := pobj.Items()
	assert.Equal(t, items[0].Property().Meta().FileName().Value(), "/abs/base.schema.json")
	assert.Equal(t, items[0].Property().Ref().IsNone(), true)
	assert.Equal(t, items[1].Property().Meta().FileName().Value(), "/abs/sub.schema.json")
	assert.Equal(t, items[1].Property().Ref().Value(), "file://sub.schema.json")
	// pobj := p.(PropertyObject)
	// keys := []PropertyObject{}
	// for _, po := range tsl.registry.Items() {
	// 	keys = append(keys, po.PropertItem().Property().(PropertyObject))
	// }
	// items := registry.Items()
	// sort.Slice(items, func(i, j int) bool {
	// 	return items[i].Property().Id() < items[j].Property().Id()
	// })
	filter := []string{
		"/abs/base.schema.json",
		"/abs/sub.schema.json",
		"/abs/wurst/sub2.schema.json",
		"/abs/wurst/sub3.schema.json",
	}
	filtered := []SchemaRegistryItem{}
	for _, key := range filter {
		item, found := ctx.Registry.registry[key]
		if found {
			filtered = append(filtered, item)
		}
	}
	base := &schemaRegistryItem{
		written:  false,
		jsonFile: BaseSchema(),
	}
	testSub := &schemaRegistryItem{
		written:  false,
		jsonFile: TestJsonSubSchema(),
	}
	sub2 := &schemaRegistryItem{
		written:  false,
		jsonFile: Sub2Schema(),
	}
	sub3 := &schemaRegistryItem{
		written:  false,
		jsonFile: Sub3Schema(),
	}
	refItems := []SchemaRegistryItem{base, testSub, sub2, sub3}
	assert.Equal(t, len(filtered), len(refItems))
	// for i, _ := range filtered {
	// 	assert.Equal(t, filtered[i].JSonFile().JSONProperty.Get("$id"), refItems[i].JSonFile().JSONProperty.Get("$id"))
	// 	assert.Equal(t, filtered[i].JSonFile().JSONProperty.Get("$ref"), refItems[i].JSonFile().FileName)
	// }
}

func TestLoaderResolveRef(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	prop := NewJSONProperty()
	prop.Set("$ref", "file://./base.schema.json")
	p := NewPropertiesBuilder(ctx).FromJson(prop).Build()
	assert.False(t, p.IsErr())

	po := p.Ok().(PropertyObject)
	items := po.Items()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "foo", items[0].Name())
	assert.Equal(t, "sub", items[1].Name())

	po1 := items[1].Property().(PropertyObject)
	assert.Equal(t, po1.Items()[0].Name(), "sub")
	assert.Equal(t, po1.Items()[0].Property().(PropertyString).Type(), "string")
	assert.Equal(t, po1.Meta().FileName().Value(), "/abs/sub.schema.json")
	po2 := po1.Items()[1].Property().(PropertyObject)
	assert.Equal(t, po2.Meta().FileName().Value(), "/abs/wurst/sub2.schema.json")

}

func TestErrorUnnamedNestedObject(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	prop := NewJSONProperty()
	prop.Set("$ref", "file://./unnamed_nested_object.schema.json")
	p := NewPropertiesBuilder(ctx).FromJson(prop).Build()
	assert.True(t, p.IsErr())
}

func TestInstances(t *testing.T) {
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	res1 := TestPayloadSchema(ctx).Ok()
	res2 := TestPayloadSchema(ctx).Ok()

	if res1 == res2 {
		t.Fatal("res1 == res2")
	}
}
