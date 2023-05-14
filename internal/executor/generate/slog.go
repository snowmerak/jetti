package generate

import (
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
	"strings"
)

func Slog(path string, slogs []check.Slog) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)

	for _, slog := range slogs {
		logCodeTemplate := func(level string) []string {
			codes := []string{
				"if obj == nil {",
				"\treturn err" + slog.Name + "IsNil",
				"}",
				"$RECEIVER$.logger." + level + "(",
				"\tmsg,",
			}
			for _, field := range slog.Fields {
				codes = append(codes, "\t\""+strings.ToLower(field.Name)+"\", "+"obj."+field.Name+",")
			}
			codes = append(codes, ")")
			codes = append(codes, "return nil")

			return codes
		}

		pkg := &model.Package{
			Name: packageName,
			Imports: []model.Import{
				{
					Path: "golang.org/x/exp/slog",
				},
				{
					Path: "errors",
				},
			},
			GlobalVariables: []model.GlobalVariable{
				{
					Name:  "err" + slog.Name + "IsNil",
					Type:  "error",
					Value: "errors.New(\"" + slog.Name + " is nil\")",
				},
			},
			Structs: []model.Struct{
				{
					Name: slog.Name + "Logger",
					Fields: []model.Field{
						{
							Name: "logger",
							Type: "*slog.Logger",
						},
					},
					Methods: []model.Method{
						{
							Name: "Debug",
							Params: []model.Field{
								{
									Name: "msg",
									Type: "string",
								},
								{
									Name: "obj",
									Type: "*" + slog.Name,
								},
							},
							Return: []model.Field{
								{
									Name: "err",
									Type: "error",
								},
							},
							Code: logCodeTemplate("Debug"),
						},
						{
							Name: "Info",
							Params: []model.Field{
								{
									Name: "msg",
									Type: "string",
								},
								{
									Name: "obj",
									Type: "*" + slog.Name,
								},
							},
							Return: []model.Field{
								{
									Name: "err",
									Type: "error",
								},
							},
							Code: logCodeTemplate("Info"),
						},
						{
							Name: "Warn",
							Params: []model.Field{
								{
									Name: "msg",
									Type: "string",
								},
								{
									Name: "obj",
									Type: "*" + slog.Name,
								},
							},
							Return: []model.Field{
								{
									Name: "err",
									Type: "error",
								},
							},
							Code: logCodeTemplate("Warn"),
						},
						{
							Name: "Error",
							Params: []model.Field{
								{
									Name: "msg",
									Type: "string",
								},
								{
									Name: "obj",
									Type: "*" + slog.Name,
								},
							},
							Return: []model.Field{
								{
									Name: "err",
									Type: "error",
								},
							},
							Code: logCodeTemplate("Error"),
						},
					},
				},
			},
			Functions: []model.Function{
				{
					Name: "New" + slog.Name + "Logger",
					Params: []model.Field{
						{
							Name: "logger",
							Type: "slog.Handler",
						},
					},
					Return: []model.Field{
						{},
					},
				},
			},
		}

		data, err := generator.GenerateFile(pkg)
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(dir, strings.ToLower(slog.Name)+".slog.go"))
		if err != nil {
			return err
		}

		if _, err := f.Write(data); err != nil {
			return err
		}
	}

	return nil
}
