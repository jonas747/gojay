package main

var structMarshalTpl = templateList{
	"def": &genTpl{
		strTpl: "\n// MarshalJSONObject implements gojay's MarshalerJSONObject" +
			"\nfunc (v *{{.StructName}}) MarshalJSONObject(enc *gojay.Encoder) {\n",
	},
	"isNil": &genTpl{
		strTpl: `
// IsNil returns wether the structure is nil value or not
func (v *{{.StructName}}) IsNil() bool { return v == nil }
`,
	},
}

func init() {
	parseTemplates(structMarshalTpl, "structMarshal")
}
