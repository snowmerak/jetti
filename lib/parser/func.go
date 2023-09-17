package parser

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"go/ast"
	"go/token"
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
		rs := ""
		for i, expr := range x.Lhs {
			if i > 0 {
				rs += ", "
			}
			rs += ParseExpr(expr)
		}
		rs += " " + x.Tok.String() + " "
		for i, expr := range x.Rhs {
			if i > 0 {
				rs += ", "
			}
			rs += ParseExpr(expr)
		}
		return rs
	case *ast.ReturnStmt:
		rs := "return "
		for i, result := range x.Results {
			if i > 0 {
				rs += ", "
			}
			rs += ParseExpr(result)
		}
		return rs
	case *ast.IfStmt:
		return "if " + ParseExpr(x.Cond) + " {" + ParseStmt(x.Body) + "}"
	case *ast.ForStmt:
		return "for " + ParseStmt(x.Body)
	case *ast.RangeStmt:
		return "for " + ParseExpr(x.Key) + ", " + ParseExpr(x.Value) + " := range " + ParseExpr(x.X) + " {" + ParseStmt(x.Body) + "}"
	case *ast.BlockStmt:
		rs := "{ "
		for _, stmt := range x.List {
			rs += ParseStmt(stmt)
		}
		return rs + " }"
	case *ast.DeclStmt:
		return ParseDecl(x.Decl)
	case *ast.IncDecStmt:
		return ParseExpr(x.X) + x.Tok.String()
	case *ast.SwitchStmt:
		return "switch " + ParseExpr(x.Tag) + " {" + ParseStmt(x.Body) + "}"
	case *ast.CaseClause:
		rs := "case "
		for i, expr := range x.List {
			if i > 0 {
				rs += ", "
			}
			rs += ParseExpr(expr)
		}
		rs += ":"
		for _, stmt := range x.Body {
			rs += ParseStmt(stmt)
		}
		return rs
	case *ast.TypeSwitchStmt:
		return "switch " + ParseStmt(x.Body)
	case *ast.TypeAssertExpr:
		return ParseExpr(x.X) + ".(" + ParseExpr(x.Type) + ")"
	case *ast.SelectStmt:
		return "select {" + ParseStmt(x.Body) + "}"
	case *ast.CommClause:
		rs := "case "
		if x.Comm != nil {
			rs += ParseStmt(x.Comm)
		}
		rs += ":"
		for _, stmt := range x.Body {
			rs += ParseStmt(stmt)
		}
		return rs
	case *ast.SendStmt:
		return ParseExpr(x.Chan) + " <- " + ParseExpr(x.Value)
	case *ast.BranchStmt:
		return x.Tok.String()
	case *ast.LabeledStmt:
		return x.Label.Name + ": " + ParseStmt(x.Stmt)
	case *ast.EmptyStmt:
		return ""
	case *ast.DeferStmt:
		return "defer " + ParseExpr(x.Call)
	case *ast.GoStmt:
		return "go " + ParseExpr(x.Call)
	}
	return ""
}

// ParseDecl parses a declaration.
func ParseDecl(node ast.Decl) string {
	switch x := node.(type) {
	case *ast.GenDecl:
		switch x.Tok {
		case token.CONST:
			rs := "const "
			for i, spec := range x.Specs {
				if i > 0 {
					rs += ", "
				}
				rs += ParseSpec(spec)
			}
			return rs
		case token.TYPE:
			rs := "type "
			for i, spec := range x.Specs {
				if i > 0 {
					rs += ", "
				}
				rs += ParseSpec(spec)
			}
			return rs
		case token.VAR:
			rs := "var "
			for i, spec := range x.Specs {
				if i > 0 {
					rs += ", "
				}
				rs += ParseSpec(spec)
			}
			return rs
		case token.IMPORT:
			rs := "import "
			for i, spec := range x.Specs {
				if i > 0 {
					rs += ", "
				}
				rs += ParseSpec(spec)
			}
		}
	}
	return ""
}

func ParseSpec(spec ast.Spec) string {
	switch x := spec.(type) {
	case *ast.ValueSpec:
		rs := ""
		for i, name := range x.Names {
			if i > 0 {
				rs += ", "
			}
			rs += name.Name
		}
		if x.Type != nil {
			rs += " " + ParseExpr(x.Type)
		}
		if x.Values != nil {
			rs += " = "
			for i, value := range x.Values {
				if i > 0 {
					rs += ", "
				}
				rs += ParseExpr(value)
			}
		}
		return rs
	case *ast.TypeSpec:
		return x.Name.Name + " " + ParseExpr(x.Type)
	}
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
		for i, arg := range x.Args {
			if i > 0 {
				v += ", "
			}
			v += ParseExpr(arg)
		}
		return v + ")"
	case *ast.BinaryExpr:
		return ParseExpr(x.X) + " " + x.Op.String() + " " + ParseExpr(x.Y)
	case *ast.UnaryExpr:
		return x.Op.String() + ParseExpr(x.X)
	case *ast.CompositeLit:
		rs := ParseExpr(x.Type) + "{"
		for i, elt := range x.Elts {
			if i > 0 {
				rs += ", "
			}
			rs += ParseExpr(elt)
		}
		return rs + "}"
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
