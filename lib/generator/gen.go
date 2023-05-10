package generator

import (
	"bytes"
	"go/format"
	"strings"

	"github.com/snowmerak/jetti/v2/lib/model"
)

const (
	RECEIVER = "$RECEIVER$"
)

func GenerateFile(pkg *model.Package) ([]byte, error) {
	rs := bytes.NewBuffer(nil)

	rs.WriteString("package ")
	rs.WriteString(pkg.Name)
	rs.WriteString("\n\n")

	for _, imp := range pkg.Imports {
		rs.WriteString("import ")
		if imp.Alias != "" {
			rs.WriteString(imp.Alias)
			rs.WriteString(" ")
		}
		rs.WriteString("\"")
		rs.WriteString(imp.Path)
		rs.WriteString("\"\n")
	}

	for _, alias := range pkg.Aliases {
		rs.WriteString("\ntype ")
		rs.WriteString(alias.Name)
		rs.WriteString(" ")
		rs.WriteString(alias.Type)
		rs.WriteString("\n")
	}

	for _, variable := range pkg.GlobalVariables {
		rs.WriteString("\nvar ")
		rs.WriteString(variable.Name)
		rs.WriteString(" ")
		rs.WriteString(variable.Type)
		if variable.Value != "" {
			rs.WriteString(" = ")
			rs.WriteString(variable.Value)
		}
	}

	for _, st := range pkg.Structs {
		rs.WriteString("\n")
		if st.Doc != "" {
			rs.WriteString(st.Doc)
			rs.WriteString("\n")
		}
		rs.WriteString("type ")
		rs.WriteString(st.Name)
		rs.WriteString(" struct {\n")
		for _, field := range st.Fields {
			rs.WriteString("\t")
			rs.WriteString(field.Name)
			rs.WriteString(" ")
			rs.WriteString(field.Type)
			if field.Tag != "" {
				rs.WriteString(" `")
				rs.WriteString(field.Tag)
				rs.WriteString("`")
			}
			rs.WriteString("\n")
		}
		rs.WriteString("}\n")

		receiver := strings.ToLower(st.Name[:1])
		for _, method := range st.Methods {
			rs.WriteString("\n")
			rs.WriteString("func (" + receiver + " *")
			rs.WriteString(st.Name)
			rs.WriteString(") ")
			rs.WriteString(method.Name)
			rs.WriteString("(")
			for i, param := range method.Params {
				if i > 0 {
					rs.WriteString(", ")
				}
				rs.WriteString(param.Name)
				rs.WriteString(" ")
				rs.WriteString(param.Type)
			}
			rs.WriteString(")")
			if len(method.Return) > 0 {
				rs.WriteString(" (")
				for i, ret := range method.Return {
					if i > 0 {
						rs.WriteString(", ")
					}
					rs.WriteString(ret.Name)
					rs.WriteString(" ")
					rs.WriteString(ret.Type)
				}
				rs.WriteString(")")
			}
			rs.WriteString(" {\n")
			for _, code := range method.Code {
				rs.WriteString("\t")
				rs.WriteString(strings.ReplaceAll(code, RECEIVER, receiver))
				rs.WriteString("\n")
			}
			rs.WriteString("}\n")
		}
	}

	for _, it := range pkg.Interfaces {
		rs.WriteString("\n")
		if it.Doc != "" {
			rs.WriteString(it.Doc)
			rs.WriteString("\n")
		}
		rs.WriteString("type ")
		rs.WriteString(it.Name)
		rs.WriteString(" interface {\n")
		for _, method := range it.Methods {
			rs.WriteString("\t")
			rs.WriteString(method.Name)
			rs.WriteString("(")
			for i, param := range method.Params {
				if i > 0 {
					rs.WriteString(", ")
				}
				rs.WriteString(param.Name)
				rs.WriteString(" ")
				rs.WriteString(param.Type)
			}
			rs.WriteString(")")
			if len(method.Return) > 0 {
				rs.WriteString(" (")
				for i, ret := range method.Return {
					if i > 0 {
						rs.WriteString(", ")
					}
					rs.WriteString(ret.Name)
					rs.WriteString(" ")
					rs.WriteString(ret.Type)
				}
				rs.WriteString(")")
			}
			rs.WriteString("\n")
		}
		rs.WriteString("}\n")
	}

	for _, fun := range pkg.Functions {
		rs.WriteString("\n")
		rs.WriteString("func ")
		rs.WriteString(fun.Name)
		rs.WriteString("(")
		for i, param := range fun.Params {
			if i > 0 {
				rs.WriteString(", ")
			}
			rs.WriteString(param.Name)
			rs.WriteString(" ")
			rs.WriteString(param.Type)
		}
		rs.WriteString(")")
		if len(fun.Return) > 0 {
			rs.WriteString(" (")
			for i, ret := range fun.Return {
				if i > 0 {
					rs.WriteString(", ")
				}
				rs.WriteString(ret.Name)
				rs.WriteString(" ")
				rs.WriteString(ret.Type)
			}
			rs.WriteString(")")
		}
		rs.WriteString(" {\n")
		for _, code := range fun.Code {
			rs.WriteString("\t")
			rs.WriteString(code)
			rs.WriteString("\n")
		}
		rs.WriteString("}\n")
	}

	for _, method := range pkg.Methods {
		rs.WriteString("\n")
		rs.WriteString("func (" + method.Receiver.Name + " *")
		rs.WriteString(method.Receiver.Type)
		rs.WriteString(") ")
		rs.WriteString(method.Name)
		rs.WriteString("(")
		for i, param := range method.Params {
			if i > 0 {
				rs.WriteString(", ")
			}
			rs.WriteString(param.Name)
			rs.WriteString(" ")
			rs.WriteString(param.Type)
		}
		rs.WriteString(")")
		if len(method.Return) > 0 {
			rs.WriteString(" (")
			for i, ret := range method.Return {
				if i > 0 {
					rs.WriteString(", ")
				}
				rs.WriteString(ret.Name)
				rs.WriteString(" ")
				rs.WriteString(ret.Type)
			}
			rs.WriteString(")")
		}
		rs.WriteString(" {\n")
		for _, code := range method.Code {
			rs.WriteString("\t")
			rs.WriteString(strings.ReplaceAll(code, RECEIVER, method.Receiver.Name))
			rs.WriteString("\n")
		}
		rs.WriteString("}\n")
	}

	fd, err := format.Source(rs.Bytes())
	if err != nil {
		return nil, err
	}

	return fd, nil
}
