package entity_generator

import (
	"encoding/json"
	"testing"

	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/assert"
)

type Item struct {
	Key string
	Val int
}

func TestOrderedMap(t *testing.T) {
	const input = `{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}`
	for i := 0; i < 100; i++ {
		out := NewJSONDict()
		err := json.Unmarshal([]byte(input), &out)
		assert.NoError(t, err)
		testOrder(t, out)
		asj, err := json.Marshal(out)
		assert.NoError(t, err)
		assert.Equal(t, input, string(asj))

	}
}

func testOrder(t *testing.T, out JSONDict) {
	order := []Item{}
	for _, k := range out.Keys() {
		v := out.Get(k)
		if k == "nested" {
			testOrder(t, v.(JSONDict))
			continue
		}
		order = append(order, Item{
			Key: k,
			Val: int(v.(float64)),
		})
	}
	for i, o := range order {
		assert.Equal(t, i+1, o.Val)
	}
}

func TestNestedOrder(t *testing.T) {
	const input = `{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6,"nested":{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}}`
	for i := 0; i < 100; i++ {
		out := orderedmap.New()
		err := json.Unmarshal([]byte(input), &out)
		assert.NoError(t, err)

		asj, err := json.Marshal(out)
		assert.NoError(t, err)
		assert.Equal(t, input, string(asj))

	}
}

func TestTypeUnmarshal(t *testing.T) {
	const input = `{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6,"nested":{"a":1,"b":2,"c":3,"d":4,"e":5,"f":6}}`
	out := NewJSONDict()
	err := json.Unmarshal([]byte(input), &out)
	assert.NoError(t, err)
	_, ok := out.Get("nested").(JSONDict)
	assert.True(t, ok)

}

func TestSerialize(t *testing.T) {
	v := TestJsonSubSchema()
	js, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, string(js), "{\"filename\":\"sub.schema.json\",\"jsonProperty\":{\"$id\":\"http://example.com/sub.schema.json\",\"$schema\":\"http://json-schema.org/draft-07/schema#\",\"title\":\"Sub\",\"type\":\"object\",\"description\":\"Sub description\",\"properties\":{\"sub\":{\"type\":\"string\"},\"sub-down\":{\"$ref\":\"file://wurst/sub2.schema.json\"}},\"required\":[\"bar\",\"sub-down\"]}}")
}

func TestJsonPropertySetMutable(t *testing.T) {
	const input = `{
		"foo": "bar",
		"prop": {
			"v1": 1,
			"v2": "v1"
		}
	}`
	out := NewJSONDict()
	err := json.Unmarshal([]byte(input), &out)
	assert.NoError(t, err)

	out.Get("prop").(JSONDict).Set("v3", "blabla")

	assert.Equal(t, "foo", out.Keys()[0])
	assert.Equal(t, "bar", out.Get("foo"))
	assert.Equal(t, "prop", out.Keys()[1])
	inProp := out.Get("prop").(JSONDict)
	assert.Equal(t, "v1", inProp.Keys()[0])
	assert.Equal(t, "v2", inProp.Keys()[1])
	assert.Equal(t, "v3", inProp.Keys()[2])

}

func TestOrderedMapSetMutable(t *testing.T) {
	const input = `{
		"foo": "bar",
		"prop": {
			"v1": 1,
			"v2": "v1"
		}
	}`
	out := orderedmap.New()
	err := json.Unmarshal([]byte(input), &out)
	assert.NoError(t, err)

	iout, _ := out.Get("prop")
	iout.(orderedmap.OrderedMap).Set("v3", "blabla")

	assert.Equal(t, "foo", out.Keys()[0])
	// assert.Equal(t, "bar", out.Get("foo"))
	assert.Equal(t, "prop", out.Keys()[1])
	inPropX, _ := out.Get("prop")
	inProp := inPropX.(orderedmap.OrderedMap)

	assert.Equal(t, "v1", inProp.Keys()[0])
	assert.Equal(t, "v2", inProp.Keys()[1])
	assert.Equal(t, "v3", inProp.Keys()[2])

}
