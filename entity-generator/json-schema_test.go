package entity_generator

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type JSONPropertyItems struct {
	Name string
	Prop JSONProperty
}

func toSorted(m map[string]JSONProperty) []JSONPropertyItems {
	out := []JSONPropertyItems{}
	for k, v := range m {
		out = append(out, JSONPropertyItems{
			Name: k,
			Prop: v,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func TestJsonSchema2Property(t *testing.T) {
	obj := TestJsonFlatSchema()
	sb := NewSchemaBuilder(NewTestSchemaLoader())
	jsprop := sb.JSON2PropertyObject(obj).Build().(PropertyObject)
	prop := TestFlatSchema(NewTestSchemaLoader()).(PropertyObject)
	pjs := JSONFromProperty(prop)
	pjsProps := toSorted(pjs.Properties)
	pjs.Properties = nil
	jsp := JSONFromProperty(jsprop)
	jspProps := toSorted(jsp.Properties)
	jsp.Properties = nil
	assert.Equal(t, pjs, jsp)
	assert.Equal(t, len(pjsProps), len(jspProps))
	for i, _ := range pjsProps {
		assert.Equal(t, pjsProps[i], jspProps[i], "Property %d:%s", i, pjsProps[i].Name)
	}
}
