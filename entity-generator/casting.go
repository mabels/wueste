package entity_generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mabels/wueste/entity-generator/rusty"
)

func interfaceToStringPtr(i interface{}) *string {
	if i == nil {
		return nil
	}
	s, found := i.(*string)
	if !found {
		return nil
	}
	return s
}

func interfaceToBoolPtr(i interface{}) *bool {
	if i == nil {
		return nil
	}
	b, found := i.(*bool)
	if !found {
		return nil
	}
	return b
}

func interfaceToIntPtr(i interface{}) *int {
	if i == nil {
		return nil
	}
	b, found := i.(*int)
	if !found {
		return nil
	}
	return b
}

func interfaceToInt8Ptr(i interface{}) *int8 {
	if i == nil {
		return nil
	}
	b, found := i.(*int8)
	if !found {
		return nil
	}
	return b
}

func interfaceToInt16Ptr(i interface{}) *int16 {
	if i == nil {
		return nil
	}
	b, found := i.(*int16)
	if !found {
		return nil
	}
	return b
}

func interfaceToInt32Ptr(i interface{}) *int32 {
	if i == nil {
		return nil
	}
	b, found := i.(*int32)
	if !found {
		return nil
	}
	return b
}

func interfaceToInt64Ptr(i interface{}) *int64 {
	if i == nil {
		return nil
	}
	b, found := i.(*int64)
	if !found {
		return nil
	}
	return b
}

func interfaceToFloat32Ptr(i interface{}) *float32 {
	if i == nil {
		return nil
	}
	b, found := i.(*float32)
	if !found {
		return nil
	}
	return b
}

func interfaceToFloat64Ptr(i interface{}) *float64 {
	if i == nil {
		return nil
	}
	b, found := i.(*float64)
	if !found {
		return nil
	}
	return b
}

func coerceString(v interface{}) rusty.Optional[string] {
	if v != nil {
		switch v.(type) {
		case string:
			return rusty.Some[string](v.(string))
		case bool:
			val := v.(bool)
			if val {
				return rusty.Some[string]("true")
			} else {
				return rusty.Some[string]("false")
			}
		case int, int8, int16, int32, int64, float32, float64:
			return rusty.Some[string](fmt.Sprintf("%v", v))
		default:
			panic(fmt.Errorf("unknown type %T", v))
		}
	}
	return rusty.None[string]()

}
func coerceBool(v interface{}) rusty.Optional[bool] {
	if v != nil {
		switch v.(type) {
		case string:
			val := strings.ToLower(v.(string))
			if val == "true" || val == "on" || val == "yes" {
				return rusty.Some[bool](true)
			}
			return rusty.Some[bool](false)
		case bool:
			val := v.(bool)
			return rusty.Some[bool](val)
		case int, int8, int16, int32, int64, float32, float64:
			val := coerceInt(v)
			if val.IsNone() {
				return rusty.None[bool]()
			}
			return rusty.Some[bool](*val.Value() != 0)
		default:
			panic(fmt.Errorf("unknown type %T", v))
		}

	}
	return rusty.None[bool]()
}

func coerceNumber[T int | int8 | int16 | int32 | int64 | float32 | float64](v interface{}) rusty.Optional[T] {
	if v != nil {
		switch v.(type) {
		case string:
			f, err := strconv.ParseFloat(v.(string), 8)
			if err != nil {
				return rusty.None[T]()
			}
			return rusty.Some[T](T(f))
		case int:
			return rusty.Some[T](T(v.(int)))
		case int8:
			return rusty.Some[T](T(v.(int8)))
		case int16:
			return rusty.Some[T](T(v.(int16)))
		case int32:
			return rusty.Some[T](T(v.(int32)))
		case int64:
			return rusty.Some[T](T(v.(int64)))
		case float32:
			return rusty.Some[T](T(v.(float32)))
		case float64:
			return rusty.Some[T](T(v.(float64)))
		default:
			panic(fmt.Errorf("unknown type %T", v))
		}
	}
	return rusty.None[T]()
}
func coerceInt(v interface{}) rusty.Optional[int] {
	return coerceNumber[int](v)
}
func coerceInt8(v interface{}) rusty.Optional[int8] {
	return coerceNumber[int8](v)
}
func coerceInt16(v interface{}) rusty.Optional[int16] {
	return coerceNumber[int16](v)
}
func coerceInt32(v interface{}) rusty.Optional[int32] {
	return coerceNumber[int32](v)
}
func coerceInt64(v interface{}) rusty.Optional[int64] {
	return coerceNumber[int64](v)
}
func coerceFloat32(v interface{}) rusty.Optional[float32] {
	return coerceNumber[float32](v)
}
func coerceFloat64(v interface{}) rusty.Optional[float64] {
	return coerceNumber[float64](v)
}

func getFromAttributeString(js JSONProperty, attr string) string {
	defVal, found := js[attr]
	if found {
		return *coerceString(defVal).Value()
	}
	return ""
	// panic("no " + attr + " found")
}

func getFromAttributeOptionalString(js JSONProperty, attr string) rusty.Optional[string] {
	format := rusty.None[string]()
	formatVal, found := js[attr].(string)
	if found {
		format = coerceString(formatVal)
	}
	return format
}

func getFromAttributeOptionalBoolean(js JSONProperty, attr string) rusty.Optional[bool] {
	format := rusty.None[bool]()
	formatVal, found := js[attr].(bool)
	if found {
		format = coerceBool(formatVal)
	}
	return format
}

func getFromAttributeOptionalInt(js JSONProperty, attr string) rusty.Optional[int] {
	format := rusty.None[int]()
	formatVal, found := js[attr].(int)
	if found {
		format = coerceInt(formatVal)
	}
	return format
}

func getFromAttributeOptionalFloat32(js JSONProperty, attr string) rusty.Optional[float32] {
	format := rusty.None[float32]()
	formatVal, found := js[attr].(float32)
	if found {
		format = coerceFloat32(formatVal)
	}
	return format
}
func getFromAttributeOptionalFloat64(js JSONProperty, attr string) rusty.Optional[float64] {
	format := rusty.None[float64]()
	formatVal, found := js[attr].(float64)
	if found {
		format = coerceFloat64(formatVal)
	}
	return format
}

func getFromAttributeOptionalInt8(js JSONProperty, attr string) rusty.Optional[int8] {
	format := rusty.None[int8]()
	formatVal, found := js[attr].(int8)
	if found {
		format = coerceInt8(formatVal)
	}
	return format
}
func getFromAttributeOptionalInt16(js JSONProperty, attr string) rusty.Optional[int16] {
	format := rusty.None[int16]()
	formatVal, found := js[attr].(int16)
	if found {
		format = coerceInt16(formatVal)
	}
	return format
}
func getFromAttributeOptionalInt32(js JSONProperty, attr string) rusty.Optional[int32] {
	format := rusty.None[int32]()
	formatVal, found := js[attr].(int32)
	if found {
		format = coerceInt32(formatVal)
	}
	return format
}
func getFromAttributeOptionalInt64(js JSONProperty, attr string) rusty.Optional[int64] {
	format := rusty.None[int64]()
	formatVal, found := js[attr].(int64)
	if found {
		format = coerceInt64(formatVal)
	}
	return format
}
