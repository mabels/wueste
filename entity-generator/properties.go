package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertiesBuilder struct {
	// items   []PropertyItem
	property Property
	_ctx     PropertyCtx
	// _loader  SchemaLoader
}

func NewPropertiesBuilder(run PropertyCtx) *PropertiesBuilder {
	return &PropertiesBuilder{
		_ctx: run,
	}
}

func (b *PropertiesBuilder) FileName(fname string) *PropertiesBuilder {
	b.property.Runtime().SetFileName(fname)
	return b
}

func (b *PropertiesBuilder) Resolve(rt PropertyRuntime, prop Property) rusty.Result[Property] {
	if prop.Ref().IsNone() {
		return rusty.Ok(prop)
	}
	refVal := prop.Ref().Value()

	return b._ctx.Registry.EnsureSchema(refVal, rt, func(fname string, rt PropertyRuntime) rusty.Result[Property] {
		return loadSchema(fname, b._ctx, func(abs string, prop JSONProperty) rusty.Result[Property] {
			return rusty.Ok(NewPropertiesBuilder(b._ctx).FromJson(rt, prop).FileName(abs).Build())
		})
	})
}

func (b *PropertiesBuilder) BuildObject() *PropertyObjectParam {
	return &PropertyObjectParam{
		// Runtime: b._runtime,
		Ctx:  b._ctx,
		Type: OBJECT,
	}
}
func (b *PropertiesBuilder) BuildArray() *PropertyArrayParam {
	return &PropertyArrayParam{
		// Runtime: b._runtime,
		Ctx:  b._ctx,
		Type: ARRAY,
	}
}

func (b *PropertiesBuilder) BuildString() *PropertyStringParam {
	return &PropertyStringParam{
		// Runtime: b._runtime,
		Ctx:  b._ctx,
		Type: STRING,
	}
}

func (b *PropertiesBuilder) BuildBoolean() *PropertyBooleanParam {
	return &PropertyBooleanParam{
		// Runtime: b._runtime,
		Ctx:  b._ctx,
		Type: BOOLEAN,
	}
}

func (b *PropertiesBuilder) BuildInteger() *PropertyIntegerParam {
	return &PropertyIntegerParam{
		// Runtime: b._runtime,
		Ctx:  b._ctx,
		Type: INTEGER,
	}
}

func (b *PropertiesBuilder) BuildNumber() *PropertyNumberParam {
	return &PropertyNumberParam{
		// Runtime: b._runtime,
		Ctx:  b._ctx,
		Type: NUMBER,
	}
}

func isOptional(name string, req []string) bool {
	for _, r := range req {
		if r == name {
			return false
		}
	}
	return true
}

func JSONsetString(js JSONProperty, key string, value string) {
	js.Set(key, value)
}

func JSONsetId(jsp JSONProperty, p Property) {
	if p.Id() != "" {
		jsp.Set("$id", p.Id())
	}
}

func JSONsetOptionalString(js JSONProperty, key string, value rusty.Optional[string]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func JSONsetOptionalBoolean(js JSONProperty, key string, value rusty.Optional[bool]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func JSONsetOptionalFloat64(js JSONProperty, key string, value rusty.Optional[float64]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func JSONsetOptionalInt(js JSONProperty, key string, value rusty.Optional[int]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func (b *PropertiesBuilder) FromJson(rt PropertyRuntime, js JSONProperty) *PropertiesBuilder {
	ref, found := js.Lookup("$ref")
	if found {
		// b.property = b.Resolve(NewProperty(PropertyParam{Ref: rusty.Some(ref)}))
		res := b.Resolve(rt, NewProperty(PropertyParam{Ref: coerceString(ref)}))
		if res.IsErr() {
			panic(res.Err())
		}
		b.property = res.Ok()
		return b
	}
	_typ, found := js.Lookup("type")
	if !found {
		panic("no type found")
	}
	typ := coerceString(_typ).Value()
	switch typ {
	case OBJECT:
		b.property = b.BuildObject().FromJson(rt, js).Build()
	case STRING:
		b.property = b.BuildString().FromJson(rt, js).Build()
	case NUMBER:
		b.property = b.BuildNumber().FromJson(rt, js).Build()
	case INTEGER:
		b.property = b.BuildInteger().FromJson(rt, js).Build()
	case BOOLEAN:
		b.property = b.BuildBoolean().FromJson(rt, js).Build()
	case ARRAY:
		b.property = b.BuildArray().FromJson(rt, js).Build()
	default:
		panic("unknown type:" + typ)
	}
	return b
}

func PropertyToJson(iprop Property) JSONProperty {
	switch prop := iprop.(type) {
	case PropertyString:
		return PropertyStringToJson(prop)
	case PropertyArray:
		return PropertyArrayToJson(prop)
	case PropertyBoolean:
		return PropertyBooleanToJson(prop)
	case PropertyInteger:
		return PropertyIntegerToJson(prop)
	case PropertyNumber:
		return PropertyNumberToJson(prop)
	case PropertyObject:
		return PropertyObjectToJson(prop)
	default:
		panic("unknown type: " + prop.(Property).Type())
	}

}

func (b *PropertiesBuilder) Build() Property {
	return b.property
}

// func (b *PropertiesBuilder) X(js JSONProperty, req []string) *PropertiesBuilder {
// 	for name, p := range js {
// 		// falsch !!!
// 		var pn Property
// 		switch p.Type {
// 		case "string":
// 			pn = NewPropertyString(PropertyStringParam{
// 				PropertyParam: PropertyParam{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.OptionalFromPtr(p.Format),
// 					Optional:    isOptional(name, req),
// 				},
// 				Default: rusty.OptionalFromPtr(interfaceToStringPtr(p.Default)),
// 				Format:  rusty.OptionalFromPtr(p.Format),
// 			})
// 		case "boolean":
// 			pn = NewPropertyBoolean(PropertyBooleanParam{
// 				Id:          p.Id,
// 				Type:        p.Type,
// 				Description: rusty.OptionalFromPtr(p.Description),
// 				Optional:    isOptional(name, req),
// 				Default:     rusty.OptionalFromPtr(p.Default.(*bool)),
// 			})
// 		case "integer":
// 			var format string
// 			if p.Format == nil {
// 				format = "int"
// 			} else {
// 				format = *p.Format
// 			}
// 			switch format {
// 			case "int":
// 				pn = NewPropertyInteger[int](PropertyIntegerParam[int]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("int"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*int)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*int)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*int)),
// 				})
// 			case "int8":
// 				pn = NewPropertyInteger[int8](PropertyIntegerParam[int8]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("int8"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*int8)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*int8)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*int8)),
// 				})
// 			case "int16":
// 				pn = NewPropertyInteger[int16](PropertyIntegerParam[int16]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("float32"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*int16)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*int16)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*int16)),
// 				})
// 			case "int32":
// 				pn = NewPropertyInteger[int32](PropertyIntegerParam[int32]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("float32"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*int32)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*int32)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*int32)),
// 				})
// 			case "int64":
// 				pn = NewPropertyInteger[int64](PropertyIntegerParam[int64]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("float32"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*int64)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*int64)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*int64)),
// 				})
// 			default:
// 				panic("unknown format")
// 			}
// 		case "number":
// 			var format string
// 			if p.Format == nil {
// 				format = "float64"
// 			} else {
// 				format = *p.Format
// 			}
// 			switch format {
// 			case "float32":
// 				pn = NewPropertyNumber[float32](PropertyNumberParam[float32]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("float32"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*float32)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*float32)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*float32)),
// 				})
// 			case "float64":
// 				pn = NewPropertyNumber[float64](PropertyNumberParam[float64]{
// 					Id:          p.Id,
// 					Type:        p.Type,
// 					Description: rusty.OptionalFromPtr(p.Description),
// 					Format:      rusty.Some("float32"),
// 					Optional:    isOptional(name, req),
// 					Default:     rusty.OptionalFromPtr(p.Default.(*float64)),
// 					Maximum:     rusty.OptionalFromPtr(p.Maximum.(*float64)),
// 					Minimum:     rusty.OptionalFromPtr(p.Minimum.(*float64)),
// 				})
// 			default:
// 				panic("unknown format")
// 			}
// 		case "object":
// 			pn = NewPropertyObject(PropertyObjectParam{
// 				FileName:    p.FileName,
// 				Id:          p.Id,
// 				Title:       p.Title,
// 				Schema:      p.Schema,
// 				Description: rusty.OptionalFromPtr(p.Description),
// 				Properties:  NewPropertiesBuilder(b._loader).FromJson(p.Properties, p.Required).Build(),
// 				Required:    p.Required,
// 				Ref:         rusty.OptionalFromPtr(p.Ref),
// 			})
// 		default:
// 			panic("unknown type")
// 		}
// 		b.Add(NewPropertyItem(name, pn))
// 	}
// 	return b
// }

// func (b *PropertiesBuilder) Build() PropertiesObject {
// 	return b
// }
