package parser

import (
	"go/ast"
	"strings"
)

func ParseName(node ast.Node) string {
	switch x := node.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + ParseName(x.X)
	case *ast.SelectorExpr:
		return ParseName(x.X) + "." + ParseName(x.Sel)
	case *ast.ArrayType:
		return "[]" + ParseName(x.Elt)
	case *ast.MapType:
		return "map[" + ParseName(x.Key) + "]" + ParseName(x.Value)
	case *ast.InterfaceType:
		return "interface{" + ParseFieldList(x.Methods, "; ") + "}"
	case *ast.ChanType:
		switch x.Dir {
		case ast.SEND:
			return "chan<- " + ParseName(x.Value)
		case ast.RECV:
			return "<-chan " + ParseName(x.Value)
		default:
			return "chan " + ParseName(x.Value)
		}
	case *ast.FuncType:
		return "func(" + ParseFieldList(x.Params, ", ") + ") (" + ParseFieldList(x.Results, ", ") + ")"
	case *ast.Ellipsis:
		return "..." + ParseName(x.Elt)
	}
	return ""
}

func ParseFieldList(fl *ast.FieldList, sep string) string {
	if fl == nil {
		return ""
	}
	fields := []string(nil)
	for _, field := range fl.List {
		typ := ParseName(field.Type)
		if len(field.Names) == 0 {
			fields = append(fields, typ)
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, name.Name+" "+typ)
		}
	}
	return strings.Join(fields, sep)
}
