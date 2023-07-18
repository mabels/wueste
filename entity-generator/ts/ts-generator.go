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

var reSplitNonAllowed = regexp.MustCompile("[^a-zA-Z0-9_]+")

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

func (l *tsLang) AsType(p eg.Property, withs ...string) string {
	switch p.Type() {
	case eg.OBJECT:
		if hasWith(WithAddCoerce(), withs) {
			return l.OrType(l.Name(p.(eg.PropertyObject).Title()), l.Name(p.(eg.PropertyObject).Title(), "Param"))
		}
		return l.Name(p.(eg.PropertyObject).Title())
	case eg.STRING:
		p := p.(eg.PropertyString)
		if p.Format().IsSome() {
			switch *p.Format().Value() {
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
		return l.AsType(p.(eg.PropertyArray).Items(), withs...) + "[]"
	default:
		panic(fmt.Sprintf("unknown type %s", p.Type()))
	}
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

func (l *tsLang) AsTypeNullable(p eg.Property, withs ...string) string {
	if hasWith("isOptional", withs) {
		return l.AsType(p, withs...) + "|undefined"
	}
	return l.AsType(p, withs...)
}

func (l *tsLang) Line(line string, tails ...string) string {
	tail := ""
	if len(tails) > 0 {
		tail = tails[0]
	}
	return line + ";" + tail
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
		for _, pi := range prop.Properties().Items() {
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
		for _, prop := range prop.Properties().Items() {
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

func (g *tsGenerator) generateClass(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(prop.Title()), prop, func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(prop.Property())), g.lang.AsType(prop.Property()), prop.Property())
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(prop.Name()), prop.Property().Optional())),
				g.lang.AsTypeNullable(prop.Property()))))
	})
	g.bodyWriter.WriteLine()

	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(prop.Title(), "Param"), prop, func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(pi.Property())), g.lang.AsType(pi.Property()), pi.Property())
		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(pi.Property())), g.lang.PublicName(g.lang.AsType(pi.Property()), "Param"), pi.Property())
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.PublicType(pi.Name()),
					pi.Property().Optional() || hasDefault(pi.Property()))),
				g.addWuestenType(g.lang.AsTypeNullable(pi.Property(), WithAddCoerce())))))
	})
	g.bodyWriter.WriteLine()

}

func (g *tsGenerator) generateJson(prop eg.PropertyObject) {
	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(prop.Title(), "Object"), prop, func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(prop.Property())), g.lang.AsType(prop.Property()), prop.Property())
		wr.WriteLine(g.lang.Line(
			g.lang.ReturnType(
				g.lang.Readonly(g.lang.Type(g.lang.Quote(prop.Name()), prop.Property().Optional())),
				g.lang.AsTypeNullable(prop.Property())), g.lang.JsonAnnotation(prop.Property())))
	})
	g.bodyWriter.WriteLine()
}

// func (g *tsGenerator) generateParam(prop eg.PropertyObject) {
// 	g.lang.Interface(g.bodyWriter, "export ", g.lang.PublicType(prop.Title(), "Param"), prop, func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
// 		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(prop.Property())), g.lang.AsType(prop.Property()), prop.Property())
// 		wr.WriteLine(g.lang.Line(
// 			g.lang.ReturnType(
// 				g.lang.Readonly(g.lang.Type(g.lang.PublicType(prop.Name()), prop.Property().Optional())),
// 				g.lang.AsTypeNullable(prop.Property()))))
// 	})

// 	g.bodyWriter.WriteLine()
// }

// func (g *tsGenerator) generateCloneFunc() {
// 	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) Clone() %s",
// 		g.lang.PrivateName(prop.Title(), "Impl"), g.lang.PublicName(prop.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {
// 		wr.WriteBlock("return", fmt.Sprintf("&%s", g.lang.PrivateName(prop.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {
// 			for _, prop := range prop.Properties().Items() {
// 				wr.FormatLine("%s: my.%s,", g.lang.PrivateName(prop.Name()), g.lang.PrivateName(prop.Name()))
// 			}
// 		})
// 	})
// 	g.bodyWriter.WriteLine()
// }

// func (g *tsGenerator) generateImpl(prop eg.PropertyObject) {
// g.lang.Class(g.bodyWriter, g.lang.PrivateType(prop.Title(), "Impl"), prop, func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
// g.bodyWriter.WriteBlock(g.lang.Class(g.lang.PrivateName(prop.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {

// 	wr.WriteBlock(g.lang.Class(g.lang.Type(g.lang.PublicType(prop.Name()), prop.Property().Optional())),
// 			g.lang.AsTypeNullable(prop.Property()))))

// 	// Class Attributes
// 	for _, prop := range prop.Properties().Items() {
// 		wr.FormatLine("%s %s%s", g.lang.Private(g.lang.PrivateName(prop.Name())), g.lang.AsTypeOptional(prop.Property()), g.lang.Semicolon())
// 	}
// })
// g.bodyWriter.WriteLine()

// for _, prop := range prop.Properties().Items() {
// 	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) %s() %s",
// 		g.lang.PrivateName(prop.Title(), "Impl"), g.lang.PublicName(prop.Name()), g.lang.AsTypeOptional(prop.Property())), func(wr *eg.ForIfWhileLangWriter) {
// 		wr.FormatLine("return my.%s", g.lang.PrivateName(prop.Name()))
// 	})
// 	g.bodyWriter.WriteLine()
// }
// }

// func arrayCompare(prop eg.PropertyArray) string {
// 	switch prop.Items().Type() {
// 	case eg.STRING:
// 		return "wueste.ArrayLessString"
// 	case eg.INTEGER:
// 		return "wueste.ArrayLessInteger"
// 	case eg.NUMBER:
// 		return "wueste.ArrayLessNumber"
// 	case eg.BOOLEAN:
// 		return "wueste.ArrayLessBoolean"
// 	case eg.ARRAY:

// 		return "wueste.ArrayCompare"
// 		//panic("not implemented")
// 		// p := p.(PropertyArray)
// 		// return optional(p.Optional(), "[]"+g.lang.AsTypeOptional(p.Items()))
// 	case eg.OBJECT:
// 		// p := p.(Schema)
// 		// required(p, publicClassName(p.Title()))
// 		panic("not implemented")
// 	default:
// 		panic("not implemented")

// 	}
// }

// func (g *tsGenerator) generateArrayLessBlock(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter, blockFn func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter)) {
// i := fmt.Sprintf("i%d", len(mo.other))
// wr.WriteIf(fmt.Sprintf("%s < %s", g.lang.ArrayLength(mo.my), g.lang.ArrayLength(mo.other)), func(wr *eg.ForIfWhileLangWriter) {
// 	wr.WriteLine("return true")
// }, func(wr *eg.ForIfWhileLangWriter) {
// 	wr.WriteBlock("if", fmt.Sprintf("%s > %s", g.lang.ArrayLength(mo.my), g.lang.ArrayLength(mo.other)), func(wr *eg.ForIfWhileLangWriter) {
// 		wr.WriteLine("return false")
// 	})
// })
// wr.WriteBlock("for", fmt.Sprintf("%s := 0; %s < %s; %s++", i, i, g.lang.ArrayLength(mo.my), i), func(wr *eg.ForIfWhileLangWriter) {
// 	blockFn(myOther{
// 		other: fmt.Sprintf("%s[%s]", mo.other, i),
// 		my:    fmt.Sprintf("%s[%s]", mo.my, i),
// 	}, item, wr)
// })
// // wr.WriteLine("return false")
// }

// func generateArrayLessFuncLiteral(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
// nextMo := myOther{
// 	my:    fmt.Sprintf("m_%s", mo.my),
// 	other: fmt.Sprintf("v_%s", mo.other),
// }
// wr.FormatLine("%s := %s", nextMo.my, mo.my)
// wr.FormatLine("%s := %s", nextMo.other, mo.other)

// if mo.optional {
// 	wr.WriteIf(fmt.Sprintf("%s.IsSome() && %s.IsSome()", nextMo.my, nextMo.other), func(wr *eg.ForIfWhileLangWriter) {
// 		my := fmt.Sprintf("%s_%d", nextMo.my, len(nextMo.my))
// 		other := fmt.Sprintf("%s_%d", nextMo.other, len(nextMo.my))
// 		wr.FormatLine("%s := *(%s.Value())", my, nextMo.my)
// 		wr.FormatLine("%s := *(%s.Value())", other, nextMo.other)
// 		generateArrayLessBlock(myOther{
// 			other: other,
// 			my:    my,
// 		}, item, wr, generateArrayLessFunc)
// 	}, func(wr *eg.ForIfWhileLangWriter) {
// 		wr.Indent().WriteBlock("if", fmt.Sprintf("%s.IsSome() && %s.IsNone()", nextMo.my, nextMo.other), func(wr *eg.ForIfWhileLangWriter) {
// 			wr.WriteLine("return true")
// 		})
// 	})
// } else {
// 	generateArrayLessBlock(nextMo, item, wr, generateArrayLessFunc)
// }
// }

// type myOther struct {
// 	other    string
// 	my       string
// 	optional bool
// }

// func generateArrayLessFunc(mo myOther, prop eg.Property, wr *eg.ForIfWhileLangWriter) {
// 	switch prop.Type() {
// 	case eg.STRING:
// 		generateLessOptional(mo, "wueste.ArrayLessString(%my%, %other%)", wr)
// 	case eg.INTEGER:
// 		generateLessOptional(mo, "wueste.ArrayLessInteger(%my%, %other%)", wr)
// 	case eg.NUMBER:
// 		generateLessOptional(mo, "wueste.ArrayLessNumber(%my%, %other%)", wr)
// 	case eg.BOOLEAN:
// 		generateLessOptional(mo, "wueste.ArrayLessBoolean(%my%, %other%)", wr)
// 	case eg.ARRAY:
// 		generateArrayLessFuncLiteral(mo, prop.(eg.PropertyArray).Items(), wr)
// 	default:
// 		panic("not implemented")
// 	}
// }

// var reMy = regexp.MustCompile(`%my%`)
// var reOther = regexp.MustCompile(`%other%`)

// func applyMyOther(mo myOther, tmpl string) string {
// 	tmpl = reMy.ReplaceAllString(tmpl, mo.my)
// 	tmpl = reOther.ReplaceAllString(tmpl, mo.other)
// 	return tmpl
// }

// func generateLessOptional(mo myOther, tmpl string, wr *eg.ForIfWhileLangWriter) {
// 	if mo.optional {
// 		wr.WriteIf(fmt.Sprintf("!%s.IsNone() && !%s.IsNone()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
// 			wr.WriteBlock("if", applyMyOther(myOther{
// 				other: fmt.Sprintf("*(%s.Value())", mo.other),
// 				my:    fmt.Sprintf("*(%s.Value())", mo.my),
// 			}, tmpl), func(wr *eg.ForIfWhileLangWriter) {
// 				wr.WriteLine("return true")
// 			})
// 		}, func(wr *eg.ForIfWhileLangWriter) {
// 			wr.Indent().WriteBlock("if", fmt.Sprintf("!%s.IsNone() && %s.IsNone()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
// 				wr.WriteLine("return true")
// 			})
// 		})
// 	} else {
// 		wr.WriteBlock("if", applyMyOther(myOther{
// 			other: mo.other,
// 			my:    mo.my,
// 		}, tmpl), func(wr *eg.ForIfWhileLangWriter) {
// 			wr.WriteLine("return true")
// 		})
// 	}
// }

// func (g *tsGenerator) generateLessBlock(mo myOther, prop eg.Property, wr *eg.ForIfWhileLangWriter) {
// 	switch prop.Type() {
// 	case eg.STRING:
// 		// g.includes["strings"] = true
// 		generateLessOptional(mo, "strings.Compare(%my%, %other%) < 0", wr)
// 	case eg.NUMBER, eg.INTEGER:
// 		generateLessOptional(mo, "%my% < %other%", wr)
// 	case eg.BOOLEAN:
// 		generateLessOptional(mo, "%my% != %other% && %my% == false", wr)
// 	case eg.ARRAY:
// 		generateArrayLessFunc(mo, prop.(eg.PropertyArray).Items(), wr)
// 	default:
// 		panic("not implemented")
// 	}
// }

// func (g *tsGenerator) generateLessFunc() {
// 	g.bodyWriter.WriteBlock("func",
// 		fmt.Sprintf("(my *%s) Less(other %s) bool", g.lang.PrivateName(prop.Title(), "Impl"), g.lang.PublicName(prop.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {

// 			for _, prop := range prop.Properties().Items() {
// 				g.generateLessBlock(myOther{
// 					other:    fmt.Sprintf("other.%s()", g.lang.PublicName(prop.Name())),
// 					my:       fmt.Sprintf("my.%s", g.lang.PrivateName(prop.Name())),
// 					optional: prop.Property().Optional(),
// 				}, prop.Property(), wr)

// 			}
// 			wr.WriteLine("return false")
// 		})
// 	g.bodyWriter.WriteLine()
// }

// func (g *tsGenerator) hashBlock(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter, lname string) {
// 	// g.includes[WUESTE] = true
// 	if prop.Property().Optional() {
// 		wr.WriteIf(fmt.Sprintf("!my.%s.IsNone()", g.lang.PrivateName(prop.Name())), func(wr *eg.ForIfWhileLangWriter) {
// 			wr.FormatLine("w.Write([]byte(%s(*my.%s.Value()).String()))", lname, g.lang.PrivateName(prop.Name()))
// 		}, func(wr *eg.ForIfWhileLangWriter) {
// 			// SECURITY important to prevent crafted hash collision
// 			wr.FormatLine("w.Write([]byte(\"-\"))")
// 		})
// 	} else {
// 		wr.FormatLine("w.Write([]byte(%s(my.%s).String()))", lname, g.lang.PrivateName(prop.Name()))
// 	}
// }

// func generateArrayHashBlock(mo myOther, item eg.Property) {

// }

// func hashBlockArray(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
// 	// if mo.optional {
// 	// 	wr.WriteIf(fmt.Sprintf("%s.IsSome() && %s.IsSome()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
// 	// 		my := fmt.Sprintf("m%d", len(mo.my))
// 	// 		other := fmt.Sprintf("o%d", len(mo.my))
// 	// 		wr.FormatLine("%s := *(%s.Value())", my, mo.my)
// 	// 		wr.FormatLine("%s := *(%s.Value())", other, mo.other)
// 	// 		generateArrayLessBlock(myOther{
// 	// 			other: other,
// 	// 			my:    my,
// 	// 		}, item, wr, generateArrayLessFunc)
// 	// 	}, func(wr *eg.ForIfWhileLangWriter) {
// 	// 		wr.Indent().WriteBlock("if", fmt.Sprintf("%s.IsSome() && %s.IsNone()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
// 	// 			wr.WriteLine("return true")
// 	// 		})
// 	// 	})
// 	// } else {
// 	// 	generateArrayLessBlock(mo, item, wr, generateArrayLessFunc)
// 	// }
// }

// func (g *tsGenerator) generateHashFunc() {
// 	// g.includes["io"] = true
// 	// g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) Hash(w io.Writer) ", g.lang.PrivateName(prop.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {
// 	// 	for _, prop := range prop.Properties().Items() {
// 	// 		switch prop.Property().Type() {
// 	// 		case eg.STRING:
// 	// 			g.hashBlock(prop, wr, "wueste.StringLiteral")
// 	// 		case eg.INTEGER:
// 	// 			g.hashBlock(prop, wr, "wueste.IntegerLiteral")
// 	// 		case eg.NUMBER:
// 	// 			g.hashBlock(prop, wr, "wueste.NumberLiteral")
// 	// 		case eg.BOOLEAN:
// 	// 			g.hashBlock(prop, wr, "wueste.BoolLiteral")
// 	// 		case eg.ARRAY:
// 	// 			hashBlockArray(myOther{
// 	// 				other: fmt.Sprintf("other.%s()", g.lang.PublicName(prop.Name())),
// 	// 				my:    fmt.Sprintf("my.%s", g.lang.PrivateName(prop.Name())),
// 	// 			}, prop.Property().(eg.PropertyArray).Items(), wr)

// 	// 		case eg.OBJECT:
// 	// 			panic("not implemented")
// 	// 		default:
// 	// 			panic("not implemented")
// 	// 		}
// 	// 	}
// 	// })
// 	// g.bodyWriter.WriteLine()
// }

// func (g *tsGenerator) generateAsMapFunc() {
// 	// g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) AsMap() map[string]interface{}", g.lang.PrivateName(prop.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {
// 	// 	wr.WriteLine("res := map[string]interface{}{}")
// 	// 	for _, prop := range prop.Properties().Items() {
// 	// 		switch prop.Property().Type() {
// 	// 		case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN, eg.ARRAY:
// 	// 			if prop.Property().Optional() {
// 	// 				wr.WriteBlock("if", fmt.Sprintf("my.%s.IsNone()", g.lang.PrivateName(prop.Name())), func(wr *eg.ForIfWhileLangWriter) {
// 	// 					wr.FormatLine("res[%s] = my.%s.Value()", wueste.QuoteString(prop.Name()), g.lang.PrivateName(prop.Name()))
// 	// 				})
// 	// 			} else {
// 	// 				wr.FormatLine("res[%s] = my.%s", wueste.QuoteString(prop.Name()), g.lang.PrivateName(prop.Name()))
// 	// 			}
// 	// 		case eg.OBJECT:
// 	// 			panic("not implemented")
// 	// 		default:
// 	// 			panic("not implemented")
// 	// 		}
// 	// 	}
// 	// 	g.bodyWriter.WriteLine("return res")
// 	// })
// 	// g.bodyWriter.WriteLine()
// }

func (g *tsGenerator) generateFactory(prop eg.PropertyObject) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenFactory").activated = true

	// export function  NewSimpleTypeFactory(): WuestenFactory<SimpleTypeBuilder, SimpleType>
	className := g.lang.PrivateName(prop.Title(), "Factory")
	g.lang.Class(g.bodyWriter, "", g.lang.Implements(className,
		g.lang.Generics("WuestenFactory", g.lang.PublicName(prop.Title(), "Builder"), g.lang.PublicName(prop.Title()))), prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
		}, func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock("Builder():", g.lang.PublicName(prop.Title(), "Builder"), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return new %s()", g.lang.PublicName(prop.Title(), "Builder"))
			})
			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result").activated = true
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("ToObject", g.lang.ReturnType("obj", g.lang.PublicName(prop.Title()))),
				g.lang.PublicName(prop.Title(), "Object")), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("const ret: Record<string, unknown> = {}")
				for _, prop := range prop.Properties().Items() {
					if !prop.Property().Optional() {
						wr.FormatLine("ret[%s] = %s", g.lang.Quote(prop.Name()), g.lang.CallDot("obj", g.lang.PublicName(prop.Name())))
						continue
					}
					wr.WriteBlock("if ", fmt.Sprintf("(typeof %s !== 'undefined')",
						g.lang.CallDot("obj", g.lang.PublicName(prop.Name()))), func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("ret[%s] = %s", g.lang.Quote(prop.Name()), g.lang.CallDot("obj", g.lang.PublicName(prop.Name())))
					})
				}
				wr.FormatLine("return ret as unknown as %s;", g.lang.PublicName(prop.Title(), "Object"))

			})

			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result").activated = true
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Clone", g.lang.ReturnType("oth", g.lang.PublicName(prop.Title()))),
				g.lang.Generics("Result", g.lang.PublicName(prop.Title()))), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("const builder = new %s();", g.lang.PublicName(prop.Title(), "Builder"))
				for _, prop := range prop.Properties().Items() {
					wr.WriteLine(
						g.lang.Line(
							g.lang.Call(
								g.lang.CallDot("builder", g.lang.PublicName(prop.Name())),
								g.lang.CallDot("oth", g.lang.PublicName(prop.Name())))))
				}
				wr.WriteLine("return builder.Get();")
			})

		})

	g.bodyWriter.WriteLine(
		g.lang.AssignDefault(
			g.lang.Export(g.lang.Const(g.lang.PublicName(prop.Title(), "Factory"))),
			g.lang.New(g.lang.PrivateName(prop.Title(), "Factory"))))
	g.bodyWriter.WriteLine()
}

// func toAttributeType[T string | int | uint | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64](lang eg.ForIfWhileLang, prop eg.Property) string {
// 	litProp, ok := prop.(eg.PropertyLiteralType[T])
// 	if !ok {
// 		panic("not implemented")
// 	}
// 	var attribute string
// 	if litProp.Optional() {
// 		if litProp.Default().IsNone() {
// 			// attribute = fmt.Sprintf("wueste.OptionalAttribute[%s]", g.lang.AsTypeOptional(litProp))
// 			attribute = fmt.Sprintf("wueste.OptionalAttribute[%s]", lang.AsTypeOptional(litProp))
// 		} else {
// 			// attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s(%s))",
// 			attribute = fmt.Sprintf("wueste.DefaultAttribute[%s]", lang.AsTypeOptional(litProp))
// 			// g.lang.AsTypeOptional(litProp, "Some"),
// 			// litProp.Default().Value().String())
// 		}
// 	} else {
// 		if !litProp.Default().IsNone() {
// 			// attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s)",
// 			attribute = fmt.Sprintf("wueste.DefaultAttribute[%s]",
// 				lang.AsTypeOptional(litProp), litProp.Default().Value().String())
// 		} else {
// 			// attribute = fmt.Sprintf("wueste.MustAttribute[%s]()", g.lang.AsTypeOptional(litProp))
// 			attribute = fmt.Sprintf("wueste.MustAttribute[%s]", lang.AsTypeOptional(litProp))
// 		}
// 	}
// 	return attribute
// }

// func toAttributeCreation[T string | int | uint | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64](lang eg.ForIfWhileLang, prop eg.Property) string {
// 	litProp, ok := prop.(eg.PropertyLiteralType[T])
// 	if !ok {
// 		panic("not implemented")
// 	}
// 	var attribute string
// 	if litProp.Optional() {
// 		if litProp.Default().IsNone() {
// 			attribute = fmt.Sprintf("wueste.OptionalAttribute[%s]", lang.AsTypeOptional(litProp))
// 		} else {
// 			attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s(%s))",
// 				lang.AsTypeOptional(litProp, "Some"),
// 				litProp.Default().Value().String())
// 		}
// 	} else {
// 		if !litProp.Default().IsNone() {
// 			attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s)",
// 				lang.AsTypeOptional(litProp), litProp.Default().Value().String())
// 		} else {
// 			attribute = fmt.Sprintf("wueste.MustAttribute[%s]()", lang.AsTypeOptional(litProp))
// 		}
// 	}
// 	return attribute
// }

// func toAttributeCreationInteger(lang eg.ForIfWhileLang, prop eg.Property) string {

// 	_, ok := prop.(eg.PropertyLiteralType[int])
// 	if ok {
// 		return toAttributeCreation[int](lang, prop)
// 	}
// 	_, ok = prop.(eg.PropertyLiteralType[uint])
// 	if ok {
// 		return toAttributeCreation[uint](lang, prop)
// 	}

// 	_, ok = prop.(eg.PropertyLiteralType[uint8])
// 	if ok {
// 		return toAttributeCreation[uint8](lang, prop)
// 	}
// 	_, ok = prop.(eg.PropertyLiteralType[int8])
// 	if ok {
// 		return toAttributeCreation[int8](lang, prop)
// 	}

// 	_, ok = prop.(eg.PropertyLiteralType[uint16])
// 	if ok {
// 		return toAttributeCreation[uint16](lang, prop)
// 	}
// 	_, ok = prop.(eg.PropertyLiteralType[int16])
// 	if ok {
// 		return toAttributeCreation[int16](lang, prop)
// 	}

// 	_, ok = prop.(eg.PropertyLiteralType[uint32])
// 	if ok {
// 		return toAttributeCreation[uint32](lang, prop)
// 	}
// 	_, ok = prop.(eg.PropertyLiteralType[int32])
// 	if ok {
// 		return toAttributeCreation[int32](lang, prop)
// 	}

// 	_, ok = prop.(eg.PropertyLiteralType[uint64])
// 	if ok {
// 		return toAttributeCreation[uint64](lang, prop)
// 	}
// 	_, ok = prop.(eg.PropertyLiteralType[int64])
// 	if ok {
// 		return toAttributeCreation[int64](lang, prop)
// 	}
// 	panic("not implemented")
// }

// func (g *tsGenerator) toAttributeArray(prop eg.PropertyArray) string {
// 	if prop.Optional() {
// 		return fmt.Sprintf("wueste.DefaultAttribute[%s](%s())",
// 			g.lang.Optional(true, fmt.Sprintf("[]%s", g.lang.AsTypeOptional(prop.Items()))),
// 			g.lang.Optional(true, fmt.Sprintf("[]%s", g.lang.AsTypeOptional(prop.Items())), "None"))
// 	}
// 	return fmt.Sprintf("wueste.DefaultAttribute[[]%s]([]%s{})", g.lang.AsType(prop.Items()), g.lang.AsType(prop.Items()))
// }
// func generateBuilderSetter(mo myOther, prop eg.Property, wr *eg.ForIfWhileLangWriter) {
// 	switch prop.Type() {
// 	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
// 		if mo.optional {
// 			wr.WriteBlock("if", fmt.Sprintf("%s.IsSome()", mo.my), func(wr *eg.ForIfWhileLangWriter) {
// 				wr.FormatLine("*(%s.Value()).Set(v)", mo.my)
// 			})
// 		} else {
// 			wr.FormatLine("%s.Set(v)", mo.my)
// 		}
// 	case eg.ARRAY:
// 		generateArraySetBlock := func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter, blockFn func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter)) {
// 			loopI := fmt.Sprintf("i%d", len(mo.my))
// 			wr.WriteBlock("for", fmt.Sprintf("%s := 0; %s < len(%s); %s++", loopI, loopI, mo.my, loopI), func(wr *eg.ForIfWhileLangWriter) {
// 				blockFn(myOther{
// 					// other: fmt.Sprintf("[i]",
// 					my: fmt.Sprintf("%s[%s]", mo.my, loopI),
// 				}, item, wr)
// 			})
// 		}
// 		if mo.optional {
// 			wr.WriteBlock("if", fmt.Sprintf("%s.IsSome()", mo.my), func(wr *eg.ForIfWhileLangWriter) {
// 				wr.FormatLine("v := *(%s.Value())", mo.my)
// 				generateArraySetBlock(mo, prop.(eg.PropertyArray).Items(), wr, func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
// 					generateBuilderSetter(mo, item, wr)
// 				})
// 			})
// 		} else {
// 			generateArraySetBlock(mo, prop.(eg.PropertyArray).Items(), wr, func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
// 				generateBuilderSetter(mo, item, wr)
// 			})
// 		}
// 	default:
// 		panic("not implemented")
// 	}
// }

// func (g *tsGenerator) genWeuesteAttributeCreation(prop eg.Property) string {
// 	switch prop.Type() {
// 	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
// 		return fmt.Sprintf("wueste.Attribute[%s]", g.lang.AsTypeOptional(prop))
// 	case eg.ARRAY:
// 		return fmt.Sprintf("wueste.Attribute[%s]", g.genWuesteAttributeCreation(prop.(eg.PropertyArray).Items()))
// 	default:
// 		panic("not implemented")
// 	}
// }

// func (g *tsGenerator) genWeuesteAttributeType(prop eg.Property) string {
// 	switch prop.Type() {
// 	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
// 		return fmt.Sprintf("wueste.Attribute[%s]", g.lang.AsTypeOptional(prop))
// 	case eg.ARRAY:
// 		return fmt.Sprintf("wueste.Attribute[%s]", g.genWeuesteAttributeType(prop.(eg.PropertyArray).Items()))
// 	default:
// 		panic("not implemented")
// 	}
// }

func getDefaultForProperty(prop eg.Property) *string {
	{
		p, ok := prop.(eg.PropertyLiteralType[string])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[bool])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[float32])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[float64])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[int])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[uint])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[int8])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[uint8])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[int16])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[uint16])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[int32])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[uint32])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[int64])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	{
		p, ok := prop.(eg.PropertyLiteralType[uint64])
		if ok {
			if p.Default().IsSome() {
				return p.Default().Value().String()
			} else {
				return nil
			}
		}
	}
	return nil
}

// func (g *tsGenerator) genWuesteBuilderAttributeType(prop eg.Property) string {
// 	returnFn := func(s string) string {
// 		return g.lang.Generics("ReturnType", fmt.Sprintf("typeof %s", s))
// 	}
// 	switch prop.Type() {
// 	case eg.STRING:
// 		g.includes.AddType(g.cfg.FromWueste, "wuesten").activated = true
// 		p := prop.(eg.PropertyString)
// 		if p.Format().IsSome() {
// 			switch *p.Format().Value() {
// 			case eg.DATE_TIME:
// 				if prop.Optional() {
// 					return returnFn("wuesten.AttributeDateTimeOptional")
// 				} else {
// 					return returnFn("wuesten.AttributeDateTime")
// 				}
// 			default:
// 			}
// 		}
// 		if prop.Optional() {
// 			return returnFn("wuesten.AttributeStringOptional")
// 		} else {
// 			return returnFn("wuesten.AttributeString")
// 		}
// 	case eg.INTEGER:
// 		g.includes.AddType(g.cfg.FromWueste, "wuesten").activated = true
// 		if prop.Optional() {
// 			return returnFn("wuesten.AttributeIntegerOptional")
// 		} else {
// 			return returnFn("wuesten.AttributeInteger")
// 		}
// 	case eg.NUMBER:
// 		g.includes.AddType(g.cfg.FromWueste, "wuesten").activated = true
// 		if prop.Optional() {
// 			return returnFn("wuesten.AttributeNumberOptional")
// 		} else {
// 			return returnFn("wuesten.AttributeNumber")
// 		}
// 	case eg.BOOLEAN:
// 		g.includes.AddType(g.cfg.FromWueste, "wuesten").activated = true
// 		if prop.Optional() {
// 			return returnFn("wuesten.AttributeBooleanOptional")
// 		} else {
// 			return returnFn("wuesten.AttributeBoolean")
// 		}
// 	case eg.ARRAY:
// 		// g.includes[WUESTE] = true
// 		p := prop.(eg.PropertyArray)
// 		if prop.Optional() {
// 			return fmt.Sprintf("%s[]", g.genWuesteBuilderAttributeType(p.Items()))
// 			// return g.lang.Generics("wuesten.AttributeArrayOptional", returnFn(g.genWuesteBuilderAttributeType(p.Items())))
// 		} else {
// 			return fmt.Sprintf("%s[]", g.genWuesteBuilderAttributeType(p.Items()))
// 			// return g.lang.Generics("wuesten.AttributeArray", returnFn(g.genWuesteBuilderAttributeType(p.Items())))
// 		}
// 		// toAttributeArray(prop.Property().(PropertyArray)))
// 	case eg.OBJECT:
// 		g.includes.AddType(g.cfg.FromWueste, "wuesten").activated = true
// 		p := prop.(eg.PropertyObject)
// 		factory := g.lang.PublicName("New", p.Title(), "Factory")
// 		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(p)), factory, p)
// 		if prop.Optional() {
// 			return g.lang.Generics("wuesten.AttributeObjectOptional", g.lang.PublicName(p.Title()))
// 		} else {
// 			return g.lang.Generics("wuesten.AttributeObject", g.lang.PublicName(p.Title()))
// 		}
// 	default:
// 		panic("not implemented")
// 	}
// }

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

func (g *tsGenerator) genWuesteBuilderAttribute(name string, prop eg.Property, paramFns ...func() string) string {
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
			switch *p.Format().Value() {
			case eg.DATE_TIME:
				if prop.Optional() {
					return g.lang.Call("wuesten.AttributeDateTimeOptional", paramFn())
				} else {
					return g.lang.Call("wuesten.AttributeDateTime", paramFn())
				}
			default:
			}
		}
		if prop.Optional() {
			return g.lang.Call("wuesten.AttributeStringOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeString", paramFn())
		}
	case eg.INTEGER:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		if prop.Optional() {
			return g.lang.Call("wuesten.AttributeIntegerOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeInteger", paramFn())
		}
	case eg.NUMBER:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true

		if prop.Optional() {
			return g.lang.Call("wuesten.AttributeNumberOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeNumber", paramFn())
		}
	case eg.BOOLEAN:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		if prop.Optional() {
			return g.lang.Call("wuesten.AttributeBooleanOptional", paramFn())
		} else {
			return g.lang.Call("wuesten.AttributeBoolean", paramFn())
		}
	case eg.ARRAY:
		if prop.Optional() {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttrOptional").activated = true
			return g.lang.New("WuestenAttrOptional", g.lang.New(g.lang.PublicName(name, "Attribute"), paramFn()))
		} else {
			return g.lang.New(g.lang.PublicName(name, "Attribute"), paramFn())
		}
	case eg.OBJECT:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "wuesten").activated = true
		p := prop.(eg.PropertyObject)
		factory := g.lang.PublicName(p.Title(), "Factory")
		g.includes.AddProperty(g.lang.FileName(g.lang.AsType(p)), factory, p)
		if prop.Optional() {
			return g.lang.Call("wuesten.AttributeObjectOptional", paramFn(), factory)
		} else {
			return g.lang.Call("wuesten.AttributeObject", paramFn(), factory)
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
		wr.FormatLine("return Result.Err(\"it's not iterable on level %d\")", level)
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
						g.lang.Generics("attrRes as unknown as Result", returnType)))
				})
				wr.WriteLine(g.lang.Call(fmt.Sprintf("r%d.push", level), g.lang.Call("attrRes.unwrap")))
			}
			wr.FormatLine("c%d++", level)
		})
}

func (g *tsGenerator) generateBuilder(prop eg.PropertyObject) {
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttribute").activated = true

	for _, p := range prop.Properties().Items() {
		pa, ok := p.Property().(eg.PropertyArray)
		if ok {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttribute").activated = true
			className := g.lang.PublicName(p.Name(), "Attribute")
			g.lang.Class(g.bodyWriter, "", g.lang.Extends(className,
				g.lang.Generics("WuestenAttr",
					g.lang.AsTypeNullable(p.Property(), WithOptional(p.Property().Optional())))),
				prop,
				func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
				}, func(wr *eg.ForIfWhileLangWriter) {
					attrib := g.lang.ReturnType(g.lang.OptionalParam(g.lang.PrivateName("value"), p.Property().Optional()), g.lang.AsType(p.Property()))
					if !p.Property().Optional() {
						attrib += " = []"
					}
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeParameter").activated = true

					wr.WriteBlock("",
						g.lang.Call("constructor",
							g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeParameter", g.lang.AsType(getItemType(pa)))),
						), func(wr *eg.ForIfWhileLangWriter) {
							attr := g.genWuesteBuilderAttribute("ARRAY", getItemType(pa), func() string { return "param" })
							wr.WriteLine(g.lang.AssignDefault(g.lang.Const("itemAttr"), attr))

							wr.WriteBlock("", "super({jsonname: param.jsonname, varname: param.varname, base: param.base}, (t0: unknown) => ", func(wr *eg.ForIfWhileLangWriter) {
								g.generateArrayCoerce(0, "t0", g.lang.AsType(pa), pa, wr)
								wr.WriteLine(g.lang.Return(g.lang.Call("Result.Ok", "r0")))
							}, " {", "})")
						})
				})
			g.bodyWriter.WriteLine()
		}
	}

	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttr").activated = true

	resultsClassName := g.lang.PublicName(prop.Title(), "Results")

	g.lang.Interface(g.bodyWriter, "", resultsClassName, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine(g.lang.Readonly(
				g.lang.ReturnType(
					g.lang.PrivateName(pi.Name()),
					g.lang.Generics("Result", g.lang.AsTypeNullable(pi.Property(), WithOptional(pi.Property().Optional()))))))
		})

	attrsClassName := g.lang.PublicName(prop.Title(), "Attributes")
	g.lang.Class(g.bodyWriter, "", attrsClassName, prop,
		func(pi eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {},
		func(wr *eg.ForIfWhileLangWriter) {
			for _, prop := range prop.Properties().Items() {
				// wr.FormatLine("readonly %s = %s;", g.lang.PrivateName(prop.Name()), g.genWuesteBuilderAttribute(prop.Name(), prop.Property()))
				wr.WriteLine(
					g.lang.Readonly(g.lang.ReturnType(g.lang.PrivateName(prop.Name()), g.lang.Generics("WuestenAttribute",
						g.lang.AsTypeNullable(prop.Property(), WithOptional(prop.Property().Optional())),
						g.lang.AsTypeNullable(prop.Property(), WithAddCoerce(), WithOptional(prop.Property().Optional())),
					))))
			}

			wr.WriteBlock("", g.lang.Call("constructor",
				g.lang.AssignDefault("param",
					fmt.Sprintf("{jsonname: %s, varname: %s, base: \"\"}",
						g.lang.Quote(prop.Title()), g.lang.Quote(g.lang.PublicName(prop.Title()))))), func(wr *eg.ForIfWhileLangWriter) {
				if len(prop.Properties().Items()) > 0 {
					g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeName").activated = true
					wr.WriteLine(g.lang.AssignDefault(g.lang.Const("baseName"), g.lang.Call("WuestenAttributeName", "param")))
					for _, prop := range prop.Properties().Items() {
						wr.WriteLine(g.lang.AssignDefault(
							g.lang.CallDot("this", g.lang.PrivateName(prop.Name())), g.genWuesteBuilderAttribute(prop.Name(), prop.Property())))
					}
				}
			})

			wr.WriteBlock("",
				g.lang.AssignDefault(g.lang.Readonly("Coerce"),
					g.lang.ReturnType(
						g.lang.Call("", g.lang.ReturnType("value", "unknown")),
						g.lang.Generics("Result", g.lang.PublicName(prop.Title())))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("if", "(!(typeof value === 'object' && value !== null))", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return Result.Err(Error('expected object'));")
					})
					wr.WriteBlock("return", "this._fromResults", func(wr *eg.ForIfWhileLangWriter) {
						for _, prop := range prop.Properties().Items() {
							wr.FormatLine("%s: this.%s.CoerceAttribute(value),", g.lang.PrivateName(prop.Name()), g.lang.PrivateName(prop.Name()))
						}
					}, "({", "});")
				}, " => {")

			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("Get", ""),
					g.lang.Generics("Result", g.lang.PublicName(prop.Title()))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteBlock("return", "this._fromResults", func(wr *eg.ForIfWhileLangWriter) {
						for _, prop := range prop.Properties().Items() {
							wr.FormatLine("%s: this.%s.Get(),", g.lang.PrivateName(prop.Name()), g.lang.PrivateName(prop.Name()))
						}
					}, "({", "});")
				})

			wr.WriteBlock("",
				g.lang.ReturnType(
					g.lang.Call("_fromResults", g.lang.ReturnType("results", resultsClassName)),
					g.lang.Generics("Result", g.lang.PublicName(prop.Title()))),
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine("const errors: string[] = [];")
					for _, prop := range prop.Properties().Items() {
						val := g.lang.CallDot("results", g.lang.PrivateName(prop.Name()))
						wr.WriteBlock("if", "("+val+".is_err())", func(wr *eg.ForIfWhileLangWriter) {
							wr.FormatLine("errors.push(%s.unwrap_err().message);", val)
						})
					}
					// wr.FormatLine("const errors = Object.values(results).filter(r => r.is_err()).map(r => r.unwrap_err().message)")
					wr.WriteBlock("if", "(errors.length > 0)", func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("return Result.Err(Error(errors.join('\\n')));")
					})
					wr.WriteBlock("return",
						g.lang.Generics("Result.Ok", g.lang.PublicName(prop.Title())),
						func(wr *eg.ForIfWhileLangWriter) {
							for _, prop := range prop.Properties().Items() {
								wr.FormatLine("%s: results.%s.unwrap(),", g.lang.PublicName(prop.Name()), g.lang.PrivateName(prop.Name()))
							}
						}, "({", "});")
				})
		})

	className := g.lang.PublicName(prop.Title(), "Builder")
	g.lang.Class(g.bodyWriter, "export ", g.lang.Extends(className,
		g.lang.Generics("WuestenAttr", g.lang.PublicName(prop.Title()),
			g.lang.OrType(
				g.lang.PublicName(prop.Title()),
				g.lang.PublicName(prop.Title(), "Param"),
				g.lang.PublicName(prop.Title(), "Object")))), prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock(
				g.lang.ReturnType(
					g.lang.Call(g.lang.Type(g.lang.PublicType(prop.Name()), false),
						g.lang.ReturnType(g.lang.OptionalParam("v", prop.Property().Optional()), g.lang.Type(g.lang.AsType(prop.Property(), WithAddCoerce()), false))),
					className),
				"", func(wr *eg.ForIfWhileLangWriter) {
					wr.FormatLine("this._attr.%s.Coerce(v);", g.lang.PrivateName(prop.Name()))
					wr.FormatLine("return this;")
				})
		}, func(wr *eg.ForIfWhileLangWriter) {

			g.includes.AddType(g.cfg.EntityCfg.FromResult, "Result").activated = true

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
						g.lang.Quote(prop.Title()), g.lang.Quote(g.lang.PublicName(prop.Title()))))), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Const(g.lang.AssignDefault("attr", g.lang.New(attrsClassName, "param"))))
				wr.WriteLine("super(param, attr.Coerce);")
				wr.WriteLine(g.lang.AssignDefault("this._attr", "attr"))
			})
			wr.WriteBlock("", g.lang.ReturnType(
				g.lang.Call("Get", ""), g.lang.Generics("Result", g.lang.PublicName(prop.Title()))), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Return(g.lang.Call("this._attr.Get", "")))
			})

		})
	// g.bodyWriter.WriteBlock("type", g.lang.PublicName(prop.Title(), "Builder")+" struct", func(wr *eg.ForIfWhileLangWriter) {
	// 	// g.includes[WUESTE] = true
	// 	for _, prop := range prop.Properties().Items() {
	// 		wr.FormatLine("%s %s", g.lang.PrivateName(prop.Name()), "WTF") // g.genWeuesteAttributeType(prop.Property()))
	// 	}
	// })
	// g.bodyWriter.WriteLine()
	// // in languages like TS we could pass a literal here.
	// // TS Allows to type define Required Types
	// g.bodyWriter.WriteBlock("func", fmt.Sprintf("New%s() *%s", g.lang.PublicName(prop.Title(), "Builder"), g.lang.PublicName(prop.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
	// 	wr.WriteBlock(fmt.Sprintf("return &%s", g.lang.PublicName(prop.Title(), "Builder")), "", func(wr *eg.ForIfWhileLangWriter) {
	// 		for _, prop := range prop.Properties().Items() {
	// 			wr.FormatLine("%s: %s,", g.lang.PrivateName(prop.Name()), "WTF1") // g.genWuesteAttributeCreation(prop.Property()))
	// 		}
	// 	})
	// 	wr.WriteLine()
	// })

	// for _, pi := range prop.Properties().Items() {
	// 	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(b *%s) %s(v %s) *%s",
	// 		g.lang.PublicName(prop.Title(), "Builder"), g.lang.PublicName(pi.Name()), g.lang.AsTypeOptional(pi.Property()), g.lang.PublicName(prop.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
	// 		// generateBuilderSetter(myOther{
	// 		// 	my:       "v", // fmt.Sprintf("s", g.lang.PrivateName(prop.Name())),
	// 		// 	optional: prop.Property().Optional(),
	// 		// }, prop.Property(), wr)
	// 		wr.WriteLine("return b")
	// 	})
	// 	g.bodyWriter.WriteLine()
	// }

	// g.bodyWriter.WriteBlock("func", fmt.Sprintf("(b *%s) IsValid() rusty.Optional[error]",
	// 	g.lang.PublicName(prop.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
	// 	// g.includes["fmt"] = true
	// 	// g.includes[RUSTY] = true
	// 	for _, pi := range prop.Properties().Items() {
	// 		wr.WriteBlock("if", fmt.Sprintf("valid := b.%s.IsValid(); !valid.IsNone()", g.lang.PrivateName(pi.Name())), func(wr *eg.ForIfWhileLangWriter) {
	// 			wr.FormatLine("return rusty.Some[error](fmt.Errorf(\"%s.%s:%%s\", (*valid.Value()).Error()))",
	// 				g.lang.PublicName(prop.Title(), "Class"),
	// 				g.lang.PublicName(pi.Name()))
	// 		})
	// 	}
	// 	wr.FormatLine("return rusty.None[error]()")
	// })
	// g.bodyWriter.WriteBlock("func", fmt.Sprintf("(b *%s) ToClass() rusty.Result<%s>",
	// 	g.lang.PublicName(prop.Title(), "Builder"), g.lang.PublicName(prop.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {
	// 	wr.WriteBlock("if", "valid := b.IsValid(); !valid.IsNone()", func(wr *eg.ForIfWhileLangWriter) {
	// 		wr.FormatLine("return rusty.Err[%s](*valid.Value())", g.lang.PublicName(prop.Title(), "Class"))
	// 	})
	// 	wr.WriteBlock(fmt.Sprintf("return rusty.Ok[%s](&"+g.lang.PrivateName(prop.Title(), "Impl"), g.lang.PublicName(prop.Title(), "Class")), "", func(wr *eg.ForIfWhileLangWriter) {
	// 		for _, prop := range prop.Properties().Items() {
	// 			wr.FormatLine("%s: b.%s.Get(),", g.lang.PrivateName(prop.Name()), g.lang.PrivateName(prop.Name()))
	// 		}
	// 	}, "{", "})")
	// })
	// g.bodyWriter.WriteLine()
}

type tsGenerator struct {
	cfg        *eg.GeneratorConfig
	lang       tsLang
	includes   *externalTypes
	bodyWriter *eg.ForIfWhileLangWriter
}

func (g *tsGenerator) generatePropertyObject(prop eg.PropertyObject, sl eg.SchemaLoader) {
	g.generateClass(prop)
	g.generateJson(prop)
	g.generateBuilder(prop)
	g.generateFactory(prop)

	os.MkdirAll(g.cfg.OutputDir, 0755)

	fname := filepath.Join(g.cfg.OutputDir, g.lang.FileName(prop.Title()))
	tmpFname := filepath.Join(g.cfg.OutputDir, "."+g.lang.FileName(prop.Title())+uuid.New().String()+".ts")

	fmt.Printf("Generate: %s -> %s\n", prop.FileName(), fname)
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
			if len(include.Types()) <= 3 {
				header.FormatLine("import { %s } from %s;", strings.Join(include.Types(), ", "),
					g.lang.Quote(g.lang.RemoveFileExt(include.fileName)))
			} else {
				header.WriteBlock("", "import", func(wr *eg.ForIfWhileLangWriter) {
					for idx, t := range include.Types() {
						comma := ","
						if idx == len(include.Types())-1 {
							comma = ""
						}
						wr.FormatLine("%s%s", t, comma)
					}
				}, "{", fmt.Sprintf("} from %s;", g.lang.Quote(g.lang.RemoveFileExt(include.fileName))))
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
	sl.SchemaRegistry().SetWritten(prop)
}

type externalType struct {
	toGenerate bool
	activated  bool
	// prefix     string
	fileName string
	types    map[string]struct{}
	property eg.PropertyObject
}

func (et *externalType) Types() []string {
	types := make([]string, 0, len(et.types))
	for k, _ := range et.types {
		types = append(types, k)
	}
	sort.Strings(types)
	return types
}

type externalTypes struct {
	types map[string]*externalType
}

func (g *externalTypes) AddProperty(fileName, typ string, prop eg.Property) {
	po, ok := prop.(eg.PropertyObject)
	if ok {
		t := g.AddType(fileName, typ)
		t.property = po
		t.activated = true
	}
}

func (g *externalTypes) AddType(fileName, typeName string) *externalType {
	if fileName == "" {
		return nil
	}
	et, ok := g.types[fileName]
	if !ok {
		et = &externalType{
			toGenerate: false,
			fileName:   fileName,
			types:      make(map[string]struct{}),
		}
		g.types[fileName] = et
	}
	et.types[typeName] = struct{}{}
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

func TsGenerator(cfg *eg.GeneratorConfig, schema eg.Property, sl eg.SchemaLoader) {
	po, found := schema.(eg.PropertyObject)
	if !found {
		panic("TsGenerator not a property object")
	}
	if sl.SchemaRegistry().IsWritten(po) {
		return
	}
	g := &tsGenerator{
		cfg:        cfg,
		includes:   newExternalTypes(),
		bodyWriter: eg.NewForIfWhileLangWriter(eg.ForIfWhileLangWriter{OfsIndent: cfg.EntityCfg.Indent}),
	}

	po, _ok := schema.(eg.PropertyObject)
	if _ok {
		g.generatePropertyObject(po, sl)
		for _, prop := range g.includes.ActiveTypes() {
			if prop.property != nil {
				TsGenerator(cfg, prop.property, sl)
			}
		}
	}
}
