package golang

import (
	"fmt"
	"regexp"
	"strings"

	eg "github.com/mabels/wueste/entity-generator"
)

type ForIfWhileLang struct {
	KeyWords map[string]bool
}

func (x *ForIfWhileLang) KeyWordFilter(prefix, name string) string {
	_, ok := x.KeyWords[name]
	if ok {
		return prefix + name
	}
	return name
}

var reSplitNonAllowed = regexp.MustCompile("[^a-zA-Z0-9]+")

func (x *ForIfWhileLang) PublicName(name string, add ...string) string {
	for _, a := range add {
		name = name + a
	}
	splitted := reSplitNonAllowed.Split(name, -1)
	clear := make([]string, 0, len(splitted))
	for _, s := range splitted {
		if len(s) != 0 {
			if len(s) == 1 {
				s = strings.ToUpper(s)
			} else {
				s = strings.ToUpper(s[0:1]) + s[1:]
			}
			clear = append(clear, s)
		}
	}
	name = strings.Join(clear, "")
	if len(clear) > 0 && !('A' <= name[0] && name[0] <= 'Z') {
		name = "X_" + name
	}
	name = x.KeyWordFilter("X_", name)
	return name
}

func (x *ForIfWhileLang) PrivateName(name string, add ...string) string {
	name = x.PublicName(name, add...)
	if len(name) > 1 {
		return x.KeyWordFilter("_", strings.ToLower(name[0:1])+name[1:])
	} else {
		return x.KeyWordFilter("_", strings.ToLower(name))
	}
}

func (x *ForIfWhileLang) Optional(optional bool, typeStr string, opts ...string) string {
	if optional {
		opt := "Optional"
		if len(opts) > 0 {
			opt = opts[0]
		}
		return fmt.Sprintf("rusty.%s[%s]", opt, typeStr)
	}
	return typeStr
}

func (x *ForIfWhileLang) Ptr(p eg.Property, typeStr string) string {
	if p.Optional() {
		return fmt.Sprintf("*%s", typeStr)
	}
	return typeStr
}

func (x *ForIfWhileLang) AsType(p eg.Property) string {
	switch p.Type() {
	case eg.STRING:
		return "string"
	case eg.INTEGER:
		return "int64"
	case eg.NUMBER:
		return "float64"
	case eg.BOOLEAN:
		return "bool"
	case eg.ARRAY:
		p := p.(eg.PropertyArray)
		return "[]" + x.AsTypePtr(p.Items())
	case eg.OBJECT:
		// p := p.(Schema)
		// required(p, publicClassName(p.Title()))
		panic("not implemented")
	default:
		panic("not implemented")
	}
}

func (x *ForIfWhileLang) AsTypePtr(p eg.Property) string {
	switch p.Type() {
	case eg.STRING:
		return x.Ptr(p, "string")
	case eg.INTEGER:
		return x.Ptr(p, "int64")
	case eg.NUMBER:
		return x.Ptr(p, "float64")
	case eg.BOOLEAN:
		return x.Ptr(p, "bool")
	case eg.ARRAY:
		p := p.(eg.PropertyArray)
		return x.Ptr(p, "[]"+x.AsTypePtr(p.Items()))
	case eg.OBJECT:
		// p := p.(Schema)
		// required(p, publicClassName(p.Title()))
		panic("not implemented")
	default:
		panic("not implemented")
	}
}

func (x *ForIfWhileLang) AsTypeOptional(p eg.Property, opts ...string) string {
	switch p.Type() {
	case eg.STRING:
		return x.Optional(p.Optional(), "string", opts...)
	case eg.INTEGER:
		return x.Optional(p.Optional(), "int64", opts...)
	case eg.NUMBER:
		return x.Optional(p.Optional(), "float64", opts...)
	case eg.BOOLEAN:
		return x.Optional(p.Optional(), "bool", opts...)
	case eg.ARRAY:
		p := p.(eg.PropertyArray)
		return x.Optional(p.Optional(), "[]"+x.AsTypeOptional(p.Items()))
	case eg.OBJECT:
		// p := p.(Schema)
		// required(p, publicClassName(p.Title()))
		panic("not implemented")
	default:
		panic("not implemented")
	}
}
