package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
	"go/token"
)

func ParseImports(node ast.Node) []model.Import {
	rs := []model.Import(nil)
	switch x := any(node).(type) {
	case *ast.File:
		for _, imp := range x.Imports {
			rs = append(rs, ParseImport(imp))
		}
	case *ast.GenDecl:
		if x.Tok != token.IMPORT {
			return nil
		}
		for _, spec := range x.Specs {
			rs = append(rs, ParseImports(spec)...)
		}
	case *ast.ImportSpec:
		return []model.Import{ParseImport(x)}
	}
	return nil
}

func ParseImport(node *ast.ImportSpec) model.Import {
	if node == nil {
		return model.Import{}
	}
	if node.Name != nil {
		return model.Import{
			Alias: node.Name.Name,
			Path:  node.Path.Value,
		}
	}
	return model.Import{
		Path: node.Path.Value,
	}
}
