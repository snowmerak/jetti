package executor

import (
	"bytes"
	"github.com/snowmerak/jetti/lib/generator"
	"github.com/snowmerak/jetti/lib/model"
	"github.com/snowmerak/jetti/lib/strcase"
	"github.com/twpayne/go-jsonstruct/v2"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

func convertJson(path string) error {
	packageName := filepath.Base(filepath.Dir(path))
	fileName := filepath.Base(path)
	fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
	structName := strcase.SnakeToPascal(fileName)

	gen := jsonstruct.NewGenerator(jsonstruct.WithTypeName(structName))

	if err := gen.ObserveJSONFile(path); err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)

	pkg := &model.Package{
		Name: packageName,
		Imports: []model.Import{
			{
				Path: "github.com/goccy/go-json",
			},
			{
				Path: "io",
			},
			{
				Path: "os",
			},
		},
		Functions: []model.Function{
			{
				Name: structName + "FromJSON",
				Params: []model.Field{
					{
						Name: "data",
						Type: "[]byte",
					},
				},
				Return: []model.Field{
					{
						Type: "*" + structName,
					},
					{
						Type: "error",
					},
				},
				Code: []string{
					"v := new(" + structName + ")",
					"if err := json.Unmarshal(data, &v); err != nil {",
					"\treturn nil, err",
					"}",
					"return v, nil",
				},
			},
			{
				Name: structName + "FromFile",
				Params: []model.Field{
					{
						Name: "path",
						Type: "string",
					},
				},
				Return: []model.Field{
					{
						Type: "*" + structName,
					},
					{
						Type: "error",
					},
				},
				Code: []string{
					"f, err := os.ReadFile(path)",
					"if err != nil {",
					"\treturn nil, err",
					"}",
					"return " + structName + "FromJSON(f)",
				},
			},
		},
		Methods: []model.Method{
			{
				Name: "Marshal2JSON",
				Receiver: model.Field{
					Name: strings.ToLower(structName),
					Type: structName,
				},
				Return: []model.Field{
					{
						Type: "[]byte",
					},
					{
						Type: "error",
					},
				},
				Code: []string{
					"return json.Marshal($RECEIVER$)",
				},
			},
			{
				Name: "Encode2JSON",
				Receiver: model.Field{
					Name: strings.ToLower(structName),
					Type: structName,
				},
				Params: []model.Field{
					{
						Name: "w",
						Type: "io.Writer",
					},
				},
				Return: []model.Field{
					{
						Type: "error",
					},
				},
				Code: []string{
					"return json.NewEncoder(w).Encode($RECEIVER$)",
				},
			},
		},
	}

	data, err := generator.GenerateFile(pkg)
	if err != nil {
		return err
	}

	buffer.Write(data)

	data, err = gen.Generate()
	if err != nil {
		return err
	}

	data = bytes.ReplaceAll(data, []byte("package main"), []byte{})

	buffer.Write(data)

	data, err = format.Source(buffer.Bytes())
	if err != nil {
		return err
	}

	f, err := os.Create(path + ".go")
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func ModelYaml(path string) {

}
