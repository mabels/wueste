package ts

import (
	"fmt"

	eg "github.com/mabels/wueste/entity-generator"
)

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
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflection")

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
	// attrType := []string{
	// 	g.lang.AsTypeNullable(pi.Property() /*WithOptional(pi.Optional())*/),
	// 	g.lang.AsTypeNullable(pi.Property(), WithAddCoerce() /*WithOptional(pi.Optional())*/),
	// }
	typ := g.lang.AsTypeNullable(pi.Property())
	coerceTyp := g.lang.AsTypeNullable(pi.Property(), WithAddCoerce())
	objectTyp := g.lang.AsTypeNullable(pi.Property())
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenBuilder")
	g.lang.Class(g.bodyWriter, "export ", g.lang.Implements(className,
		g.lang.Generics("WuestenBuilder", typ, coerceTyp, objectTyp)),
		prop,
		func(prop eg.PropertyItem, wr *eg.ForIfWhileLangWriter) {},
		func(wr *eg.ForIfWhileLangWriter) {

			attrib := g.lang.ReturnType(g.lang.OptionalParam(g.lang.PrivateName("value"), pi.Optional()), g.lang.AsType(pi.Property()))
			if !pi.Optional() {
				attrib += " = []"
			}
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenAttributeBase")

			wr.WriteLine("// WuestenAttributeParameter jojo")
			// boolean[][][][], WuesteCoerceTypeboolean[][][][]>

			wr.WriteLine(g.lang.Readonly(g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeBase", coerceTyp))))

			wr.WriteBlock("",
				g.lang.Call("constructor",
					g.lang.ReturnType("param", g.lang.Generics("WuestenAttributeBase",
						g.lang.AsType(getItemType(pa), WithAddCoerce()),
					)),
				), func(wr *eg.ForIfWhileLangWriter) {
					// pi := eg.NewPropertyArrayItem("ARRAY", rusty.Ok(getItemType(pa)), false).Ok()
					// attr := g.genWuesteBuilderAttribute("ARRAY", pi /*, func() string { return "param" } */)

					wr.WriteLine(g.lang.AssignDefault("this.param", "param"))

					// wr.WriteLine(g.lang.AssignDefault(g.lang.Const("itemAttr"), attr))

					// wr.WriteBlock("", "super({jsonname: param.jsonname, varname: param.varname, base: param.base}, {coerce: (c0: unknown) => ", func(wr *eg.ForIfWhileLangWriter) {
					// 	g.generateArrayCoerce(0, "c0", g.lang.AsType(pa), pa, wr)
					// 	wr.WriteLine(g.lang.Return(g.lang.Call("WuesteResult.Ok", "s0")))
					// }, " {", "}})")

				})

			wr.WriteBlock("",
				g.lang.ReturnType(g.lang.Call("Get"), g.lang.Generics("WuesteResult", typ)), func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(`throw new Error("Method not implemented.")`)
				})

			wr.WriteBlock("",
				g.lang.ReturnType(g.lang.Call("CoerceAttribute", g.lang.ReturnType("val", "unknown")), g.lang.Generics("WuesteResult", typ)), func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(`throw new Error("Method not implemented.")`)
				})

			wr.WriteBlock("",
				g.lang.ReturnType(g.lang.Call("Coerce", g.lang.ReturnType("val", coerceTyp)), className), func(wr *eg.ForIfWhileLangWriter) {
					// walk(val, (item) => {
					// 	if (Array.isArray(item)) {
					// 	  return
					// 	}
					// 	WuestenFormatter.Boolean.Formatter.Coerce(item)
					//   })
					//   return this
					// }
				})

			wr.WriteBlock("",
				g.lang.ReturnType(g.lang.Call("ToObject"), g.lang.Generics("WuesteResult", objectTyp)), func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(`throw new Error("Method not implemented.")`)
				})

			wr.WriteBlock("",
				g.lang.ReturnType(g.lang.Call("ToPayload"), g.lang.Generics("WuesteResult", "WuestePayload")), func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(`throw new Error("Method not implemented.")`)
				})
		})
	g.bodyWriter.WriteLine()
}

func (g *tsGenerator) generateSchemaExport(prop eg.Property, baseName string) {
	g.generateReflectionGetter(prop, baseName)
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflection")
	g.bodyWriter.WriteBlock(g.lang.Export(g.lang.Const(g.lang.ReturnType(
		g.lang.PublicName(baseName, "Schema"),
		"WuestenReflection"))), "", func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteLine(g.lang.Comma(g.lang.ReturnType(
			"ref", g.lang.Quote(g.lang.PublicName(baseName)))))
		g.writeSchema(wr, prop)
	}, " = {")
	// g.generateToObject(prop, baseName)
}
