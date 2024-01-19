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

func (l *tsLang) addCoerceType(typ string, withs ...withResult) string {
	if hasWith(WithAddCoerce(), withs) {
		switch typ {
		case "number", "string", "boolean", "Date":
			typ = fmt.Sprintf("WuesteCoerceType%s", typ)
		}
	}
	res := hasWithResult(WithAddType(func(typ string, prop eg.Property) {}), withs)
	prop := hasWithResult(WithIsCoerceType(nil), withs)
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
			return l.Generics("Record", "string", "unknown")
		}
		name := getObjectName(p)
		if hasWith(WithAddCoerce(), withs) {
			return l.addCoerceType(l.Name(name, "CoerceType"),
				append(withs, WithIsCoerceType(p))...)
		}
		ret := l.Name(name)
		res := hasWithResult(WithAddType(func(typ string, prop eg.Property) {}), withs)
		if res != nil {
			if isNamedType(p) {
				res.addType(ret, p)
			}
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

func (l *tsLang) New(line string, retType ...string) string {
	return "new " + l.Call(line, retType...)
}

func (l *tsLang) Return(line string) string {
	return "return " + line
}

func (l *tsLang) Call(line string, params ...string) string {
	return line + "(" + strings.Join(params, ", ") + ")"
}

func (l *tsLang) Interface(
	wr *eg.ForIfWhileLangWriter,
	prefix, name string,
	prop eg.PropertyObject,
	itemFn func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter),
	onceFn ...func(wr *eg.ForIfWhileLangWriter)) string {
	return wr.WriteBlock(prefix+"interface ", name, func(wr *eg.ForIfWhileLangWriter) {
		for _, once := range onceFn {
			once(wr)
		}
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
	if strings.HasPrefix(prop, "[") {
		return str + prop
	}
	if prop == "" {
		return str
	}
	return str + "." + prop
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

func hasDefault(prop eg.Property) bool {
	return getDefaultForProperty(prop) != nil
}

// var reOrArray = regexp.MustCompile(`[(\|)(\[\])]+`)

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
		return true
	default:
	}
	return false
}

func (g *tsGenerator) generateClass(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop)), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		out := []string{}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()), pi.Optional())),
				g.lang.AsTypeNullable(pi.Property(),
					WithAddType(func(typ string, addProp eg.Property) {
						out = append(out, typ)
						if addProp == nil {
							g.includes.AddType(g.cfg.EntityCfg.FromWueste, typ)
						} else {
							g.includes.AddProperty(typ, addProp)
						}
					}),
				))))
		wr.FormatLine("// %v", out)

	})
	g.bodyWriter.WriteLine()

	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Param"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		if isNamedType(pi.Property()) {
			paramName := getObjectName(pi.Property())
			g.includes.AddProperty(paramName, pi.Property())
			// g.includes.AddProperty(g.lang.PublicName(paramName, "Param"), pi.Property())
		}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()),
					pi.Optional() || hasDefault(pi.Property()))),
				g.lang.AsTypeNullable(pi.Property(),
					WithAddType(func(typ string, prop eg.Property) {
						if prop == nil {
							g.includes.AddType(g.cfg.EntityCfg.FromWueste, typ)
						} else {
							g.includes.AddProperty(typ, prop)
						}
					}),
					WithAddCoerce()))))
	})
	g.bodyWriter.WriteLine()

}

func (g *tsGenerator) generateJSONDict(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Object"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		typ := g.lang.AsTypeNullable(pi.Property())
		if isNamedType(pi.Property()) {
			typ = g.lang.PublicName(g.lang.AsType(pi.Property()), "Object")
			g.includes.AddProperty(typ, pi.Property())
			// if pi.Optional() {
			// 	typ = g.lang.OrType(typ, "undefined")
			// }
		}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.Quote(pi.Name()), pi.Optional())),
				typ), g.lang.JsonAnnotation(pi.Property())))
	})
	g.bodyWriter.WriteLine()
}

func (g *tsGenerator) generatePayload(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(getObjectName(prop), "Payload"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
	}, func(wr *eg.ForIfWhileLangWriter) {
		names := g.getNames(prop)
		quotedNames := []string{}
		for _, name := range names.names {
			quotedNames = append(quotedNames, g.lang.Quote(name))
		}
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type("Type", false)),
				strings.Join(quotedNames, "|"))))
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type("Data", false)),
				g.lang.PublicType(getObjectName(prop), "Object"))))
	})
	g.bodyWriter.WriteLine()
}

type Names struct {
	title   string
	names   []string
	varname string
}

func (g *tsGenerator) getNames(prop eg.PropertyObject) Names {
	title := prop.Title()
	if title == "" {
		title = prop.Id()
	}
	names := []string{prop.Id()}
	if prop.Id() != prop.Title() {
		names = append(names, prop.Title())
	}
	varname := g.lang.PublicName(title)
	if varname != prop.Id() && varname != title {
		names = append(names, varname)
	}
	return Names{title: title, names: names, varname: varname}
}

func (g *tsGenerator) generateFactory(prop eg.PropertyObject) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFactory")

	// export function  NewSimpleTypeFactory(): WuestenFactory<SimpleTypeBuilder, SimpleType>
	className := g.lang.PublicName(getObjectName(prop), "FactoryImpl")
	partialType := g.lang.PublicName(getObjectName(prop), "CoerceType")
	// g.lang.OrType(
	// 	g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
	// 	g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Param")),
	// 	g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")))
	g.lang.Class(g.bodyWriter, "export ", g.lang.Implements(className,
		g.lang.Generics("WuestenFactory", g.lang.PublicName(getObjectName(prop)), partialType,
			g.lang.PublicName(getObjectName(prop), "Object"))), prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		}, func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine(g.lang.Readonly(
				g.lang.AssignDefault("T",
					fmt.Sprintf("undefined as unknown as %s", g.lang.PublicName(getObjectName(prop))))))

			wr.WriteLine(g.lang.Readonly(
				g.lang.AssignDefault("I",
					fmt.Sprintf("undefined as unknown as %s", partialType))))

			wr.WriteLine(g.lang.Readonly(
				g.lang.AssignDefault("O",
					fmt.Sprintf("undefined as unknown as %s", g.lang.PublicName(getObjectName(prop), "Object")))))

			wr.WriteBlock("Builder():", g.lang.PublicName(getObjectName(prop), "Builder"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return new %s()", g.lang.PublicName(getObjectName(prop), "Builder"))
			})
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteResult")
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("ToObject", g.lang.ReturnType("obj", g.lang.PublicName(getObjectName(prop)))),
				g.lang.PublicName(getObjectName(prop), "Object")), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return %s", g.lang.Call(g.lang.PublicName(getObjectName(prop), "ToObject"), "obj"))
			})

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenJSONPassThroughDecoder")
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenNames")

			wr.WriteBlock("", g.lang.ReturnType(g.lang.Call("Names"), "WuestenNames"), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteBlock("return", "", func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("id: %s,", g.lang.Quote(prop.Id()))
					names := g.getNames(prop)
					wr.FormatLine("title: %s,", g.lang.Quote(names.title))
					jsonBytes, _ := json.Marshal(names.names)
					wr.FormatLine("names: %s,", string(jsonBytes))
					wr.FormatLine("varname: %s", g.lang.Quote(names.varname))
				})
			})
			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("FromPayload", g.lang.ReturnType("val", g.lang.PublicName(getObjectName(prop), "Payload")), "decoder = WuestenJSONPassThroughDecoder"),
					g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("if", "(!this.Names().names.includes(val.Type))", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(new Error(`%s Type mismatch:[${this.Names().names.join(',')}] != ${val.Type}`));", g.lang.PublicName(getObjectName(prop), "Payload"))
					})
					// <Partial<SimpleTypeParam>>
					wr.FormatLine("const data = %s", g.lang.Call("decoder", "val.Data"))
					wr.WriteBlock("if", "(data.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteLine("return WuesteResult.Err(data.unwrap_err());")
					})

					wr.WriteLine(
						g.lang.AssignDefault("const builder",
							g.lang.New(g.lang.PublicName(getObjectName(prop), "Builder"))))
					wr.FormatLine("return builder.Coerce(data.unwrap() as %s);", partialType)
				})

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenJSONPassThroughEncoder")
			name := g.lang.PublicName(getObjectName(prop))
			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("ToPayload", g.lang.ReturnType("val", g.lang.OrType(g.lang.Generics("WuesteResult", name), name)),
						"encoder = WuestenJSONPassThroughEncoder"), g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop), "Payload"))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("let toEncode: %s;", name)
					wr.WriteIf("(WuesteResult.Is(val))", func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteIf("(val.is_err())", func(wr *eg.ForIfWhileLangWriter) {
							wr.FormatLine("return WuesteResult.Err(val.unwrap_err());")
						})
						wr.FormatLine("toEncode = val.unwrap();")
					}, func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("toEncode = val;")
					})
					wr.FormatLine("const data = encoder(%s.ToObject(toEncode))", g.lang.PublicName(getObjectName(prop), "Factory"))
					wr.WriteBlock("if", "(data.is_err())", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(data.unwrap_err());")
					})
					wr.WriteBlock("return", "WuesteResult.Ok", func(wr *eg.ForIfWhileLangWriter) {
						id := prop.Id()
						if id == "" {
							id = prop.Title()
						}
						wr.FormatLine("Type: %s,", g.lang.Quote(id))
						wr.FormatLine("Data: data.unwrap() as unknown as %s", g.lang.PublicName(getObjectName(prop), "Object"))
					}, "({", "});")
				})

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteResult")
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

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionObject")
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Schema"), "WuestenReflectionObject"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return %s as WuestenReflectionObject;", g.lang.PublicName(getObjectName(prop), "Schema"))
			})
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionValue")
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Getter", g.lang.ReturnType("typ", g.lang.PublicName(getObjectName(prop))),
					g.lang.ReturnType("base", "WuestenReflectionValue[] = []")), "WuestenGetterBuilder"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return %s", g.lang.Call(g.lang.PublicName(getObjectName(prop), "Getter"), "typ", "base"))
			})

		})

	g.bodyWriter.WriteLine(
		g.lang.AssignDefault(
			g.lang.Export(g.lang.Const(g.lang.PublicName(getObjectName(prop), "Factory"))),
			g.lang.New(g.lang.PublicName(getObjectName(prop), "FactoryImpl"))))
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenTypeRegistry")
	g.bodyWriter.WriteLine(g.lang.CallDot("WuestenTypeRegistry", g.lang.Call("Register", g.lang.PublicName(getObjectName(prop), "Factory"))))
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
	if prop.Description().IsSome() {
		wr.WriteLine(g.lang.Comma(g.lang.ReturnType("description", g.lang.Quote(prop.Description().Value()))))
	}
	if prop.XProperties() != nil {
		for xk, xv := range prop.XProperties() {
			jsonXk, err := json.Marshal(xk)
			if err != nil {
				panic(err)
			}
			jsonXv, err := json.Marshal(xv)
			if err != nil {
				panic(err)
			}

			wr.WriteLine(g.lang.Comma(g.lang.ReturnType(string(jsonXk), string(jsonXv))))
		}
	}

	switch prop.Type() {
	case eg.BOOLEAN:
	case eg.STRING, eg.INTEGER, eg.NUMBER:
		pi := prop.(eg.PropertyFormat)
		if pi.Format().IsSome() {
			wr.WriteLine(g.lang.Comma(g.lang.ReturnType("format", g.lang.Quote(pi.Format().Value()))))
		}
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
						if pi.Property().Type() == eg.OBJECT && isNamedType(pi.Property()) {
							reflection := g.lang.PublicName(getObjectName(pi.Property()), "Schema")
							g.includes.AddProperty(reflection, pi.Property())
							// builderName := g.lang.Call(g.lang.PublicType(pi.Name(), "Builder"))
							// if pi.Optional() {
							// 	builderName = g.lang.CallDot(builderName, "typ")
							// }
							wr.WriteLine(g.lang.ReturnType("property", reflection))
						} else if pi.Property().Type() == eg.ARRAY {
							reflection := g.lang.PublicName(getObjectName(pi.Property(), []string{pi.Name()}), "Schema")
							// builderName := g.lang.Call(g.lang.PublicType(pi.Name(), "Builder"))
							// if pi.Optional() {
							// 	builderName = g.lang.CallDot(builderName, "typ")
							// }
							wr.WriteLine(g.lang.ReturnType("property", reflection))
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
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten")
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
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten")
		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeIntegerOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeInteger", paramFn())
		}
	case eg.NUMBER:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten")

		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeNumberOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeNumber", paramFn())
		}
	case eg.BOOLEAN:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten")
		if pi.Optional() {
			return g.lang.Call("wuesten.AttributeBooleanOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeBoolean", paramFn())
		}
	case eg.ARRAY:
		baseName := g.lang.PublicName(getObjectName(pi.Property(), []string{name}), "Builder")
		if pi.Optional() {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
			// return g.lang.New(baseName, paramFn())
			return g.lang.New("WuestenObjectOptional", g.lang.New(baseName, paramFn()))
		} else {
			return g.lang.New(baseName, paramFn())
		}
	case eg.OBJECT:
		po := prop.(eg.PropertyObject)
		objName := getObjectName(po)
		if pi.Optional() {
			if !isNamedType(po) {
				// return g.lang.Call(g.lang.Generics("wuesten.AttributeObjectOptional", generics()...), paramFn(), factory)
				factory := "WuestenObjectFactory"
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, factory)
				return g.lang.Call("wuesten.AttributeObjectOptional", paramFn(), factory)
			} else {
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
				name := g.lang.PublicName(objName, "Builder")
				g.includes.AddProperty(name, po)
				return g.lang.New("WuestenObjectOptional",
					g.lang.New(name, paramFn()))
				// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttrOptional").activated = true
				// return g.lang.New("WuestenAttrOptional",
				// 	g.lang.New(g.lang.PublicName(objName, "Attributes"), paramFn()))
			}
		} else {
			if !isNamedType(po) {
				// return g.lang.Call(g.lang.Generics("wuesten.AttributeObject", generics()...), paramFn(), factory)
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten")
				factory := "WuestenObjectFactory"
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, factory)
				return g.lang.Call("wuesten.AttributeObject", paramFn(), factory)
			} else {
				name := g.lang.PublicName(objName, "Builder")
				g.includes.AddProperty(name, po)
				return g.lang.New(name, paramFn())
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
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteToIterator")
	resIter := fmt.Sprintf("r%d", level)
	wr.FormatLine(g.lang.AssignDefault(
		g.lang.Const(resIter),
		g.lang.Call(g.lang.Generics("WuesteToIterator", g.lang.AsType(prop.Items())), rootArray)))
	wr.WriteIf(g.lang.RoundBrackets(g.lang.CallDot(resIter, g.lang.Call("is_err"))), func(wr *eg.ForIfWhileLangWriter) {
		wr.FormatLine("return WuesteResult.Err(`it's not iterable on level %d:${%s}`)", level, g.lang.CallDot(resIter, g.lang.Call("unwrap_err")))
	})
	resultFmt := "s%d"
	result := fmt.Sprintf(resultFmt, level)
	wr.WriteLine(
		g.lang.AssignDefault(
			g.lang.Const(
				g.lang.ReturnType(result, g.lang.AsType(prop))),
			"[]"))
	iter := fmt.Sprintf("t%d", level)
	inc := fmt.Sprintf("i%d", level)
	wr.FormatLine("const %s = %s.unwrap()", iter, resIter)
	wr.FormatLine("let %s = %s.next()", inc, iter)
	wr.WriteBlock("for",
		g.lang.RoundBrackets(
			"; !"+g.lang.CallDot(inc, "done")+"; "+inc+" = "+g.lang.CallDot(iter, g.lang.Call("next"))),
		func(wr *eg.ForIfWhileLangWriter) {
			p, ok := prop.Items().(eg.PropertyArray)
			if ok {
				g.generateArrayCoerce(level+1, g.lang.CallDot(inc, "value"), returnType, p, wr)
				wr.WriteLine(g.lang.CallDot(result, g.lang.Call("push", fmt.Sprintf(resultFmt, level+1))))
			} else {
				// param := []string{}
				// for i := 0; i <= level; i++ {
				// 	param = append(param, fmt.Sprintf("c%d", i))
				// }
				// wr.WriteLine(g.lang.Call("itemAttr.SetNameSuffix", strings.Join(param, ", ")))
				wr.WriteLine(g.lang.AssignDefault(
					g.lang.Const("attrRes"), g.lang.Call("itemAttr.Coerce", g.lang.CallDot(inc, "value"))))
				wr.WriteIf(g.lang.RoundBrackets("attrRes.is_err()"), func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(g.lang.Return(
						g.lang.Generics("attrRes as unknown as WuesteResult", returnType)))
				})
				wr.WriteLine(g.lang.CallDot(result, g.lang.Call("push", g.lang.Call("attrRes.unwrap"))))
			}
		})
}

func (g *tsGenerator) generateLocalArrays(prop eg.PropertyObject, pa eg.PropertyArray, pi eg.PropertyItem) {
	baseName := getObjectName(pi.Property(), []string{pi.Name()})
	// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionArray")

	g.generateSchemaExport(pi.Property(), baseName)
	// g.bodyWriter.WriteLine("// eslint-disable-next-line @typescript-eslint/no-unused-vars")
	// g.bodyWriter.WriteBlock("function ", g.lang.ReturnType(
	// 	g.lang.Call(g.lang.PublicName(baseName, "Reflection"),
	// 		g.lang.ReturnType("param", g.lang.PublicName(baseName, "Builder"))), "WuestenReflection"),
	// 	func(wr *eg.ForIfWhileLangWriter) {
	// 		wr.WriteBlock("return", "", func(wr *eg.ForIfWhileLangWriter) {
	// 			g.writeSchema(wr, pi.Property())
	// 		})
	// 	})
	className := g.lang.PublicName(baseName, "Builder")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttr")
	g.lang.Class(g.bodyWriter, "export ", g.lang.Extends(className,
		g.lang.Generics("WuestenAttr",
			// g.lang.AsType(pi.Property()),
			g.lang.AsTypeNullable(pi.Property() /*WithOptional(pi.Optional())*/),
			g.lang.AsTypeNullable(pi.Property(), WithAddCoerce() /*WithOptional(pi.Optional())*/))),
		prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {},
		func(wr *eg.ForIfWhileLangWriter) {
			attrib := g.lang.ReturnType(g.lang.OptionalParam(g.lang.PrivateName("value"), pi.Optional()), g.lang.AsType(pi.Property()))
			if !pi.Optional() {
				attrib += " = []"
			}
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter")

			wr.WriteBlock("",
				g.lang.Call("constructor",
					g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeParameter", g.lang.AsType(getItemType(pa)))),
				), func(wr *eg.ForIfWhileLangWriter) {
					pi := eg.NewPropertyArrayItem("ARRAY", rusty.Ok(getItemType(pa)), false).Ok()
					attr := g.genWuesteBuilderAttribute("ARRAY", pi, func() string { return "param" })
					wr.WriteLine(g.lang.AssignDefault(g.lang.Const("itemAttr"), attr))

					wr.WriteBlock("", "super({jsonname: param.jsonname, varname: param.varname, base: param.base}, {coerce: (c0: unknown) => ", func(wr *eg.ForIfWhileLangWriter) {
						g.generateArrayCoerce(0, "c0", g.lang.AsType(pa), pa, wr)
						wr.WriteLine(g.lang.Return(g.lang.Call("WuesteResult.Ok", "s0")))
					}, " {", "}})")
				})
		})
	g.bodyWriter.WriteLine()
}

func (g *tsGenerator) generateSchemaExport(prop eg.Property, baseName string) {
	g.generateReflectionGetter(propertyValue{prop: prop}, baseName)
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflection")
	g.bodyWriter.WriteBlock(g.lang.Export(g.lang.Const(g.lang.ReturnType(
		g.lang.PublicName(baseName, "Schema"),
		"WuestenReflection"))), "", func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteLine(g.lang.Comma(g.lang.ReturnType(
			"ref", g.lang.Quote(g.lang.PublicName(baseName)))))
		g.writeSchema(wr, prop)
	}, " = {")
	g.generateToObject(prop, baseName)
}

func (g *tsGenerator) generateAttributesClass(prop eg.PropertyObject, resultsClassName string) string {
	attrsClassName := g.lang.PublicName(getObjectName(prop), "Attributes")
	cname := attrsClassName
	g.lang.Class(g.bodyWriter, "export ", cname, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {},
		func(wr *eg.ForIfWhileLangWriter) {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter")
			wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("param",
				g.lang.Generics("WuestenAttributeParameter", g.lang.PublicName(getObjectName(prop))))))
			for _, pi := range prop.Items() {
				// wr.FormatLine("readonly %s = %s;", g.lang.PrivateName(prop.Name()), g.genWuesteBuilderAttribute(prop.Name(), prop.Property()))
				if pi.Property().Type() == eg.OBJECT && isNamedType(pi.Property()) {
					piAttrClassName := g.lang.PublicName(getObjectName(pi.Property()), "Builder")
					g.includes.AddProperty(piAttrClassName, pi.Property())
					if pi.Optional() {
						coerceType := g.lang.PublicName(getObjectName(pi.Property()), "CoerceType")
						g.includes.AddProperty(coerceType, pi.Property())
						g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
								g.lang.Generics("WuestenObjectOptional",
									piAttrClassName,
									// g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()),
									g.lang.PublicName(getObjectName(pi.Property())),
									g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
								))))
					} else {
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()), piAttrClassName)))
					}
				} else if pi.Property().Type() == eg.ARRAY {
					name := getObjectName(pi.Property(), []string{pi.Name()})
					if pi.Optional() {
						piAttrClassName := g.lang.PublicName(name, "Builder")
						coerceType := g.lang.PublicName(getObjectName(pi.Property()), "CoerceType")
						g.includes.AddProperty(coerceType, pi.Property())
						g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
								g.lang.Generics("WuestenObjectOptional",
									piAttrClassName,
									// g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()),
									g.lang.AsType(pi.Property()),
									g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
								))))
					} else {
						wr.WriteLine(
							g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(pi.Name()),
								g.lang.PublicName(name, "Builder"))))
					}
				} else {
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttribute")
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
				// wr.WriteLine(
				// 	g.lang.Call("super", "param", g.lang.CurlyBrackets(
				// 		g.lang.ReturnType("coerce", g.lang.CallDot("(p) => "+attrsClassName,
				// 			g.lang.Call("_coerce", "this", "p"))))))
				if len(prop.Items()) > 0 {
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeName")
					wr.WriteLine(g.lang.AssignDefault(g.lang.CallDot("this", "param"), "param"))
					wr.WriteLine(g.lang.AssignDefault(g.lang.Const("baseName"), g.lang.Call("WuestenAttributeName", "param")))
					for _, pi := range prop.Items() {
						wr.WriteLine(g.lang.AssignDefault(
							g.lang.CallDot("this", g.lang.PrivateName(pi.Name())), g.genWuesteBuilderAttribute(pi.Name(), pi)))
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

			wr.WriteBlock("readonly ", g.lang.ReturnType(
				g.lang.Call("Coerce = ",
					// g.lang.ReturnType("bound", attrsClassName),
					g.lang.ReturnType("value", "unknown")),
				g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("if", "(!(typeof value === 'object' && value !== null))", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return WuesteResult.Err(Error('expected object'));")
					})
					wr.WriteBlock("return", g.lang.CallDot(attrsClassName, "_fromResults"), func(wr *eg.ForIfWhileLangWriter) {
						for _, pi := range prop.Items() {
							wr.FormatLine("%s: this.%s.CoerceAttribute(value),", g.lang.PrivateName(pi.Name()), g.lang.PrivateName(pi.Name()))
						}
					}, "({", "});")
				}, " => {", "}")

			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("Get", ""),
					g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("return", g.lang.CallDot(attrsClassName, "_fromResults"), func(wr *eg.ForIfWhileLangWriter) {
						for _, pi := range prop.Items() {
							wr.FormatLine("%s: this.%s.Get(),", g.lang.PrivateName(pi.Name()), g.lang.PrivateName(pi.Name()))
						}
					}, "({", "});")
				})

			wr.WriteBlock("static ",
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
	return attrsClassName
}

func (g *tsGenerator) generateFunctionHandler(wr *eg.ForIfWhileLangWriter, pi eg.PropertyItem) {
	wr.WriteIf(g.lang.RoundBrackets("typeof v === 'function'"), func(wr *eg.ForIfWhileLangWriter) {
		switch pi.Property().Type() {
		case eg.STRING, eg.BOOLEAN, eg.NUMBER, eg.INTEGER:
			wr.WriteLine(g.lang.Const(
				g.lang.AssignDefault("val",
					g.lang.Call(
						g.lang.CallDot(
							g.lang.CallDot("this._attr", g.lang.PrivateName(pi.Name())),
							"Get")))))
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("ret", g.lang.Call("v",
				g.lang.Trinary("val.is_ok()", "val.unwrap()", "undefined")))))
		case eg.ARRAY:
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("ret", g.lang.Call("v",
				g.lang.CallDot("this._attr", g.lang.PrivateName(pi.Name()))))))
		case eg.OBJECT:
			wr.WriteLine(g.lang.Const(g.lang.AssignDefault("ret", g.lang.Call("v",
				g.lang.CallDot("this._attr", g.lang.PrivateName(pi.Name()))))))
		default:
			panic("what is this?")
		}
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenRetValType")
		wr.WriteIf("(!(ret instanceof WuestenRetValType))", func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine("return this")
		})
		wr.WriteLine(g.lang.AssignDefault("v", fmt.Sprintf("ret.Val as %s ", g.lang.AsType(pi.Property()))))
	})
}

func (g *tsGenerator) generateBuilder(prop eg.PropertyObject) {

	for _, pi := range prop.Items() {
		pa, ok := pi.Property().(eg.PropertyArray)
		if ok {
			g.generateLocalArrays(prop, pa, pi)
		}
	}

	resultsClassName := g.lang.PublicName(getObjectName(prop), "Results")

	g.generateSchemaExport(prop, getObjectName(prop))

	g.lang.Interface(g.bodyWriter, "", resultsClassName, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine(g.lang.Readonly(
				g.lang.ReturnType(
					g.lang.PrivateName(pi.Name()),
					g.lang.Generics("WuesteResult", g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Optional()))))))
		})

	attrsClassName := g.generateAttributesClass(prop, resultsClassName)

	className := g.lang.PublicName(getObjectName(prop), "Builder")
	// extends WuestenAttr<Sub, Partial<Sub>|Partial<SubParam>|Partial<SubObject>>
	// implements WuestenBuilder<Sub, Partial<Sub>|Partial<SubParam>|Partial<SubObject>>

	partialType := g.lang.OrType(
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop))),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Param")),
		g.lang.Generics("Partial", g.lang.PublicName(getObjectName(prop), "Object")))
	coerceType := g.lang.PublicType(getObjectName(prop), "CoerceType")
	genericType := []string{g.lang.PublicName(getObjectName(prop)), coerceType}
	g.bodyWriter.FormatLine("export type %s = %s", coerceType, partialType)

	// type SX = (b: SimpleType$PayloadBuilder)=>void
	// g.bodyWriter.FormatLine("export type %s = (%s) => void",
	// 	g.lang.PublicType(getObjectName(prop), "FnGetBuilder"),
	// 	g.lang.ReturnType("b", g.lang.PublicType(getObjectName(prop), "Builder")))

	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenBuilder")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttr")
	g.lang.Class(g.bodyWriter, "export ", g.lang.Implements(
		g.lang.Extends(className, g.lang.Generics("WuestenAttr", genericType...)),
		g.lang.Generics("WuestenBuilder", genericType[0], genericType[1])), prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			paramTyp := g.lang.Type(g.lang.AsType(pi.Property(), WithAddCoerce()), false)
			fnGetBuilderType := paramTyp
			if pi.Property().Type() == eg.ARRAY || pi.Property().Type() == eg.OBJECT {
				if pi.Property().Type() == eg.ARRAY || isNamedType(pi.Property()) {
					baseName := getObjectName(pi.Property())
					typeName := g.lang.PublicName(baseName)
					if pi.Property().Type() == eg.ARRAY {
						baseName = getObjectName(pi.Property(), []string{pi.Name()})
						typeName = g.lang.AsType(pi.Property(), WithAddCoerce())
					}
					if pi.Optional() {
						// WuestenObjectOptional<SimpleType$Payload, SimpleType$PayloadCoerceType|undefined>
						g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenObjectOptional")
						fnGetBuilderType =
							g.lang.Generics("WuestenObjectOptional",
								g.lang.PublicName(baseName, "Builder"),
								typeName,
								g.lang.AsTypeNullable(pi.Property(), WithAddCoerce(), WithOptional(pi.Optional())),
							)
					} else {
						fnGetBuilderType = g.lang.PublicName(baseName, "Builder")
					}
				} else if pi.Property().Type() == eg.OBJECT {
					recordType := g.lang.Generics("Record", "string", "unknown")
					if pi.Optional() {
						recordType = g.lang.OrType(recordType, "undefined")
					}
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttribute")
					fnGetBuilderType = g.lang.Generics("WuestenAttribute", recordType, recordType)
				}
				// paramTyp = g.lang.OrType(g.lang.AsType(pi.Property(), WithAddCoerce()), paramTyp)
			}
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFNGetBuilder")
			paramTyp = g.lang.OrType(paramTyp, g.lang.Generics("WuestenFNGetBuilder", fnGetBuilderType))
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
					wr.FormatLine("this._attr.%s.Coerce(v);", g.lang.PrivateName(pi.Name()))
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

			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuesteResult")

			wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("_attr", attrsClassName)))
			wr.WriteBlock("", g.lang.Call("constructor",
				g.lang.AssignDefault("param",
					fmt.Sprintf("{jsonname: %s, varname: %s, base: \"\"}",
						g.lang.Quote(getObjectName(prop)), g.lang.Quote(g.lang.PublicName(getObjectName(prop)))))), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Const(g.lang.AssignDefault("attr", g.lang.New(attrsClassName, "param"))))
				wr.WriteLine("super(param, {coerce: attr.Coerce});")
				wr.WriteLine(g.lang.AssignDefault("this._attr", "attr"))
			})

			// wr.WriteBlock("", g.lang.ReturnType(
			// 	g.lang.Call("Reflection", ""), "WuestenReflection"), func(wr *eg.ForIfWhileLangWriter) {
			// 	wr.WriteLine(g.lang.Return(g.lang.Call(
			// 		g.lang.PublicName(getObjectName(prop), "Reflection"), "this")))
			// })

			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Get", ""), g.lang.Generics("WuesteResult", g.lang.PublicName(getObjectName(prop)))), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Return(g.lang.Call("this._attr.Get", "")))
			})
			// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestePayload")
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
	g.generatePayload(prop)
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
