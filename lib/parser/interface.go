package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
)

func ParseInterface(node ast.Node) *model.Interface {
	switch x := node.(type) {
	case *ast.TypeSpec:
		if _, ok := x.Type.(*ast.InterfaceType); !ok {
			return nil
		}
		return &model.Interface{
			Name:    x.Name.Name,
			Methods: ParseMethods(x.Type.(*ast.InterfaceType).Methods),
		}
	}
	return nil
}

func ParseMethods(fl *ast.FieldList) []model.Method {
	if fl == nil {
		return nil
	}
	methods := []model.Method(nil)
	for _, field := range fl.List {
		funcType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		params := ParseFields(funcType.Params)
		results := ParseFields(funcType.Results)

		if len(field.Names) == 0 {
			methods = append(methods, model.Method{
				Params: params,
				Return: results,
			})
			continue
		}
		for _, name := range field.Names {
			methods = append(methods, model.Method{
				Name:   name.Name,
				Params: params,
				Return: results,
			})
		}
	}
	return methods
}
