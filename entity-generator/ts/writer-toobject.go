package ts

import (
	"fmt"

	eg "github.com/mabels/wueste/entity-generator"
)

func (g *tsGenerator) writeArrayToObject(wr *eg.ForIfWhileLangWriter, prop eg.Property, l int) {
	switch prop.Type() {
	case eg.OBJECT:
		if l == 0 {
			g.writeObjectToObject(wr, prop.(eg.PropertyObject))
		} else {
			name := g.lang.PublicName(getObjectName(prop), "ToObject")
			g.includes.AddProperty(name, prop)
			wr.FormatLine(g.lang.Const(g.lang.AssignDefault(
				fmt.Sprintf("o%d", l), g.lang.Call(name, fmt.Sprintf("v%d", l)))))
		}
	case eg.ARRAY:
		pa := prop.(eg.PropertyArray)
		tname := g.lang.AsType(pa)
		if pa.Items().Type() == eg.OBJECT {
			g.includes.AddProperty(g.lang.AsType(pa, WithTypeSuffix("Object")), prop)
			tname = g.lang.AsType(pa, WithTypeSuffix("Object"))
		}
		wr.FormatLine(
			g.lang.AssignDefault(
				g.lang.Const(
					g.lang.ReturnType(fmt.Sprintf("o%d", l), tname)),
				" []"))
		wr.WriteBlock("for",
			fmt.Sprintf("(let i%d = 0; i%d < v%d.length; i%d++)", l+1, l+1, l, l+1),
			func(wr *eg.ForIfWhileLangWriter) {
				wr.FormatLine("const v%d = v%d[i%d]", l+1, l, l+1)
				g.writeArrayToObject(wr, prop.(eg.PropertyArray).Items(), l+1)
				wr.FormatLine(
					g.lang.Call(
						fmt.Sprintf("o%d.push", l), fmt.Sprintf("o%d", l+1)))
			})
		if l == 0 {
			wr.FormatLine("return o0")
		}
	case eg.STRING, eg.BOOLEAN, eg.INTEGER, eg.NUMBER:
		wr.FormatLine(g.lang.Const(g.lang.AssignDefault(
			fmt.Sprintf("o%d", l), fmt.Sprintf("v%d", l))))
	default:
		panic("writeArrayToObject not implemented")
	}
}

func (g *tsGenerator) generateToObject(prop eg.Property, baseName string) {
	var rType string
	switch prop.Type() {
	case eg.OBJECT:
		rType = g.lang.PublicName(g.lang.AsType(prop), "Object")
	case eg.ARRAY:
		pa := prop.(eg.PropertyArray)
		rType = g.lang.AsType(prop)
		if pa.Items().Type() == eg.OBJECT {
			g.includes.AddProperty(g.lang.AsType(pa.Items(), WithTypeSuffix("Object")), pa.Items())
			rType = g.lang.AsType(pa, WithTypeSuffix("Object"))
		}
	default:
		panic("generateToObject not implemented")
	}
	g.bodyWriter.WriteBlock(g.lang.Export(
		g.lang.ReturnType(
			g.lang.Call("function "+g.lang.PublicName(baseName, "ToObject"),
				g.lang.ReturnType("v0", g.lang.AsType(prop))), rType)), "", func(wr *eg.ForIfWhileLangWriter) {
		// wr.FormatLine(
		// 	g.lang.AssignDefault(
		// 		g.lang.Const(g.lang.ReturnType("o0", g.lang.AsType(prop))), " []"))
		g.writeArrayToObject(wr, prop, 0)
		// wr.FormatLine("return o0")
		//   return out
	})

}

func (g *tsGenerator) generateObjectToObject(pi eg.PropertyItem) string {
	switch pi.Property().Type() {
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		return g.lang.CallDot("v0", g.lang.PublicName(pi.Name()))
	case eg.ARRAY:
		name := g.lang.PublicName(getObjectName(pi.Property(), []string{pi.Name()}), "ToObject")
		g.includes.AddProperty(name, pi.Property())
		return g.lang.Call(name,
			g.lang.CallDot("v0", g.lang.PublicName(pi.Name()))+"/* Choose Factory */")
	case eg.OBJECT:
		if pi.Property().(eg.PropertyObject).Title() != "" {
			name := g.lang.PublicName(getObjectName(pi.Property()), "Factory")
			g.includes.AddProperty(name, pi.Property())
			return g.lang.CallDot(name, g.lang.Call("ToObject",
				g.lang.CallDot("v0", g.lang.PublicName(pi.Name()))))
		} else {
			// TODO OpenObject
			return g.lang.CallDot("v0", g.lang.PublicName(pi.Name()))
		}
	default:
		panic("not implemented")
	}
}

func (g *tsGenerator) writeObjectToObject(wr *eg.ForIfWhileLangWriter, prop eg.PropertyObject) {
	wr.WriteLine("const ret: Record<string, unknown> = {}")
	for _, pi := range prop.Items() {
		if !pi.Optional() {
			wr.FormatLine("ret[%s] = %s", g.lang.Quote(pi.Name()), g.generateObjectToObject(pi))
			continue
		}
		wr.WriteBlock("if ", fmt.Sprintf("(typeof %s !== 'undefined')",
			g.lang.CallDot("v0", g.lang.PublicName(pi.Name()))), func(wr *eg.ForIfWhileLangWriter) {
			wr.FormatLine("ret[%s] = %s", g.lang.Quote(pi.Name()), g.generateObjectToObject(pi))
		})
	}
	wr.FormatLine("return ret as unknown as %s;", g.lang.PublicName(getObjectName(prop), "Object"))
}
