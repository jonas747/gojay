package main

var structUnmarshalTpl = templateList{
	"def": &genTpl{
		strTpl: "\n// UnmarshalJSONObject implements gojay's UnmarshalerJSONObject" +
			"\nfunc (v *{{.StructName}}) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {\n",
	},
	"nKeys": &genTpl{
		strTpl: `
// NKeys returns the number of keys to unmarshal
func (v *{{.StructName}}) NKeys() int { return {{.NKeys}} }
`,
	},
	"case": &genTpl{
		strTpl: "\tcase \"{{.Key}}\":\n",
	},
}

func init() {
	parseTemplates(structUnmarshalTpl, "structUnmarshal")
}
