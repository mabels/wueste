package entity_generator

import (
	"fmt"

	"github.com/iancoleman/orderedmap"
)

// type JSONProperty map[string]interface{}
type JSONProperty interface {
	// orderedmap.OrderedMap
	Set(key string, value any)
	Get(key string) any
	Keys() []string
	Len() int
	Lookup(key string) (any, bool)
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)
}

type jsonProperty struct {
	omap orderedmap.OrderedMap
}

// Len implements JSONProperty.
func (j *jsonProperty) Len() int {
	return len(j.omap.Keys())
}

// Get implements JSONProperty.
func (j *jsonProperty) Get(key string) any {
	val, found := j.Lookup(key)
	if !found {
		panic(fmt.Sprintf("key[%s] not found", key))
	}
	return val
}

// Lookup implements JSONProperty.
func (j *jsonProperty) Lookup(key string) (any, bool) {
	out, found := j.omap.Get(key)
	if !found {
		return nil, false
	}
	isOmap, found := out.(orderedmap.OrderedMap)
	if found {
		return &jsonProperty{
			omap: isOmap,
		}, true
	}
	return out, true

}

// MarshalJSON implements JSONProperty.
func (j *jsonProperty) MarshalJSON() ([]byte, error) {
	return j.omap.MarshalJSON()
}

// Set implements JSONProperty.
func (j *jsonProperty) Set(key string, value any) {
	j.omap.Set(key, value)
}

// UnmarshalJSON implements JSONProperty.
func (j *jsonProperty) UnmarshalJSON(b []byte) error {
	return j.omap.UnmarshalJSON(b)
}

func (j *jsonProperty) Keys() []string {
	return j.omap.Keys()
}

func NewJSONProperty() JSONProperty {
	return &jsonProperty{
		omap: orderedmap.New(),
	}
}
