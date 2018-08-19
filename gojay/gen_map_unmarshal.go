package main

import (
	"fmt"
	"go/ast"
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

func (g *Gen) mapGetType(mapV *ast.MapType) (string, error) {
	// check if has hide tag
	switch t := mapV.Value.(type) {
	case *ast.Ident:
		return t.String(), nil
	case *ast.StarExpr:
		switch ptrExp := t.X.(type) {
		case *ast.Ident:
			return "*" + ptrExp.String(), nil
		case *ast.SelectorExpr:
			pkgName, err := g.getNameFromAstExpr(ptrExp.X)
			if err != nil {
				return "", err
			}
			return "*" + formatType(ptrExp.Sel.String(), pkgName), nil
		default:
			return "", fmt.Errorf("Unknown type %T", ptrExp)
		}
	case *ast.SelectorExpr:
		pkgName, err := g.getNameFromAstExpr(t.X)
		if err != nil {
			return "", err
		}
		return formatType(t.Sel.String(), pkgName), nil
	}
	return "", ErrUnknownType
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
	typeName, err := g.mapGetType(s)
	if err != nil {
		return err
	}
	err = g.mapGenUnmarshalIdent(typeName)
	if err != nil {
		return err
	}
	_, err = g.b.Write(structUnmarshalClose)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gen) mapGenUnmarshalIdent(typeName string) error {
	t := typeName
	if t[0] == '*' {
		t = t[1:]
	}
	if v, ok := genTypes[typeName]; ok {
		var err = v.mapTpl.unmarshalTpl.Execute(
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
		var err = genTypes[structType].mapTpl.unmarshalTpl.Execute(
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
