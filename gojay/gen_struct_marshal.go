package main

import (
	"go/ast"
)

func (g *Gen) structGenIsNil(n string) error {
	err := structMarshalTpl["isNil"].tpl.Execute(g.b, struct {
		StructName string
	}{
		StructName: n,
	})
	return err
}

func (g *Gen) structGenMarshalObj(n string, s *ast.StructType) (int, error) {
	err := structMarshalTpl["def"].tpl.Execute(g.b, struct {
		StructName string
	}{
		StructName: n,
	})
	if err != nil {
		return 0, err
	}
	keys := 0
	if len(s.Fields.List) > 0 {
		// TODO:  check tags
		for _, field := range s.Fields.List {
			// check if has hide tag
			var omitEmpty string
			if field.Tag != nil {
				if hasTagMarshalHide(field.Tag) {
					continue
				}
				if hasTagOmitEmpty(field.Tag) {
					omitEmpty = omitEmptyFuncName
				}
			}
			typeName, ptr, err := g.structGetTypeFromField(field)
			if err != nil {
				return 0, err
			}
			keys, err = g.structGenMarshalIdent(field, typeName, keys, omitEmpty, ptr)
			if err != nil {
				return 0, err
			}
		}
	}
	_, err = g.b.Write([]byte("}\n"))
	if err != nil {
		return 0, err
	}
	return keys, nil
}

func (g *Gen) structGenMarshalIdent(field *ast.Field, typeName string, keys int, omitEmpty string, ptr bool) (int, error) {
	var keyV = getStructFieldJSONKey(field)
	if v, ok := genTypes[typeName]; ok {
		var s, err = v.structTpl.marshalFunc(g, field, keyV)
		if err != nil {
			return 0, err
		}
		err = v.structTpl.marshalTpl.Execute(
			g.b,
			s,
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
			var err = g.structMarshalNonPrim(field, keyV, sp, typeName)
			if err != nil {
				return 0, err
			}
			keys++
		}
	}
	return keys, nil
}

func (g *Gen) structMarshalNonPrim(field *ast.Field, keyV string, sp *ast.TypeSpec, typeName string) error {
	switch sp.Type.(type) {
	case *ast.StructType:
		var t = "struct"
		if typeName[0] == '*' {
			t = "*" + t
		}
		var s, err = genTypes[t].structTpl.marshalFunc(g, field, keyV, sp.Name.String())
		if err != nil {
			return err
		}
		err = genTypes[t].structTpl.marshalTpl.Execute(
			g.b,
			s,
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
		var s, err = genTypes[t].structTpl.marshalFunc(g, field, keyV)
		if err != nil {
			return err
		}
		err = genTypes[t].structTpl.marshalTpl.Execute(
			g.b,
			s,
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
		var s, err = genTypes[t].structTpl.marshalFunc(g, field, keyV)
		if err != nil {
			return err
		}
		err = genTypes[t].structTpl.marshalTpl.Execute(
			g.b,
			s,
		)
		if err != nil {
			return err
		}
		return nil
	}
}
