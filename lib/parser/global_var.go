package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
	"go/token"
)

func ParseGlobalVariable(node ast.Node) *model.GlobalVariable {
	switch x := node.(type) {
	case *ast.GenDecl:
		if x.Tok != token.VAR {
			return nil
		}
		if len(x.Specs) < 1 {
			return nil
		}
		vs, ok := x.Specs[0].(*ast.ValueSpec)
		if !ok {
			return nil
		}
		name := ""
		if len(vs.Names) > 0 {
			name = vs.Names[0].Name
		}
		value := ""
		if len(vs.Values) > 0 {
			value = ParseName(vs.Values[0])
		}
		return &model.GlobalVariable{
			Doc:   x.Doc.Text(),
			Name:  name,
			Type:  ParseName(vs.Type),
			Value: value,
		}
	}
	return nil
}
