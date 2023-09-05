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
		return items[i].PropertItem().Name() < items[j].PropertItem().Name()
	})
	abs := func(path string) string {
		p, _ := tsl.Abs(path)
		return p
	}
	refItems := []SchemaRegistryItem{
		&schemaRegistryItem{
			written: false,
			prop: NewPropertyItem("http://example.com/base.schema.json",
				NewPropertiesBuilder(tsl).BuildObject().FromJson(BaseSchema).
					fileName(abs(BaseSchema["fileName"].(string))).Build()),
		},
		&schemaRegistryItem{
			written: false,
			prop: NewPropertyItem("http://example.com/sub.schema.json",
				NewPropertiesBuilder(tsl).BuildObject().FromJson(SubSchema).
					fileName(abs(SubSchema["fileName"].(string))).Build()),
		},
		&schemaRegistryItem{
			written: false,
			prop: NewPropertyItem("http://example.com/sub2.schema.json",
				NewPropertiesBuilder(tsl).BuildObject().FromJson(Sub2Schema).
					fileName(abs(Sub2Schema["fileName"].(string))).Build()),
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
	items := po.Items()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "foo", items[0].Name())
	assert.Equal(t, "sub", items[1].Name())

	po1 := items[1].Property().(PropertyObject)
	assert.Equal(t, po1.Items()[0].Name(), "sub")
	assert.Equal(t, po1.Items()[0].Property().(PropertyString).Type(), "string")
	assert.Equal(t, po1.FileName(), "/abs/sub.schema.json")
	po2 := po1.Items()[1].Property().(PropertyObject)
	assert.Equal(t, po2.FileName(), "/abs/sub2.schema.json")

}
