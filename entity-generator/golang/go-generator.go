package golang

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	eg "github.com/mabels/wueste/entity-generator"
	"github.com/mabels/wueste/entity-generator/wueste"
)

const RUSTY = "github.com/mabels/wueste/entity-generator/src/rusty"
const WUESTE = "github.com/mabels/wueste/entity-generator/src/wueste"

type ObjectType[T any] interface {
	Clone() T
	Equals(other T) bool
	AsMap() map[string]interface{}
}

var KeyWords = map[string]bool{
	"break":     true,
	"default":   true,
	"func":      true,
	"interface": true,
	"select":    true,
	"case":      true,
	"defer":     true,
	"go":        true,
	"map":       true,
	"struct":    true,
	"chan":      true,
	"else":      true,

	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,
	"error":       true,
	"string":      true,
	"bool":        true,
	"uint8":       true,
	"uint16":      true,
	"uint32":      true,
	"uint64":      true,

	"int8":  true,
	"int16": true,
	"int32": true,
	"int64": true,

	"float32": true,
	"float64": true,

	"complex64":  true,
	"complex128": true,

	"byte": true,
	"rune": true,
}

func jsonTag(prop eg.PropertyItem) string {
	name := prop.Name()
	if prop.Optional() {
		return fmt.Sprintf("`json:%s`", wueste.QuoteString(fmt.Sprintf("%s,omitempty", name)))
	}
	return fmt.Sprintf("`json:%s`", wueste.QuoteString(prop.Name()))
}

var keyWords = map[string]bool{
	"break":      true,
	"as":         true,
	"any":        true,
	"switch":     true,
	"case":       true,
	"if":         true,
	"throw	":     true,
	"else":       true,
	"var":        true,
	"number":     true,
	"string":     true,
	"get":        true,
	"module":     true,
	"type":       true,
	"instanceof": true,
	"typeof":     true,
	"public	":    true,
	"private":    true,
	"enum":       true,
	"export":     true,
	"finally":    true,
	"for":        true,
	"while	":     true,
	"void":       true,
	"null":       true,
	"super	":     true,
	"this":       true,
	"new":        true,
	"in	":        true,
	"return	":    true,
	"true":       true,
	"false":      true,
	"any	":       true,
	"extends	":   true,
	"static	":    true,
	"let":        true,
	"package":    true,
	"implements": true,
	"interface	": true,
	"function":   true,
	"try":        true,
	"yield	":     true,
	"const":      true,
	"continue	":  true,
	"do	":        true,
	"catch":      true,
}

func (g *goGenerator) generateClass() {
	g.bodyWriter.WriteBlock("type", g.lang.PublicName(g.schema.Title(), "Class")+" interface", func(wr *eg.ForIfWhileLangWriter) {
		for _, prop := range g.schema.Items() {
			wr.FormatLine("%s() %s", g.lang.PublicName(prop.Name()), g.lang.AsTypeOptional(prop))
			// String() string
		}
	})
	g.bodyWriter.WriteLine()
}

func (g *goGenerator) generateJson() {
	g.bodyWriter.WriteBlock("type", g.lang.PublicName(g.schema.Title(), "Json")+" struct", func(wr *eg.ForIfWhileLangWriter) {
		for _, prop := range g.schema.Items() {
			wr.FormatLine("%s %s %s", g.lang.PublicName(prop.Name()), g.lang.AsTypePtr(prop.Property()), jsonTag(prop))
		}
	})
	g.bodyWriter.WriteLine()
}

func (g *goGenerator) generateParam() {
	g.bodyWriter.WriteBlock("type", g.lang.PublicName(g.schema.Title(), "Param")+" struct", func(wr *eg.ForIfWhileLangWriter) {
		for _, prop := range g.schema.Items() {
			wr.FormatLine("%s %s", g.lang.PublicName(prop.Name()), g.lang.AsTypeOptional(prop))
		}
	})
	g.bodyWriter.WriteLine()
}

func (g *goGenerator) generateCloneFunc() {
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) Clone() %s",
		g.lang.PrivateName(g.schema.Title(), "Impl"), g.lang.PublicName(g.schema.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteBlock("return", fmt.Sprintf("&%s", g.lang.PrivateName(g.schema.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {
			for _, prop := range g.schema.Items() {
				wr.FormatLine("%s: my.%s,", g.lang.PrivateName(prop.Name()), g.lang.PrivateName(prop.Name()))
			}
		})
	})
	g.bodyWriter.WriteLine()
}

func arrayCompare(prop eg.PropertyArray) string {
	switch prop.Items().Type() {
	case eg.STRING:
		return "wueste.ArrayLessString"
	case eg.INTEGER:
		return "wueste.ArrayLessInteger"
	case eg.NUMBER:
		return "wueste.ArrayLessNumber"
	case eg.BOOLEAN:
		return "wueste.ArrayLessBoolean"
	case eg.ARRAY:

		return "wueste.ArrayCompare"
		//panic("not implemented")
		// p := p.(PropertyArray)
		// return optional(p.Optional(), "[]"+g.lang.AsTypeOptional(p.Items()))
	case eg.OBJECT:
		// p := p.(Schema)
		// required(p, publicClassName(p.Title()))
		panic("not implemented")
	default:
		panic("not implemented")

	}
}

func generateArrayLessBlock(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter, blockFn func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter)) {
	i := fmt.Sprintf("i%d", len(mo.other))
	wr.WriteIf(fmt.Sprintf("len(%s) < len(%s)", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteLine("return true")
	}, func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteBlock("if", fmt.Sprintf("len(%s) > len(%s)", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine("return false")
		})
	})
	wr.WriteBlock("for", fmt.Sprintf("%s := 0; %s < len(%s); %s++", i, i, mo.my, i), func(wr *eg.ForIfWhileLangWriter) {
		blockFn(myOther{
			other: fmt.Sprintf("%s[%s]", mo.other, i),
			my:    fmt.Sprintf("%s[%s]", mo.my, i),
		}, item, wr)
	})
	// wr.WriteLine("return false")
}

func generateArrayLessFuncLiteral(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
	nextMo := myOther{
		my:    fmt.Sprintf("m_%s", mo.my),
		other: fmt.Sprintf("v_%s", mo.other),
	}
	wr.FormatLine("%s := %s", nextMo.my, mo.my)
	wr.FormatLine("%s := %s", nextMo.other, mo.other)

	if mo.optional {
		wr.WriteIf(fmt.Sprintf("%s.IsSome() && %s.IsSome()", nextMo.my, nextMo.other), func(wr *eg.ForIfWhileLangWriter) {
			my := fmt.Sprintf("%s_%d", nextMo.my, len(nextMo.my))
			other := fmt.Sprintf("%s_%d", nextMo.other, len(nextMo.my))
			wr.FormatLine("%s := *(%s.Value())", my, nextMo.my)
			wr.FormatLine("%s := *(%s.Value())", other, nextMo.other)
			generateArrayLessBlock(myOther{
				other: other,
				my:    my,
			}, item, wr, generateArrayLessFunc)
		}, func(wr *eg.ForIfWhileLangWriter) {
			wr.Indent().WriteBlock("if", fmt.Sprintf("%s.IsSome() && %s.IsNone()", nextMo.my, nextMo.other), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("return true")
			})
		})
	} else {
		generateArrayLessBlock(nextMo, item, wr, generateArrayLessFunc)
	}
}

type myOther struct {
	other    string
	my       string
	optional bool
}

func generateArrayLessFunc(mo myOther, prop eg.Property, wr *eg.ForIfWhileLangWriter) {
	switch prop.Type() {
	case eg.STRING:
		generateLessOptional(mo, "wueste.ArrayLessString(%my%, %other%)", wr)
	case eg.INTEGER:
		generateLessOptional(mo, "wueste.ArrayLessInteger(%my%, %other%)", wr)
	case eg.NUMBER:
		generateLessOptional(mo, "wueste.ArrayLessNumber(%my%, %other%)", wr)
	case eg.BOOLEAN:
		generateLessOptional(mo, "wueste.ArrayLessBoolean(%my%, %other%)", wr)
	case eg.ARRAY:
		generateArrayLessFuncLiteral(mo, prop.(eg.PropertyArray).Items(), wr)
	default:
		panic("not implemented")
	}

}

var reMy = regexp.MustCompile(`%my%`)
var reOther = regexp.MustCompile(`%other%`)

func applyMyOther(mo myOther, tmpl string) string {
	tmpl = reMy.ReplaceAllString(tmpl, mo.my)
	tmpl = reOther.ReplaceAllString(tmpl, mo.other)
	return tmpl
}

func generateLessOptional(mo myOther, tmpl string, wr *eg.ForIfWhileLangWriter) {
	if mo.optional {
		wr.WriteIf(fmt.Sprintf("!%s.IsNone() && !%s.IsNone()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock("if", applyMyOther(myOther{
				other: fmt.Sprintf("*(%s.Value())", mo.other),
				my:    fmt.Sprintf("*(%s.Value())", mo.my),
			}, tmpl), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("return true")
			})
		}, func(wr *eg.ForIfWhileLangWriter) {
			wr.Indent().WriteBlock("if", fmt.Sprintf("!%s.IsNone() && %s.IsNone()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("return true")
			})
		})
	} else {
		wr.WriteBlock("if", applyMyOther(myOther{
			other: mo.other,
			my:    mo.my,
		}, tmpl), func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine("return true")
		})
	}
}

func (g *goGenerator) generateLessBlock(mo myOther, prop eg.Property, wr *eg.ForIfWhileLangWriter) {
	switch prop.Type() {
	case eg.STRING:
		g.includes["strings"] = true
		generateLessOptional(mo, "strings.Compare(%my%, %other%) < 0", wr)
	case eg.NUMBER, eg.INTEGER:
		generateLessOptional(mo, "%my% < %other%", wr)
	case eg.BOOLEAN:
		generateLessOptional(mo, "%my% != %other% && %my% == false", wr)
	case eg.ARRAY:
		generateArrayLessFunc(mo, prop.(eg.PropertyArray).Items(), wr)
	default:
		panic("not implemented")
	}
}

func (g *goGenerator) generateLessFunc() {
	g.bodyWriter.WriteBlock("func",
		fmt.Sprintf("(my *%s) Less(other %s) bool", g.lang.PrivateName(g.schema.Title(), "Impl"), g.lang.PublicName(g.schema.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {

			for _, prop := range g.schema.Items() {
				g.generateLessBlock(myOther{
					other:    fmt.Sprintf("other.%s()", g.lang.PublicName(prop.Name())),
					my:       fmt.Sprintf("my.%s", g.lang.PrivateName(prop.Name())),
					optional: prop.Optional(),
				}, prop.Property(), wr)

			}
			wr.WriteLine("return false")
		})
	g.bodyWriter.WriteLine()
}

func (g *goGenerator) hashBlock(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter, lname string) {
	g.includes[WUESTE] = true
	if prop.Optional() {
		wr.WriteIf(fmt.Sprintf("!my.%s.IsNone()", g.lang.PrivateName(prop.Name())), func(wr *eg.ForIfWhileLangWriter) {
			wr.FormatLine("w.Write([]byte(%s(*my.%s.Value()).String()))", lname, g.lang.PrivateName(prop.Name()))
		}, func(wr *eg.ForIfWhileLangWriter) {
			// SECURITY important to prevent crafted hash collision
			wr.FormatLine("w.Write([]byte(\"-\"))")
		})
	} else {
		wr.FormatLine("w.Write([]byte(%s(my.%s).String()))", lname, g.lang.PrivateName(prop.Name()))
	}
}

func generateArrayHashBlock(mo myOther, item eg.Property) {

}

func hashBlockArray(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
	if mo.optional {
		wr.WriteIf(fmt.Sprintf("%s.IsSome() && %s.IsSome()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
			my := fmt.Sprintf("m%d", len(mo.my))
			other := fmt.Sprintf("o%d", len(mo.my))
			wr.FormatLine("%s := *(%s.Value())", my, mo.my)
			wr.FormatLine("%s := *(%s.Value())", other, mo.other)
			generateArrayLessBlock(myOther{
				other: other,
				my:    my,
			}, item, wr, generateArrayLessFunc)
		}, func(wr *eg.ForIfWhileLangWriter) {
			wr.Indent().WriteBlock("if", fmt.Sprintf("%s.IsSome() && %s.IsNone()", mo.my, mo.other), func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine("return true")
			})
		})
	} else {
		generateArrayLessBlock(mo, item, wr, generateArrayLessFunc)
	}
}

func (g *goGenerator) generateHashFunc() {
	g.includes["io"] = true
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) Hash(w io.Writer) ", g.lang.PrivateName(g.schema.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {
		for _, prop := range g.schema.Items() {
			switch prop.Property().Type() {
			case eg.STRING:
				g.hashBlock(prop, wr, "wueste.StringLiteral")
			case eg.INTEGER:
				g.hashBlock(prop, wr, "wueste.IntegerLiteral")
			case eg.NUMBER:
				g.hashBlock(prop, wr, "wueste.NumberLiteral")
			case eg.BOOLEAN:
				g.hashBlock(prop, wr, "wueste.BoolLiteral")
			case eg.ARRAY:
				hashBlockArray(myOther{
					other: fmt.Sprintf("other.%s()", g.lang.PublicName(prop.Name())),
					my:    fmt.Sprintf("my.%s", g.lang.PrivateName(prop.Name())),
				}, prop.Property().(eg.PropertyArray).Items(), wr)

			case eg.OBJECT:
				panic("not implemented")
			default:
				panic("not implemented")
			}
		}
	})
	g.bodyWriter.WriteLine()
}

func (g *goGenerator) generateAsMapFunc() {
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) AsMap() map[string]interface{}", g.lang.PrivateName(g.schema.Title(), "Impl")), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteLine("res := map[string]interface{}{}")
		for _, prop := range g.schema.Items() {
			switch prop.Property().Type() {
			case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN, eg.ARRAY:
				if prop.Optional() {
					wr.WriteBlock("if", fmt.Sprintf("my.%s.IsNone()", g.lang.PrivateName(prop.Name())), func(wr *eg.ForIfWhileLangWriter) {
						wr.FormatLine("res[%s] = my.%s.Value()", wueste.QuoteString(prop.Name()), g.lang.PrivateName(prop.Name()))
					})
				} else {
					wr.FormatLine("res[%s] = my.%s", wueste.QuoteString(prop.Name()), g.lang.PrivateName(prop.Name()))
				}
			case eg.OBJECT:
				panic("not implemented")
			default:
				panic("not implemented")
			}
		}
		g.bodyWriter.WriteLine("return res")
	})
	g.bodyWriter.WriteLine()
}

func (g *goGenerator) generateImpl() {
	g.bodyWriter.WriteBlock("type", g.lang.PrivateName(g.schema.Title(), "Impl")+" struct", func(wr *eg.ForIfWhileLangWriter) {
		for _, prop := range g.schema.Items() {
			wr.FormatLine("%s %s", g.lang.PrivateName(prop.Name()), g.lang.AsTypeOptional(prop))
		}
	})
	g.bodyWriter.WriteLine()

	for _, prop := range g.schema.Items() {
		g.bodyWriter.WriteBlock("func", fmt.Sprintf("(my *%s) %s() %s",
			g.lang.PrivateName(g.schema.Title(), "Impl"), g.lang.PublicName(prop.Name()), g.lang.AsTypeOptional(prop)), func(wr *eg.ForIfWhileLangWriter) {
			wr.FormatLine("return my.%s", g.lang.PrivateName(prop.Name()))
		})
		g.bodyWriter.WriteLine()
	}
}

func (g *goGenerator) generateFactory() {
	g.bodyWriter.WriteBlock("type", g.lang.PublicName(g.schema.Title(), "Factory")+" struct", func(wr *eg.ForIfWhileLangWriter) {
	})
	g.bodyWriter.WriteLine()
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("New%s() *%s", g.lang.PublicName(g.schema.Title(), "Factory"), g.lang.PublicName(g.schema.Title(), "Factory")), func(wr *eg.ForIfWhileLangWriter) {
		wr.FormatLine("return &%s{}", g.lang.PublicName(g.schema.Title(), "Factory"))
	})
	g.bodyWriter.WriteLine()

	g.includes[RUSTY] = true
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(f *%s) Builder() *%s", g.lang.PublicName(g.schema.Title(), "Factory"),
		g.lang.PublicName(g.schema.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteLine(fmt.Sprintf("return New%s()", g.lang.PublicName(g.schema.Title(), "Builder")))
		// panic("not implemented")
		// wr.WriteBlock("if", "len(ps) == 0", func(wr *entity_generator.ForIfWhileLangWriter) {
		// 	w
		// })
		// for _, prop := range g.schema.Properties().Items() {
		// }
	})
	g.bodyWriter.WriteLine()

	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(f *%s) FromMap(m map[string]interface{}) rusty.Result[%s]", g.lang.PublicName(g.schema.Title(), "Factory"), g.lang.PublicName(g.schema.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {
		wr.FormatLine("my := f.Builder()")
		wr.FormatLine("var val interface{}")
		wr.FormatLine("var found bool")
		for _, prop := range g.schema.Items() {
			wr.FormatLine("val, found = m[%s]", wueste.QuoteString(prop.Name()))
			wr.WriteBlock("if", "found", func(wr *eg.ForIfWhileLangWriter) {
				if prop.Optional() {
					g.includes[RUSTY] = true
					wr.WriteLine("my." + g.lang.PublicName(prop.Name()) + fmt.Sprintf("(rusty.Some[%s](*val.(%s)))",
						g.lang.AsType(prop.Property()), g.lang.AsTypePtr(prop.Property())))
				} else {
					wr.WriteLine("my." + g.lang.PublicName(prop.Name()) + fmt.Sprintf("(val.(%s))", g.lang.AsTypePtr(prop.Property())))
				}
			})
		}
		wr.WriteLine("return my.ToClass()")
	})
	g.bodyWriter.WriteLine()
}

func toAttributeType[T string | int | uint | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64](lang ForIfWhileLang, prop eg.Property) string {
	// litProp, ok := prop.(eg.PropertyInteger[T])
	// if !ok {
	panic("not implemented")
	// }
	// attribute := fmt.Sprintf("wueste.Attribute[%s]", lang.AsTypeOptional(litProp))
	// if litProp.Optional() {
	// 	if litProp.Default().IsNone() {
	// 		// attribute = fmt.Sprintf("wueste.OptionalAttribute[%s]", g.lang.AsTypeOptional(litProp))
	// 		attribute = fmt.Sprintf("wueste.OptionalAttribute[%s]", lang.AsTypeOptional(litProp))
	// 	} else {
	// 		// attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s(%s))",
	// 		attribute = fmt.Sprintf("wueste.DefaultAttribute[%s]", lang.AsTypeOptional(litProp))
	// 		// g.lang.AsTypeOptional(litProp, "Some"),
	// 		// litProp.Default().Value().String())
	// 	}
	// } else {
	// 	if !litProp.Default().IsNone() {
	// 		// attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s)",
	// 		attribute = fmt.Sprintf("wueste.DefaultAttribute[%s]",
	// 			lang.AsTypeOptional(litProp))
	// 	} else {
	// 		// attribute = fmt.Sprintf("wueste.MustAttribute[%s]()", g.lang.AsTypeOptional(litProp))
	// 		attribute = fmt.Sprintf("wueste.MustAttribute[%s]", lang.AsTypeOptional(litProp))
	// 	}
	// }
	// return attribute
}

func toAttributeCreation[T string | int | uint | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64 | bool | float32 | float64](lang ForIfWhileLang, prop eg.Property) string {
	// litProp, ok := prop.(eg.PropertyLiteralType[T])
	// if !ok {
	panic("not implemented")
	// }
	// var attribute string
	// if litProp.Optional() {
	// 	if litProp.Default().IsNone() {
	// 		attribute = fmt.Sprintf("wueste.OptionalAttribute[%s]()", lang.AsTypeOptional(litProp))
	// 	} else {
	// 		attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s(%s))",
	// 			lang.AsTypeOptional(litProp),
	// 			lang.AsTypeOptional(litProp, "rusty.Some"),
	// 			*litProp.Default().Value().String())
	// 	}
	// } else {
	// 	if !litProp.Default().IsNone() {
	// 		attribute = fmt.Sprintf("wueste.DefaultAttribute[%s](%s)",
	// 			lang.AsTypeOptional(litProp), *litProp.Default().Value().String())
	// 	} else {
	// 		attribute = fmt.Sprintf("wueste.MustAttribute[%s]()", lang.AsTypeOptional(litProp))
	// 	}
	// }
	// return attribute
}

func toAttributeCreationInteger(lang ForIfWhileLang, prop eg.Property) string {

	_, ok := prop.(eg.PropertyInteger)
	if ok {
		return toAttributeCreation[int](lang, prop)
	}
	panic("not implemented")
}

func toAttributeCreationNumber(lang ForIfWhileLang, prop eg.Property) string {
	_, ok := prop.(eg.PropertyNumber)
	if ok {
		return toAttributeCreation[float64](lang, prop)
	}
	panic("not implemented")
}

func (g *goGenerator) toAttributeArray(prop eg.PropertyArray) string {
	// if prop.Optional() {
	// 	return fmt.Sprintf("wueste.DefaultAttribute[%s](%s())",
	// 		g.lang.Optional(true, fmt.Sprintf("[]%s", g.lang.AsTypeOptional(prop.Items()))),
	// 		g.lang.Optional(true, fmt.Sprintf("[]%s", g.lang.AsTypeOptional(prop.Items())), "None"))
	// }
	return fmt.Sprintf("wueste.DefaultAttribute[[]%s]([]%s{})", g.lang.AsType(prop.Items()), g.lang.AsType(prop.Items()))
}
func generateBuilderSetter(mo myOther, prop eg.Property, wr *eg.ForIfWhileLangWriter) {
	switch prop.Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		if mo.optional {
			wr.WriteBlock("if", fmt.Sprintf("%s.IsSome()", mo.my), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("*(%s.Value()).Set(v)", mo.my)
			})
		} else {
			wr.FormatLine("%s.Set(v)", mo.my)
		}
	case eg.ARRAY:
		generateArraySetBlock := func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter, blockFn func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter)) {
			loopI := fmt.Sprintf("i%d", len(mo.my))
			wr.WriteBlock("for", fmt.Sprintf("%s := 0; %s < len(%s); %s++", loopI, loopI, mo.my, loopI), func(wr *eg.ForIfWhileLangWriter) {
				blockFn(myOther{
					// other: fmt.Sprintf("[i]",
					my: fmt.Sprintf("%s[%s]", mo.my, loopI),
				}, item, wr)
			})
		}
		if mo.optional {
			wr.WriteBlock("if", fmt.Sprintf("%s.IsSome()", mo.my), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("v := *(%s.Value())", mo.my)
				generateArraySetBlock(mo, prop.(eg.PropertyArray).Items(), wr, func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
					generateBuilderSetter(mo, item, wr)
				})
			})
		} else {
			generateArraySetBlock(mo, prop.(eg.PropertyArray).Items(), wr, func(mo myOther, item eg.Property, wr *eg.ForIfWhileLangWriter) {
				generateBuilderSetter(mo, item, wr)
			})
		}
	default:
		panic("not implemented")
	}
}

func (g *goGenerator) genWeuesteAttributeCreation(prop eg.PropertyItem) string {
	switch prop.Property().Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		return fmt.Sprintf("wueste.Attribute[%s]", g.lang.AsTypeOptional(prop))
	case eg.ARRAY:
		return fmt.Sprintf("wueste.Attribute[%s]", g.genWuesteAttributeCreation(prop.(eg.PropertyArray).Items()))
	default:
		panic("not implemented")
	}
}

func (g *goGenerator) genWeuesteAttributeType(prop eg.PropertyItem) string {
	switch prop.Property().Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		return fmt.Sprintf("wueste.Attribute[%s]", g.lang.AsTypeOptional(prop))
	case eg.ARRAY:
		panic("not implemented")
		// return fmt.Sprintf("wueste.Attribute[%s]", g.genWeuesteAttributeType(prop.(eg.PropertyArray).Items()))
	default:
		panic("not implemented")
	}
}

func (g *goGenerator) genWuesteAttributeCreation(prop eg.Property) string {
	switch prop.Type() {
	case eg.STRING:
		g.includes[WUESTE] = true
		return toAttributeCreation[string](g.lang, prop.(eg.PropertyString))
	case eg.INTEGER:
		g.includes[WUESTE] = true
		return toAttributeCreationInteger(g.lang, prop)
	case eg.NUMBER:
		g.includes[WUESTE] = true
		return toAttributeCreationNumber(g.lang, prop)
	case eg.BOOLEAN:
		g.includes[WUESTE] = true
		return toAttributeCreation[bool](g.lang, prop.(eg.PropertyBoolean))
	case eg.ARRAY:
		// g.includes[WUESTE] = true
		panic("not implemented")
		// return g.genWeuesteAttributeCreation(prop.(eg.PropertyArray).Items())
		// toAttributeArray(prop.Property().(PropertyArray)))

	case eg.OBJECT:
		panic("not implemented")
	default:
		panic("not implemented")
	}
}

func (g *goGenerator) generateBuilder() {
	g.bodyWriter.WriteBlock("type", g.lang.PublicName(g.schema.Title(), "Builder")+" struct", func(wr *eg.ForIfWhileLangWriter) {
		g.includes[WUESTE] = true
		for _, prop := range g.schema.Items() {
			wr.FormatLine("%s %s", g.lang.PrivateName(prop.Name()), g.genWeuesteAttributeType(prop))
		}
	})
	g.bodyWriter.WriteLine()
	// in languages like TS we could pass a literal here.
	// TS Allows to type define Required Types
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("New%s() *%s", g.lang.PublicName(g.schema.Title(), "Builder"), g.lang.PublicName(g.schema.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteBlock(fmt.Sprintf("return &%s", g.lang.PublicName(g.schema.Title(), "Builder")), "", func(wr *eg.ForIfWhileLangWriter) {
			for _, prop := range g.schema.Items() {
				wr.FormatLine("%s: %s,", g.lang.PrivateName(prop.Name()), g.genWuesteAttributeCreation(prop.Property()))
			}
		})
		wr.WriteLine()
	})

	for _, prop := range g.schema.Items() {
		g.bodyWriter.WriteBlock("func", fmt.Sprintf("(b *%s) %s(v %s) *%s",
			g.lang.PublicName(g.schema.Title(), "Builder"), g.lang.PublicName(prop.Name()), g.lang.AsTypeOptional(prop), g.lang.PublicName(g.schema.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
			generateBuilderSetter(myOther{
				my:       "v", // fmt.Sprintf("s", g.lang.PrivateName(prop.Name())),
				optional: prop.Optional(),
			}, prop.Property(), wr)
			wr.WriteLine("return b")
		})
		g.bodyWriter.WriteLine()
	}

	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(b *%s) IsValid() rusty.Optional[error]",
		g.lang.PublicName(g.schema.Title(), "Builder")), func(wr *eg.ForIfWhileLangWriter) {
		g.includes["fmt"] = true
		g.includes[RUSTY] = true
		for _, prop := range g.schema.Items() {
			wr.WriteBlock("if", fmt.Sprintf("valid := b.%s.IsValid(); !valid.IsNone()", g.lang.PrivateName(prop.Name())), func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("return rusty.Some[error](fmt.Errorf(\"%s.%s:%%s\", (*valid.Value()).Error()))",
					g.lang.PublicName(g.schema.Title(), "Class"),
					g.lang.PublicName(prop.Name()))
			})
		}
		wr.FormatLine("return rusty.None[error]()")
	})
	g.bodyWriter.WriteBlock("func", fmt.Sprintf("(b *%s) ToClass() rusty.Result[%s]",
		g.lang.PublicName(g.schema.Title(), "Builder"), g.lang.PublicName(g.schema.Title(), "Class")), func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteBlock("if", "valid := b.IsValid(); !valid.IsNone()", func(wr *eg.ForIfWhileLangWriter) {
			wr.FormatLine("return rusty.Err[%s](*valid.Value())", g.lang.PublicName(g.schema.Title(), "Class"))
		})
		wr.WriteBlock(fmt.Sprintf("return rusty.Ok[%s](&"+g.lang.PrivateName(g.schema.Title(), "Impl"), g.lang.PublicName(g.schema.Title(), "Class")), "", func(wr *eg.ForIfWhileLangWriter) {
			for _, prop := range g.schema.Items() {
				wr.FormatLine("%s: b.%s.Get(),", g.lang.PrivateName(prop.Name()), g.lang.PrivateName(prop.Name()))
			}
		}, "{", "})")
	})
	g.bodyWriter.WriteLine()
}

type goGenerator struct {
	cfg    *eg.Config
	schema eg.PropertyObject
	lang   ForIfWhileLang
	// headerWriter *entity_generator.ForIfWhileLangWriter
	includes   map[string]bool
	bodyWriter *eg.ForIfWhileLangWriter
}

type bufLineWriter struct {
	lines []string
}

func (w *bufLineWriter) Write(p []byte) (n int,
	err error) {
	w.lines = append(w.lines,
		string(p))
	return len(p),
		nil
}

var reReplaceCaps = regexp.MustCompile(`[A-Z]+`)
var reReplaceNoAlpha = regexp.MustCompile(`[^a-zA-Z0-9]+`)
var reTrimNoAlpha = regexp.MustCompile(`^[^a-zA-Z0-9]+`)

func FileName(fname, suffix string) string {
	fname = reReplaceCaps.ReplaceAllString(fname, "_$0")
	fname = reTrimNoAlpha.ReplaceAllString(fname, "")
	fname = reReplaceNoAlpha.ReplaceAllString(fname, "_")
	return strings.ToLower(fname) + suffix
}

func GoGenerator(cfg *eg.Config, schema eg.PropertyObject, writer io.Writer) {
	// bodyLines := bufLineWriter{}
	g := &goGenerator{
		cfg:    cfg,
		schema: schema,
		// headerWriter: &goWriter{
		// 	writer: os.Stdout,
		// },
		includes: make(map[string]bool),
		bodyWriter: &eg.ForIfWhileLangWriter{
			OfsIndent: cfg.Indent,
			// Writer:    &bodyLines,
		},
	}

	g.generateClass()
	g.generateParam()
	g.generateJson()
	g.generateImpl()

	g.generateBuilder()

	g.generateCloneFunc()
	g.generateLessFunc()
	g.generateHashFunc()
	g.generateAsMapFunc()

	g.generateFactory()

	file := &eg.ForIfWhileLangWriter{
		// Writer: writer
	}

	file.FormatLine("package %s", cfg.PackageName)
	file.WriteLine()
	if len(g.includes) > 0 {
		gincs := make([]string, 0, len(g.includes))
		for include, _ := range g.includes {
			gincs = append(gincs, include)
		}
		sort.Strings(gincs)
		file.WriteLine("import (")
		for _, include := range gincs {
			file.FormatLine("%s\"%s\"", cfg.Indent, include)
		}
		file.WriteLine(")")
		file.WriteLine()
	}
	// for _, line := range bodyLines.lines {
	// file.Writer.Write([]byte(line))
	// }

}
