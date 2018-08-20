package main

import (
	"go/ast"
	"log"
)

func (g *Gen) mapGenIsNil(n string) error {
	err := mapMarshalTpl["isNil"].tpl.Execute(g.b, struct {
		StructName string
	}{
		StructName: n,
	})
	return err
}

func (g *Gen) mapGenMarshalObj(n string, s *ast.MapType) error {
	err := mapMarshalTpl["def"].tpl.Execute(g.b, struct {
		StructName string
	}{
		StructName: n,
	})
	if err != nil {
		return err
	}
	typeName, err := g.mapGetType(s)
	if err != nil {
		return err
	}
	err = g.mapGenMarshalIdent(typeName)
	if err != nil {
		return err
	}
	_, err = g.b.Write([]byte("}\n"))
	if err != nil {
		return err
	}
	return nil
}

func (g *Gen) mapGenMarshalIdent(typeName string) error {
	t := typeName
	if t[0] == '*' {
		t = t[1:]
	}
	if v, ok := genTypes[typeName]; ok {
		var err = v.mapTpl.marshalTpl.Execute(
			g.b,
			struct {
				TypeName string
			}{t},
		)
		if err != nil {
			return err
		}
		return nil
	}
	var structType = "struct"
	if typeName[0] == '*' {
		structType = "*" + structType
	}
	if _, ok := g.genTypes[t]; ok {
		var err = genTypes[structType].mapTpl.marshalTpl.Execute(
			g.b,
			struct {
				TypeName string
			}{t},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gen) mapMarshalNonPrim(sp *ast.TypeSpec, ptr bool) {
	switch sp.Type.(type) {
	case *ast.StructType:
		g.mapMarshalStruct(sp, ptr)
	case *ast.ArrayType:
		g.mapMarshalArr(sp, ptr)
	}
}

func (g *Gen) mapMarshalStruct(st *ast.TypeSpec, ptr bool) {
	ptrStr := ""
	if ptr {
		ptrStr = "*"
	}
	var err = mapMarshalTpl["struct"].tpl.Execute(g.b, struct {
		Ptr string
	}{ptrStr})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Gen) mapMarshalArr(st *ast.TypeSpec, ptr bool) {
	ptrStr := ""
	if ptr {
		ptrStr = "*"
	}
	var err = mapMarshalTpl["arr"].tpl.Execute(g.b, struct {
		Ptr string
	}{ptrStr})
	if err != nil {
		log.Fatal(err)
	}
}
