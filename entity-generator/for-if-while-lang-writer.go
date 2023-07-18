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
