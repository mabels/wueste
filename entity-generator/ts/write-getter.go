package ts

import (
	"fmt"
	"regexp"
	"strings"

	eg "github.com/mabels/wueste/entity-generator"
	"github.com/mabels/wueste/entity-generator/rusty"
)

type vname struct {
	init  string
	local int
}

var removeIndex = regexp.MustCompile(`\[[^\]]+\]`)

func (v *vname) newContext(init string) vname {
	v.local = v.local + 1
	return vname{
		init:  init,
		local: v.local,
	}
}

func (v vname) contextVar() string {
	return v.init
}

func (v *vname) newVar() string {
	ret := fmt.Sprintf("%s%d", removeIndex.ReplaceAllString(v.init, ""), v.local)
	v.local++
	return ret
}

// func (v *vname) get() string {
// 	if v.idx == 0 {
// 		return "v"
// 	}
// 	return fmt.Sprintf("_v%d", v.idx)
// }

// func (v *vname) inc() *vname {
// 	v.idx++
// 	return v
// }

func (g *tsGenerator) generateReflectionGetter(prop propertyValue, baseName string) {
	// g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenGetterFn").activated = true
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenGetterBuilder")
	g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionValue")
	g.bodyWriter.WriteBlock("export function ",
		g.lang.ReturnType(
			g.lang.Call(
				g.lang.PublicName(baseName, "Getter"),
				g.lang.ReturnType("v", g.lang.AsTypeNullable(prop.prop /*WithOptional(pi.Optional())*/)),
				g.lang.ReturnType("base", "WuestenReflectionValue[] = []")),
			"WuestenGetterBuilder"),
		func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteBlock("return new WuestenGetterBuilder((fn) => ", "", func(wr *eg.ForIfWhileLangWriter) {
				g.writeReflectionGetter(wr, baseName, vname{
					init: "v",
				}, prop, []propertyValue{prop})
			}, "{", "})")
		})
}

func (g *tsGenerator) getLastPath(level []propertyValue) string {
	// for (let i = 0; i < level.length; i++) {
	if len(level) >= 1 {
		return g.getPath(level[len(level)-1])
	}
	return ""
}

func (g *tsGenerator) getLastName(level []propertyValue) string {
	// for (let i = 0; i < level.length; i++) {
	if len(level) >= 1 {
		pi, ok := level[len(level)-1].prop.(eg.PropertyItem)
		if ok {
			switch pi.Type() {
			case eg.ARRAYITEM:
				return pi.Name()
			default:
				return g.lang.PublicName(pi.Name())
			}
		}
	}
	return ""
}

func (g *tsGenerator) getPath(c propertyValue) string {
	out := []string{}
	pi, ok := c.prop.(eg.PropertyItem)
	if ok {
		if strings.HasPrefix(pi.Name(), "[") {
			out = append(out, pi.Name())
		} else {
			out = append(out, g.lang.PublicName(pi.Name()))
		}
	}
	return strings.Join(out, "")

}

// func (g *tsGenerator) getFullPath(level []eg.Property) string {
// 	out := make([]string, 0, len(level))
// 	for i := 1; i < len(level); i++ {
// 		out = append(out, g.getPath(level[i]))
// 	}
// 	return strings.Join(out, "")
// }

// func curVname(level []eg.Property) string {
// 	if len(level) <= 1 {
// 		return "v"
// 	}
// 	lastLevel := level[len(level)-1]
// 	if pi, ok := lastLevel.(eg.PropertyItem); ok {
// 		return fmt.Sprintf("_v%d_%d", len(level)-1, pi.Idx())
// 	}
// 	return fmt.Sprintf("_v%d", len(level)-1)
// }

// func nextVname(level []eg.Property) string {
// 	if len(level) <= 1 {
// 		return "_v1"
// 	}
// 	lastLevel := level[len(level)-1]
// 	if pi, ok := lastLevel.(eg.PropertyItem); ok {
// 		return fmt.Sprintf("_v%d_%d", len(level), pi.Idx())
// 	}
// 	return fmt.Sprintf("_v%d", len(level))
// }

func (g *tsGenerator) toPropLevel(baseName string, ps []propertyValue) string {
	if len(ps) < 1 {
		return ""
	}
	p := ps[len(ps)-1].prop
	tail := ps[0 : len(ps)-1]
	if len(tail) == 0 {
		return g.lang.PublicName(baseName, "Schema")
	}
	switch p.Type() {
	case eg.OBJECT:
		po := p.(eg.PropertyObject)
		if po.Properties().Len() > 0 {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionObject")
			return fmt.Sprintf("(%s as WuestenReflectionObject)", g.toPropLevel(baseName, tail))
		}
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenRecordGetter")
		return fmt.Sprintf("WuestenRecordGetter(%s)", g.toPropLevel(baseName, tail))
	case eg.OBJECTITEM:
		pi := p.(eg.PropertyItem)
		property := ""
		// if len(ps) > 2 {
		// 	property = ".property"
		// }
		return fmt.Sprintf("(%s%s as WuestenReflectionObject).properties![%d]",
			g.toPropLevel(baseName, tail), property, pi.Idx())
	case eg.ARRAYITEM:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionArray")
		return fmt.Sprintf("(%s as WuestenReflectionArray).items", g.toPropLevel(baseName, tail))
	case eg.ARRAY:
		g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionArray")
		return fmt.Sprintf("(%s as WuestenReflectionArray)", g.toPropLevel(baseName, tail))
	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		return fmt.Sprintf("%s.%s", g.toPropLevel(baseName, tail), p.Type())
	}
	panic("not implemented")
}

type propertyValue struct {
	prop    eg.Property
	varname string
}

func (g *tsGenerator) writePath(wr *eg.ForIfWhileLangWriter, baseName string, path []propertyValue) {
	wr.WriteBlock("", "", func(wr *eg.ForIfWhileLangWriter) {
		wr.WriteLine(g.lang.Comma("...base"))
		for i, l := range path {
			wr.WriteBlock("", "", func(wr *eg.ForIfWhileLangWriter) {
				wr.WriteLine(g.lang.Comma(g.lang.ReturnType("schema", g.toPropLevel(baseName, path[:i+1]))))
				wr.WriteLine(g.lang.ReturnType("value", l.varname))
			}, "{", "},")
		}
	}, "[", "]")
}

func toPropertValue(level []propertyValue, varname string) []propertyValue {
	varnames := []string{varname}
	ret := make([]propertyValue, 0, len(level))
	ret = append(ret, level...)
	startidx := len(level) - len(varnames)
	if startidx < 0 {
		panic("not enough levels")
	}
	for i := 0; i < len(varnames); i++ {
		ret[startidx+i].varname = varnames[i]
	}
	return ret
}

func (g *tsGenerator) writeReflectionGetter(wr *eg.ForIfWhileLangWriter, baseName string, vname vname, pv propertyValue, path []propertyValue) {
	prop := pv.prop
	// isOptional := false
	// if len(level) > 0 {
	// 	la := level[len(level)-1]
	// 	if pi, ok := la.(eg.PropertyItem); ok && pi.Optional() {
	// 		isOptional = true
	// 	}
	// }
	switch prop.Type() {
	case eg.OBJECT:
		nextWithVar := vname.newVar()
		po := prop.(eg.PropertyObject)
		if po.Properties().Len() > 0 {
			for _, l := range path {
				name := "NoName"
				if pi, ok := l.prop.(eg.PropertyItem); ok {
					name = pi.Name()
				}
				wr.FormatLine("/* %s:%v:%s */", l.prop.Type(), name, g.getLastName(path))
			}
			if len(path) == 0 {
				panic("there must be a level")
			}
			lastLevel := path[len(path)-1]
			if lastLevel.prop.Type() == eg.OBJECTITEM {
				// SimpleType$PayloadGetter(v0.sub, [])
				getterName := g.lang.PublicName(getObjectName(prop), "Getter")
				g.includes.AddProperty(getterName, prop)
				// wr.WriteLine(
				// 	g.lang.Call(getterName,
				// 		g.lang.CallDot(vname.contextVar(), g.getLastPath(level))))

				wr.WriteBlock(
					getterName,
					"",
					func(wr *eg.ForIfWhileLangWriter) {
						wr.WriteLine(g.lang.Comma(vname.contextVar()))
						g.writePath(wr, baseName, toPropertValue(path, vname.contextVar()))
					}, "(", ").Apply(fn)")

			} else {
				// if level[0].Type() != eg.OBJECT {
				// 	wr.FormatLine("/* ZZZ %s:%s */", g.getLastName(level), level[0].Type())
				// 	getterName := g.lang.PublicName(getObjectName(prop), "Getter")
				// 	g.includes.AddProperty(getterName, prop)
				// 	wr.FormatLine("// !Object %s", getterName)
				// 	wr.WriteLine(
				// 		g.lang.Call(
				// 			g.lang.CallDot(
				// 				g.lang.Call(getterName,
				// 					g.lang.CallDot(vname.contextVar(), g.getLastPath(level)), "base"), "Apply"), "fn"))
				// 	return
				// }
				wr.FormatLine("const %s = %s", nextWithVar, vname.contextVar())
				path[len(path)-1].varname = nextWithVar
				g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenReflectionObject")
				for i, pi := range po.Items() {
					pc := eg.NewPropertyObjectItem(pi.Name(), rusty.Ok(pi.Property()), i, pi.Optional()).Ok()
					idx := g.lang.CallDot(nextWithVar, g.lang.PublicName(pi.Name()))
					cp := propertyValue{prop: pi.Property(), varname: idx}
					wrapOptional(pi.Optional(), g.lang.CallDot(nextWithVar, g.lang.PublicName(pi.Name())), wr, func(wr *eg.ForIfWhileLangWriter) {
						g.writeReflectionGetter(wr, baseName, vname.newContext(idx), cp, append(path, propertyValue{prop: pc, varname: idx}))
					})
					//   out.push(gen(pi.property, vname.inc(), `${nextWithVar}`, [...level, { ...prop.properties[i], type: 'objectItem'}]));
				}
			}
		} else {
			g.includes.AddType(g.cfg.EntityCfg.FromWueste, "WuestenRecordGetter")
			wr.WriteBlock("WuestenRecordGetter(", "fn, ", func(wr *eg.ForIfWhileLangWriter) {
				g.writePath(wr, baseName, toPropertValue(path, vname.contextVar()))
			}, "", ")")
		}

	case eg.ARRAY:
		// wrapOptional(isOptional, g.lang.CallDot(vname.contextVar(), g.getLastPath(level)), wr, func(wr *eg.ForIfWhileLangWriter) {
		// getObjectName(pi.Property(), []string{pi.Name()})
		if path[0].prop.Type() != eg.ARRAY {
			// wr.FormatLine("/* XX %s:%s */", prop.Type(), level[0].Type())
			getterName := g.lang.PublicName(getObjectName(prop, []string{g.getLastName(path)}), "Getter")
			// if len(level) > 1 {
			// my := level[len(level)-2]
			g.includes.AddProperty(getterName, prop)
			wr.FormatLine("// Array %s --- %s --- %v", getterName, getObjectName(prop), prop.Meta().Parent().Value().Type())
			wr.WriteBlock(
				getterName,
				"",
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine(g.lang.Comma(vname.contextVar()))
					g.writePath(wr, baseName, toPropertValue(path, vname.contextVar()))
				}, "(", ").Apply(fn)")
			return
		}
		nextWithVar := vname.newVar()
		wr.FormatLine("const %s = %s", nextWithVar, vname.contextVar())
		path[len(path)-1].varname = nextWithVar
		forLine := fmt.Sprintf("for (let i%s = 0; i%s < %s.length; i%s++)", nextWithVar, nextWithVar, nextWithVar, nextWithVar)
		wr.WriteBlock(forLine, "", func(wr *eg.ForIfWhileLangWriter) {
			idx := fmt.Sprintf("%s[i%s]", nextWithVar, nextWithVar)
			// wr.FormatLine("const %s = %s /* pathlen:%d */", nextWithVar, idx, len(path))
			pa := prop.(eg.PropertyArray)
			pi := propertyValue{prop: eg.NewPropertyArrayItem(idx, rusty.Ok(pa.Items())).Ok(), varname: idx}
			cp := propertyValue{prop: pa.Items(), varname: idx}
			// pi := propertyValue{prop: pa.Items(), varname: idx}
			wr.FormatLine("/* HO %s:%s */", pi.prop.Type(), idx)
			g.writeReflectionGetter(wr, baseName, vname.newContext(idx), cp, append(path, pi))
		})
		// })
	// case eg.ARRAYITEM:
	// 	wr.WriteLine("/* ARRAYITEM */")
	// 	pai := prop.(eg.PropertyItem)
	// 	// path[len(path)-1].prop = pai.Property()
	// 	g.writeReflectionGetter(wr, baseName, vname, propertyValue{prop: pai.Property(), varname: pv.varname}, path)

	case eg.STRING, eg.INTEGER, eg.NUMBER, eg.BOOLEAN:
		wr.WriteBlock("", "fn(", func(wr *eg.ForIfWhileLangWriter) {
			path[len(path)-1].varname = vname.contextVar()
			g.writePath(wr, baseName, path)
		}, "", ")")
	default:
		panic("writeReflectionGetter not implemented:" + prop.Type())
	}

}

func wrapOptional(opt bool, varName string, wr *eg.ForIfWhileLangWriter, fn func(wr *eg.ForIfWhileLangWriter)) {
	if opt {
		wr.WriteBlock("if (", fmt.Sprintf("!!%s && %s !== null)", varName, varName), fn)
	} else {
		fn(wr)
	}
}
