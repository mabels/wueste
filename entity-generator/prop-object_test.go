package entity_generator

import (
	"testing"

	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/stretchr/testify/assert"
)

func TestResolveRef(t *testing.T) {
	rt := PropertyRuntime{}
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}

	prop := NewProperty(PropertyParam{
		Ref: rusty.Some("file://./base.schema.json"),
		// Runtime: rt,
	})
	po := NewPropertiesBuilder(ctx).Resolve(rt, prop).Ok().Runtime().ToPropertyObject()
	assert.True(t, po.IsOk())
	assert.Equal(t, po.Ok().Runtime().FileName.Value(), "/abs/base.schema.json")

	direct := NewPropertiesBuilder(ctx).Resolve(rt, NewProperty(PropertyParam{
		Ref: rusty.Some("file://./base.schema.json"),
	}))
	if direct.IsErr() {
		t.Fatal(direct.Err())
	}
	assert.Equal(t, PropertyToJson(po.Ok()), PropertyToJson(direct.Ok()))
}

func TestDeref(t *testing.T) {
	rt := PropertyRuntime{}
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	err := NewPropertiesBuilder(ctx).Resolve(rt, NewProperty(PropertyParam{
		Ref: rusty.Some("file://./base.schema.json"),
	}))
	if err.IsErr() {
		t.Fatal(err.Err())
	}
	pobj := err.Ok().Runtime().ToPropertyObject().Ok()
	assert.Equal(t, pobj.Runtime().FileName.Value(), "/abs/base.schema.json")
	assert.Equal(t, pobj.Runtime().Ref.Value(), "file://./base.schema.json")
	items := pobj.Items()
	assert.Equal(t, items[0].Property().Runtime().FileName.Value(), "/abs/base.schema.json")
	assert.Equal(t, items[0].Property().Runtime().Ref.Value(), "file://./base.schema.json")
	assert.Equal(t, items[1].Property().Runtime().FileName.Value(), "/abs/sub.schema.json")
	assert.Equal(t, items[1].Property().Runtime().Ref.Value(), "file://sub.schema.json")
	// pobj := p.(PropertyObject)
	// keys := []PropertyObject{}
	// for _, po := range tsl.registry.Items() {
	// 	keys = append(keys, po.PropertItem().Property().(PropertyObject))
	// }
	// items := registry.Items()
	// sort.Slice(items, func(i, j int) bool {
	// 	return items[i].Property().Id() < items[j].Property().Id()
	// })
	abs := func(path string) string {
		p, _ := ctx.Registry.loader.Abs(path)
		return p
	}
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
		written: false,
		prop: // NewPropertyItem("http://example.com/base.schema.json",
		NewPropertiesBuilder(ctx).BuildObject().FromJson(rt, BaseSchema).
			fileName(abs(BaseSchema.Get("fileName").(string))).Build(),
	}
	testSub := &schemaRegistryItem{
		written: false,
		prop: // NewPropertyItem("http://example.com/sub.schema.json",
		NewPropertiesBuilder(ctx).BuildObject().FromJson(
			*base.prop.Runtime().ToPropertyObject().Ok().PropertyByName("sub").Property().Runtime(), TestJsonSubSchema()).
			fileName(abs(TestJsonSubSchema().Get("fileName").(string))).Build(),
	}
	sub2 := &schemaRegistryItem{
		written: false,
		prop: // NewPropertyItem("http://example.com/sub2.schema.json",
		NewPropertiesBuilder(ctx).BuildObject().
			FromJson(*testSub.prop.Runtime().ToPropertyObject().Ok().PropertyByName("sub-down").Property().Runtime(), Sub2Schema).
			fileName(abs(Sub2Schema.Get("fileName").(string))).Build(),
	}
	sub3 := &schemaRegistryItem{
		written: false,
		prop: // NewPropertyItem("http://example.com/sub2.schema.json",
		NewPropertiesBuilder(ctx).BuildObject().
			FromJson(*sub2.prop.Runtime().ToPropertyObject().Ok().PropertyByName("bar").Property().Runtime(), Sub3Schema).
			fileName(abs(Sub3Schema.Get("fileName").(string))).Build(),
	}
	refItems := []SchemaRegistryItem{base, testSub, sub2, sub3}
	assert.Equal(t, len(filtered), len(refItems))
	for i, _ := range filtered {
		assert.Equal(t, filtered[i].Property().Id(), refItems[i].Property().Id())
		assert.Equal(t, filtered[i].Property().Runtime().FileName.Value(), refItems[i].Property().Runtime().FileName.Value())
		assert.Equal(t, PropertyToJson(filtered[i].Property()), PropertyToJson(refItems[i].Property()))
	}
}

func TestLoaderResolveRef(t *testing.T) {
	rt := PropertyRuntime{}
	ctx := PropertyCtx{
		Registry: NewSchemaRegistry(&TestSchemaLoader{}),
	}
	p := NewPropertiesBuilder(ctx).Resolve(rt, NewProperty(PropertyParam{
		Ref: rusty.Some("file://./base.schema.json"),
	}))
	if p.IsErr() {
		t.Fatal(p.Err())
	}

	po := p.Ok().Runtime().ToPropertyObject().Ok()
	items := po.Items()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "foo", items[0].Name())
	assert.Equal(t, "sub", items[1].Name())

	po1 := items[1].Property().(PropertyObject)
	assert.Equal(t, po1.Items()[0].Name(), "sub")
	assert.Equal(t, po1.Items()[0].Property().(PropertyString).Type(), "string")
	assert.Equal(t, po1.Runtime().FileName.Value(), "/abs/sub.schema.json")
	po2 := po1.Items()[1].Property().(PropertyObject)
	assert.Equal(t, po2.Runtime().FileName.Value(), "/abs/wurst/sub2.schema.json")

}
