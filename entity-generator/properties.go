package entity_generator

import (
	"fmt"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertiesBuilder struct {
	parentFileName rusty.Optional[string]
	property       rusty.Optional[Property]
	filename       rusty.Optional[string]
	// ref      rusty.Optional[string]
	errors []error
	ctx    PropertyCtx
}

func NewPropertiesBuilder(run PropertyCtx) *PropertiesBuilder {
	return &PropertiesBuilder{
		ctx: run,
		// meta: NewPropertyMeta(),
	}
}

func (b *PropertiesBuilder) FileName() rusty.Optional[string] {
	if b.filename.IsSome() {
		return b.filename
	}
	return b.parentFileName
}

// func (b *PropertiesBuilder) SetFileName(fname string) *PropertiesBuilder {
// 	b.filename = rusty.Some(fname)
// 	return b
// }

// func (b *PropertiesBuilder) SetFileName(fname string) *PropertiesBuilder {
// 	if b.property.IsNone() {
// 		b.errors = append(b.errors, fmt.Errorf("FileName Set with no property set"))
// 		return b
// 	}
// 	if b.property.Value().Meta().FileName().IsSome() {
// 		b.errors = append(b.errors, fmt.Errorf("double FileName Set"))
// 		return b
// 	}
// 	b.property.Value().Meta().SetFileName(fname)
// 	return b
// }

func isOptional(name string, req []string) bool {
	for _, r := range req {
		if r == name {
			return false
		}
	}
	return true
}

func JSONsetString(js JSONDict, key string, value string) {
	js.Set(key, value)
}

func JSONsetId(jsp JSONDict, p Property) {
	if p.Id() != "" {
		jsp.Set("$id", p.Id())
	}
}

func JSONsetOptionalString(js JSONDict, key string, value rusty.Optional[string]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func JSONsetXProperties(js JSONDict, value map[string]interface{}) {
	for k, v := range value {
		js.Set(k, v)
	}
}

func JSONsetOptionalBoolean(js JSONDict, key string, value rusty.Optional[bool]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func JSONsetOptionalFloat64(js JSONDict, key string, value rusty.Optional[float64]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

func JSONsetOptionalInt(js JSONDict, key string, value rusty.Optional[int]) {
	if !value.IsNone() {
		js.Set(key, value.Value())
	}
}

// func (b *PropertiesBuilder) FromProperty(prop Property, optParent ...PropertyMeta) *PropertiesBuilder {
// 	propMeta := NewPropertyMeta()
// 	if len(optParent) > 0 {
// 		propMeta = optParent[0]
// 	}
// 	rProp := b.Resolve(propMeta, prop)
// 	if rProp.IsErr() {
// 		b.errors = append(b.errors, rProp.Err())
// 		return b
// 	}
// 	switch prop.Type() {
// 	case OBJECT:
// 		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
// 			return NewPropertyObjectBuilder(b).FromProperty(rProp.Ok()).Build()
// 		})
// 	case STRING:
// 		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
// 			return NewPropertyStringBuilder(b).FromProperty(rProp.Ok()).Build()
// 		})
// 	case NUMBER:
// 		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
// 			return NewPropertyNumberBuilder(b).FromProperty(rProp.Ok()).Build()
// 		})
// 	case INTEGER:
// 		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
// 			return NewPropertyIntegerBuilder(b).FromProperty(rProp.Ok()).Build()
// 		})
// 	case BOOLEAN:
// 		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
// 			return NewPropertyBooleanBuilder(b).FromProperty(rProp.Ok()).Build()
// 		})
// 	case ARRAY:
// 		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
// 			return NewPropertyArrayBuilder(b).FromProperty(rProp.Ok()).Build()
// 		})
// 	default:
// 		b.errors = append(b.errors, fmt.Errorf("unknown type: %s", prop.Type()))
// 	}
// 	return b
// }

func (b *PropertiesBuilder) MergeJson(parentFname rusty.Optional[string], ref string, js JSONDict) rusty.Result[JSonFile] {
	// refVal := b.ref.Value()
	rJson := b.ctx.Registry.EnsureJSONProperty(parentFname, ref)
	if rJson.IsErr() {
		return rusty.Err[JSonFile](rJson.Err())
	}
	fjs := rJson.Ok().JSONProperty
	for _, k := range fjs.Keys() {
		v, _ := fjs.Lookup(k)
		js.Set(k, v)
	}
	return rusty.Ok[JSonFile](JSonFile{
		FileName:     rJson.Ok().FileName,
		JSONProperty: js,
	})
}

func (b *PropertiesBuilder) FromJson(js JSONDict) *PropertiesBuilder {
	ref, found := js.Lookup("$ref")
	if found {
		// if b.property.IsSome() {
		// 	b.fixlename = b.property.Value().Meta().FileName()
		// }
		refStr := coerceString(ref)
		if refStr.IsNone() {
			b.errors = append(b.errors, fmt.Errorf("ref not a string"))
			return b
		}
		rJs := b.MergeJson(b.parentFileName, refStr.Value(), js)
		if rJs.IsErr() {
			b.errors = append(b.errors, rJs.Err())
			return b
		}
		js = rJs.Ok().JSONProperty
		b.filename = rusty.Some(rJs.Ok().FileName)
	}
	_typ, found := js.Lookup("type")
	if !found {
		b.errors = append(b.errors, fmt.Errorf("no type"))
		return b
	}
	typ := coerceString(_typ).Value()
	switch typ {
	case OBJECT:
		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
			return NewPropertyObjectBuilder(b).FromJson(js).Build()
		})
	case STRING:
		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
			return NewPropertyStringBuilder(b).FromJson(js).Build()
		})
	case NUMBER:
		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
			return NewPropertyNumberBuilder(b).FromJson(js).Build()
		})
	case INTEGER:
		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
			return NewPropertyIntegerBuilder(b).FromJson(js).Build()
		})
	case BOOLEAN:
		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
			return NewPropertyBooleanBuilder(b).FromJson(js).Build()
		})
	case ARRAY:
		b.assignProperty(func(b *PropertiesBuilder) rusty.Result[Property] {
			return NewPropertyArrayBuilder(b).FromJson(js).Build()
		})
	default:
		panic("unknown type:" + typ)
	}
	return b
}

func (b *PropertiesBuilder) assignProperty(fn func(b *PropertiesBuilder) rusty.Result[Property]) *PropertiesBuilder {
	p := fn(b)
	if p.IsErr() {
		b.errors = append(b.errors, p.Err())
	} else {
		if b.property.IsSome() {
			b.errors = append(b.errors, fmt.Errorf("property already set"))
			return b
		}
		b.property = rusty.Some(p.Ok())
	}
	return b
}

func PropertyToJson(iprop Property) JSONDict {
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
		panic("PropertyToJson unknown type: " + prop.Type())
	}
}

// func (b *PropertiesBuilder) Resolve(meta PropertyMeta, prop Property) rusty.Result[Property] {
// 	if prop.Ref().IsSome() && prop.Meta().FileName().IsSome() {
// 		return rusty.Ok(prop)
// 	}
// 	refVal := prop.Ref().Value()
// 	return b.ctx.Registry.EnsureSchema(refVal, meta.FileName(), func(fname string) rusty.Result[Property] {
// 		return loadSchema(fname, b.ctx, func(abs string, prop JSONProperty) rusty.Result[Property] {
// 			return NewPropertiesBuilder(b.ctx).FromJson(prop).SetFileName(abs).Build()
// 		})
// 	})
// }

func (b *PropertiesBuilder) Build() rusty.Result[Property] {
	if len(b.errors) > 0 {
		str := ""
		for _, v := range b.errors {
			str += v.Error() + "\n"
		}
		return rusty.Err[Property](fmt.Errorf(str))
	}

	if b.property.IsNone() {
		b.errors = append(b.errors, fmt.Errorf("no property set"))
		return rusty.Err[Property](fmt.Errorf("no property set"))
	}
	if b.filename.IsSome() {
		b.property.Value().Meta().SetFileName(b.filename.Value())
	}

	return rusty.Ok[Property](b.property.Value())
}
