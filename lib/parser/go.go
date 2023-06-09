package parser

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func ParseFile(path string) (*model.Package, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	packageName := ""
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
			packageName = x.Name.Name
		case *ast.ImportSpec:
			imp := model.Import{
				Path: x.Path.Value,
			}
			if x.Name != nil {
				imp.Alias = x.Name.Name
			}
			imports = append(imports, imp)
		case *ast.FuncDecl:
			fun := model.Function{
				Name: x.Name.Name,
			}
			if x.Doc != nil {
				fun.Doc = x.Doc.Text()
			}
			if x.Recv != nil {
				switch v := x.Recv.List[0].Type.(type) {
				case *ast.StarExpr:
					fun.Receiver = v.X.(*ast.Ident).Name
				case *ast.Ident:
					fun.Receiver = v.Name
				}
			}
			if x.Type != nil {
				if x.Type.Params != nil {
					for _, param := range x.Type.Params.List {
						typ, ok := param.Type.(*ast.Ident)
						if !ok {
							continue
						}
						f := model.Field{
							Type: typ.Name,
						}
						if len(param.Names) > 0 {
							f.Name = param.Names[0].Name
						}
						fun.Params = append(fun.Params, f)
					}
				}
				if x.Type.Results != nil {
					for _, result := range x.Type.Results.List {
						typ, ok := result.Type.(*ast.Ident)
						if !ok {
							continue
						}
						f := model.Field{
							Type: typ.Name,
						}
						if len(result.Names) > 0 {
							f.Name = result.Names[0].Name
						}
						fun.Return = append(fun.Return, f)
					}
				}
			}
			functions = append(functions, fun)
		case *ast.GenDecl:
			doc := x.Doc.Text()
			for _, spec := range x.Specs {
				switch p := spec.(type) {
				case *ast.TypeSpec:
					switch e := p.Type.(type) {
					case *ast.StructType:
						st := model.Struct{
							Doc:  doc,
							Name: p.Name.Name,
						}
						for _, field := range e.Fields.List {
							typ, ok := field.Type.(*ast.Ident)
							if !ok {
								continue
							}
							name := ""
							if len(field.Names) > 0 {
								name = field.Names[0].Name
							}
							fi := model.Field{
								Name: name,
								Type: typ.Name,
							}
							if field.Tag != nil {
								fi.Tag = field.Tag.Value
							}
							st.Fields = append(st.Fields, fi)
						}
						structs = append(structs, st)
					case *ast.InterfaceType:
						it := model.Interface{
							Doc:  doc,
							Name: p.Name.Name,
						}
						for _, method := range e.Methods.List {
							fun, ok := method.Type.(*ast.FuncType)
							if !ok {
								continue
							}
							m := model.Method{
								Name: method.Names[0].Name,
							}
							if fun.Params != nil {
								for _, param := range fun.Params.List {
									typ, ok := param.Type.(*ast.Ident)
									if !ok {
										continue
									}
									f := model.Field{
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
									f := model.Field{
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
					case *ast.ArrayType:
						length := ""
						if p.Type.(*ast.ArrayType).Len != nil {
							length = p.Type.(*ast.ArrayType).Len.(*ast.BasicLit).Value
						}
						arr := model.Alias{
							Doc:  doc,
							Name: p.Name.Name,
							Type: fmt.Sprintf("[%s]%s", length, p.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name),
						}
						aliases = append(aliases, arr)
					case *ast.MapType:
						m := model.Alias{
							Doc:  doc,
							Name: p.Name.Name,
							Type: fmt.Sprintf("map[%s]%s", p.Type.(*ast.MapType).Key.(*ast.Ident).Name, p.Type.(*ast.MapType).Value.(*ast.Ident).Name),
						}
						aliases = append(aliases, m)
					case *ast.Ident:
						ali := model.Alias{
							Doc:  doc,
							Name: p.Name.Name,
							Type: p.Type.(*ast.Ident).Name,
						}
						aliases = append(aliases, ali)
					case *ast.FuncType:
						params := []string(nil)
						if e.Params != nil {
							for _, param := range e.Params.List {
								typ := ""
								switch param.Type.(type) {
								case *ast.Ident:
									typ = param.Type.(*ast.Ident).Name
								case *ast.StarExpr:
									typ = "*" + param.Type.(*ast.StarExpr).X.(*ast.Ident).Name
								}
								params = append(params, typ)
							}
						}
						returns := []string(nil)
						if e.Results != nil {
							for _, result := range e.Results.List {
								typ := ""
								switch result.Type.(type) {
								case *ast.Ident:
									typ = result.Type.(*ast.Ident).Name
								case *ast.StarExpr:
									typ = "*" + result.Type.(*ast.StarExpr).X.(*ast.Ident).Name
								}
								returns = append(returns, typ)
							}
						}
						fun := model.Alias{
							Doc:  doc,
							Name: p.Name.Name,
							Type: fmt.Sprintf("func(%s) (%s)", strings.Join(params, ", "), strings.Join(returns, ", ")),
						}
						aliases = append(aliases, fun)
					}
				}
			}
		}
		return true
	})

	return &model.Package{
		Name:       packageName,
		Structs:    structs,
		Interfaces: interfaces,
		Imports:    imports,
		Functions:  functions,
		Aliases:    aliases,
	}, nil
}
