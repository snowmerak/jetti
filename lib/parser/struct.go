package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
)

func ParseStruct(node ast.Node) *model.Struct {
	switch x := node.(type) {
	case *ast.TypeSpec:
		if _, ok := x.Type.(*ast.StructType); !ok {
			return nil
		}
		return &model.Struct{
			Name:   x.Name.Name,
			Fields: ParseFields(x.Type.(*ast.StructType).Fields),
			Doc:    x.Doc.Text(),
		}
	}
	return nil
}

func ParseFields(fl *ast.FieldList) []model.Field {
	if fl == nil {
		return nil
	}
	fields := []model.Field(nil)
	for _, field := range fl.List {
		typ := ParseName(field.Type)
		if len(field.Names) == 0 {
			fields = append(fields, model.Field{
				Type: typ,
			})
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, model.Field{
				Name: name.Name,
				Type: typ,
			})
		}
	}
	return fields
}
