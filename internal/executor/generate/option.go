package generate

import (
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
	"strings"
)

func Option(path string, opts []string) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)

	for _, opt := range opts {
		typeName := strings.ReplaceAll(opt, "*", "")

		err := func() error {
			someFunc := "Some"
			if packageName != strings.ToLower(typeName) {
				someFunc += typeName
			}

			noneFunc := "None"
			if packageName != strings.ToLower(typeName) {
				noneFunc += typeName
			}

			pkg := &model.Package{
				Name: packageName,
				Structs: []model.Struct{
					{
						Name: "Optional" + typeName,
						Fields: []model.Field{
							{
								Name: "value",
								Type: opt,
							},
							{
								Name: "valid",
								Type: "bool",
							},
						},
						Methods: []model.Method{
							{
								Name: "Unwrap",
								Return: []model.Field{
									{
										Type: opt,
									},
								},
								Code: []string{
									"if !o.valid {",
									"\tpanic(\"unwrap a none value\")",
									"}",
									"return o.value",
								},
							},
							{
								Name: "IsSome",
								Return: []model.Field{
									{
										Type: "bool",
									},
								},
								Code: []string{
									"return o.valid",
								},
							},
							{
								Name: "IsNone",
								Return: []model.Field{
									{
										Type: "bool",
									},
								},
								Code: []string{
									"return !o.valid",
								},
							},
							{
								Name: "UnwrapOr",
								Params: []model.Field{
									{
										Name: "defaultValue",
										Type: opt,
									},
								},
								Return: []model.Field{
									{
										Type: opt,
									},
								},
								Code: []string{
									"if !o.valid {",
									"\treturn defaultValue",
									"}",
									"return o.value",
								},
							},
						},
					},
				},
				Functions: []model.Function{
					{
						Name: someFunc,
						Params: []model.Field{
							{
								Name: "value",
								Type: opt,
							},
						},
						Return: []model.Field{
							{
								Type: "Optional" + typeName,
							},
						},
						Code: []string{
							"return Optional" + typeName + "{",
							"\tvalue: value,",
							"\tvalid: true,",
							"}",
						},
					},
					{
						Name: noneFunc,
						Return: []model.Field{
							{
								Type: "Optional" + typeName,
							},
						},
						Code: []string{
							"return Optional" + typeName + "{",
							"\tvalid: false,",
							"}",
						},
					},
				},
			}

			data, err := generator.GenerateFile(pkg)
			if err != nil {
				return err
			}

			f, err := os.Create(filepath.Join(dir, strings.ToLower(typeName)+".option.go"))
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.Write(data); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
