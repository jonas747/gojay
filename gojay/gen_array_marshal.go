package main

import (
	"go/ast"
	"log"
)

func init() {}

func (g *Gen) arrGenIsNil(n string) error {
	err := arrMarshalTpl["isNil"].tpl.Execute(g.b, struct {
		TypeName string
	}{
		TypeName: n,
	})
	return err
}

func (g *Gen) arrGenMarshal(n string, s *ast.ArrayType) error {
	err := arrMarshalTpl["def"].tpl.Execute(g.b, struct {
		TypeName string
	}{
		TypeName: n,
	})
	if err != nil {
		return err
	}
	// determine type of element in array
	t, err := g.getArrayType(s)
	if err != nil {
		return err
	}
	err = g.arrGenMarshalIdent(n, t)
	if err != nil {
		return err
	}
	_, err = g.b.Write([]byte("}\n"))
	if err != nil {
		return err
	}
	return err
}

func (g *Gen) arrGenMarshalIdent(n, eltName string) error {
	t := eltName
	if t[0] == '*' {
		t = t[1:]
	}
	if v, ok := genTypes[t]; ok {
		var err = v.arrTpl.marshalTpl.Execute(
			g.b,
			struct {
				TypeName string
			}{n},
		)
		if err != nil {
			return err
		}
		return nil
	}
	var structType = "struct"
	if t[0] == '*' {
		structType = "*" + structType
	}
	if _, ok := g.genTypes[t]; ok {
		var err = genTypes[structType].arrTpl.marshalTpl.Execute(
			g.b,
			struct {
				TypeName string
			}{n},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gen) arrMarshalNonPrim(sp *ast.TypeSpec, ptr bool) error {
	switch sp.Type.(type) {
	case *ast.StructType:
		g.arrMarshalStruct(sp, ptr)
	case *ast.ArrayType:
		g.arrMarshalArr(sp, ptr)
	}
	return nil
}

func (g *Gen) arrMarshalStruct(st *ast.TypeSpec, ptr bool) {
	if ptr {
		err := arrMarshalTpl["structPtr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := arrMarshalTpl["struct"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) arrMarshalArr(st *ast.TypeSpec, ptr bool) {
	if ptr {
		err := arrMarshalTpl["arrPtr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := arrMarshalTpl["arr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	}
}
