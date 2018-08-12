package main

var arrMarshalTpl = templateList{
	"def": &genTpl{
		strTpl: "\n// MarshalJSONArray implements gojay's MarshalerJSONArray" +
			"\nfunc (v *{{.TypeName}}) MarshalJSONArray(enc *gojay.Encoder) {\n",
	},
	"isNil": &genTpl{
		strTpl: "\n// IsNil implements gojay's MarshalerJSONArray" +
			"\nfunc (v *{{.TypeName}}) IsNil() bool {\n" +
			"\treturn *v == nil || len(*v) == 0\n" +
			"}\n",
	},
	"string": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.String(s)\n" +
			"\t}\n",
	},
	"bool": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Bool(s)\n" +
			"\t}\n",
	},
	"int": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Int{{.IntLen}}(s)\n" +
			"\t}\n",
	},
	"uint": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Uint{{.IntLen}}(s)\n" +
			"\t}\n",
	},
	"float": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Float{{.IntLen}}(s)\n" +
			"\t}\n",
	},
	"struct": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Object(s)\n" +
			"\t}\n",
	},
	"structPtr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Object(s)\n" +
			"\t}\n",
	},
	"arr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Array(s)\n" +
			"\t}\n",
	},
	"arrPtr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.Array(s)\n" +
			"\t}\n",
	},
	"sqlNullString": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullString(&s)\n" +
			"\t}\n",
	},
	"sqlNullStringPtr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullString(s)\n" +
			"\t}\n",
	},
	"sqlNullInt64": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullInt64(&s)\n" +
			"\t}\n",
	},
	"sqlNullInt64Ptr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullInt64(s)\n" +
			"\t}\n",
	},
	"sqlNullFloat64": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullFloat64(&s)\n" +
			"\t}\n",
	},
	"sqlNullFloat64Ptr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullFloat64(s)\n" +
			"\t}\n",
	},
	"sqlNullBool": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullBool(&s)\n" +
			"\t}\n",
	},
	"sqlNullBoolPtr": &genTpl{
		strTpl: "\tfor _, s := range *v {\n" +
			"\t\tenc.SQLNullBool(s)\n" +
			"\t}\n",
	},
}

func init() {
	parseTemplates(arrMarshalTpl, "arrMarshal")
}
