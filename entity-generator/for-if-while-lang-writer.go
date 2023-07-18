package entity_generator

import (
	"fmt"
	"strings"
)

type BufLineWriter struct {
	Lines []string
}

func (w *BufLineWriter) Write(p []byte) (n int,
	err error) {
	w.Lines = append(w.Lines,
		string(p))
	return len(p),
		nil
}

type ForIfWhileLangWriter struct {
	OfsIndent string
	writer    *BufLineWriter
	current   string
	// ExpressionBracketOpen  string
	// ExpressionBracketClose string
}

func NewForIfWhileLangWriter(f ForIfWhileLangWriter) *ForIfWhileLangWriter {
	return &ForIfWhileLangWriter{
		writer:    &BufLineWriter{},
		OfsIndent: f.OfsIndent,
		current:   "",
	}
}

func (w *ForIfWhileLangWriter) Lines() []string {
	return w.writer.Lines
}

func (w *ForIfWhileLangWriter) Indent() *ForIfWhileLangWriter {
	return &ForIfWhileLangWriter{
		writer:    w.writer,
		OfsIndent: w.OfsIndent,
		current:   w.current + w.OfsIndent,
		// ExpressionBracketOpen:  w.ExpressionBracketOpen,
		// ExpressionBracketClose: w.ExpressionBracketClose,
	}
}

func (w *ForIfWhileLangWriter) WriteLine(lines ...string) {
	for _, line := range lines {
		my := w.current + line
		if len(strings.TrimSpace(my)) != 0 {
			fmt.Fprintln(w.writer, my)
		} else {
			fmt.Fprintln(w.writer)
		}
	}
	if len(lines) == 0 {
		fmt.Fprintln(w.writer)
	}
}

func (w *ForIfWhileLangWriter) FormatLine(format string, args ...interface{}) string {
	out := fmt.Sprintln(w.current + fmt.Sprintf(format, args...))
	w.writer.Write([]byte(out))
	return out
}

func (w *ForIfWhileLangWriter) WriteBlock(keyword, expression string, fn func(wr *ForIfWhileLangWriter), opts ...string) string {
	opening := " {"
	if len(opts) > 0 {
		opening = opts[0]
	}
	closing := "}"
	if len(opts) > 1 {
		closing = opts[1]
	}
	space := ""
	if len(keyword) != 0 && len(expression) != 0 {
		space = " "
	}
	out := w.FormatLine("%s%s%s%s", strings.TrimSpace(keyword), space, strings.TrimSpace(expression), opening)
	fn(w.Indent())
	w.WriteLine(closing)
	return out
}

func (w *ForIfWhileLangWriter) WriteIf(expression string, fnTrue func(wr *ForIfWhileLangWriter), fnFalses ...func(wr *ForIfWhileLangWriter)) {
	w.FormatLine("if %s {", expression)
	fnTrue(w.Indent())
	if len(fnFalses) > 0 {
		w.WriteLine("} else {")
		fnFalses[0](w.Indent())
	}
	w.WriteLine("}")
}

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

// var reSplitNonAllowed = regexp.MustCompile("[^a-zA-Z0-9]+")

// func (x *ForIfWhileLang) PublicName(name string, add ...string) string {
// 	for _, a := range add {
// 		name = name + a
// 	}
// 	splitted := reSplitNonAllowed.Split(name, -1)
// 	clear := make([]string, 0, len(splitted))
// 	for _, s := range splitted {
// 		if len(s) != 0 {
// 			if len(s) == 1 {
// 				s = strings.ToUpper(s)
// 			} else {
// 				s = strings.ToUpper(s[0:1]) + s[1:]
// 			}
// 			clear = append(clear, s)
// 		}
// 	}
// 	name = strings.Join(clear, "")
// 	if len(clear) > 0 && !('A' <= name[0] && name[0] <= 'Z') {
// 		name = "X_" + name
// 	}
// 	name = x.KeyWordFilter("X_", name)
// 	return name
// }

// func (x *ForIfWhileLang) PrivateName(name string, add ...string) string {
// 	name = x.PublicName(name, add...)
// 	if len(name) > 1 {
// 		return x.KeyWordFilter("_", strings.ToLower(name[0:1])+name[1:])
// 	} else {
// 		return x.KeyWordFilter("_", strings.ToLower(name))
// 	}
// }

// func (x *ForIfWhileLang) Optional(optional bool, typeStr string, opts ...string) string {
// 	if optional {
// 		opt := "Optional"
// 		if len(opts) > 0 {
// 			opt = opts[0]
// 		}
// 		return fmt.Sprintf("rusty.%s[%s]", opt, typeStr)
// 	}
// 	return typeStr
// }

// func (x *ForIfWhileLang) Ptr(p Property, typeStr string) string {
// 	if p.Optional() {
// 		return fmt.Sprintf("*%s", typeStr)
// 	}
// 	return typeStr
// }

// func (x *ForIfWhileLang) AsType(p Property) string {
// 	switch p.Type() {
// 	case STRING:
// 		return "string"
// 	case INTEGER:
// 		return "int64"
// 	case NUMBER:
// 		return "float64"
// 	case BOOLEAN:
// 		return "bool"
// 	case ARRAY:
// 		p := p.(PropertyArray)
// 		return "[]" + x.AsTypePtr(p.Items())
// 	case OBJECT:
// 		// p := p.(Schema)
// 		// required(p, publicClassName(p.Title()))
// 		panic("not implemented")
// 	default:
// 		panic("not implemented")
// 	}
// }

// func (x *ForIfWhileLang) AsTypePtr(p Property) string {
// 	switch p.Type() {
// 	case STRING:
// 		return x.Ptr(p, "string")
// 	case INTEGER:
// 		return x.Ptr(p, "int64")
// 	case NUMBER:
// 		return x.Ptr(p, "float64")
// 	case BOOLEAN:
// 		return x.Ptr(p, "bool")
// 	case ARRAY:
// 		p := p.(PropertyArray)
// 		return x.Ptr(p, "[]"+x.AsTypePtr(p.Items()))
// 	case OBJECT:
// 		// p := p.(Schema)
// 		// required(p, publicClassName(p.Title()))
// 		panic("not implemented")
// 	default:
// 		panic("not implemented")
// 	}
// }

// func (x *ForIfWhileLang) AsTypeOptional(p Property, opts ...string) string {
// 	switch p.Type() {
// 	case STRING:
// 		return x.Optional(p.Optional(), "string", opts...)
// 	case INTEGER:
// 		return x.Optional(p.Optional(), "int64", opts...)
// 	case NUMBER:
// 		return x.Optional(p.Optional(), "float64", opts...)
// 	case BOOLEAN:
// 		return x.Optional(p.Optional(), "bool", opts...)
// 	case ARRAY:
// 		p := p.(PropertyArray)
// 		return x.Optional(p.Optional(), "[]"+x.AsTypeOptional(p.Items()))
// 	case OBJECT:
// 		// p := p.(Schema)
// 		// required(p, publicClassName(p.Title()))
// 		panic("not implemented")
// 	default:
// 		panic("not implemented")
// 	}
// }
