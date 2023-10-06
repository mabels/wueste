package entity_generator

import (
	"fmt"

	"github.com/iancoleman/orderedmap"
)

// type JSONDict map[string]interface{}
type JSONDict interface {
	// orderedmap.OrderedMap
	Set(key string, value any)
	Get(key string) any
	Keys() []string
	Len() int
	Lookup(key string) (any, bool)
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)
}

type jsonDict struct {
	omap orderedmap.OrderedMap
}

// Len implements JSONProperty.
func (j *jsonDict) Len() int {
	return len(j.omap.Keys())
}

// Get implements JSONProperty.
func (j *jsonDict) Get(key string) any {
	val, found := j.Lookup(key)
	if !found {
		panic(fmt.Sprintf("key[%s] not found", key))
	}
	return val
}

// Lookup implements JSONProperty.
func (j *jsonDict) Lookup(key string) (any, bool) {
	out, found := j.omap.Get(key)
	if !found {
		return nil, false
	}
	isOmap, found := out.(orderedmap.OrderedMap)
	if found {
		return &jsonDict{
			omap: isOmap,
		}, true
	}
	return out, true

}

// MarshalJSON implements JSONProperty.
func (j *jsonDict) MarshalJSON() ([]byte, error) {
	return j.omap.MarshalJSON()
}

// Set implements JSONProperty.
func (j *jsonDict) Set(key string, value any) {
	j.omap.Set(key, value)
}

// UnmarshalJSON implements JSONProperty.
func (j *jsonDict) UnmarshalJSON(b []byte) error {
	return j.omap.UnmarshalJSON(b)
}

func (j *jsonDict) Keys() []string {
	return j.omap.Keys()
}

func NewJSONDict() JSONDict {
	return &jsonDict{
		omap: orderedmap.New(),
	}
}
