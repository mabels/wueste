package entity_generator

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeref(t *testing.T) {
	tsl := &TestSchemaLoader{
		registry: NewSchemaRegistry(),
	}
	_, err := LoadSchema("base.schema.json", tsl)
	assert.NoError(t, err)
	// pobj := p.(PropertyObject)
	// keys := []PropertyObject{}
	// for _, po := range tsl.registry.Items() {
	// 	keys = append(keys, po.PropertItem().Property().(PropertyObject))
	// }
	items := tsl.registry.Items()
	sort.Slice(items, func(i, j int) bool {
		return items[i].PropertItem().Property().Id() < items[j].PropertItem().Property().Id()
	})
	abs := func(path string) string {
		p, _ := tsl.Abs(path)
		return p
	}
	refItems := []SchemaRegistryItem{
		&schemaRegistryItem{
			written: false,
			prop: NewPropertyItem("http://example.com/base.schema.json",
				NewSchemaBuilder(tsl).JSON2PropertyObject(BaseSchema).FileName(abs(BaseSchema.FileName)).Build().(PropertyObject)),
		},
		&schemaRegistryItem{
			written: false,
			prop: NewPropertyItem("http://example.com/sub.schema.json",
				NewSchemaBuilder(tsl).JSON2PropertyObject(SubSchema).FileName(abs(SubSchema.FileName)).Build().(PropertyObject)),
		},
		&schemaRegistryItem{
			written: false,
			prop: NewPropertyItem("http://example.com/sub2.schema.json",
				NewSchemaBuilder(tsl).JSON2PropertyObject(Sub2Schema).FileName(abs(Sub2Schema.FileName)).Build().(PropertyObject)),
		},
	}
	assert.Equal(t, items, refItems)
}

func TestResolveRef(t *testing.T) {
	tsl := &TestSchemaLoader{
		registry: NewSchemaRegistry(),
	}
	p, err := LoadSchema("base.schema.json", tsl)
	assert.NoError(t, err)

	po := p.(PropertyObject)
	items := po.Properties().Items()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "foo", items[0].Name())
	assert.Equal(t, "sub", items[1].Name())

	po1 := items[1].Property().(PropertyObject)
	assert.Equal(t, po1.Properties().Items()[0].Name(), "sub")
	assert.Equal(t, po1.Properties().Items()[0].Property().(PropertyString).Type(), "string")
	assert.Equal(t, po1.FileName(), "/abs/sub.schema.json")
	po2 := po1.Properties().Items()[1].Property().(PropertyObject)
	assert.Equal(t, po2.FileName(), "/abs/sub2.schema.json")

}
