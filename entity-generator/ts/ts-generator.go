package ts

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/google/uuid"
	eg "github.com/mabels/wueste/entity-generator"
	"github.com/mabels/wueste/entity-generator/rusty"
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

func (l *tsLang) addCoerceType(typ string, withs ...withResult) string {
	var prop *withResult
	if hasWith(WithAddInputType(), withs) {
		switch typ {
		case "number":
			typ = "WuestenFormatter.Number.Type"
		case "integer":
			typ = "WuestenFormatter.Integer.Type"
		case "string":
			typ = "WuestenFormatter.String.Type"
		case "boolean":
			typ = "WuestenFormatter.Boolean.Type"
		case "Date":
			typ = "WuestenFormatter.Date.Type"
		case "any":
			typ = "WuestenFormatter.Any.Type"
		}
		prop = hasWithResult(WithIsInputType(nil), withs)
	}
	if hasWith(WithAddCoerce(), withs) {
		switch typ {
		case "number":
			typ = "WuestenFormatter.Number.CoerceType"
		case "integer":
			typ = "WuestenFormatter.Integer.CoerceType"
		case "string":
			typ = "WuestenFormatter.String.CoerceType"
		case "boolean":
			typ = "WuestenFormatter.Boolean.CoerceType"
		case "Date":
			typ = "WuestenFormatter.Date.CoerceType"
		case "any":
			typ = "WuestenFormatter.Any.CoerceType"
		}
		prop = hasWithResult(WithIsCoerceType(nil), withs)
	}
	if hasWith(WithAddObject(), withs) {
		switch typ {
		case "number":
			typ = "WuestenFormatter.Number.ObjectType"
		case "integer":
			typ = "WuestenFormatter.Integer.ObjectType"
		case "string":
			typ = "WuestenFormatter.String.ObjectType"
		case "boolean":
			typ = "WuestenFormatter.Boolean.ObjectType"
		case "Date":
			typ = "WuestenFormatter.Date.ObjectType"
		case "any":
			typ = "WuestenFormatter.Any.ObjectType"
		}
		prop = hasWithResult(WithIsObjectType(nil), withs)
	}
	res := hasWithResult(WithAddType(func(typ string, prop eg.Property) {}), withs)
	if res != nil && prop == nil && hasWith(WithAddCoerce(), withs) {
		res.addType(typ, nil)
	}
	if res != nil && prop != nil {
		res.addType(typ, prop.prop)
	}
	return typ
}

// func handleAddType(typ string, prop eg.Property, withs ...withResult) string {
// 	res := hasWithResult(WithAddType(func(typ string, prop eg.Property) {}), withs)
// 	if res != nil {
// 		res.addType(typ, prop)
// 	}
// 	return typ
// }

func (l *tsLang) AsTypeHelper(p eg.Property, withs ...withResult) string {
	switch p.Type() {
	case eg.OBJECT:
		po := p.(eg.PropertyObject)
		if po.Properties() == nil || po.Properties().Len() == 0 {
			return l.addCoerceType("any", withs...)
		}
		name := getObjectName(p)
		if hasWith(WithAddCoerce(), withs) {
			return l.addCoerceType(l.Name(name, "CoerceType"),
				append(withs, WithIsCoerceType(p))...)
		}

		if hasWith(WithAddObject(), withs) {
			return l.addCoerceType(l.Name(name, "Object"), append(withs, WithIsObjectType(p))...)
		}
		if hasWith(WithAddInputType(), withs) {
			return l.addCoerceType(l.Name(name, ""), append(withs, WithIsInputType(p))...)
		}

		ret := l.Name(name)
		res := hasWithResult(WithAddType(func(typ string, prop eg.Property) {}), withs)
		if res != nil {
			res.addType(ret, p)
		}
		return ret
	case eg.STRING:
		p := p.(eg.PropertyString)
		if p.Format().IsSome() {
			switch p.Format().Value() {
			case eg.DATE_TIME:
				return l.addCoerceType("Date", withs...)
			case eg.DATE, eg.TIME:
				panic("not implemented")
			default:
				if strings.HasPrefix(p.Format().Value(), "#") {
					return l.addCoerceType("string", withs...)
				}
				dynUrl, err := url.Parse(p.Format().Value())
				if err == nil {
					if dynUrl.Scheme == "file" {
						return l.addCoerceType("string", withs...)
					}
				}
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

func (l *tsLang) AsType(p eg.Property, withs ...withResult) string {
	out := l.AsTypeHelper(p, withs...)
	if hasWith(WithPartial(), withs) {
		out = l.Generics("Partial", out)
	}
	return out
}

func (l *tsLang) AsTypeOptional(p eg.Property) string {
	return "rusty.Optional<" + l.AsType(p) + ">"
}

func hasWith(str withResult, withs []withResult) bool {
	return hasWithResult(str, withs) != nil
}

func hasWithResult(str withResult, withs []withResult) *withResult {
	for _, with := range withs {
		if with.action != "" && str.action == with.action {
			return &with
		}
	}
	return nil
}

type withResult struct {
	action  string
	addType func(typ string, prop eg.Property)
	prop    eg.Property
}

func WithIsCoerceType(prop eg.Property) withResult {
	return withResult{action: "isCoerceType", prop: prop}
}

func WithIsObjectType(prop eg.Property) withResult {
	return withResult{action: "isObjectType", prop: prop}
}

func WithIsInputType(prop eg.Property) withResult {
	return withResult{action: "isInputType", prop: prop}
}

// func WithIsType(prop eg.Property) withResult {
// 	return withResult{action: "isType", prop: prop}
// }

func WithAddType(addType func(typ string, prop eg.Property)) withResult {
	return withResult{action: "addType", addType: addType}
}

func WithOptional(optional bool) withResult {
	if optional {
		return withResult{action: "isOptional"}
	}
	return withResult{action: "isNotOptional"}
}

func WithAddCoerce() withResult {
	return withResult{action: "addCoerceType"}
}

func WithAddObject() withResult {
	return withResult{action: "addObjectType"}
}

func WithAddInputType() withResult {
	return withResult{action: "addInputType"}
}

func WithPartial() withResult {
	return withResult{action: "isPartial"}
}

func (l *tsLang) AsTypeNullable(p eg.Property, withs ...withResult) string {
	res := l.AsType(p, withs...)

	if hasWith(WithOptional(true), withs) {
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

func (l *tsLang) New(line string, args ...string) string {
	return "new " + l.Call(line, args...)
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

func callDot(prop string) (string, string) {
	if strings.HasPrefix(prop, "\"") {
		return "", "[" + prop + "]"
	}
	if strings.HasPrefix(prop, "[") {
		return "", prop
	}
	return ".", prop
}

func (l *tsLang) CallDot(props ...string) string {
	out := ""
	first := true
	for _, prop := range props {
		sep, dot := callDot(prop)
		if dot == "" {
			continue
		}
		if first {
			sep = ""
		}
		out = out + sep + dot
		first = false
	}
	return out
}

func (l *tsLang) OrType(str ...string) string {
	return strings.Join(str, "|")
}

func (l *tsLang) Const(str string) string {
	return "const " + str
}

func (l *tsLang) Trinary(condition, trueish, falsisch string) string {
	return condition + " ? " + trueish + " : " + falsisch
}
func (l *tsLang) Let(str string) string {
	return "let " + str
}

func (l *tsLang) RoundBrackets(str string) string {
	return "(" + str + ")"
}

func (l *tsLang) CurlyBrackets(str string) string {
	return "{" + str + "}"
}

func (l *tsLang) Cast(str string, asType string) string {
	return l.RoundBrackets(str + " as " + asType)
}

func (l *tsLang) Index(str string, index string) string {
	return str + "[" + index + "]"
}

func hasDefault(prop eg.Property) bool {
	return getDefaultForProperty(prop) != nil
}

// var reOrArray = regexp.MustCompile(`[(\|)(\[\])]+`)

// func isNamedType(p eg.Property) bool {
// 	switch p.Type() {
// 	case eg.OBJECT:
// 		po := p.(eg.PropertyObject)
// 		if po.Id() == "" || po.Title() == "" {
// 			return false
// 		}
// 		if po.Properties() == nil || po.Properties().Len() == 0 {
// 			return false
// 		}
// 		return true
// 	default:
// 	}
// 	return false
// }

func handleAddType(g *tsGenerator) func(string, eg.Property) {
	return func(typ string, prop eg.Property) {
		if prop == nil {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFormatter")
		} else {
			g.includes.AddProperty(typ, prop)
		}
	}
}

func (g *tsGenerator) generateClass(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop)), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		// out := []string{}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()), pi.Optional())),
				g.lang.AsTypeNullable(pi.Property(),
					WithAddType(handleAddType(g)), WithAddInputType()))))
		// wr.FormatLine("// %v", out)

	})
	g.bodyWriter.WriteLine()

	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Param"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		// if isNamedType(pi.Property()) {
		// 	paramName := getObjectName(pi.Property())
		// 	g.includes.AddProperty(paramName, pi.Property())
		// 	// g.includes.AddProperty(g.lang.PublicName(paramName, "Param"), pi.Property())
		// }
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()),
					pi.Optional() || hasDefault(pi.Property()))),
				g.lang.AsTypeNullable(pi.Property(),
					WithAddType(handleAddType(g)),
					WithAddCoerce()))))
	})
	g.bodyWriter.WriteLine()

}

func (g *tsGenerator) generateJSONDict(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Object"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		// typ := g.lang.AsTypeNullable(pi.Property())
		// if isNamedType(pi.Property()) {
		// 	typ = g.lang.PublicName(g.lang.AsType(pi.Property()), "Object")
		// 	g.includes.AddProperty(typ, pi.Property())
		// 	// if pi.Optional() {
		// 	// 	typ = g.lang.OrType(typ, "undefined")
		// 	// }
		// }
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.Quote(pi.Name()), pi.Optional())),
				g.lang.AsTypeNullable(pi.Property(),
					WithAddType(handleAddType(g)),
					WithAddObject())))) //, g.lang.JsonAnnotation(pi.Property()))))
	})
	g.bodyWriter.WriteLine()
}

func (g *tsGenerator) generateFactory(prop eg.PropertyObject) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFactory")
	// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFactoryParam")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenNames")

	names := g.lang.PrivateName(getObjectName(prop), "Names")
	g.bodyWriter.WriteBlock(
		g.lang.Const(
			g.lang.ReturnType(names, "WuestenNames")), " = ",
		func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock(g.lang.Call("", "id", "title", "varname"), "=>", func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteBlock("return", "", func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("id, title, names: [id, title], varname")
				})
			})
		}, "(", g.lang.Call(")",
			g.lang.Quote(prop.Id()),
			g.lang.Quote(prop.Title()),
			g.lang.Quote(g.lang.PublicName(getObjectName(prop)))))
	g.bodyWriter.WriteLine()

	coerceType := g.lang.PublicType(getObjectName(prop), "CoerceType")
	typ := g.lang.PublicType(getObjectName(prop))
	objectType := g.lang.PublicType(getObjectName(prop), "Object")

	// export function  NewSimpleTypeFactory(): WuestenFactory<SimpleTypeBuilder, SimpleType>
	className := g.lang.PublicName(getObjectName(prop), "FactoryImpl")
	// partialType := g.lang.OrType(
	// 	g.lang.PublicName(getObjectName(prop), "Object"),
	// 	g.lang.PublicName(getObjectName(prop), "CoerceType"))
	// fullType := g.lang.PublicName(getObjectName(prop))
	// g.lang.OrType(
	// 	g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
	// 	g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Param")),
	// 	g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")))
	g.lang.Class(g.bodyWriter, "export ", g.lang.Extends(className,
		g.lang.Generics("WuestenFactory", typ, coerceType, objectType)), prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		}, func(wr *eg.ForIfWhileLangWriter) {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenConstructionParams")
			fpArg := g.lang.ReturnType("params", g.lang.Generics("WuestenConstructionParams", coerceType))
			// wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("params", g.lang.Generics("WuestenAttributeParameter", fullType))))
			wr.WriteBlock(g.lang.Call("constructor", fpArg), "", func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteBlock("super", "", func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(g.lang.Comma(g.lang.ReturnType("jsonname", g.lang.CallDot(names, "title"))))
					wr.WriteLine(g.lang.Comma(g.lang.ReturnType("varname", g.lang.CallDot(names, "varname"))))
					wr.WriteLine(g.lang.Comma(g.lang.ReturnType("base", g.lang.Quote(""))))
					// wr.WriteLine(g.lang.Comma(g.lang.ReturnType("dynamic", "{}")))
					// wr.WriteLine(g.lang.Comma(g.lang.ReturnType("encoder", "WuestenJSONPassThroughEncoder")))
					// wr.WriteLine(g.lang.Comma(g.lang.ReturnType("decoder", "WuestenJSONPassThroughDecoder")))
					wr.WriteLine("...params")
				}, "({", "})")
			})

			wr.WriteBlock(
				g.lang.ReturnType(
					g.lang.Call("Names"), "WuestenNames"), "",
				func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("return %s", g.lang.PrivateName(getObjectName(prop), "Names"))
				})

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter")
			wr.WriteBlock(
				g.lang.ReturnType(
					g.lang.Call("Builder",
						// base: WuestenAttributeBase<unknown>, params: WuestenFactoryParam<C>)
						g.lang.ReturnType("base?", "WuestenAttributeBase<unknown>"),
						g.lang.ReturnType("params?", g.lang.Generics("WuestenAttributeParameter", coerceType))),
					g.lang.PublicName(getObjectName(prop), "Builder")), "",
				func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("return new %s(this._params, base, params)", g.lang.PublicName(getObjectName(prop), "Builder"))
				})
			// g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult")
			// wr.WriteBlock("", g.lang.ReturnType(
			// 	g.lang.Call("ToObject", g.lang.ReturnType("obj", g.lang.PublicName(getObjectName(prop)))),
			// 	g.lang.PublicName(getObjectName(prop), "Object")), func(wr *eg.ForIfWhileLangWriter) {
			// 	wr.FormatLine("return %s", g.lang.Call(g.lang.PublicName(getObjectName(prop), "ToObject"), "obj"))
			// })
			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult")
			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("FromPayload", g.lang.ReturnType("val", "WuestePayload"), "decoder = this._params.decoder"),
					g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop), "Builder"))),
				func(wr *eg.ForIfWhileLangWriter) {
					// ids := []string{getObjectName(prop)}
					// if prop.Id() != "" {
					// 	ids = append(ids, prop.Id())
					// }
					// if prop.Title() != "" {
					// 	ids = append(ids, prop.Title())
					// }
					// conditions := []string{}
					// for _, id := range ids {
					// 	conditions = append(conditions, fmt.Sprintf("val.Type === %s", g.lang.Quote(id)))
					// }
					// !_SimpleTypeNames.names.includes(val.Type)
					wr.WriteBlock("if",
						fmt.Sprintf("(!%s.names.includes(val.Type))", g.lang.PrivateName(getObjectName(prop), "Names")),
						func(wr *eg.ForIfWhileLangWriter) {
							wr.FormatLine("return WuesteResult.Err(new Error(`WuestePayload Type mismatch:[${%s.names.join(',')}] != ${val.Type}`));",
								g.lang.PrivateName(getObjectName(prop), "Names"))
						})
					// <Partial<SimpleTypeParam>>
					wr.FormatLine("const data = %s", g.lang.Call("decoder", "val.Data"))
					wr.WriteBlock("if", "(data.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteLine("return WuesteResult.Err(data.unwrap_err());")
					})
					// (new DynamicDefaultBuilder(this._params)).Coerce(data.unwrap() as DynamicDefaultCoerceType)
					wr.WriteLine(
						g.lang.Return(g.lang.Call(g.lang.CallDot("WuesteResult", "Ok"),
							g.lang.Call(
								g.lang.CallDot(
									g.lang.RoundBrackets(g.lang.Call(g.lang.CallDot("this", "Builder"))),
									"Coerce"), g.lang.Cast("data.unwrap()", coerceType)))))
				})

			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult")
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Clone", g.lang.ReturnType("oth", g.lang.PublicName(getObjectName(prop)))),
				g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("const builder = new %s(this._params);", g.lang.PublicName(getObjectName(prop), "Builder"))
				for _, pi := range prop.Items() {
					wr.WriteLine(
						g.lang.Line(
							g.lang.Call(
								g.lang.CallDot("builder", g.lang.PublicName(pi.Name())),
								g.lang.CallDot("oth", g.lang.PublicName(pi.Name())))))
				}
				wr.WriteLine("return builder.Get();")
			})

			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Schema"), "WuestenReflection"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return %s;", g.lang.PublicName(getObjectName(prop), "Schema"))
			})
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Getter", g.lang.ReturnType("typ", g.lang.PublicName(getObjectName(prop))),
					g.lang.ReturnType("base", "WuestenReflection[] = []")), "WuestenGetterBuilder"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return %s", g.lang.Call(g.lang.PublicName(getObjectName(prop), "Getter"), "typ", "base"))
			})

		})

	// WuestenTypeRegistry.Register(n
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenTypeRegistry")
	g.bodyWriter.WriteLine(
		g.lang.AssignDefault(
			g.lang.Export(g.lang.Const(g.lang.PublicName(getObjectName(prop), "Factory"))),
			g.lang.Call("WuestenTypeRegistry.Register",
				g.lang.New(g.lang.PublicName(getObjectName(prop), "FactoryImpl"),
					g.lang.Call("WuestenTypeRegistry.cloneAttributeBase")))))
	g.bodyWriter.WriteLine()
}

type propFormat interface {
	Format() rusty.Optional[string]
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
		po := prop.(propFormat)
		if po.Format().IsSome() {
			wr.WriteLine(g.lang.Comma(g.lang.ReturnType("format", g.lang.Quote(po.Format().Value()))))
		}
		// nimpl := "() => { throw new Error('not implemented') }"
		// if len(itemNames) > 0 {
		// 	nimpl = fmt.Sprintf("(val) => %s", g.lang.Call(g.lang.CallDot("my",
		// 		g.lang.PublicName(itemNames[0])), "val"))
		// }
		// wr.WriteLine(g.lang.Comma(g.lang.ReturnType("coerceFromString", nimpl)))
		// nimpl = "() => { throw new Error('not implemented') }"
		// if len(itemNames) > 0 {
		// 	nimpl = fmt.Sprintf("() => %s", g.lang.Call(
		// 		"JSON.stringify",
		// 		g.lang.CallDot("my.Get().unwrap()",
		// 			g.lang.PublicName(itemNames[0]))))
		// }
		// wr.WriteLine(g.lang.Comma(g.lang.ReturnType("getAsString", nimpl)))
	case eg.OBJECT:
		po := prop.(eg.PropertyObject)
		if po.Schema() != "" {
			wr.WriteLine(g.lang.Comma(g.lang.ReturnType("schema", g.lang.Quote(po.Schema()))))
		}
		if po.Title() != "" {
			wr.WriteLine(g.lang.Comma(g.lang.ReturnType("title", g.lang.Quote(po.Title()))))
		}
		if len(po.Required()) > 0 {
			wr.WriteBlock("required:", "", func(wr *eg.ForIfWhileLangWriter) {
				for _, req := range po.Required() {
					wr.WriteLine(g.lang.Comma(g.lang.Quote(req)))
				}
			}, "[", "],")
		}
		items := po.Items()
		if len(items) > 0 {
			wr.WriteBlock("properties:", "", func(wr *eg.ForIfWhileLangWriter) {
				for _, pi := range po.Items() {
					wr.WriteBlock("", "", func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteLine(g.lang.Comma(g.lang.ReturnType("type", g.lang.Quote("objectitem"))))
						wr.WriteLine(g.lang.Comma(g.lang.ReturnType("name", g.lang.Quote(pi.Name()))))
						if pi.Property().Type() == eg.OBJECT /* && isNamedType(pi.Property()) */ {
							reflection := g.lang.PublicName(getObjectName(pi.Property()), "Schema")
							g.includes.AddProperty(reflection, pi.Property())
							// builderName := g.lang.Call(g.lang.PublicType(pi.Name(), "Builder"))
							// if pi.Optional() {
							// 	builderName = g.lang.CallDot(builderName, "typ")
							// }
							wr.WriteLine(g.lang.ReturnType("property", reflection))
							// } else if pi.Property().Type() == eg.ARRAY {
							// 	reflection := g.lang.PublicName(getObjectName(pi.Property(), []string{pi.Name()}), "Schema")
							// 	// builderName := g.lang.Call(g.lang.PublicType(pi.Name(), "Builder"))
							// 	// if pi.Optional() {
							// 	// 	builderName = g.lang.CallDot(builderName, "typ")
							// 	// }
							// 	wr.WriteLine(g.lang.ReturnType("property", reflection))
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
		}
	case eg.ARRAY:
		pa := prop.(eg.PropertyArray)
		wr.WriteBlock("items:", "", func(wr *eg.ForIfWhileLangWriter) {
			g.writeSchema(wr, pa.Items())
		}, "{", "},")
	}
}

type propDefault interface {
	Default() rusty.Optional[any]
}

func getDefaultForProperty(prop eg.Property) *string {
	pDef, ok := prop.(propDefault)
	if !ok {
		return nil
	}
	if pDef.Default().IsNone() {
		return nil
	}
	bStr, err := json.Marshal(pDef.Default().Value())
	if err != nil {
		panic(err)
	}
	str := string(bStr)
	return &str
}

func genDefaultWuestenAttribute[T string](g *tsGenerator, name string, prop eg.Property) []string {
	defStr := getDefaultForProperty(prop)
	var param string
	format := ""
	switch prop.Type() {
	case eg.STRING:
		ps := prop.(eg.PropertyString)
		if ps.Format().IsSome() {
			format = g.lang.ReturnType(",format", g.lang.Quote(ps.Format().Value()))
		}
	case eg.NUMBER:
		ps := prop.(eg.PropertyNumber)
		if ps.Format().IsSome() {
			format = g.lang.ReturnType(",format", g.lang.Quote(ps.Format().Value()))
		}
	case eg.INTEGER:
		ps := prop.(eg.PropertyInteger)
		if ps.Format().IsSome() {
			format = g.lang.ReturnType(",format", g.lang.Quote(ps.Format().Value()))
		}
	case eg.BOOLEAN:
		ps := prop.(eg.PropertyBoolean)
		if ps.Format().IsSome() {
			format = g.lang.ReturnType(",format", g.lang.Quote(ps.Format().Value()))
		}
	}
	if defStr != nil {
		param = fmt.Sprintf("{jsonname: %s, varname: %s, base: baseName, default: %s%s}", g.lang.Quote(name), g.lang.Quote(g.lang.PublicName(name)), *defStr, format)
	} else {
		param = fmt.Sprintf("{jsonname: %s, varname: %s, base: baseName%s}", g.lang.Quote(name), g.lang.Quote(g.lang.PublicName(name)), format)
	}
	// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeFactory")
	return []string{"this.param", param}
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

func getArrayItemProperty(prop eg.Property) eg.Property {
	switch prop.Type() {
	case eg.ARRAY:
		return getArrayItemProperty(prop.(eg.PropertyArray).Items())
	default:
		return prop
	}
}

func (g *tsGenerator) genWuesteBuilderAttribute(name string, pi eg.PropertyItem /*, paramFns ...func() string */) string {
	prop := pi.Property()
	// paramFn := func() string {
	// 	return genDefaultWuestenAttribute(g, name, prop)
	// }
	// if len(paramFns) > 0 {
	// 	paramFn = paramFns[0]
	// }
	typ := g.lang.AsTypeNullable(pi.Property(), WithAddInputType())
	coerceTyp := g.lang.AsTypeNullable(pi.Property(), WithAddCoerce())
	objectTyp := g.lang.AsTypeNullable(pi.Property(), WithAddObject())
	// params := genDefaultWuestenAttribute(g, name, prop)
	switch prop.Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		if pi.Optional() {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeFactoryOptional")
			return g.lang.Call(g.lang.Generics("WuestenAttributeFactoryOptional", typ, coerceTyp, objectTyp))
		} else {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeFactory")
			return g.lang.Call(g.lang.Generics("WuestenAttributeFactory", typ, coerceTyp, objectTyp))
		}
	case eg.ARRAY:
		pa := eg.NewPropertyArrayItem(name, rusty.Ok(getArrayItemProperty(pi.Property())), false).Ok()
		// baseName := g.lang.PublicName(getObjectName(pi.Property(), []string{name}), "Builder")
		itemProp := getArrayItemProperty(pi.Property().(eg.PropertyArray).Items())

		ityp := g.lang.AsTypeNullable(itemProp, WithAddInputType())
		icoerceTyp := g.lang.AsTypeNullable(itemProp, WithAddCoerce())
		iobjectTyp := g.lang.AsTypeNullable(itemProp, WithAddObject())
		if pi.Optional() {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenArrayFactoryOptional")
			return g.lang.New(g.lang.Generics("WuestenArrayFactoryOptional", typ, coerceTyp, objectTyp, ityp, icoerceTyp, iobjectTyp), g.genWuesteBuilderAttribute(name, pa))
		} else {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenArrayFactory")
			return g.lang.New(g.lang.Generics("WuestenArrayFactory", typ, coerceTyp, objectTyp, ityp, icoerceTyp, iobjectTyp), g.genWuesteBuilderAttribute(name, pa))
		}
	case eg.OBJECT:
		po := prop.(eg.PropertyObject)
		objName := getObjectName(po)
		if pi.Optional() {
			// if !isNamedType(po) {
			// 	// return g.lang.Call(g.lang.Generics("wuesten.AttributeObjectOptional", generics()...), paramFn(), factory)
			// 	factory := "WuestenObjectFactory"
			// 	g.includes.AddType(g.cfg.EntityCfg.FromWueste, factory)
			// 	return g.lang.Call("WuestenObjectFactory", factory, "KAPUT")
			// } else {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
			name := g.lang.PublicName(objName, "Factory")
			g.includes.AddProperty(name, po)
			return g.lang.Call("WuestenObjectOptional", name)
			// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttrOptional").activated = true
			// return g.lang.New("WuestenAttrOptional",
			// 	g.lang.New(g.lang.PublicName(objName, "Attributes"), paramFn()))
			// }
		} else {
			// if !isNamedType(po) {
			// 	// return g.lang.Call(g.lang.Generics("wuesten.AttributeObject", generics()...), paramFn(), factory)
			// 	// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten")
			// 	factory := "WuestenObjectFactory"
			// 	g.includes.AddType(g.cfg.EntityCfg.FromWueste, factory)
			// 	return g.lang.Call("WuestenObjectFactory", factory, "KAPUT")
			// } else {
			name := g.lang.PublicName(objName, "Factory")
			g.includes.AddProperty(name, po)
			return name // g.lang.Call(name, params...)
			// }
		}
	default:
		panic("not implemented")
	}
}

type attributeClassParam struct {
	resultsClassName string
	coerceTyp        string
	typ              string
	objectTyp        string
	attrsClassName   string
}

type objectProperty struct {
	prop              eg.Property
	coerceTyp         string
	builderTyp        string
	builderCall       string
	fnGetBuilderTypes []string
}

func (g *tsGenerator) propertyFactory(pin eg.PropertyItem) objectProperty {
	po, found := pin.Property().(eg.PropertyObject)

	coerceTyp := g.lang.Type(g.lang.AsType(pin.Property(), WithAddCoerce()), false)
	builderTyp := g.lang.Type(g.lang.AsType(pin.Property(), WithAddInputType()), false)

	if found {
		if po.Type() == eg.OBJECT && (po.Properties() == nil || po.Properties().Len() == 0) {
			// po = eg.NewPropertyObject(eg.PropertyObjectBuilder{
			// 	Type:        po.Type(),
			// 	Description: po.Description(),

			// 	Id:         "WuestenFormatter.Any",
			// 	Title:      "WuestenFormatter.Any",
			// 	Schema:     po.Schema(),
			// 	Properties: po.Properties(),
			// 	Required:   po.Required(),
			// 	Ref:        po.Ref(),

			// 	// Errors: po.Errors(),
			// }).Ok()
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFormatter")
			return objectProperty{
				prop:              pin,
				coerceTyp:         "WuestenFormatter.Any.CoerceType",
				builderTyp:        "WuestenFormatter.Any.Builder",
				builderCall:       "new WuestenFormatter.Any.Builder",
				fnGetBuilderTypes: []string{"WuestenFormatter.Any.Builder", "WuestenFormatter.Any.CoerceType"},
			}
		}
		builderTyp = g.lang.PublicName(getObjectName(pin.Property()), "Builder")
	}
	g.includes.AddProperty(g.lang.PublicName(getObjectName(pin.Property()), "Builder"), pin.Property())

	return objectProperty{
		prop:              pin,
		coerceTyp:         coerceTyp,
		builderCall:       g.lang.CallDot(g.genWuesteBuilderAttribute(pin.Name(), pin), "Builder"),
		builderTyp:        builderTyp,
		fnGetBuilderTypes: []string{builderTyp, coerceTyp},
	}
}

func (g *tsGenerator) generateAttributesClass(wr *eg.ForIfWhileLangWriter, prop eg.PropertyObject, acp attributeClassParam) {
	// attrsClassName := g.lang.PublicName(getObjectName(prop), "Attributes")
	// cname := attrsClassName
	// // implements WuestenFormatter<>
	// typ := g.lang.PublicType(getObjectName(prop))
	// coerceTyp := g.lang.PublicType(getObjectName(prop), "CoerceType")
	// objectTyp := g.lang.PublicType(getObjectName(prop), "Object")
	// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFormatter")
	// g.lang.Class(g.bodyWriter, "export ",
	// 	g.lang.Implements(cname, g.lang.Generics("WuestenFormatter", typ, coerceTyp, objectTyp)),
	// 	prop,
	// 	func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {},
	// 	func(wr *eg.ForIfWhileLangWriter) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeBase")
	fpArgs := g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeBase", acp.coerceTyp))
	wr.WriteLine(g.lang.Readonly(fpArgs))
	g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult")
	// wr.WriteLine(g.lang.AssignDefault("_valueError", "WuesteResult.Ok(0)"))
	for _, pi := range prop.Items() {
		// wr.FormatLine("readonly %s = %s;", g.lang.PrivateName(prop.Name()), g.genWuesteBuilderAttribute(prop.Name(), prop.Property()))
		if pi.Property().Type() == eg.OBJECT /* && isNamedType(pi.Property()) */ {
			// piAttrClassName := g.lang.PublicName(getObjectName(pi.Property()), "Builder")
			// g.includes.AddProperty(piAttrClassName, pi.Property())
			// if pi.Optional() {
			// 	coerceType := g.lang.PublicName(getObjectName(pi.Property()), "CoerceType")
			// 	g.includes.AddProperty(coerceType, pi.Property())
			// 	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
			// 	wr.WriteLine(
			// 		g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
			// 			g.lang.Generics("WuestenObjectOptional",
			// 				piAttrClassName,
			// 				// g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()),
			// 				g.lang.PublicName(getObjectName(pi.Property())),
			// 				g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
			// 			))))
			// } else {
			// if pi.Optional() {
			// 	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenBuilder")
			// 	wr.WriteLine(
			// 		g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), g.lang.Generics("WuestenBuilder",
			// 			g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()), WithAddType(handleAddType(g)), WithAddInputType()),
			// 			g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional()), WithAddType(handleAddType(g))),
			// 			g.lang.AsTypeNullable(pi.Property(), WithAddObject(), WithOptional(pi.Optional()), WithAddType(handleAddType(g))),
			// 		))))
			// } else {
			wr.WriteLine(
				g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), g.propertyFactory(pi).builderTyp)))
			// }
			// wr.WriteLine(
			// g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), piAttrClassName)))
			// }
		} else if pi.Property().Type() == eg.ARRAY {
			t := g.lang.AsTypeNullable(pi.Property(), WithAddInputType(), WithOptional(pi.Optional()))
			c := g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional()))
			o := g.lang.AsTypeNullable(pi.Property(), WithAddObject(), WithOptional(pi.Optional()))
			it := g.lang.AsTypeNullable(getArrayItemProperty(pi.Property()), WithAddInputType(), WithOptional(pi.Optional()))
			ic := g.lang.AsTypeNullable(getArrayItemProperty(pi.Property()), WithAddCoerce(), WithOptional(pi.Optional()))
			io := g.lang.AsTypeNullable(getArrayItemProperty(pi.Property()), WithAddObject(), WithOptional(pi.Optional()))
			if pi.Optional() {
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenArrayBuilderOptional")
				wr.WriteLine(
					g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
						g.lang.Generics("WuestenArrayBuilderOptional", t, c, o, it, ic, io))))
				// piAttrClassName := g.lang.PublicName(name, "Builder")
				// coerceType := g.lang.PublicName(getObjectName(pi.Property()), "CoerceType")
				// g.includes.AddProperty(coerceType, pi.Property())
				// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
				// wr.WriteLine(
				// 	g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
				// 		g.lang.Generics("WuestenObjectOptional",
				// 			piAttrClassName,
				// 			// g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()),
				// 			g.lang.AsType(pi.Property()),
				// 			g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
				// 		))))
			} else {
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenArrayBuilder")
				wr.WriteLine(
					g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
						g.lang.Generics("WuestenArrayBuilder", t, c, o, it, ic, io))))
			}
		} else {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenBuilder")
			wr.WriteLine(
				g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), g.lang.Generics("WuestenBuilder",
					g.lang.AsTypeNullable(pi.Property(), WithAddInputType(), WithOptional(pi.Optional())),
					g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
					g.lang.AsTypeNullable(pi.Property(), WithAddObject(), WithOptional(pi.Optional())),
				))))
		}
	}

	// constructor(factory: WuestenFormatFactory, param: WuestenAttributeParameter<SimpleType$PayloadCoerceType>) {
	// 	this.param = WuestenMergeAttributeBase(factory, param)

	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenMergeAttributeBase")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeBase")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter")
	wr.WriteBlock("", g.lang.Call("constructor",
		g.lang.ReturnType("factory", g.lang.Generics("WuestenAttributeBase", "unknown")),
		g.lang.ReturnType("...params", g.lang.Index(
			g.lang.RoundBrackets(
				g.lang.OrType(g.lang.Generics("WuestenAttributeParameter", "unknown"), "undefined")), ""))),
		func(wr *eg.ForIfWhileLangWriter) {
			if len(prop.Items()) > 0 {
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeName")
				wr.WriteLine(g.lang.AssignDefault(g.lang.CallDot("this", "param"),
					g.lang.Call(g.lang.Generics("WuestenMergeAttributeBase", acp.coerceTyp), "factory", "...params")))
				wr.WriteLine(g.lang.AssignDefault(g.lang.Const("baseName"), g.lang.Call("WuestenAttributeName", "this.param")))
				for _, pi := range prop.Items() {
					params := genDefaultWuestenAttribute(g, pi.Name(), pi.Property())
					pin := g.propertyFactory(pi)
					wr.WriteLine(g.lang.AssignDefault(
						g.lang.CallDot("this", g.lang.PrivateName(pi.Name())),
						g.lang.Call(pin.builderCall, params...)))
				}
			}
		})

	// wr.WriteLine("// eslint-disable-next-line @typescript-eslint/no-unused-vars")
	// wr.WriteBlock("", g.lang.ReturnType(g.lang.Call("SetNameSuffix",
	// 	g.lang.ReturnType("...idxs", "number[]")), "void"), func(wr *eg.ForIfWhileLangWriter) {
	// 	wr.WriteLine("throw new Error(\"SetNameSuffix:Method not implemented.\")")
	// })
	// wr.WriteLine("// eslint-disable-next-line @typescript-eslint/no-unused-vars")
	// wr.WriteBlock("", g.lang.ReturnType(g.lang.Call("CoerceAttribute",
	// 	g.lang.ReturnType("val", "unknown")), g.lang.Generics("WuesteResult",
	// 	g.lang.PublicName(getObjectName(prop)), "Error")), func(wr *eg.ForIfWhileLangWriter) {
	// 	wr.WriteLine("throw new Error(\"CoerceAttribute:Method not implemented.\")")
	// })

	wr.WriteBlock("", g.lang.ReturnType(
		g.lang.Call("CoerceAttribute",
			g.lang.ReturnType("value", "unknown")), g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
		func(wr *eg.ForIfWhileLangWriter) {

			wr.WriteIf(g.lang.RoundBrackets("typeof value !== 'object' || value === null"), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("return WuesteResult.Err(Error(`Attribute[${WuestenAttributeName(this.param)}] is required`))")
			})
			wr.WriteLine(
				g.lang.Const(
					g.lang.AssignDefault("jsVal",
						g.lang.Index(g.lang.Cast("value", "Record<string, unknown>"), "this.param.jsonname")+
							" || "+
							g.lang.Index(g.lang.Cast("value", "Record<string, unknown>"), "this.param.varname"))))
			wr.WriteIf(g.lang.RoundBrackets("typeof jsVal !== 'object' || jsVal === null"), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("return WuesteResult.Err(Error(`Attribute[${WuestenAttributeName(this.param)}] is required`));")
			})
			wr.WriteLine("return this.Coerce(jsVal).Get()")
		})

	wr.WriteBlock("", g.lang.ReturnType(
		g.lang.Call("Coerce",
			g.lang.ReturnType("value?", acp.coerceTyp)),
		g.lang.PublicName(getObjectName(prop), "Builder")),
		func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock("if", "(!(typeof value === 'object' && value !== null))", func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return this")
			})
			for _, pi := range prop.Items() {
				wr.FormatLine("this.%s.CoerceAttribute(value)", g.lang.PrivateName(pi.Name()))
			}
			wr.WriteLine("return this")
		})

	wr.WriteBlock("",
		g.lang.ReturnType(
			g.lang.Call("Get", ""),
			g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
		func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock("return", g.lang.CallDot(acp.attrsClassName, "_fromResults"), func(wr *eg.ForIfWhileLangWriter) {
				for _, pi := range prop.Items() {
					wr.FormatLine("%s: this.%s.Get(),", g.lang.PrivateName(pi.Name()), g.lang.PrivateName(pi.Name()))
				}
			}, "({", "});")
		})

	// ToObject(value: T): O;
	wr.WriteBlock("", g.lang.ReturnType(
		g.lang.Call("ToObject"), //, g.lang.ReturnType("value", typ)),

		g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop), "Object"))),
		func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine("const errors: string[] = [];")
			for _, pi := range prop.Items() {
				r_ := "r_" + g.lang.PrivateName(pi.Name())
				wr.WriteLine(g.lang.Const(g.lang.AssignDefault(
					r_, g.lang.Call(
						g.lang.CallDot("this", g.lang.PrivateName(pi.Name()), "ToObject")))))
				wr.WriteBlock("if", "("+r_+".is_err())", func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("errors.push(%s.unwrap_err().message);", r_)
				})
			}
			// wr.FormatLine("const errors = Object.values(results).filter(r => r.is_err()).map(r => r.unwrap_err().message)")
			wr.WriteBlock("if", "(errors.length > 0)", func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return WuesteResult.Err(Error(errors.join('\\n')));")
			})
			wr.WriteBlock("return",
				g.lang.Generics("WuesteResult.Ok", g.lang.PublicName(getObjectName(prop), "Object")),
				func(wr *eg.ForIfWhileLangWriter) {
					for _, pi := range prop.Items() {
						r_ := "r_" + g.lang.PrivateName(pi.Name())
						wr.WriteLine(g.lang.Comma(g.lang.ReturnType(
							g.lang.Quote(pi.Name()),
							g.lang.Call(g.lang.CallDot(r_, "unwrap")))))
					}
				}, "({", "});")
		})

	wr.WriteBlock("static ",
		g.lang.ReturnType(
			g.lang.Call("_fromResults", g.lang.ReturnType("results", acp.resultsClassName)),
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
	// })
	// return attrsClassName
}

func (g *tsGenerator) generateFunctionHandler(wr *eg.ForIfWhileLangWriter, pi eg.PropertyItem) {
	wr.WriteIf(g.lang.RoundBrackets("typeof v === 'function'"), func(wr *eg.ForIfWhileLangWriter) {
		switch pi.Property().Type() {
		case eg.STRING, eg.BOOLEAN, eg.NUMBER, eg.INTEGER:
			wr.WriteLine(g.lang.Const(
				g.lang.AssignDefault("val",
					g.lang.Call(
						g.lang.CallDot(
							g.lang.CallDot("this", g.lang.PrivateName(pi.Name())),
							"Get")))))
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("ret", g.lang.Call("v",
				g.lang.Trinary("val.is_ok()", "val.unwrap()", "undefined")))))
		case eg.ARRAY:
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("rval", g.lang.Call(g.lang.CallDot("this", g.lang.PrivateName(pi.Name()), "Get")))))
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("ret", g.lang.Call("v",
				g.lang.Trinary("rval.isOk()", "rval.unwrap()", "undefined")))))
		case eg.OBJECT:
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("ret", g.lang.Call("v",
				g.lang.CallDot("this", g.lang.PrivateName(pi.Name()))))))
		default:
			panic("what is this?")
		}
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenRetValType")
		wr.WriteIf("(!(ret instanceof WuestenRetValType))", func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine("return this")
		})
		wr.WriteLine(g.lang.AssignDefault("v", fmt.Sprintf("ret.Val as %s ", g.lang.AsType(pi.Property(), WithAddInputType()))))
	})
}

func (g *tsGenerator) generateBuilder(prop eg.PropertyObject) {
	// for _, pi := range prop.Items() {
	// 	pa, ok := pi.Property().(eg.PropertyArray)
	// 	if ok {
	// 		g.generateLocalArrays(prop, pa, pi)
	// 	}
	// }

	resultsClassName := g.lang.PublicName(getObjectName(prop), "Results")

	g.generateSchemaExport(prop, getObjectName(prop))

	g.lang.Interface(g.bodyWriter, "", resultsClassName, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine(g.lang.Readonly(
				g.lang.ReturnType(
					g.lang.PrivateName(pi.Name()),
					g.lang.Generics("WuesteResult", g.lang.AsTypeNullable(pi.Property(), WithAddInputType(), WithOptional(pi.Optional()))))))
		})

	className := g.lang.PublicName(getObjectName(prop), "Builder")
	// extends WuestenAttr<Sub, Partial<Sub>|Partial<SubParam>|Partial<SubObject>>
	// implements WuestenBuilder<Sub, Partial<Sub>|Partial<SubParam>|Partial<SubObject>>

	partialType := g.lang.OrType(
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Param")),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")))
	coerceType := g.lang.PublicType(getObjectName(prop), "CoerceType")
	// fullCoerceType := g.lang.OrType(
	// 	coerceType,
	// 	g.lang.PublicType(getObjectName(prop), "Object"),
	// )
	typ := g.lang.PublicType(getObjectName(prop))
	g.bodyWriter.FormatLine("export type %s = %s", coerceType, partialType)
	objectType := g.lang.PublicType(getObjectName(prop), "Object")

	// type SX = (b: SimpleType$PayloadBuilder)=>void
	// g.bodyWriter.FormatLine("export type %s = (%s) => void",
	// 	g.lang.PublicType(getObjectName(prop), "FnGetBuilder"),
	// 	g.lang.ReturnType("b", g.lang.PublicType(getObjectName(prop), "Builder")))

	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenBuilder")
	g.lang.Class(g.bodyWriter, "export ", g.lang.Implements(className,
		g.lang.Generics("WuestenBuilder", typ, coerceType, objectType)), prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {

			pin := g.propertyFactory(pi)

			// WuestenFormatter.Any.Builder, WuestenFormatter.Any.CoerceType
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFNGetBuilder")
			paramTyp := g.lang.OrType(pin.coerceTyp, g.lang.Generics("WuestenFNGetBuilder", pin.fnGetBuilderTypes...))
			wr.WriteBlock(
				g.lang.ReturnType(
					g.lang.Call(g.lang.Type(g.lang.PublicType(pi.Name()), false),
						g.lang.ReturnType(g.lang.OptionalParam("v", pi.Optional()),
							paramTyp)),
					className),
				"", func(wr *eg.ForIfWhileLangWriter) {
					//if pi.Property().Type() == eg.OBJECT && isNamedType(pi.Property()) {
					// }
					g.generateFunctionHandler(wr, pi)
					wr.FormatLine("this.%s.Coerce(v);", g.lang.PrivateName(pi.Name()))
					wr.FormatLine("return this;")
				})
			// if pi.Property().Type() == eg.OBJECT && isNamedType(pi.Property()) {
			// 	retType := g.lang.PublicName(getObjectName(pi.Property()), "Builder")
			// 	if pi.Optional() {
			// 		retType = g.lang.Generics("WuestenObjectOptional",
			// 			// g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()),
			// 			g.lang.PublicName(getObjectName(pi.Property()), "Builder"),
			// 			g.lang.PublicName(getObjectName(pi.Property())),
			// 			g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
			// 		)
			// 	}
			// 	wr.WriteBlock(
			// 		g.lang.ReturnType(
			// 			g.lang.Call(g.lang.Type(g.lang.PublicType(pi.Name(), "Builder"), false)), retType),
			// 		"", func(wr *eg.ForIfWhileLangWriter) {
			// 			wr.FormatLine("return this._attr.%s;", g.lang.PrivateName(pi.Name()))
			// 		})

			// }
			// if pi.Property().Type() == eg.ARRAY {
			// 	retType := g.lang.PublicType(getObjectName(pi.Property(),
			// 		[]string{g.lang.PublicName(pi.Name())}), "Builder")
			// 	if pi.Optional() {
			// 		retType = g.lang.Generics("WuestenObjectOptional", retType,
			// 			g.lang.AsType(pi.Property()),
			// 			g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())))
			// 		// readonly _opt_arrayarrayBool: WuestenObjectOptional<NestedType$opt_arrayarrayBoolBuilder,
			// 		// boolean[][][][], WuesteCoerceTypeboolean[][][][]|undefined>
			// 	}
			// 	wr.WriteBlock(
			// 		g.lang.ReturnType(
			// 			g.lang.Call(g.lang.Type(g.lang.PublicType(pi.Name(), "Builder"), false)), retType),
			// 		"", func(wr *eg.ForIfWhileLangWriter) {
			// 			wr.FormatLine("return this._attr.%s;", g.lang.PrivateName(pi.Name()))
			// 		})
			// }
		}, func(wr *eg.ForIfWhileLangWriter) {

			g.generateAttributesClass(wr, prop, attributeClassParam{
				resultsClassName: resultsClassName,
				coerceTyp:        coerceType,
				objectTyp:        objectType,
				typ:              typ,
				attrsClassName:   className,
			})

			// g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result", "WuesteResult")
			// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeBase")

			// wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("_attr", attrsClassName)))
			// wr.WriteBlock("", g.lang.Call("constructor",
			// 	g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeBase", coerceType))),
			// 	func(wr *eg.ForIfWhileLangWriter) {
			// 		wr.WriteLine(g.lang.AssignDefault("this._attr", g.lang.New(attrsClassName, "param")))
			// 		// wr.WriteLine("super(param, attr);")
			// 		// wr.WriteLine(g.lang.AssignDefault("this._attr", "attr"))
			// 	})

			// wr.WriteBlock("", g.lang.ReturnType(
			// 	g.lang.Call("Reflection", ""), "WuestenReflection"), func(wr *eg.ForIfWhileLangWriter) {
			// 	wr.WriteLine(g.lang.Return(g.lang.Call(
			// 		g.lang.PublicName(getObjectName(prop), "Reflection"), "this")))
			// })
			// wr.WriteBlock("", g.lang.ReturnType(
			// 	g.lang.Call("Coerce", g.lang.ReturnType("value", coerceType)),
			// 	g.lang.PublicName(getObjectName(prop), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
			// 	wr.WriteLine(g.lang.Call("this._attr.Coerce", "value"))
			// 	wr.WriteLine(g.lang.Return("this"))
			// })
			// wr.WriteBlock("", g.lang.ReturnType(
			// 	g.lang.Call("ToObject", ""), g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop), "Object"))), func(wr *eg.ForIfWhileLangWriter) {
			// 	wr.WriteLine(g.lang.Return(g.lang.Call("this._attr.ToObject", "")))
			// })

			// name := g.lang.PublicName(getObjectName(prop))

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestePayload")
			wr.WriteBlock("",
				// this._attr.param.encoder
				g.lang.ReturnType(
					g.lang.Call("ToPayload", "encoder = this.param.encoder"), g.lang.Generics("WuesteResult", "WuestePayload")),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(g.lang.Const(g.lang.AssignDefault("robj", g.lang.Call(g.lang.CallDot("this", "ToObject")))))
					wr.WriteIf(g.lang.RoundBrackets("robj.is_err()"), func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteLine(g.lang.Return("WuesteResult.Err(robj.unwrap_err())"))
					})
					wr.WriteLine("const rdata = encoder(robj.unwrap())")
					wr.WriteBlock("if", "(rdata.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(rdata.unwrap_err());")
					})
					wr.WriteBlock("return", "WuesteResult.Ok", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("Type: %s,", g.lang.CallDot(g.lang.PrivateName(getObjectName(prop), "Names"), "id"))
						wr.FormatLine("Data: rdata.unwrap() as unknown as Record<string, unknown>")
					}, "({", "});")
				})

			// wr.WriteBlock("", g.lang.ReturnType(
			// 	g.lang.Call("Get", ""), g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))), func(wr *eg.ForIfWhileLangWriter) {
			// 	wr.WriteLine(g.lang.Return(g.lang.Call("this._attr.Get", "")))
			// })
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
	g.generateJSONDict(prop)
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
		myFname := getObjectFileName(prop)
		for _, externalTyp := range g.includes.ActiveTypes() {
			filename := externalTyp.tsFileName
			if filename == myFname {
				continue
			}
			// if externalTyp.activated && externalTyp.property.IsSome() {
			// 	filename = getObjectFileName(externalTyp.property.Value())
			// }
			// getObjectName(include.property)
			if len(externalTyp.Types()) <= 3 {
				header.FormatLine("import { %s } from %s;", strings.Join(externalTyp.Types(), ", "),
					g.lang.Quote(filename))
			} else {
				header.WriteBlock("", "import", func(wr *eg.ForIfWhileLangWriter) {
					for idx, t := range externalTyp.Types() {
						comma := ","
						if idx == len(externalTyp.Types())-1 {
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
	// toGenerate bool
	// activated  bool
	// prefix     string
	schemaFileName string
	tsFileName     string
	types          map[string]*string
	fileProperty   rusty.Optional[eg.Property]
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
	// if !o_ok {
	// 	parent := prop.Meta().Parent()
	// 	if parent.IsNone() {
	// 		panic("no parent")
	// 	}
	// 	prop = parent.Value()
	// 	if prop.Type() != eg.OBJECT {
	// 		panic("no parent")
	// 	}
	// 	// po = prop.(eg.PropertyObject)
	// 	o_ok = true
	// }
	fileName := getObjectFileName(prop)
	et, ok := g.types[fileName]
	if !ok {
		et = &externalType{
			// toGenerate:     true,
			schemaFileName: prop.Meta().FileName().Value(),
			tsFileName:     fileName,
			types:          make(map[string]*string),
		}
		g.types[fileName] = et
	}
	et.types[typ] = nil
	po, ok := prop.(eg.PropertyObject)
	if ok && et.fileProperty.IsNone() {
		// fmt.Printf("TOP-AddProperty: %s -> %s\n", typ, fileName)
		et.fileProperty = rusty.Some[eg.Property](po)
	} else {
		// fmt.Printf("AddProperty: %s -> %s\n", typ, fileName)
	}
}

func (g *externalTypes) AddType(fileName, typeName string, optAlias ...string) *externalType {
	var alias *string = nil
	if len(optAlias) > 0 {
		alias = &optAlias[0]
	}
	if fileName == "" {
		return nil
	}
	et, ok := g.types[fileName]
	if !ok {
		et = &externalType{
			// toGenerate: false,
			tsFileName: fileName,
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
		// if et.activated {
		atyp = append(atyp, et)
		// }
	}
	sort.Slice(atyp, func(i, j int) bool {
		return strings.Compare(atyp[i].tsFileName, atyp[j].tsFileName) < 0
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

	g.generatePropertyObject(prop.(eg.PropertyObject), sl)
	for _, at := range g.includes.ActiveTypes() {
		if at.fileProperty.IsSome() &&
			at.fileProperty.Value().Id() != prop.Id() {
			TsGenerator(cfg, at.fileProperty.Value(), sl)
		}
	}
}
