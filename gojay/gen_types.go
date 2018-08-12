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
	mapTpl    MapTpl
	structTpl *StructTpl
	arrTpl    ArrTpl
}

var genTypes = map[string]T{
	"string": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*string": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"int": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*int": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"int64": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*int64": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"int32": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*int32": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"int16": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*int16": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"int8": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*int8": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"uint64": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*uint64": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"uint32": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*uint32": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"uint16": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*uint16": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"uint8": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*uint8": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"float64": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*float64": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"float32": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*float32": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"bool": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*bool": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"arr": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*arr": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"struct": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"*struct": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"time.Time": {
		mapTpl: MapTpl{},
		structTpl: &StructTpl{
			unmarshalStr: `		return dec.Time(&v.{{.Field}}, {{.Format}})
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
		arrTpl: ArrTpl{},
	},
	"*time.Time": {
		mapTpl: MapTpl{},
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
		arrTpl: ArrTpl{},
	},
	"any": {
		mapTpl: MapTpl{},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Any(&v.{{.Field}})\n",
			marshalStr:   "\tenc.AnyKey(\"{{.Key}}\", v.{{.Field}})\n",
		},
		arrTpl: ArrTpl{},
	},
	"*any": {
		mapTpl: MapTpl{},
		structTpl: &StructTpl{
			unmarshalStr: "\t\treturn dec.Any(v.{{.Field}})\n",
			marshalStr:   "\tenc.AnyKey(\"{{.Key}}\", *v.{{.Field}})\n",
		},
		arrTpl: ArrTpl{},
	},
	"sql.NullString": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"*sql.NullString": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"sql.NullInt64": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"*sql.NullInt64": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"sql.NullFloat64": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"*sql.NullFloat64": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"sql.NullBool": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
	"*sql.NullBool": {
		mapTpl:    MapTpl{},
		structTpl: &StructTpl{},
		arrTpl:    ArrTpl{},
	},
}

func init() {
	for typeName, genType := range genTypes {
		// map tpl
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
		log.Print(typeName, " ", genType.structTpl.unmarshalTpl)
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
