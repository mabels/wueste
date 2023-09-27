package ts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/google/uuid"
	eg "github.com/mabels/wueste/entity-generator"
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/mabels/wueste/entity-generator/wueste"
)

var keyWords = map[string]bool{
	"break":      true,
	"as":         true,
	"any":        true,
	"switch":     true,
	"case":       true,
	"if":         true,
	"throw":      true,
	"else":       true,
	"var":        true,
	"number":     true,
	"string":     true,
	"get":        true,
	"module":     true,
	"type":       true,
	"instanceof": true,
	"typeof":     true,
	"public":     true,
	"private":    true,
	"enum":       true,
	"export":     true,
	"finally":    true,
	"for":        true,
	"while":      true,
	"void":       true,
	"null":       true,
	"super":      true,
	"this":       true,
	"new":        true,
	"in":         true,
	"return":     true,
	"true":       true,
	"false":      true,
	"extends":    true,
	"static":     true,
	"let":        true,
	"package":    true,
	"implements": true,
	"interface":  true,
	"function":   true,
	"try":        true,
	"yield":      true,
	"const":      true,
	"continue":   true,
	"do":         true,
	"catch":      true,
}

type tsLang struct {
}

var reSplitNonAllowed = regexp.MustCompile("[^a-zA-Z0-9_$]+")

var reReplaceCaps = regexp.MustCompile(`[A-Z]+`)
var reReplaceNoAlpha = regexp.MustCompile(`[^a-zA-Z0-9]+`)
var reTrimNoAlpha = regexp.MustCompile(`^[^a-zA-Z0-9]+`)

func (l *tsLang) FileName(fname string) string {
	fname = reReplaceCaps.ReplaceAllString(fname, "_$0")
	fname = reTrimNoAlpha.ReplaceAllString(fname, "")
	fname = reReplaceNoAlpha.ReplaceAllString(fname, "_")
	return "./" + strings.ToLower(fname) + ".ts"
}

func (l *tsLang) keyWordFilter(name string) string {
	_, ok := keyWords[name]
	if ok {
		return l.Quote(name)
	}
	return name
}

func (l *tsLang) RemoveFileExt(fname string) string {
	return strings.TrimSuffix(fname, ".ts")
}

func (l *tsLang) PublicName(name string, opts ...string) string {
	re := reSplitNonAllowed.ReplaceAllString(name+strings.Join(opts, ""), "_")
	return l.keyWordFilter(strings.TrimLeft(re, "_"))
}

func (l *tsLang) PrivateName(name string, opts ...string) string {
	re := reSplitNonAllowed.ReplaceAllString(name+strings.Join(opts, ""), "_")
	return "_" + re
}

func (l *tsLang) Name(name string, opts ...string) string {
	re := reSplitNonAllowed.ReplaceAllString(name+strings.Join(opts, ""), "_")
	return l.keyWordFilter(re)
}

func (l *tsLang) Type(name string, opt bool) string {
	optStr := ""
	if opt {
		optStr = "?"
	}
	return name + optStr
}

func (l *tsLang) addCoerceType(typ string, withs ...string) string {
	if hasWith(WithAddCoerce(), withs) {
		switch typ {
		case "number", "string", "boolean", "Date":
			return fmt.Sprintf("WuesteCoerceType%s", typ)
		}
	}
	return typ
}

func (l *tsLang) AsTypeHelper(p eg.Property, withs ...string) string {
	switch p.Type() {
	case eg.OBJECT:
		po := p.(eg.PropertyObject)
		if po.Properties() == nil || po.Properties().Len() == 0 {
			return l.Generics("Record", "string", "unknown")
		}
		name := getObjectName(p)
		if hasWith(WithAddCoerce(), withs) {
			return l.OrType(l.Name(name), l.Name(name, "Param"))
		}
		return l.Name(name)
	case eg.STRING:
		p := p.(eg.PropertyString)
		if p.Format().IsSome() {
			switch p.Format().Value() {
			case eg.DATE_TIME:
				return l.addCoerceType("Date", withs...)
			case eg.DATE, eg.TIME:
				panic("not implemented")
			default:
				panic("not implemented")
			}
		}
		return l.addCoerceType("string", withs...)
	case eg.NUMBER:
		return l.addCoerceType("number", withs...)
	case eg.INTEGER:
		return l.addCoerceType("number", withs...)
	case eg.BOOLEAN:
		return l.addCoerceType("boolean", withs...)
	case eg.ARRAY:
		return l.AsTypeHelper(p.(eg.PropertyArray).Items(), withs...) + "[]"
	default:
		panic(fmt.Sprintf("unknown type %s", p.Type()))
	}
}

func (l *tsLang) AsType(p eg.Property, withs ...string) string {
	out := l.AsTypeHelper(p, withs...)
	if hasWith("isPartial", withs) {
		out = l.Generics("Partial", out)
	}
	return out
}

func (l *tsLang) AsTypeOptional(p eg.Property) string {
	return "rusty.Optional<" + l.AsType(p) + ">"
}

func hasWith(str string, withs []string) bool {
	for _, with := range withs {
		if str == with {
			return true
		}
	}
	return false
}

func WithOptional(optional bool) string {
	if optional {
		return "isOptional"
	}
	return ""
}

func WithAddCoerce() string {
	return "addCoerceType"
}

func WithPartial() string {
	return "isPartial"
}

func (l *tsLang) AsTypeNullable(p eg.Property, withs ...string) string {
	res := l.AsType(p, withs...)

	if hasWith("isOptional", withs) {
		res = res + "|undefined"
	}
	return res
}

func (l *tsLang) Line(line string, tails ...string) string {
	tail := strings.Join(tails, "")
	return line + ";" + tail
}

func (l *tsLang) Comma(line string, tails ...string) string {
	tail := strings.Join(tails, "")
	return line + "," + tail
}

func (l *tsLang) ReturnType(line string, retType string) string {
	return line + ": " + retType
}

func (l *tsLang) AssignDefault(line string, retType string) string {
	return line + " = " + retType
}

func (l *tsLang) New(line string, retType ...string) string {
	return "new " + l.Call(line, retType...)
}

func (l *tsLang) Return(line string) string {
	return "return " + line
}

func (l *tsLang) Call(line string, params ...string) string {
	return line + "(" + strings.Join(params, ", ") + ")"
}

func (l *tsLang) Interface(wr *eg.ForIfWhileLangWriter, prefix, name string, prop eg.PropertyObject, itemFn func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter)) string {
	return wr.WriteBlock(prefix+"interface ", name, func(wr *eg.ForIfWhileLangWriter) {
		for _, pi := range prop.Items() {
			itemFn(pi, wr)
		}
	})
}

func (l *tsLang) Class(
	wr *eg.ForIfWhileLangWriter,
	prefix, name string,
	prop eg.PropertyObject,
	itemFn func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter),
	onceFn ...func(wr *eg.ForIfWhileLangWriter)) string {
	return wr.WriteBlock(prefix+"class ", name, func(wr *eg.ForIfWhileLangWriter) {
		for _, once := range onceFn {
			once(wr)
		}
		for _, prop := range prop.Items() {
			itemFn(prop, wr)
		}
	})
}

func (l *tsLang) JsonAnnotation(prop eg.Property) string {
	return ""
}

func (l *tsLang) PublicType(str string, suffix ...string) string {
	// remove non alpha numeric characters - _
	return l.PublicName(str + strings.Join(suffix, ""))
}

func (l *tsLang) Quote(str string) string {
	bytes, _ := json.Marshal(str)
	return string(bytes)
}

func (l *tsLang) Readonly(str string) string {
	return "readonly " + str
}

func (l *tsLang) Export(str string) string {
	return "export " + str
}

func (l *tsLang) Implements(str string, interfaces ...string) string {
	return str + " implements " + strings.Join(interfaces, ", ")
}

func (l *tsLang) Extends(str string, interfaces ...string) string {
	return str + " extends " + strings.Join(interfaces, ", ")
}

func (l *tsLang) Generics(str string, typ ...string) string {
	return str + "<" + strings.Join(typ, ", ") + ">"
}

func (l *tsLang) OptionalParam(str string, opt bool) string {
	if opt {
		return str + "?"
	}
	return str
}

func (l *tsLang) CallDot(str string, prop string) string {
	if strings.HasPrefix(prop, "\"") {
		return str + "[" + prop + "]"
	}
	return str + "." + prop
}

func (l *tsLang) OrType(str ...string) string {
	return strings.Join(str, "|")
}

func (l *tsLang) Const(str string) string {
	return "const " + str
}
func (l *tsLang) Let(str string) string {
	return "let " + str
}

func (l *tsLang) RoundBrackets(str string) string {
	return "(" + str + ")"
}

func hasDefault(prop eg.Property) bool {
	return getDefaultForProperty(prop) != nil
}

var reOrArray = regexp.MustCompile(`[(\|)(\[\])]+`)

func (g *tsGenerator) addWuestenType(typ string) string {
	split := reOrArray.Split(typ, -1)
	for _, s := range split {
		if strings.HasPrefix(s, "WuesteCoerceType") {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, s).activated = true
		}
	}
	return typ
}

func isNamedType(p eg.Property) bool {
	switch p.Type() {
	case eg.OBJECT:
		po := p.(eg.PropertyObject)
		if po.Id() == "" || po.Title() == "" {
			return false
		}
		if po.Properties() == nil || po.Properties().Len() == 0 {
			return false
		}
	}
	return true
}

func (g *tsGenerator) generateClass(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop)), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		if isNamedType(pi.Property()) {
			g.includes.AddProperty(getObjectName(pi.Property()), pi.Property())
		}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()), pi.Optional())),
				g.lang.AsTypeNullable(pi.Property()))))
	})
	g.bodyWriter.WriteLine()

	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Param"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		if isNamedType(pi.Property()) {
			paramName := getObjectName(pi.Property())
			g.includes.AddProperty(paramName, pi.Property())
			g.includes.AddProperty(g.lang.PublicName(paramName, "Param"), pi.Property())
		}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()),
					pi.Optional() || hasDefault(pi.Property()))),
				g.addWuestenType(g.lang.AsTypeNullable(pi.Property(), WithAddCoerce())))))
	})
	g.bodyWriter.WriteLine()

}

func (g *tsGenerator) generateJson(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Object"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		if isNamedType(pi.Property()) {
			g.includes.AddProperty(g.lang.AsType(pi.Property()), pi.Property())
		}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.Quote(pi.Name()), pi.Optional())),
				g.lang.AsTypeNullable(pi.Property())), g.lang.JsonAnnotation(pi.Property())))
	})
	g.bodyWriter.WriteLine()
}

func (g *tsGenerator) genToObject(pi eg.PropertyItem) string {
	switch pi.Property().Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		return g.lang.CallDot("obj", g.lang.PublicName(pi.Name()))
	case eg.ARRAY:
		return g.lang.CallDot("obj", g.lang.PublicName(pi.Name())) + "// Choose Factory"
	case eg.OBJECT:
		if pi.Property().(eg.PropertyObject).Title() != "" {
			return g.lang.CallDot(
				g.lang.PublicName(getObjectName(pi.Property()) /*pi.Property().Meta().ToPropertyObject().Ok().Title(), */, "Factory"),
				g.lang.Call("ToObject",
					g.lang.CallDot("obj", g.lang.PublicName(pi.Name()))))
		} else {
			// TODO OpenObject
			return g.lang.CallDot("obj", g.lang.PublicName(pi.Name()))
		}
	default:
		panic("not implemented")
	}
}

func (g *tsGenerator) generateFactory(prop eg.PropertyObject) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFactory").activated = true

	// export function  NewSimpleTypeFactory(): WuestenFactory<SimpleTypeBuilder, SimpleType>
	className := g.lang.PrivateName(getObjectName(prop), "Factory")
	partialType := g.lang.OrType(
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Param")),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")))
	g.lang.Class(g.bodyWriter, "", g.lang.Implements(className,
		g.lang.Generics("WuestenFactory", g.lang.PublicName(getObjectName(prop)), partialType,
			g.lang.PublicName(getObjectName(prop), "Object"))), prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		}, func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock("Builder():", g.lang.PublicName(getObjectName(prop), "Builder"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return new %s()", g.lang.PublicName(getObjectName(prop), "Builder"))
			})
			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult").activated = true
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("ToObject", g.lang.ReturnType("obj", g.lang.PublicName(getObjectName(prop)))),
				g.lang.PublicName(getObjectName(prop), "Object")), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("const ret: Record<string, unknown> = {}")
				for _, pi := range prop.Items() {
					if !pi.Optional() {
						wr.FormatLine("ret[%s] = %s", g.lang.Quote(pi.Name()), g.genToObject(pi))
						continue
					}
					wr.WriteBlock("if ", fmt.Sprintf("(typeof %s !== 'undefined')",
						g.lang.CallDot("obj", g.lang.PublicName(pi.Name()))), func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("ret[%s] = %s", g.lang.Quote(pi.Name()), g.genToObject(pi))
					})
				}
				wr.FormatLine("return ret as unknown as %s;", g.lang.PublicName(getObjectName(prop), "Object"))

			})

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteJsonDecoder").activated = true
			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("FromPayload", g.lang.ReturnType("val", "WuestePayload"),
						g.lang.Generics("decoder = WuesteJsonDecoder", partialType)),
					g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
				func(wr *eg.ForIfWhileLangWriter) {
					ids := []string{getObjectName(prop)}
					if prop.Id() != "" {
						ids = append(ids, prop.Id())
					}
					if prop.Title() != "" {
						ids = append(ids, prop.Title())
					}
					conditions := []string{}
					for _, id := range ids {
						conditions = append(conditions, fmt.Sprintf("val.Type === %s", g.lang.Quote(id)))
					}
					wr.WriteBlock("if", fmt.Sprintf("(!(%s))", strings.Join(conditions, " || ")), func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(new Error(`WuestePayload Type mismatch:[%s] != ${val.Type}`));", strings.Join(ids, ","))
					})
					// <Partial<SimpleTypeParam>>
					wr.FormatLine("const data = %s", g.lang.Call("decoder", "val.Data"))
					wr.WriteBlock("if", "(data.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteLine("return WuesteResult.Err(data.unwrap_err());")
					})

					wr.WriteLine(
						g.lang.AssignDefault("const builder",
							g.lang.New(g.lang.PublicName(getObjectName(prop), "Builder"))))
					wr.WriteLine("return builder.Coerce(data.unwrap());")
				})

			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult").activated = true
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Clone", g.lang.ReturnType("oth", g.lang.PublicName(getObjectName(prop)))),
				g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("const builder = new %s();", g.lang.PublicName(getObjectName(prop), "Builder"))
				for _, pi := range prop.Items() {
					wr.WriteLine(
						g.lang.Line(
							g.lang.Call(
								g.lang.CallDot("builder", g.lang.PublicName(pi.Name())),
								g.lang.CallDot("oth", g.lang.PublicName(pi.Name())))))
				}
				wr.WriteLine("return builder.Get();")
			})

		})

	g.bodyWriter.WriteLine(
		g.lang.AssignDefault(
			g.lang.Export(g.lang.Const(g.lang.PublicName(getObjectName(prop), "Factory"))),
			g.lang.New(g.lang.PrivateName(getObjectName(prop), "Factory"))))
	g.bodyWriter.WriteLine()
}

func (g *tsGenerator) writeSchema(wr *eg.ForIfWhileLangWriter, prop eg.Property, itemNames ...string) {
	// if prop.Meta().Ref.IsSome() {
	// wr.WriteLine("xxxxx Ref")
	// return
	// }
	if prop.Id() != "" {
		wr.WriteLine(g.lang.Comma(g.lang.ReturnType("id", g.lang.Quote(prop.Id()))))
	}
	wr.WriteLine(g.lang.Comma(g.lang.ReturnType("type", g.lang.Quote(prop.Type()))))

	switch prop.Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		nimpl := "() => { throw new Error('not implemented') }"
		if len(itemNames) > 0 {
			nimpl = fmt.Sprintf("(val) => my.%s.Coerce(val)", g.lang.PrivateName(itemNames[0]))
		}
		wr.WriteLine(g.lang.Comma(g.lang.ReturnType("coerceFromString", nimpl)))
	case eg.OBJECT:
		po := prop.(eg.PropertyObject)
		if po.Schema() != "" {
			wr.WriteLine(g.lang.Comma(g.lang.ReturnType("schema", g.lang.Quote(po.Schema()))))
		}
		if po.Title() != "" {
			wr.WriteLine(g.lang.Comma(g.lang.ReturnType("title", g.lang.Quote(po.Title()))))
		}
		wr.WriteBlock("properties:", "", func(wr *eg.ForIfWhileLangWriter) {
			for _, pi := range po.Items() {
				wr.WriteBlock("", "", func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(g.lang.Comma(g.lang.ReturnType("name", g.lang.Quote(pi.Name()))))
					if pi.Property().Type() == eg.OBJECT && isNamedType(pi.Property()) {
						reflection := g.lang.PublicName(getObjectName(pi.Property()), "Reflection")
						g.includes.AddProperty(reflection, pi.Property())
						wr.WriteLine(g.lang.ReturnType("property",
							g.lang.Call(reflection, g.lang.CallDot("my", g.lang.PrivateName(pi.Name())))))
					} else if pi.Property().Type() == eg.ARRAY {
						reflection := g.lang.PublicName(getObjectName(pi.Property(), []string{pi.Name()}), "Reflection")
						wr.WriteLine(g.lang.ReturnType("property",
							g.lang.Call(reflection, g.lang.CallDot("my", g.lang.PrivateName(pi.Name())))))
					} else {
						wr.WriteBlock("property:", "", func(wr *eg.ForIfWhileLangWriter) {
							if pi.Property().Meta().Parent().IsSome() &&
								pi.Property().Meta().Parent().Value().Type() == eg.OBJECT &&
								pi.Property().Meta().Parent().Value().Meta().Parent().IsSome() &&
								pi.Property().Meta().Parent().Value().Meta().Parent().Value().Type() == eg.ARRAY {
								g.writeSchema(wr, pi.Property())
							} else {
								g.writeSchema(wr, pi.Property(), pi.Name())
							}
						})
					}
				}, "{", "},")
			}
		}, "[", "],")
	case eg.ARRAY:
		pa := prop.(eg.PropertyArray)
		wr.WriteBlock("items:", "", func(wr *eg.ForIfWhileLangWriter) {
			g.writeSchema(wr, pa.Items())
		}, "{", "},")
	}
}

func getDefaultForProperty(prop eg.Property) *string {
	{
		p, ok := prop.(eg.PropertyString)
		if ok {
			if p.Default().IsSome() {
				return wueste.StringLiteral(p.Default().Value()).String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyBoolean)
		if ok {
			if p.Default().IsSome() {
				return wueste.BoolLiteral(p.Default().Value()).String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyNumber)
		if ok {
			if p.Default().IsSome() {
				return wueste.NumberLiteral[float64](p.Default().Value()).String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyInteger)
		if ok {
			if p.Default().IsSome() {
				return wueste.IntegerLiteral[int](p.Default().Value()).String()
			} else {
				return nil
			}
		}
	}
	return nil
}

func genDefaultWuestenAttribute[T string](lang tsLang, name string, prop eg.Property) string {
	defStr := getDefaultForProperty(prop)
	var param string
	if defStr != nil {
		param = fmt.Sprintf("{jsonname: %s, varname: %s, base: baseName, default: %s}", lang.Quote(name), lang.Quote(lang.PublicName(name)), *defStr)
	} else {
		param = fmt.Sprintf("{jsonname: %s, varname: %s, base: baseName}", lang.Quote(name), lang.Quote(lang.PublicName(name)))
	}
	return param
}

func reverse[S ~[]E, E any](s S) S {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func getObjectName(p eg.Property, namess ...[]string) string {
	names := []string{}
	if len(namess) > 0 {
		names = namess[0]
	}
	if p.Meta().Parent().IsNone() {
		name := ""
		if p.Type() == eg.OBJECT {
			name = p.(eg.PropertyObject).Title()
		}
		// fmt.Printf("getObjectName NoParent %s\n", name)
		return strings.Join(reverse(append(names, name)), "$")
	}
	if p.Type() != eg.OBJECT {
		return getObjectName(p.Meta().Parent().Value(), names)
	}
	title := p.(eg.PropertyObject).Title()
	// fmt.Printf("getObjectName Parent %s\n", title)
	// if p.Meta().Ref.IsSome() {
	// return strings.Join(append(names, title), "$")
	// }
	return getObjectName(p.Meta().Parent().Value(), append(names, title))
}

func getObjectFileName(prop eg.Property) string {
	fname := getObjectName(prop)
	fname = fmt.Sprintf("./%s", strings.ToLower(fname))
	return fname
}

func (g *tsGenerator) genWuesteBuilderAttribute(name string, pi eg.PropertyItem, paramFns ...func() string) string {
	prop := pi.Property()
	paramFn := func() string {
		return genDefaultWuestenAttribute(g.lang, name, prop)
	}
	if len(paramFns) > 0 {
		paramFn = paramFns[0]
	}
	switch prop.Type() {
	case eg.STRING:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		p := prop.(eg.PropertyString)
		if p.Format().IsSome() {
			switch p.Format().Value() {
			case eg.DATE_TIME:
				if pi.Optional() {
					return g.lang.Call("wuesten.AttributeDateTimeOptional", paramFn())
				} else {
					return g.lang.Call("wuesten.AttributeDateTime", paramFn())
				}
			default:
			}
		}
		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeStringOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeString", paramFn())
		}
	case eg.INTEGER:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeIntegerOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeInteger", paramFn())
		}
	case eg.NUMBER:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true

		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeNumberOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeNumber", paramFn())
		}
	case eg.BOOLEAN:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeBooleanOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeBoolean", paramFn())
		}
	case eg.ARRAY:
		baseName := g.lang.PublicName(getObjectName(pi.Property(), []string{name}), "Attributes")
		if pi.Optional() {
			return g.lang.New(baseName, paramFn())
			// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttrOptional").activated = true
			// return g.lang.New("WuestenAttrOptional", g.lang.New(baseName, paramFn()))
		} else {
			return g.lang.New(baseName, paramFn())
		}
	case eg.OBJECT:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		p := prop.(eg.PropertyObject)
		objName := getObjectName(p)
		var factory string
		if !isNamedType(p) {
			factory = "WuestenObjectFactory"
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, factory).activated = true
			// g.includes.AddProperty(WUE, factory, p)
		} else {
			factory = g.lang.PublicName(objName, "Factory")
			g.includes.AddProperty(factory, p)
		}
		// param := g.lang.PublicName(objName, "Param")
		// object := g.lang.PublicName(objName, "Object")
		// plain := g.lang.PublicName(objName)
		// obj := g.lang.PublicName(objName, "Object")
		// typ := g.lang.PublicName(objName)
		// var generics func() []string
		if !isNamedType(p) {
			// typ = g.lang.Generics("Record", "string", "unknown")
			// param = g.lang.Generics("Record", "string", "unknown")
			// object = g.lang.Generics("Record", "string", "unknown")
			// plain = g.lang.Generics("Record", "string", "unknown")
			// obj = g.lang.Generics("Record", "string", "unknown")
			// generics = func() []string {
			// 	return []string{obj}
			// }
		} else {
			// generics = func() []string {
			// 	g.includes.AddProperty(object, p)
			// 	g.includes.AddProperty(param, p)
			// 	g.includes.AddProperty(plain, p)
			// 	return []string{typ,
			// 		g.lang.OrType(
			// 			g.lang.Generics("Partial", plain),
			// 			g.lang.Generics("Partial", param),
			// 			g.lang.Generics("Partial", object)),
			// 		obj,
			// 	}
			// }
		}
		if pi.Optional() {
			if !isNamedType(p) {
				// return g.lang.Call(g.lang.Generics("wuesten.AttributeObjectOptional", generics()...), paramFn(), factory)
				return g.lang.Call("wuesten.AttributeObjectOptional", paramFn(), factory)
			} else {
				return g.lang.New(g.lang.PublicName(objName, "Attributes"), paramFn())
				// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttrOptional").activated = true
				// return g.lang.New("WuestenAttrOptional",
				// 	g.lang.New(g.lang.PublicName(objName, "Attributes"), paramFn()))
			}
		} else {
			if !isNamedType(p) {
				// return g.lang.Call(g.lang.Generics("wuesten.AttributeObject", generics()...), paramFn(), factory)
				return g.lang.Call("wuesten.AttributeObject", paramFn(), factory)
			} else {
				return g.lang.New(g.lang.PublicName(objName, "Attributes"), paramFn())
			}
		}
	default:
		panic("not implemented")
	}
}

func getItemType(pa eg.PropertyArray) eg.Property {
	switch pa.Items().Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN, eg.OBJECT:
		return pa.Items()
	case eg.ARRAY:
		return getItemType(pa.Items().(eg.PropertyArray))
	default:
		panic("not implemented")
	}
}

func (g *tsGenerator) generateArrayCoerce(level int, rootArray, returnType string, prop eg.PropertyArray, wr *eg.ForIfWhileLangWriter) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteIterable").activated = true
	wr.FormatLine(g.lang.AssignDefault(
		g.lang.Const(fmt.Sprintf("a%d", level)),
		g.lang.Call(g.lang.Generics("WuesteIterable", g.lang.AsType(prop.Items())), rootArray)))
	wr.WriteIf(fmt.Sprintf("(!a%d)", level), func(wr *eg.ForIfWhileLangWriter) {
		wr.FormatLine("return WuesteResult.Err(\"it's not iterable on level %d\")", level)
	})
	wr.WriteLine(
		g.lang.AssignDefault(
			g.lang.Const(
				g.lang.ReturnType(fmt.Sprintf("r%d", level), g.lang.AsType(prop))),
			"[]"))
	wr.WriteLine(
		g.lang.AssignDefault(
			g.lang.Let(fmt.Sprintf("c%d", level)), "0"))
	root := fmt.Sprintf("i%d", level)
	wr.WriteBlock("for",
		g.lang.RoundBrackets(
			g.lang.Const(root)+" of "+fmt.Sprintf("a%d", level)), func(wr *eg.ForIfWhileLangWriter) {
			p, ok := prop.Items().(eg.PropertyArray)
			if ok {
				g.generateArrayCoerce(level+1, root, returnType, p, wr)
				wr.WriteLine(g.lang.Call(fmt.Sprintf("r%d.push", level), fmt.Sprintf("r%d", level+1)))
			} else {
				param := []string{}
				for i := 0; i <= level; i++ {
					param = append(param, fmt.Sprintf("c%d", i))
				}
				wr.WriteLine(g.lang.Call("itemAttr.SetNameSuffix", strings.Join(param, ", ")))
				wr.WriteLine(g.lang.AssignDefault(
					g.lang.Const("attrRes"), g.lang.Call("itemAttr.Coerce", root)))
				wr.WriteIf(g.lang.RoundBrackets("attrRes.is_err()"), func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(g.lang.Return(
						g.lang.Generics("attrRes as unknown as WuesteResult", returnType)))
				})
				wr.WriteLine(g.lang.Call(fmt.Sprintf("r%d.push", level), g.lang.Call("attrRes.unwrap")))
			}
			wr.FormatLine("c%d++", level)
		})
}

func (g *tsGenerator) generateBuilder(prop eg.PropertyObject) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttribute").activated = true

	for _, pi := range prop.Items() {
		pa, ok := pi.Property().(eg.PropertyArray)
		if ok {
			baseName := getObjectName(pi.Property(), []string{pi.Name()})
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflection").activated = true
			g.bodyWriter.WriteLine("// eslint-disable-next-line @typescript-eslint/no-unused-vars")
			g.bodyWriter.WriteBlock("function ", g.lang.ReturnType(
				g.lang.Call(g.lang.PublicName(baseName, "Reflection"),
					g.lang.ReturnType("param", g.lang.PublicName(baseName, "Attributes"))), "WuestenReflection"),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("return", "", func(wr *eg.ForIfWhileLangWriter) {
						g.writeSchema(wr, pi.Property())
					})
				})
			className := g.lang.PublicName(baseName, "Attributes")
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttr").activated = true
			g.lang.Class(g.bodyWriter, "", g.lang.Extends(className,
				g.lang.Generics("WuestenAttr",
					g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional())),
					g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())))),
				prop,
				func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
				}, func(wr *eg.ForIfWhileLangWriter) {
					attrib := g.lang.ReturnType(g.lang.OptionalParam(g.lang.PrivateName("value"), pi.Optional()), g.lang.AsType(pi.Property()))
					if !pi.Optional() {
						attrib += " = []"
					}
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter").activated = true

					wr.WriteBlock("",
						g.lang.Call("constructor",
							g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeParameter", g.lang.AsType(getItemType(pa)))),
						), func(wr *eg.ForIfWhileLangWriter) {
							pi := eg.NewPropertyItem("ARRAY", rusty.Ok(getItemType(pa)), false).Ok()
							attr := g.genWuesteBuilderAttribute("ARRAY", pi, func() string { return "param" })
							wr.WriteLine(g.lang.AssignDefault(g.lang.Const("itemAttr"), attr))

							wr.WriteBlock("", "super({jsonname: param.jsonname, varname: param.varname, base: param.base}, {coerce: (t0: unknown) => ", func(wr *eg.ForIfWhileLangWriter) {
								g.generateArrayCoerce(0, "t0", g.lang.AsType(pa), pa, wr)
								wr.WriteLine(g.lang.Return(g.lang.Call("WuesteResult.Ok", "r0")))
							}, " {", "}})")
						})
				})
			g.bodyWriter.WriteLine()
		}
	}

	resultsClassName := g.lang.PublicName(getObjectName(prop), "Results")

	g.lang.Interface(g.bodyWriter, "", resultsClassName, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine(g.lang.Readonly(
				g.lang.ReturnType(
					g.lang.PrivateName(pi.Name()),
					g.lang.Generics("WuesteResult", g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()))))))
		})

	attrsClassName := g.lang.PublicName(getObjectName(prop), "Attributes")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflection").activated = true
	g.bodyWriter.WriteBlock("export function ", g.lang.ReturnType(
		g.lang.Call(g.lang.PublicName(getObjectName(prop), "Reflection"),
			g.lang.ReturnType("my", attrsClassName)), "WuestenReflection"), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteBlock("return", "", func(wr *eg.ForIfWhileLangWriter) {
			g.writeSchema(wr, prop)
		})
	})
	cname := g.lang.Implements(attrsClassName,
		g.lang.Generics("WuestenAttribute", g.lang.PublicName(getObjectName(prop)),
			g.lang.OrType(
				g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
				g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")),
			)))

	g.lang.Class(g.bodyWriter, "export ", cname, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {},
		func(wr *eg.ForIfWhileLangWriter) {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter").activated = true
			wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("param",
				g.lang.Generics("WuestenAttributeParameter", g.lang.PublicName(getObjectName(prop))))))
			for _, pi := range prop.Items() {
				// wr.FormatLine("readonly %s = %s;", g.lang.PrivateName(prop.Name()), g.genWuesteBuilderAttribute(prop.Name(), prop.Property()))
				if pi.Property().Type() == eg.OBJECT && isNamedType(pi.Property()) {
					piAttrClassName := g.lang.PublicName(getObjectName(pi.Property()), "Attributes")
					g.includes.AddProperty(piAttrClassName, pi.Property())
					if pi.Optional() {
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), piAttrClassName)))
						// wr.WriteLine(
						// 	g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
						// 		g.lang.Generics("WuestenAttrOptional",
						// 			g.lang.PublicName(getObjectName(pi.Property())),
						// 			g.lang.OrType(
						// 				g.lang.PublicName(getObjectName(pi.Property())),
						// 				g.lang.PublicName(getObjectName(pi.Property()), "Param"),
						// 				"undefined",
						// 			)))))
					} else {
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), piAttrClassName)))
					}
				} else if pi.Property().Type() == eg.ARRAY {
					name := getObjectName(pi.Property(), []string{pi.Name()})
					if pi.Optional() {
						// WuestenAttrOptional<boolean[][][][], (WuesteCoerceTypeboolean|undefined)[][][][]>
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
								g.lang.PublicName(name, "Attributes"))))
						// wr.WriteLine(
						// 	g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
						// 		g.lang.Generics("WuestenAttrOptional",
						// 			g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional())),
						// 			g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
						// 		))))
					} else {
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
								g.lang.PublicName(name, "Attributes"))))
					}
				} else {
					wr.WriteLine(
						g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), g.lang.Generics("WuestenAttribute",
							g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional())),
							g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
						))))
				}
			}

			wr.WriteBlock("", g.lang.Call("constructor",
				g.lang.AssignDefault("param",
					fmt.Sprintf("{jsonname: %s, varname: %s, base: \"\"}",
						g.lang.Quote(getObjectName(prop)), g.lang.Quote(g.lang.PublicName(getObjectName(prop)))))), func(wr *eg.ForIfWhileLangWriter) {
				if len(prop.Items()) > 0 {
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeName").activated = true
					wr.WriteLine(g.lang.AssignDefault(g.lang.CallDot("this", "param"), "param"))
					wr.WriteLine(g.lang.AssignDefault(g.lang.Const("baseName"), g.lang.Call("WuestenAttributeName", "param")))
					for _, pi := range prop.Items() {
						wr.WriteLine(g.lang.AssignDefault(
							g.lang.CallDot("this", g.lang.PrivateName(pi.Name())), g.genWuesteBuilderAttribute(pi.Name(), pi)))
					}
				}
			})

			wr.WriteLine("// eslint-disable-next-line @typescript-eslint/no-unused-vars")
			wr.WriteBlock("", g.lang.ReturnType(g.lang.Call("SetNameSuffix",
				g.lang.ReturnType("...idxs", "number[]")), "void"), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("throw new Error(\"SetNameSuffix:Method not implemented.\")")
			})
			wr.WriteLine("// eslint-disable-next-line @typescript-eslint/no-unused-vars")
			wr.WriteBlock("", g.lang.ReturnType(g.lang.Call("CoerceAttribute",
				g.lang.ReturnType("val", "unknown")), g.lang.Generics("WuesteResult",
				g.lang.PublicName(getObjectName(prop)), "Error")), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("throw new Error(\"CoerceAttribute:Method not implemented.\")")
			})

			wr.WriteBlock("",
				g.lang.AssignDefault(g.lang.Readonly("Coerce"),
					g.lang.ReturnType(
						g.lang.Call("", g.lang.ReturnType("value", "unknown")),
						g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop))))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("if", "(!(typeof value === 'object' && value !== null))", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(Error('expected object'));")
					})
					wr.WriteBlock("return", "this._fromResults", func(wr *eg.ForIfWhileLangWriter) {
						for _, pi := range prop.Items() {
							wr.FormatLine("%s: this.%s.CoerceAttribute(value),", g.lang.PrivateName(pi.Name()), g.lang.PrivateName(pi.Name()))
						}
					}, "({", "});")
				}, " => {")

			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("Get", ""),
					g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("return", "this._fromResults", func(wr *eg.ForIfWhileLangWriter) {
						for _, pi := range prop.Items() {
							wr.FormatLine("%s: this.%s.Get(),", g.lang.PrivateName(pi.Name()), g.lang.PrivateName(pi.Name()))
						}
					}, "({", "});")
				})

			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("_fromResults", g.lang.ReturnType("results", resultsClassName)),
					g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine("const errors: string[] = [];")
					for _, pi := range prop.Items() {
						val := g.lang.CallDot("results", g.lang.PrivateName(pi.Name()))
						wr.WriteBlock("if", "("+val+".is_err())", func(wr *eg.ForIfWhileLangWriter) {
							wr.FormatLine("errors.push(%s.unwrap_err().message);", val)
						})
					}
					// wr.FormatLine("const errors = Object.values(results).filter(r => r.is_err()).map(r => r.unwrap_err().message)")
					wr.WriteBlock("if", "(errors.length > 0)", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(Error(errors.join('\\n')));")
					})
					wr.WriteBlock("return",
						g.lang.Generics("WuesteResult.Ok", g.lang.PublicName(getObjectName(prop))),
						func(wr *eg.ForIfWhileLangWriter) {
							for _, pi := range prop.Items() {
								wr.FormatLine("%s: results.%s.unwrap(),", g.lang.PublicName(pi.Name()), g.lang.PrivateName(pi.Name()))
							}
						}, "({", "});")
				})
		})

	className := g.lang.PublicName(getObjectName(prop), "Builder")
	// extends WuestenAttr<Sub, Partial<Sub>|Partial<SubParam>|Partial<SubObject>>
	// implements WuestenBuilder<Sub, Partial<Sub>|Partial<SubParam>|Partial<SubObject>>

	partialType := g.lang.OrType(
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Param")),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")))
	genericType := []string{g.lang.PublicName(getObjectName(prop)), partialType}

	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenBuilder").activated = true
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttr").activated = true
	g.lang.Class(g.bodyWriter, "export ", g.lang.Implements(
		g.lang.Extends(className, g.lang.Generics("WuestenAttr", genericType...)),
		g.lang.Generics("WuestenBuilder", genericType[0], genericType[1], g.lang.PublicName(getObjectName(prop), "Object"))), prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock(
				g.lang.ReturnType(
					g.lang.Call(g.lang.Type(g.lang.PublicType(prop.Name()), false),
						g.lang.ReturnType(g.lang.OptionalParam("v", prop.Optional()), g.lang.Type(g.lang.AsType(prop.Property(), WithAddCoerce()), false))),
					className),
				"", func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("this._attr.%s.Coerce(v);", g.lang.PrivateName(prop.Name()))
					wr.FormatLine("return this;")
				})
		}, func(wr *eg.ForIfWhileLangWriter) {

			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult").activated = true

			// Coerce(value: unknown): Result<NestedType, Error> {
			// 	throw new Error("Method not implemented.");
			//   }
			//   Get(): Result<NestedType, Error> {
			// 	throw new Error("Method not implemented.");
			//   }

			wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("_attr", attrsClassName)))
			wr.WriteBlock("", g.lang.Call("constructor",
				g.lang.AssignDefault("param",
					fmt.Sprintf("{jsonname: %s, varname: %s, base: \"\"}",
						g.lang.Quote(getObjectName(prop)), g.lang.Quote(g.lang.PublicName(getObjectName(prop)))))), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Const(g.lang.AssignDefault("attr", g.lang.New(attrsClassName, "param"))))
				wr.WriteLine("super(param, {coerce: attr.Coerce});")
				wr.WriteLine(g.lang.AssignDefault("this._attr", "attr"))
			})
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Get", ""), g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Return(g.lang.Call("this._attr.Get", "")))
			})
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "Payload", "WuestePayload").activated = true
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteJsonEncoder").activated = true
			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("AsPayload",
						g.lang.Generics("encoder = WuesteJsonEncoder", g.lang.PublicName(getObjectName(prop), "Object"))),
					g.lang.Generics("WuesteResult", "WuestePayload")),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("const val = this.Get();")
					wr.WriteBlock("if", "(val.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(val.unwrap_err());")
					})
					wr.FormatLine("const data = encoder(%s.ToObject(val.unwrap()))", g.lang.PublicName(getObjectName(prop), "Factory"))
					wr.WriteBlock("if", "(data.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(data.unwrap_err());")
					})
					wr.WriteBlock("return", "WuesteResult.Ok", func(wr *eg.ForIfWhileLangWriter) {
						id := prop.Id()
						if id == "" {
							id = prop.Title()
						}
						wr.FormatLine("Type: %s,", g.lang.Quote(id))
						wr.FormatLine("Data: data.unwrap()")
					}, "({", "});")
				})

		})

}

type tsGenerator struct {
	cfg        *eg.GeneratorConfig
	lang       tsLang
	includes   *externalTypes
	bodyWriter *eg.ForIfWhileLangWriter
}

func (g *tsGenerator) generatePropertyObject(prop eg.PropertyObject, sl eg.PropertyCtx) {
	g.generateClass(prop)
	g.generateJson(prop)
	g.generateBuilder(prop)
	g.generateFactory(prop)

	os.MkdirAll(g.cfg.OutputDir, 0755)

	fname := filepath.Join(g.cfg.OutputDir, getObjectFileName(prop)+".ts")
	tmpFname := filepath.Join(g.cfg.OutputDir, "."+getObjectFileName(prop)+uuid.New().String()+".ts")

	fmt.Printf("Generate: %s -> %s\n", prop.Meta().FileName().Value(), fname)
	wr, err := os.OpenFile(tmpFname, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		wr.Close()
		os.Remove(fname)
		os.Rename(tmpFname, fname)
	}()

	header := eg.NewForIfWhileLangWriter(eg.ForIfWhileLangWriter{OfsIndent: g.cfg.EntityCfg.Indent})
	if len(g.includes.ActiveTypes()) > 0 {
		for _, include := range g.includes.ActiveTypes() {
			filename := include.fileName
			if include.activated && include.property.IsSome() {
				filename = getObjectFileName(include.property.Value())
			}
			// getObjectName(include.property)
			if len(include.Types()) <= 3 {
				header.FormatLine("import { %s } from %s;", strings.Join(include.Types(), ", "),
					g.lang.Quote(filename))
			} else {
				header.WriteBlock("", "import", func(wr *eg.ForIfWhileLangWriter) {
					for idx, t := range include.Types() {
						comma := ","
						if idx == len(include.Types())-1 {
							comma = ""
						}
						wr.FormatLine("%s%s", t, comma)
					}
				}, " {", fmt.Sprintf("} from %s;", g.lang.Quote(filename)))
			}
		}
		header.WriteLine()
	}

	for _, line := range header.Lines() {
		wr.Write([]byte(line))
	}

	for _, line := range g.bodyWriter.Lines() {
		wr.Write([]byte(line))
	}
	// sl.Registry.SetWritten(prop)
}

type externalType struct {
	toGenerate bool
	activated  bool
	// prefix     string
	fileName string
	types    map[string]*string
	property rusty.Optional[eg.PropertyObject]
}

type ImportType struct {
	Alias *string
	Type  string
}

func (et *externalType) Types() []string {
	types := make([]string, 0, len(et.types))
	for k, alias := range et.types {
		str := k
		if alias != nil {
			str = fmt.Sprintf("%s as %s", k, *alias)
		}
		types = append(types, str)
	}
	sort.Strings(types)
	return types
}

type externalTypes struct {
	types map[string]*externalType
}

func (g *externalTypes) AddProperty(typ string, prop eg.Property) {
	po, ok := prop.(eg.PropertyObject)
	if ok {
		t := g.AddType(prop.Meta().FileName().Value(), typ)
		t.property = rusty.Some(po)
		t.activated = true
	}
}

func (g *externalTypes) AddType(fileName, typeName string, optAlias ...string) *externalType {
	var alias *string = nil
	if len(optAlias) > 0 {
		alias = &optAlias[0]
	}
	if typeName == "SimpleType$PayloadObject" {
		fmt.Println("SimpleType$PayloadObject")
	}
	if fileName == "" {
		return nil
	}
	et, ok := g.types[fileName]
	if !ok {
		et = &externalType{
			toGenerate: false,
			fileName:   fileName,
			types:      make(map[string]*string),
		}
		g.types[fileName] = et
	}
	et.types[typeName] = alias
	return et
}

func (g *externalTypes) ActiveTypes() []*externalType {
	atyp := make([]*externalType, 0, len(g.types))
	for _, et := range g.types {
		if et.activated {
			atyp = append(atyp, et)
		}
	}
	sort.Slice(atyp, func(i, j int) bool {
		return strings.Compare(atyp[i].fileName, atyp[j].fileName) < 0
	})
	return atyp
}

func newExternalTypes() *externalTypes {
	return &externalTypes{
		types: make(map[string]*externalType),
	}
}

func TsGenerator(cfg *eg.GeneratorConfig, prop eg.Property, sl eg.PropertyCtx) {
	// po, found := schema.(eg.PropertyObject)
	// if !found {
	// 	panic("TsGenerator not a property object")
	// }
	// if sl.Registry.IsWritten(jsf) {
	// 	return
	// }
	g := &tsGenerator{
		cfg:        cfg,
		includes:   newExternalTypes(),
		bodyWriter: eg.NewForIfWhileLangWriter(eg.ForIfWhileLangWriter{OfsIndent: cfg.EntityCfg.Indent}),
	}

	// rpo := eg.NewPropertiesBuilder(sl).SetFileName(jsf.FileName).FromJson(jsf.JSONProperty).Build()
	// if rpo.IsErr() {
	// 	panic(rpo.Err())
	// }

	g.generatePropertyObject(prop.(eg.PropertyObject), sl)
	for _, prop := range g.includes.ActiveTypes() {
		if prop.property.IsSome() {
			TsGenerator(cfg, prop.property.Value(), sl)
		}
	}
}
