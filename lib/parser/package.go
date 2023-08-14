package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
)

func ParsePackage(node *ast.File) model.Package {
	return model.Package{
		Name:    node.Name.Name,
		Imports: ParseImports(node),
	}
}
