package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
)

func ParseAlias(node ast.Node) *model.Alias {
	switch x := node.(type) {
	case *ast.TypeSpec:
		switch x.Type.(type) {
		case *ast.InterfaceType:
			return nil
		case *ast.StructType:
			return nil
		}
		return &model.Alias{
			Name: x.Name.Name,
			Type: ParseName(x.Type),
		}
	}
	return nil
}
