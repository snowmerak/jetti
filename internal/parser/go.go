package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Package struct {
	Name       string
	Imports    []Import
	Structs    []Struct
	Interfaces []Interface
}

type Import struct {
	Alias string
	Path  string
}

type Struct struct {
	Doc    string
	Name   string
	Fields []Field
}

type Method struct {
	Name   string
	Params []Field
	Return []Field
}

type Interface struct {
	Doc     string
	Name    string
	Methods []Method
}

func ParseFile(path string) (*Package, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	packageName := ""
	structs := []Struct(nil)
	interfaces := []Interface(nil)
	imports := []Import(nil)
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch x := n.(type) {
		case *ast.File:
			packageName = x.Name.Name
		case *ast.ImportSpec:
			imp := Import{
				Path: x.Path.Value,
			}
			if x.Name != nil {
				imp.Alias = x.Name.Name
			}
			imports = append(imports, imp)
		case *ast.GenDecl:
			doc := x.Doc.Text()
			for _, spec := range x.Specs {
				switch p := spec.(type) {
				case *ast.TypeSpec:
					switch e := p.Type.(type) {
					case *ast.StructType:
						st := Struct{
							Doc:  doc,
							Name: p.Name.Name,
						}
						for _, field := range e.Fields.List {
							typ, ok := field.Type.(*ast.Ident)
							if !ok {
								continue
							}
							fi := Field{
								Name: field.Names[0].Name,
								Type: typ.Name,
							}
							if field.Tag != nil {
								fi.Tag = field.Tag.Value
							}
							st.Fields = append(st.Fields, fi)
						}
						structs = append(structs, st)
					case *ast.InterfaceType:
						it := Interface{
							Doc:  doc,
							Name: p.Name.Name,
						}
						for _, method := range e.Methods.List {
							fun, ok := method.Type.(*ast.FuncType)
							if !ok {
								continue
							}
							m := Method{
								Name: method.Names[0].Name,
							}
							if fun.Params != nil {
								for _, param := range fun.Params.List {
									typ, ok := param.Type.(*ast.Ident)
									if !ok {
										continue
									}
									f := Field{
										Type: typ.Name,
									}
									if len(param.Names) > 0 {
										f.Name = param.Names[0].Name
									}
									m.Params = append(m.Params, f)
								}
							}
							if fun.Results != nil {
								for _, result := range fun.Results.List {
									typ, ok := result.Type.(*ast.Ident)
									if !ok {
										continue
									}
									f := Field{
										Type: typ.Name,
									}
									if len(result.Names) > 0 {
										f.Name = result.Names[0].Name
									}
									m.Return = append(m.Return, f)
								}
							}
							it.Methods = append(it.Methods, m)
						}
						interfaces = append(interfaces, it)
					}
				}
			}
		}
		return true
	})

	return &Package{
		Name:       packageName,
		Structs:    structs,
		Interfaces: interfaces,
		Imports:    imports,
	}, nil
}
