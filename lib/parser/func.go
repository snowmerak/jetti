package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
)

func ParseFunction(node ast.Node) *model.Function {
	if node == nil {
		return nil
	}
	x, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}
	recv := ""
	if x.Recv != nil {
		recv = ParseFieldList(x.Recv, ", ")
	}
	return &model.Function{
		Name:     x.Name.Name,
		Receiver: recv,
		Doc:      x.Doc.Text(),
		Params:   ParseFields(x.Type.Params),
		Return:   ParseFields(x.Type.Results),
		Code:     ParseFuncBody(x.Body),
	}
}

func ParseFuncBody(body *ast.BlockStmt) []string {
	if body == nil {
		return nil
	}

	code := []string(nil)

	for _, stmt := range body.List {
		code = append(code, ParseStmt(stmt))
	}

	return code
}

func ParseStmt(node ast.Node) string {
	switch x := node.(type) {
	case *ast.ExprStmt:
		return ParseExpr(x.X)
	case *ast.AssignStmt:
		return ParseExpr(x.Lhs[0]) + " " + x.Tok.String() + " " + ParseExpr(x.Rhs[0])
	case *ast.ReturnStmt:
		return "return " + ParseExpr(x.Results[0])
	case *ast.IfStmt:
		return "if " + ParseExpr(x.Cond) + " {" + ParseStmt(x.Body) + "}"
	case *ast.ForStmt:
		return "for " + ParseStmt(x.Body)
	case *ast.RangeStmt:
		return "for " + ParseExpr(x.Key) + ", " + ParseExpr(x.Value) + " := range " + ParseExpr(x.X) + " {" + ParseStmt(x.Body) + "}"
	case *ast.BlockStmt:
		return "{" + ParseStmt(x.List[0]) + "}"
	case *ast.DeclStmt:
		return ParseDecl(x.Decl)
	case *ast.IncDecStmt:
		return ParseExpr(x.X) + x.Tok.String()
	case *ast.SwitchStmt:
		return "switch " + ParseExpr(x.Tag) + " {" + ParseStmt(x.Body) + "}"
	case *ast.CaseClause:
		return "case " + ParseExpr(x.List[0]) + ": " + ParseStmt(x.Body[0])
	case *ast.TypeSwitchStmt:
		return "switch " + ParseStmt(x.Body)
	case *ast.TypeAssertExpr:
		return ParseExpr(x.X) + ".(" + ParseExpr(x.Type) + ")"
	case *ast.SelectStmt:
		return "select {" + ParseStmt(x.Body) + "}"
	case *ast.CommClause:
		return "case " + ParseStmt(x.Body[0])
	case *ast.SendStmt:
		return ParseExpr(x.Chan) + " <- " + ParseExpr(x.Value)
	case *ast.BranchStmt:
		return x.Tok.String()
	case *ast.LabeledStmt:
		return x.Label.Name + ": " + ParseStmt(x.Stmt)
	case *ast.EmptyStmt:
		return ""
	}
	return ""
}

// ParseDecl parses a declaration.
// TODO: ParseDecl
func ParseDecl(node ast.Decl) string {
	return ""
}

func ParseExpr(node ast.Expr) string {
	switch x := node.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + ParseExpr(x.X)
	case *ast.SelectorExpr:
		return ParseExpr(x.X) + "." + ParseExpr(x.Sel)
	case *ast.ArrayType:
		return "[]" + ParseExpr(x.Elt)
	case *ast.MapType:
		return "map[" + ParseExpr(x.Key) + "]" + ParseExpr(x.Value)
	case *ast.InterfaceType:
		return "interface{" + ParseFieldList(x.Methods, "; ") + "}"
	case *ast.ChanType:
		switch x.Dir {
		case ast.SEND:
			return "chan<- " + ParseExpr(x.Value)
		case ast.RECV:
			return "<-chan " + ParseExpr(x.Value)
		default:
			return "chan " + ParseExpr(x.Value)
		}
	case *ast.FuncType:
		return "func(" + ParseFieldList(x.Params, ", ") + ") (" + ParseFieldList(x.Results, ", ") + ")"
	case *ast.Ellipsis:
		return "..." + ParseExpr(x.Elt)
	case *ast.BasicLit:
		return x.Value
	case *ast.CallExpr:
		v := ParseExpr(x.Fun) + "("
		for _, arg := range x.Args {
			v += ParseExpr(arg) + ", "
		}
		return v[:len(v)-2] + ")"
	case *ast.BinaryExpr:
		return ParseExpr(x.X) + " " + x.Op.String() + " " + ParseExpr(x.Y)
	case *ast.UnaryExpr:
		return x.Op.String() + ParseExpr(x.X)
	case *ast.CompositeLit:
		return ParseExpr(x.Type) + "{" + ParseExpr(x.Elts[0]) + "}"
	case *ast.IndexExpr:
		return ParseExpr(x.X) + "[" + ParseExpr(x.Index) + "]"
	case *ast.SliceExpr:
		return ParseExpr(x.X) + "[" + ParseExpr(x.Low) + ":" + ParseExpr(x.High) + "]"
	case *ast.KeyValueExpr:
		return ParseExpr(x.Key) + ":" + ParseExpr(x.Value)
	case *ast.TypeAssertExpr:
		return ParseExpr(x.X) + ".(" + ParseExpr(x.Type) + ")"
	case *ast.ParenExpr:
		return "(" + ParseExpr(x.X) + ")"
	}
	return ""
}
