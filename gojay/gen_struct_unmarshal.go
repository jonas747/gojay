package main

import (
	"errors"
	"fmt"
	"go/ast"
)

var ErrUnknownType = errors.New("unknown type")

var structUnmarshalSwitchOpen = []byte("\tswitch k {\n")
var structUnmarshalClose = []byte("\treturn nil\n}\n")

func (g *Gen) structGenNKeys(n string, count int) error {
	err := structUnmarshalTpl["nKeys"].tpl.Execute(g.b, struct {
		NKeys      int
		StructName string
	}{
		NKeys:      count,
		StructName: n,
	})
	return err
}

func (g *Gen) getNameFromAstExpr(expr ast.Expr) (string, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.String(), nil
	}
	return "", ErrUnknownType
}

func (g *Gen) structGetTypeFromField(field *ast.Field) (string, bool, error) {
	// check if has hide tag
	switch t := field.Type.(type) {
	case *ast.Ident:
		return t.String(), false, nil
	case *ast.StarExpr:
		switch ptrExp := t.X.(type) {
		case *ast.Ident:
			return "*" + ptrExp.String(), true, nil
		case *ast.SelectorExpr:
			pkgName, err := g.getNameFromAstExpr(ptrExp.X)
			if err != nil {
				return "", false, err
			}
			return "*" + formatType(ptrExp.Sel.String(), pkgName), true, nil
		default:
			return "", false, fmt.Errorf("Unknown type %T", ptrExp)
		}
	case *ast.SelectorExpr:
		pkgName, err := g.getNameFromAstExpr(t.X)
		if err != nil {
			return "", false, err
		}
		return formatType(t.Sel.String(), pkgName), false, nil
	}
	return "", false, ErrUnknownType
}

func (g *Gen) structGenUnmarshalObj(n string, s *ast.StructType) (int, error) {
	err := structUnmarshalTpl["def"].tpl.Execute(g.b, struct {
		StructName string
	}{
		StructName: n,
	})
	if err != nil {
		return 0, err
	}
	keys := 0
	if len(s.Fields.List) > 0 {
		// open  switch statement
		g.b.Write(structUnmarshalSwitchOpen)
		// TODO:  check tags
		// check type of field
		// add accordingly
		for _, field := range s.Fields.List {
			// check if has hide tag
			if field.Tag != nil && hasTagUnmarshalHide(field.Tag) {
				continue
			}
			typeName, _, err := g.structGetTypeFromField(field)
			if err != nil {
				return 0, err
			}
			keys, err = g.structGenUnmarshalIdent(field, typeName, keys)
			if err != nil {
				return 0, err
			}
		}
		// close  switch statement
		g.b.Write([]byte("\t}\n"))
	}
	_, err = g.b.Write(structUnmarshalClose)
	if err != nil {
		return 0, err
	}
	return keys, nil
}

func (g *Gen) structGenUnmarshalIdent(field *ast.Field, typeName string, keys int) (int, error) {
	var keyV = getStructFieldJSONKey(field)
	if v, ok := genTypes[typeName]; ok {
		err := structUnmarshalTpl["case"].tpl.Execute(g.b, struct {
			Key string
		}{keyV})
		if err != nil {
			return 0, err
		}
		err = v.structTpl.unmarshalTpl.Execute(
			g.b,
			struct {
				Field     string
				Key       string
				OmitEmpty string
				Format    string
			}{field.Names[0].Name, keyV, getOmitEmpty(field), timeFormat(field)},
		)
		if err != nil {
			return 0, err
		}
		keys++
	} else {
		t := typeName
		if t[0] == '*' {
			t = t[1:]
		}
		if sp, ok := g.genTypes[t]; ok {
			err := structUnmarshalTpl["case"].tpl.Execute(g.b, struct {
				Key string
			}{keyV})
			if err != nil {
				return 0, err
			}
			err = g.structUnmarshalNonPrim(field, keyV, sp, typeName)
			if err != nil {
				return 0, err
			}
			keys++
		}
	}
	return keys, nil
}

func (g *Gen) structUnmarshalNonPrim(field *ast.Field, keyV string, sp *ast.TypeSpec, typeName string) error {
	switch sp.Type.(type) {
	case *ast.StructType:
		var t = "struct"
		if typeName[0] == '*' {
			t = "*" + t
		}
		if typeName[0] == '*' {
			typeName = typeName[1:]
		}
		var err = genTypes[t].structTpl.unmarshalTpl.Execute(
			g.b,
			struct {
				Field      string
				Key        string
				OmitEmpty  string
				Format     string
				StructName string
			}{field.Names[0].Name, keyV, getOmitEmpty(field), timeFormat(field), typeName},
		)
		if err != nil {
			return err
		}
		return nil
	case *ast.ArrayType:
		var t = "arr"
		if typeName[0] == '*' {
			t = "*" + t
		}
		var err = genTypes[t].structTpl.unmarshalTpl.Execute(
			g.b,
			struct {
				Field     string
				Key       string
				OmitEmpty string
				Format    string
			}{field.Names[0].Name, keyV, getOmitEmpty(field), timeFormat(field)},
		)
		if err != nil {
			return err
		}
		return nil
	default:
		var t = "any"
		if typeName[0] == '*' {
			t = "*" + t
		}
		var err = genTypes[t].structTpl.unmarshalTpl.Execute(
			g.b,
			struct {
				Field     string
				Key       string
				OmitEmpty string
				Format    string
			}{field.Names[0].Name, keyV, getOmitEmpty(field), timeFormat(field)},
		)
		if err != nil {
			return err
		}
		return nil
	}
}

// func (g *Gen) structUnmarshalAny(field *ast.Field, keyV string, st *ast.TypeSpec) {
// 	key := field.Names[0].String()
// 	err := structUnmarshalTpl["case"].tpl.Execute(g.b, struct {
// 		Key string
// 	}{keyV})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if ptr {
// 		err = structUnmarshalTpl["anyPtr"].tpl.Execute(g.b, struct {
// 			Field string
// 		}{key})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	} else {
// 		err = structUnmarshalTpl["any"].tpl.Execute(g.b, struct {
// 			Field string
// 		}{key})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }
