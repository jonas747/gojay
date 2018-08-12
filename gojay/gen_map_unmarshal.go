package main

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
)

func (g *Gen) mapGenNKeys(n string, count int) error {
	err := mapUnmarshalTpl["nKeys"].tpl.Execute(g.b, struct {
		NKeys      int
		StructName string
	}{
		NKeys:      count,
		StructName: n,
	})
	return err
}

func (g *Gen) mapGetType(mapV *ast.MapType) (string, bool, error) {
	// check if has hide tag
	switch t := mapV.Value.(type) {
	case *ast.Ident:
		return t.String(), false, nil
	case *ast.StarExpr:
		switch ptrExp := t.X.(type) {
		case *ast.Ident:
			return ptrExp.String(), true, nil
		case *ast.SelectorExpr:
			pkgName, err := g.getNameFromAstExpr(ptrExp.X)
			if err != nil {
				return "", false, err
			}
			return formatType(ptrExp.Sel.String(), pkgName), true, nil
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

func (g *Gen) mapGenUnmarshalObj(n string, s *ast.MapType) error {
	err := mapUnmarshalTpl["def"].tpl.Execute(g.b, struct {
		TypeName string
	}{
		TypeName: n,
	})
	if err != nil {
		return err
	}
	typeName, ptr, err := g.mapGetType(s)
	if err != nil {
		return err
	}
	err = g.mapGenUnmarshalIdent(typeName, ptr)
	if err != nil {
		return err
	}
	_, err = g.b.Write(structUnmarshalClose)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gen) mapGenUnmarshalIdent(typeName string, ptr bool) error {
	switch typeName {
	case "string":
		g.mapUnmarshalString(ptr)
	case "bool":
		g.mapUnmarshalBool(ptr)
	case "int":
		g.mapUnmarshalInt("", ptr)
	case "int64":
		g.mapUnmarshalInt("64", ptr)
	case "int32":
		g.mapUnmarshalInt("32", ptr)
	case "int16":
		g.mapUnmarshalInt("16", ptr)
	case "int8":
		g.mapUnmarshalInt("8", ptr)
	case "uint64":
		g.mapUnmarshalUint("64", ptr)
	case "uint32":
		g.mapUnmarshalUint("32", ptr)
	case "uint16":
		g.mapUnmarshalUint("16", ptr)
	case "uint8":
		g.mapUnmarshalUint("8", ptr)
	case "float64":
		g.mapUnmarshalFloat("64", ptr)
	case "float32":
		g.mapUnmarshalFloat("32", ptr)
	default:
		// if ident is already in our spec list
		if sp, ok := g.genTypes[typeName]; ok {
			err := g.mapUnmarshalNonPrim(sp, ptr)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unknown type %s", typeName)
		}
	}
	return nil
}

func (g *Gen) mapUnmarshalNonPrim(sp *ast.TypeSpec, ptr bool) error {
	switch sp.Type.(type) {
	case *ast.StructType:
		g.mapUnmarshalStruct(sp, ptr)
		return nil
	case *ast.ArrayType:
		g.mapUnmarshalArr(sp, ptr)
		return nil
	}
	return errors.New("Unknown type")
}

func (g *Gen) mapUnmarshalString(ptr bool) {
	if ptr {
		err := mapUnmarshalTpl["string"].tpl.Execute(g.b, struct {
			Ptr string
		}{"&"})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := mapUnmarshalTpl["string"].tpl.Execute(g.b, struct {
			Ptr string
		}{""})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) mapUnmarshalBool(ptr bool) {
	if ptr {
		err := mapUnmarshalTpl["bool"].tpl.Execute(g.b, struct {
			Ptr string
		}{"&"})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := mapUnmarshalTpl["bool"].tpl.Execute(g.b, struct {
			Ptr string
		}{""})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) mapUnmarshalInt(intLen string, ptr bool) {
	if ptr {
		err := mapUnmarshalTpl["int"].tpl.Execute(g.b, struct {
			IntLen string
			Ptr    string
		}{intLen, "&"})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := mapUnmarshalTpl["int"].tpl.Execute(g.b, struct {
			IntLen string
			Ptr    string
		}{intLen, ""})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) mapUnmarshalUint(intLen string, ptr bool) {
	if ptr {
		err := mapUnmarshalTpl["uint"].tpl.Execute(g.b, struct {
			IntLen string
			Ptr    string
		}{intLen, "&"})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := mapUnmarshalTpl["uint"].tpl.Execute(g.b, struct {
			IntLen string
			Ptr    string
		}{intLen, ""})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) mapUnmarshalFloat(intLen string, ptr bool) {
	if ptr {
		err := mapUnmarshalTpl["float"].tpl.Execute(g.b, struct {
			IntLen string
			Ptr    string
		}{intLen, "&"})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := mapUnmarshalTpl["float"].tpl.Execute(g.b, struct {
			IntLen string
			Ptr    string
		}{intLen, ""})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) mapUnmarshalStruct(st *ast.TypeSpec, ptr bool) {
	if ptr {
		err := mapUnmarshalTpl["structPtr"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := mapUnmarshalTpl["struct"].tpl.Execute(g.b, struct {
			StructName string
		}{st.Name.String()})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Gen) mapUnmarshalArr(st *ast.TypeSpec, ptr bool) {
	err := mapUnmarshalTpl["arr"].tpl.Execute(g.b, struct {
		TypeName string
	}{st.Name.String()})
	if err != nil {
		log.Fatal(err)
	}
}
