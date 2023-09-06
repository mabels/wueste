package golang

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	eg "github.com/mabels/wueste/entity-generator"
	"github.com/mabels/wueste/entity-generator/rusty"
	"github.com/stretchr/testify/assert"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func TestWriterWriteLineEmpty(t *testing.T) {
	wt := bufLineWriter{}
	w := eg.ForIfWhileLangWriter{
		// Writer:    &wt,
		OfsIndent: "  ",
	}
	w.WriteLine("", "\t")
	for len(wt.lines) != 2 {
		t.Error("writeLine failed")
	}
	assert.Equal(t, []string{"\n", "\n"}, wt.lines)
}

func TestWriterWriteLine(t *testing.T) {
	wt := bufLineWriter{}
	w := (&eg.ForIfWhileLangWriter{
		// Writer:    &wt,
		OfsIndent: "\t",
	}).Indent()
	w.WriteLine("Hello",
		"World")
	for len(wt.lines) != 2 {
		t.Error("writeLine failed")
	}
	assert.Equal(t,
		[]string{"\tHello\n",
			"\tWorld\n"},
		wt.lines)
}

func TestFormatLine(t *testing.T) {
	wt := bufLineWriter{}
	w := (&eg.ForIfWhileLangWriter{
		// Writer:    &wt,
		OfsIndent: "\t",
	}).Indent()
	w.FormatLine("%s %s",
		"Hello",
		"World")
	for len(wt.lines) != 1 {
		t.Error("writeLine failed")
	}
	assert.Equal(t,
		[]string{"\tHello World\n"},
		wt.lines)
}

func TestWriteBlock(t *testing.T) {
	wt := bufLineWriter{}
	w := (&eg.ForIfWhileLangWriter{
		// Writer:    &wt,
		OfsIndent: "\t",
	}).Indent()
	w.WriteBlock("Level",
		"I",
		func(wr *eg.ForIfWhileLangWriter) {
			wr.WriteLine("Start-Level-I")
			wr.WriteBlock("Level",
				"II",
				func(wr *eg.ForIfWhileLangWriter) {
					wr.WriteLine("Inside-Level-II")
				})
			wr.WriteLine("Close-Level-I")
		})
	assert.Equal(t,
		[]string{
			"\tLevel I {\n",

			"\t\tStart-Level-I\n",

			"\t\tLevel II {\n",

			"\t\t\tInside-Level-II\n",

			"\t\t}\n",

			"\t\tClose-Level-I\n",

			"\t}\n",
		},
		wt.lines)
}

func TestKeyWordFilter(t *testing.T) {
	kf := ForIfWhileLang{KeyWords: KeyWords}
	assert.Equal(t, "_func", kf.KeyWordFilter("func", "_"))
	assert.Equal(t, "fUnc", kf.KeyWordFilter("fUnc", "_"))
	assert.Equal(t, "_uint32", kf.KeyWordFilter("uint32", "_"))
}

func TestPublicName(t *testing.T) {
	lang := ForIfWhileLang{KeyWords: KeyWords}
	p := lang.PublicName("")
	assert.Equal(t, "", p)
	p = lang.PublicName("a")
	assert.Equal(t, "A", p)
	p = lang.PublicName("ab")
	assert.Equal(t, "Ab", p)
	p = lang.PublicName("ab-cd@ef_GH.Xo")
	assert.Equal(t, "AbCdEfGHXo", p)
	p = lang.PublicName("0ab-cd@ef_GH.Xo")
	assert.Equal(t, "X_0abCdEfGHXo", p)
	p = lang.PublicName("uint32")
	assert.Equal(t, "Uint32", p)
}

func TestPrivateName(t *testing.T) {
	lang := ForIfWhileLang{KeyWords: KeyWords}
	p := lang.PrivateName("")
	assert.Equal(t, "", p)
	p = lang.PrivateName("a")
	assert.Equal(t, "a", p)
	p = lang.PrivateName("ab")
	assert.Equal(t, "ab", p)
	p = lang.PrivateName("Ab-cd@ef_GH.Xo")
	assert.Equal(t, "abCdEfGHXo", p)
	p = lang.PrivateName("0ab-cd@ef_GH.Xo")
	assert.Equal(t, "_0abCdEfGHXo", p)
	p = lang.PrivateName("uint32")
	assert.Equal(t, "_uint32", p)
}

func TestSimpleTypeClazz(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema: eg.TestSchema(sl, eg.PropertyRuntime{}).(eg.PropertyObject),

		includes: make(map[string]bool),

		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateClass()
	assert.Equal(t,
		[]string{
			"type SimpleTypeClass interface {\n",
			"\tString() string\n",
			"\tOptionalString() rusty.Optional[string]\n",
			"\tFloat64() float64\n",
			"\tOptionalFloat32() rusty.Optional[float64]\n",
			"\tInt64() int64\n",
			"\tOptionalInt32() rusty.Optional[int64]\n",
			"\tBool() bool\n",
			"\tOptionalBool() rusty.Optional[bool]\n",
			"}\n",
			"\n"},
		wt.lines)
}

type SimpleTypeObject interface {
	String() string
	OptionalString() rusty.Optional[string]
	Float64() float64
	OptionalFloat64() rusty.Optional[float64]
	Int64() int64
	OptionalInt64() rusty.Optional[int64]
	Uint64() uint64
	OptionalUint64() rusty.Optional[uint64]
	Bool() bool
	OptionalBool() rusty.Optional[bool]
}

func TestSimpleTypeParam(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema: eg.TestSchema(sl, eg.PropertyRuntime{}).(eg.PropertyObject),

		includes: make(map[string]bool),

		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateParam()
	assert.Equal(t, []string{
		"type SimpleTypeParam struct {\n",
		"\tString string `json:\"string\"`\n",
		"\tOptionalString rusty.Optional[string] `json:\"optional-string\",omitempty`\n",
		"\tFloat64 float64 `json:\"float64\"`\n",
		"\tOptionalFloat32 rusty.Optional[float64] `json:\"optional-float32\",omitempty`\n",
		"\tInt64 int64 `json:\"int64\"`\n",
		"\tOptionalInt32 rusty.Optional[int64] `json:\"optional-int32\",omitempty`\n",
		"\tBool bool `json:\"bool\"`\n",
		"\tOptionalBool rusty.Optional[bool] `json:\"optional-bool\",omitempty`\n",
		"}\n",
		"\n",
	}, wt.lines)
}

type SimpleType struct {
	String          string
	OptionalString  rusty.Optional[string]
	Float64         float64
	OptionalFloat64 rusty.Optional[float64]
	Int64           int64
	OptionalInt64   rusty.Optional[int64]
	Uint64          uint64
	OptionalUint64  rusty.Optional[uint64]
	Bool            bool
	OptionalBool    rusty.Optional[bool]
}

func TestSimpleTypeImpl(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema:   eg.TestSchema(sl).(eg.PropertyObject),
		includes: make(map[string]bool),
		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateImpl()
	assert.Equal(t,
		[]string{
			"type simpleTypeImpl struct {\n",
			"\tstring string\n",
			"\toptionalString rusty.Optional[string]\n",
			"\tfloat64 float64\n",
			"\toptionalFloat32 rusty.Optional[float64]\n",
			"\tint64 int64\n",
			"\toptionalInt32 rusty.Optional[int64]\n",
			"\tbool bool\n",
			"\toptionalBool rusty.Optional[bool]\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) String() string {\n",
			"\treturn my.string\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) OptionalString() rusty.Optional[string] {\n",
			"\treturn my.optionalString\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) Float64() float64 {\n",
			"\treturn my.float64\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) OptionalFloat32() rusty.Optional[float64] {\n",
			"\treturn my.optionalFloat32\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) Int64() int64 {\n",
			"\treturn my.int64\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) OptionalInt32() rusty.Optional[int64] {\n",
			"\treturn my.optionalInt32\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) Bool() bool {\n",
			"\treturn my.bool\n",
			"}\n",
			"\n",
			"func (my *simpleTypeImpl) OptionalBool() rusty.Optional[bool] {\n",
			"\treturn my.optionalBool\n",
			"}\n",
			"\n"}, wt.lines)
}

func TestSimpleTypeImplClone(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema:   eg.TestSchema(sl).(eg.PropertyObject),
		includes: make(map[string]bool),
		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateCloneFunc()
	assert.Equal(t, []string{
		"func (my *simpleTypeImpl) Clone() *SimpleTypeClass {\n",
		"\treturn &simpleTypeImpl {\n",
		"\t\tstring: my.string,\n",
		"\t\toptionalString: my.optionalString,\n",
		"\t\tfloat64: my.float64,\n",
		"\t\toptionalFloat32: my.optionalFloat32,\n",
		"\t\tint64: my.int64,\n",
		"\t\toptionalInt32: my.optionalInt32,\n",
		"\t\tbool: my.bool,\n",
		"\t\toptionalBool: my.optionalBool,\n",
		"\t}\n",
		"}\n",
		"\n"}, wt.lines)
}

func TestSimpleTypeImplLess(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema:   eg.TestSchema(sl).(eg.PropertyObject),
		includes: make(map[string]bool),
		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateLessFunc()
	assert.Equal(t,
		[]string{
			"func (my *simpleTypeImpl) Less(other *SimpleTypeClass) bool {\n",
			"\tif strings.Compare(my.string, other.String()) < 0 {\n",
			"\t\treturn true\n",
			"\t}\n",
			"\tif my.optionalString != nil && other.OptionalString() != nil {\n",
			"\t\tif strings.Compare(*my.optionalString, *other.OptionalString()) < 0 {\n",
			"\t\t\treturn true\n",
			"\t\t}\n",
			"\t} else {\n",
			"\t\t\tif my.optionalString != nil && other.OptionalString() == nil {\n",
			"\t\t\t\treturn true\n",
			"\t\t\t}\n",
			"\t}\n",
			"\tif my.float64 < other.Float64() {\n",
			"\t\treturn true\n",
			"\t}\n",
			"\tif my.optionalFloat32 != nil && other.OptionalFloat32() != nil {\n",
			"\t\tif *my.optionalFloat32 < *other.OptionalFloat32() {\n",
			"\t\t\treturn true\n",
			"\t\t}\n",
			"\t} else {\n",
			"\t\t\tif my.optionalFloat32 != nil && other.OptionalFloat32() == nil {\n",
			"\t\t\t\treturn true\n",
			"\t\t\t}\n",
			"\t}\n",
			"\tif my.int64 < other.Int64() {\n",
			"\t\treturn true\n",
			"\t}\n",
			"\tif my.optionalInt32 != nil && other.OptionalInt32() != nil {\n",
			"\t\tif *my.optionalInt32 < *other.OptionalInt32() {\n",
			"\t\t\treturn true\n",
			"\t\t}\n",
			"\t} else {\n",
			"\t\t\tif my.optionalInt32 != nil && other.OptionalInt32() == nil {\n",
			"\t\t\t\treturn true\n",
			"\t\t\t}\n",
			"\t}\n",
			"\tif my.bool != other.Bool() && my.bool == false {\n",
			"\t\treturn true\n",
			"\t}\n",
			"\tif my.optionalBool != nil && other.OptionalBool() != nil {\n",
			"\t\tif *my.optionalBool != *other.OptionalBool() && *my.optionalBool == false {\n",
			"\t\t\treturn true\n",
			"\t\t}\n",
			"\t} else {\n",
			"\t\t\tif my.optionalBool != nil && other.OptionalBool() == nil {\n",
			"\t\t\t\treturn true\n",
			"\t\t\t}\n",
			"\t}\n",
			"\treturn false\n",
			"}\n",
			"\n"}, wt.lines)
}

func TestSimpleTypeImplHash(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema:   eg.TestSchema(sl).(eg.PropertyObject),
		includes: make(map[string]bool),
		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateHashFunc()
	assert.Equal(t, []string{
		"func (my *SimpleTypeImpl) Hash(w io.Writer)  {\n",
		"\tio.Write([]byte(my.string))\n",
		"\tif my.optionalString != nil {\n",
		"\t\tio.Write([]byte(*my.optionalString))\n",
		"\t}\n",
		"\tio.Write([]byte(strconv.FormatFloat(float64(my.float64), 'e', 64)))\n",
		"\tif my.optionalFloat32 != nil {\n",
		"\t\tio.Write([]byte(strconv.FormatFloat(float64(*my.optionalFloat32), 'e', 64)))\n",
		"\t}\n",
		"\tio.Write([]byte(strconv.FormatInt(my.int64, 10)))\n",
		"\tif my.optionalInt32 != nil {\n",
		"\t\tio.Write([]byte(strconv.FormatInt(*my.optionalInt32, 10)))\n",
		"\t}\n",
		"\tif my.bool {\n",
		"\t\tio.Write([]byte(\"true\"))\n",
		"\t} else {\n",
		"\t\tio.Write([]byte(\"false\"))\n",
		"\t}\n",
		"\tif my.optionalBool != nil {\n",
		"\t\tif *my.optionalBool {\n",
		"\t\t\tio.Write([]byte(\"true\"))\n",
		"\t\t} else {\n",
		"\t\t\tio.Write([]byte(\"false\"))\n",
		"\t\t}\n",
		"\t}\n",
		"}\n",
		"\n",
	}, wt.lines)
}

func TestSimpleTypeImplAsMap(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema: eg.TestSchema(sl).(eg.PropertyObject),

		includes: make(map[string]bool),

		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateAsMapFunc()
	assert.Equal(t,
		[]string{
			"func (my *simpleTypeImpl) AsMap() map[string]interface{} {\n",
			"\treturn map[string]interface{}{\n",
			"\t\"string\": my.string,\n",
			"\t\"optional-string\": my.optionalString,\n",
			"\t\"float64\": my.float64,\n",
			"\t\"optional-float32\": my.optionalFloat32,\n",
			"\t\"int64\": my.int64,\n",
			"\t\"optional-int32\": my.optionalInt32,\n",
			"\t\"bool\": my.bool,\n",
			"\t\"optional-bool\": my.optionalBool,\n",
			"}\n",
			"\n"}, wt.lines)
}

func TestSimpleTypeBuilder(t *testing.T) {
	sl := eg.NewTestContext()
	wt := bufLineWriter{}
	g := &goGenerator{
		schema: eg.TestSchema(sl).(eg.PropertyObject),

		includes: make(map[string]bool),

		bodyWriter: &eg.ForIfWhileLangWriter{
			// Writer:    &wt,
			OfsIndent: "\t",
		},
	}
	g.generateBuilder()
	assert.Equal(t,
		[]string{}, wt.lines)
}

func TestFileName(t *testing.T) {
	assert.Equal(t, "simple_type_impl.go", FileName("SimpleTypeImpl", ".go"))
	assert.Equal(t, "simple_type_impl.go", FileName("Simple_type__impl", ".go"))
	assert.Equal(t, "simple_type_impl.go", FileName("Simple-TYPe__impl", ".go"))
}

func TestGoGenerator(t *testing.T) {
	uuid := uuid.New().String()
	dir, err := os.MkdirTemp("./", "tst."+uuid)
	if err != nil {
		t.Fatal(err)
	}
	// defer os.RemoveAll(dir)
	sl := eg.NewTestContext()
	schema := eg.TestSchema(sl).(eg.PropertyObject)
	wr, err := os.OpenFile(filepath.Join(dir, FileName(schema.Title(), ".go")), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer wr.Close()
	GoGenerator(&eg.Config{
		Indent:      "  ",
		PackageName: "test",
	}, eg.TestSchema(sl).(eg.PropertyObject), wr)
	t.Log(dir)

	cmd := exec.Command("go", "env", "GOPATH")
	// err = cmd.Run()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	gopath, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	i := interp.New(interp.Options{
		GoPath: strings.TrimSpace(string(gopath)),
	})
	i.Use(stdlib.Symbols)

	src, err := ioutil.ReadFile(filepath.Join(dir, FileName(schema.Title(), ".go")))
	if err != nil {
		t.Fatal(err)
	}
	_, err = i.Eval(string(src))
	if err != nil {
		t.Fatal(err)
	}

	// v, err := i.Eval("foo.Bar")
	// if err != nil {
	// 	panic(err)
	// }

	// bar := v.Interface().(func(string) string)

	// r := bar("Kung")
	// println(r)

	// GOARCH=wasm GOOS=js go build -o wasm.wasm simple_type.go

	// cmd := exec.Command("go", "build",
	// 	"-o", path.Join(dir, FileName(schema.Title(), ".wasm")), path.Join(dir, FileName(schema.Title(), ".go")))
	// cmd.Env = os.Environ()
	// cmd.Env = append(cmd.Env, "GOARCH=wasm")
	// cmd.Env = append(cmd.Env, "GOOS=js")
	// err = cmd.Run()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// wasmBytes, _ := ioutil.ReadFile(path.Join(dir, FileName(schema.Title(), ".wasm")))

	// engine := wasmer.NewEngine()
	// store := wasmer.NewStore(engine)

	// // Compiles the module
	// module, _ := wasmer.NewModule(store, wasmBytes)

	// // Instantiates the module
	// importObject := wasmer.NewImportObject()
	// instance, _ := wasmer.NewInstance(module, importObject)

	// // Gets the `sum` exported function from the WebAssembly instance.
	// sum, _ := instance.Exports.GetFunction("Sum")

	// // Calls that exported function with Go standard values. The WebAssembly
	// // types are inferred and values are casted automatically.
	// result, _ := sum(5, 37)

	// fmt.Println(result) // 42!

}
