package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseFile(path string) (*model.Package, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	packageSpec := model.Package{}
	structs := []model.Struct(nil)
	interfaces := []model.Interface(nil)
	imports := []model.Import(nil)
	functions := []model.Function(nil)
	aliases := []model.Alias(nil)
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch x := n.(type) {
		case *ast.File:
			packageSpec = ParsePackage(x)
		case *ast.ImportSpec:
			if imp := ParseImports(x); imp != nil {
				imports = append(imports, imp...)
			}
		case *ast.FuncDecl:
			if fun := ParseFunction(x); fun != nil {
				functions = append(functions, *fun)
			}
		case *ast.GenDecl:
			doc := x.Doc.Text()
			for _, spec := range x.Specs {
				switch y := spec.(type) {
				case *ast.TypeSpec:
					if str := ParseStruct(y); str != nil {
						str.Doc = doc
						structs = append(structs, *str)
					}
					if inf := ParseInterface(y); inf != nil {
						inf.Doc = doc
						interfaces = append(interfaces, *inf)
					}
					if al := ParseAlias(y); al != nil {
						al.Doc = doc
						aliases = append(aliases, *al)
					}
				}
			}
		case *ast.ValueSpec:
			if gv := ParseGlobalVariable(x); gv != nil {
				packageSpec.GlobalVariables = append(packageSpec.GlobalVariables, *gv)
			}
		}
		return true
	})

	packageSpec.Structs = structs
	packageSpec.Interfaces = interfaces
	packageSpec.Imports = imports
	packageSpec.Functions = functions
	packageSpec.Aliases = aliases
	return &packageSpec, nil
}
