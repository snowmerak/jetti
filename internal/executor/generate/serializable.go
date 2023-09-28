package generate

import (
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
)

const serializableDirectory = "serializable"

func JsonSerializable(root string) error {
	genPath := filepath.Join(root, "gen", serializableDirectory)
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(MakeGeneratedFileName(genPath, "json", "serializable"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	pkg := &model.Package{
		Name: "serializable",
		Imports: []model.Import{
			{
				Path: "io",
			},
		},
		Interfaces: []model.Interface{
			{
				Name: "JsonSerializable",
				Methods: []model.Method{
					{
						Name: "Marshal2JSON",
						Return: []model.Field{
							{
								Type: "[]byte",
							},
							{
								Type: "error",
							},
						},
					},
					{
						Name: "UnmarshalFromJSON",
						Params: []model.Field{
							{
								Type: "[]byte",
							},
						},
						Return: []model.Field{
							{
								Type: "error",
							},
						},
					},
					{
						Name: "Encode2JSON",
						Params: []model.Field{
							{
								Type: "io.Writer",
							},
						},
						Return: []model.Field{
							{
								Type: "error",
							},
						},
					},
					{
						Name: "DecodeFromJSON",
						Params: []model.Field{
							{
								Type: "io.Reader",
							},
						},
						Return: []model.Field{
							{
								Type: "error",
							},
						},
					},
				},
			},
		},
	}

	data, err := generator.GenerateFile(pkg)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func YamlSerializable(root string) error {
	genPath := filepath.Join(root, "gen", serializableDirectory)
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(MakeGeneratedFileName(genPath, "yaml", "serializable"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	pkg := &model.Package{
		Name: "serializable",
		Imports: []model.Import{
			{
				Path: "io",
			},
		},
		Interfaces: []model.Interface{
			{
				Name: "YamlSerializable",
				Methods: []model.Method{
					{
						Name: "Marshal2YAML",
						Return: []model.Field{
							{
								Type: "[]byte",
							},
							{
								Type: "error",
							},
						},
					},
					{
						Name: "UnmarshalFromYAML",
						Params: []model.Field{
							{
								Type: "[]byte",
							},
						},
						Return: []model.Field{
							{
								Type: "error",
							},
						},
					},
					{
						Name: "Encode2YAML",
						Params: []model.Field{
							{
								Type: "io.Writer",
							},
						},
						Return: []model.Field{
							{
								Type: "error",
							},
						},
					},
					{
						Name: "DecodeFromYAML",
						Params: []model.Field{
							{
								Type: "io.Reader",
							},
						},
						Return: []model.Field{
							{
								Type: "error",
							},
						},
					},
				},
			},
		},
	}

	data, err := generator.GenerateFile(pkg)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}
