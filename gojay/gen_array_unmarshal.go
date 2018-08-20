package main

import (
	"fmt"
	"go/ast"
	"log"
)

func (g *Gen) getArrayType(s *ast.ArrayType) (string, error) {
	// determine type of element in array
	switch t := s.Elt.(type) {
	case *ast.Ident:
		return t.Name, nil
	case *ast.StarExpr:
		switch ptrExp := t.X.(type) {
		case *ast.Ident:
			return "*" + ptrExp.Name, nil
		default:
			return "", fmt.Errorf("Unknown type %T", ptrExp)
		}
	}
	return "", fmt.Errorf("Unknown type %T", s)
}

func (g *Gen) arrGenUnmarshal(n string, s *ast.ArrayType) error {
	err := arrUnmarshalTpl["def"].tpl.Execute(g.b, struct {
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
	err = g.arrGenUnmarshalIdent(n, t)
	if err != nil {
		return err
	}
	_, err = g.b.Write(structUnmarshalClose)
	if err != nil {
		return err
	}
	return err
}

func (g *Gen) arrGenUnmarshalIdent(n, eltName string) error {
	t := eltName
	if t[0] == '*' {
		t = t[1:]
	}
	if v, ok := genTypes[t]; ok {
		var err = v.arrTpl.unmarshalTpl.Execute(
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
		var err = genTypes[structType].arrTpl.unmarshalTpl.Execute(
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

func (g *Gen) arrUnmarshalNonPrim(sp *ast.TypeSpec, ptr bool) error {
	switch sp.Type.(type) {
	case *ast.StructType:
		g.arrUnmarshalStruct(sp, ptr)
	case *ast.ArrayType:
		g.arrUnmarshalArr(sp, ptr)
	}
	return nil
}

func (g *Gen) arrUnmarshalStruct(st *ast.TypeSpec, ptr bool) {
	if ptr {
		err := arrUnmarshalTpl["structPtr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := arrUnmarshalTpl["struct"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) arrUnmarshalArr(st *ast.TypeSpec, ptr bool) {
	if ptr {
		err := arrUnmarshalTpl["arrPtr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := arrUnmarshalTpl["arr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	}
}
