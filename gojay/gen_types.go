package main

import (
	"go/ast"
	"log"
	"text/template"
)

type StructTpl struct {
	marshalStr    string
	marshalTpl    *template.Template
	unmarshalStr  string
	unmarshalTpl  *template.Template
	marshalFunc   func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error)
	unmarshalFunc func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error)
}

type MapTpl struct {
	marshalStr   string
	marshalTpl   *template.Template
	unmarshalStr string
	unmarshalTpl *template.Template
}

type ArrTpl struct {
	marshalStr   string
	marshalTpl   *template.Template
	unmarshalStr string
	unmarshalTpl *template.Template
}

type T struct {
	mapTpl    *MapTpl
	structTpl *StructTpl
	arrTpl    *ArrTpl
}

var genTypes = map[string]T{
	"string": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var str string
	if err := dec.String(&str); err != nil {
		return err
	}
	v[k] = str
`,
			marshalStr: "\tfor k, s := range v {\n" +
				"\t\tenc.StringKey(k, s)\n" +
				"\t}\n",
		},

		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.String(&v.{{.Field}})\n",
			marshalStr:   "\tenc.StringKey{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{
			unmarshalStr: "\tvar str string" +
				"\n\tif err := dec.String(&str); err != nil {\n" +
				"\t\treturn err\n\t}\n" +
				"\t*v = append(*v, str)\n",
			marshalStr: "\tfor _, s := range *v {\n" +
				"\t\tenc.String(s)\n" +
				"\t}\n",
		},
	},
	"*string": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var str string
	if err := dec.String(&str); err != nil {
		return err
	}
	v[k] = &str
`,
			marshalStr: "\tfor k, s := range v {\n" +
				"\t\tenc.StringKey(k, *s)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.String(v.{{.Field}})\n",
			marshalStr:   "\tenc.StringKey{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{
			unmarshalStr: "\tvar str string" +
				"\n\tif err := dec.String(&str); err != nil {\n" +
				"\t\treturn err\n\t}\n" +
				"\t*v = append(*v, &str)\n",
			marshalStr: "\tfor _, s := range *v {\n" +
				"\t\tenc.String(s)\n" +
				"\t}\n",
		},
	},

	"int": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int
	if err := dec.Int(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.IntKey(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int(&v.{{.Field}})\n",
			marshalStr:   "\tenc.IntKey{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*int": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int
	if err := dec.Int(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.IntKey(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int(v.{{.Field}})\n",
			marshalStr:   "\tenc.IntKey{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"int64": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int64
	if err := dec.Int64(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int64Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int64(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Int64Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*int64": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int64
	if err := dec.Int64(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int64Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int64(v.{{.Field}})\n",
			marshalStr:   "\tenc.Int64Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"int32": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int32
	if err := dec.Int32(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int32Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int32(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Int32Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*int32": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int32
	if err := dec.Int32(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int32Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int32(v.{{.Field}})\n",
			marshalStr:   "\tenc.Int32Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"int16": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int16
	if err := dec.Int16(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int16Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int16(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Int16Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*int16": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int16
	if err := dec.Int16(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int16Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int16(v.{{.Field}})\n",
			marshalStr:   "\tenc.Int16Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"int8": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int8
	if err := dec.Int8(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int8Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int8(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Int8Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*int8": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i int8
	if err := dec.Int8(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Int8Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Int8(v.{{.Field}})\n",
			marshalStr:   "\tenc.Int8Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"uint64": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint64
	if err := dec.Uint64(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint64Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint64(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint64Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*uint64": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint64
	if err := dec.Uint64(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint64Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint64(v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint64Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"uint32": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint32
	if err := dec.Uint32(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint32Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint32(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint32Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*uint32": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint32
	if err := dec.Uint32(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint32Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint32(v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint32Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"uint16": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint16
	if err := dec.Uint16(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint16Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint16(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint16Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*uint16": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint16
	if err := dec.Uint16(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint16Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint16(v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint16Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"uint8": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint8
	if err := dec.Uint8(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint8Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint8(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint8Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*uint8": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i uint8
	if err := dec.Uint8(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Uint8Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Uint8(v.{{.Field}})\n",
			marshalStr:   "\tenc.Uint8Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"float64": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i float64
	if err := dec.Float64(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Float64Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Float64(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Float64Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*float64": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i float64
	if err := dec.Float64(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Float64Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Float(v.{{.Field}})\n",
			marshalStr:   "\tenc.FloatKey{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"float32": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i float32
	if err := dec.Float32(&i); err != nil {
		return err
	}
	v[k] = i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Float32Key(k, i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Float32(&v.{{.Field}})\n",
			marshalStr:   "\tenc.Float32Key{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*float32": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var i float32
	if err := dec.Float32(&i); err != nil {
		return err
	}
	v[k] = &i
`,
			marshalStr: "\tfor k, i := range v {\n" +
				"\t\tenc.Float32Key(k, *i)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Float32(v.{{.Field}})\n",
			marshalStr:   "\tenc.Float32Key{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"bool": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var b bool
	if err := dec.Bool(&b); err != nil {
		return err
	}
	v[k] = b
`,
			marshalStr: "\tfor k, b := range v {\n" +
				"\t\tenc.BoolKey(k, b)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Bool(&v.{{.Field}})\n",
			marshalStr:   "\tenc.BoolKey{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*bool": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var b bool
	if err := dec.Bool(&b); err != nil {
		return err
	}
	v[k] = &b
`,
			marshalStr: "\tfor k, b := range v {\n" +
				"\t\tenc.BoolKey(k, *b)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Bool(&v.{{.Field}})\n",
			marshalStr:   "\tenc.BoolKey{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"arr": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var a = {{.TypeName}}{}
	if err := dec.Array(&a); err != nil {
		return err
	}
	v[k] = a
`,
			marshalStr: "\tfor _, a := range v {\n" +
				"\t\tenc.Array(&a)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: `		if v.{{.Field}} == nil {
			arr := make({{.TypeName}}, 0)
			v.{{.Field}} = arr
		}
		return dec.Array(&v.{{.Field}})
`,
			marshalStr: "\tenc.ArrayKey{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*arr": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var a = {{.TypeName}}{}
	if err := dec.Array(&a); err != nil {
		return err
	}
	v[k] = &a
`,
			marshalStr: "\tfor _, a := range v {\n" +
				"\t\tenc.Array(&a)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: `		if v.{{.Field}} == nil {
			arr := make({{.TypeName}}, 0)
			v.{{.Field}} = &arr
		}
		return dec.Array(v.{{.Field}})
`,
			marshalStr: "\tenc.ArrayKey{{.OmitEmpty}}(\"{{.Key}}\", *v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field string
				}{field.Names[0].Name}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"struct": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var s = {{.TypeName}}{}
	if err := dec.Object(s); err != nil {
		return err
	}
	v[k] = s
`,
			marshalStr: "\tfor _, s := range v {\n" +
				"\t\tenc.Object(s)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: `		if v.{{.Field}} == nil {
			v.{{.Field}} = {{.StructName}}{}
		}
		return dec.Object(&v.{{.Field}})
`,
			marshalStr: "\tenc.ObjectKey{{.OmitEmpty}}(\"{{.Key}}\", &v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field      string
					StructName string
				}{field.Names[0].Name, args[0].(string)}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*struct": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var s = &{{.TypeName}}{}
	if err := dec.Object(s); err != nil {
		return err
	}
	v[k] = s
`,
			marshalStr: "\tfor _, s := range v {\n" +
				"\t\tenc.Object(s)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: `		if v.{{.Field}} == nil {
			v.{{.Field}} = &{{.StructName}}{}
		}
		return dec.Object(v.{{.Field}})
`,
			marshalStr: "\tenc.ObjectKey{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
				}{field.Names[0].Name, key, getOmitEmpty(field)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field      string
					StructName string
				}{field.Names[0].Name, args[0].(string)}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"time.Time": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var t = time.Time{}
	if err := dec.Time(&t, time.RFC3339); err != nil {
		return err
	}
	v[k] = t
`,
			marshalStr: "\tfor _, t := range v {\n" +
				"\t\tenc.Time(&t, time.RFC3339)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: `		return dec.Time(&v.{{.Field}}, {{.Format}})
`,
			marshalStr: "\tenc.TimeKey{{.OmitEmpty}}(\"{{.Key}}\", &v.{{.Field}}, {{.Format}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
					Format    string
				}{field.Names[0].Name, key, getOmitEmpty(field), timeFormat(field.Tag)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field  string
					Format string
				}{field.Names[0].Name, timeFormat(field.Tag)}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"*time.Time": {
		mapTpl: &MapTpl{
			unmarshalStr: `	var t = time.Time{}
	if err := dec.Time(&t, time.RFC3339); err != nil {
		return err
	}
	v[k] = &t
`,
			marshalStr: "\tfor _, t := range v {\n" +
				"\t\tenc.Time(&t, time.RFC3339)\n" +
				"\t}\n",
		},
		structTpl: &StructTpl{
			unmarshalStr: `		if v.{{.Field}} == nil {
			v.{{.Field}} = &time.Time{}
		}
		return dec.Time(v.{{.Field}}, {{.Format}})
`,
			marshalStr: "\tenc.TimeKey{{.OmitEmpty}}(\"{{.Key}}\", v.{{.Field}}, {{.Format}})\n",
			marshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field     string
					Key       string
					OmitEmpty string
					Format    string
				}{field.Names[0].Name, key, getOmitEmpty(field), timeFormat(field.Tag)}, nil
			},
			unmarshalFunc: func(g *Gen, field *ast.Field, key string, args ...interface{}) (interface{}, error) {
				return struct {
					Field  string
					Format string
				}{field.Names[0].Name, timeFormat(field.Tag)}, nil
			},
		},
		arrTpl: &ArrTpl{},
	},
	"any": {
		mapTpl: &MapTpl{},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Any(&v.{{.Field}})\n",
			marshalStr:   "\tenc.AnyKey(\"{{.Key}}\", v.{{.Field}})\n",
		},
		arrTpl: &ArrTpl{},
	},
	"*any": {
		mapTpl: &MapTpl{},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Any(v.{{.Field}})\n",
			marshalStr:   "\tenc.AnyKey(\"{{.Key}}\", *v.{{.Field}})\n",
		},
		arrTpl: &ArrTpl{},
	},
	"sql.NullString": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"*sql.NullString": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"sql.NullInt64": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"*sql.NullInt64": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"sql.NullFloat64": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"*sql.NullFloat64": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"sql.NullBool": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
	"*sql.NullBool": {
		mapTpl:    &MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    &ArrTpl{},
	},
}

func init() {
	for typeName, genType := range genTypes {
		// map tpl
		log.Print("prep tpl: ", typeName)
		var tpl, err = template.New(typeName + ".unmarshal.map").Parse(genType.mapTpl.unmarshalStr)
		if err != nil {
			panic(err)
		}
		genType.mapTpl.unmarshalTpl = tpl
		tpl, err = template.New(typeName + ".marshal.map").Parse(genType.mapTpl.marshalStr)
		if err != nil {
			panic(err)
		}
		genType.mapTpl.marshalTpl = tpl
		// struct tpl
		tpl, err = template.New(typeName + ".unmarshal.struct").Parse(genType.structTpl.unmarshalStr)
		if err != nil {
			panic(err)
		}
		genType.structTpl.unmarshalTpl = tpl
		tpl, err = template.New(typeName + ".marshal.struct").Parse(genType.structTpl.marshalStr)
		if err != nil {
			panic(err)
		}
		genType.structTpl.marshalTpl = tpl
		// arr tpl
		tpl, err = template.New(typeName + ".unmarshal.arr").Parse(genType.arrTpl.unmarshalStr)
		if err != nil {
			panic(err)
		}
		genType.arrTpl.unmarshalTpl = tpl
		tpl, err = template.New(typeName + ".marshal.arr").Parse(genType.arrTpl.marshalStr)
		if err != nil {
			panic(err)
		}
		genType.arrTpl.marshalTpl = tpl
	}
}
