package fp

import (
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
)

func MonadOption(genPath string) error {
	genPath = filepath.Join(genPath, "option")
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(generate.MakeGeneratedFileName(genPath, "option"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	pkg := &model.Package{
		Name: "option",
		Structs: []model.Struct{
			{
				Name: "Option[T any]",
				Fields: []model.Field{
					{
						Name: "value",
						Type: "any",
					},
				},
			},
		},
		Methods: []model.Method{
			{
				Name: "Unwrap",
				Receiver: model.Field{
					Name: "o",
					Type: "Option[T]",
				},
				Return: []model.Field{
					{
						Type: "T",
					},
				},
				Code: []string{
					"if o.value == nil {",
					"\tpanic(\"unwrap a nil value\")",
					"}",
					"r, ok := o.value.(T)",
					"if !ok {",
					"\tpanic(\"unwrap a invalid value\")",
					"}",
					"return r",
				},
			},
			{
				Name: "IsSome",
				Receiver: model.Field{
					Name: "o",
					Type: "Option[T]",
				},
				Return: []model.Field{
					{
						Type: "bool",
					},
				},
				Code: []string{
					"return o.value != nil",
				},
			},
			{
				Name: "IsNone",
				Receiver: model.Field{
					Name: "o",
					Type: "Option[T]",
				},
				Return: []model.Field{
					{
						Type: "bool",
					},
				},
				Code: []string{
					"return o.value == nil",
				},
			},
			{
				Name: "UnwrapOr",
				Receiver: model.Field{
					Name: "o",
					Type: "Option[T]",
				},
				Params: []model.Field{
					{
						Name: "defaultValue",
						Type: "T",
					},
				},
				Return: []model.Field{
					{
						Type: "T",
					},
				},
				Code: []string{
					"if o.value == nil {",
					"\treturn defaultValue",
					"}",
					"r, ok := o.value.(T)",
					"if !ok {",
					"\treturn defaultValue",
					"}",
					"return r",
				},
			},
		},
		Functions: []model.Function{
			{
				Name: "Some[T any]",
				Params: []model.Field{
					{
						Name: "value",
						Type: "T",
					},
				},
				Return: []model.Field{
					{
						Type: "Option[T]",
					},
				},
				Code: []string{
					"return Option[T]{",
					"\tvalue: value,",
					"}",
				},
			},
			{
				Name: "None[T any]",
				Return: []model.Field{
					{
						Type: "Option[T]",
					},
				},
				Code: []string{
					"return Option[T]{",
					"\tvalue: nil,",
					"}",
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

func MonadResult(genPath string) error {
	genPath = filepath.Join(genPath, "result")
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(generate.MakeGeneratedFileName(genPath, "result"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	pkg := &model.Package{
		Name: "result",
		Structs: []model.Struct{
			{
				Name: "Result[T, R any]",
				Fields: []model.Field{
					{
						Name: "value",
						Type: "any",
					},
					{
						Name: "isOk",
						Type: "bool",
					},
				},
			},
		},
		Methods: []model.Method{
			{
				Name: "Unwrap",
				Receiver: model.Field{
					Name: "r",
					Type: "Result[T, R]",
				},
				Return: []model.Field{
					{
						Type: "T",
					},
				},
				Code: []string{
					"if !r.isOk {",
					"\tpanic(\"unwrap a error value\")",
					"}",
					"v, ok := r.value.(T)",
					"if !ok {",
					"\tpanic(\"unwrap a invalid value\")",
					"}",
					"return v",
				},
			},
			{
				Name: "IsOk",
				Receiver: model.Field{
					Name: "r",
					Type: "Result[T, R]",
				},
				Return: []model.Field{
					{
						Type: "bool",
					},
				},
				Code: []string{
					"return r.isOk",
				},
			},
			{
				Name: "IsErr",
				Receiver: model.Field{
					Name: "r",
					Type: "Result[T, R]",
				},
				Return: []model.Field{
					{
						Type: "bool",
					},
				},
				Code: []string{
					"return !r.isOk",
				},
			},
			{
				Name: "UnwrapOr",
				Receiver: model.Field{
					Name: "r",
					Type: "Result[T, R]",
				},
				Params: []model.Field{
					{
						Name: "defaultValue",
						Type: "T",
					},
				},
				Return: []model.Field{
					{
						Type: "T",
					},
				},
				Code: []string{
					"if !r.isOk {",
					"\treturn defaultValue",
					"}",
					"v, ok := r.value.(T)",
					"if !ok {",
					"\treturn defaultValue",
					"}",
					"return v",
				},
			},
			{
				Name: "UnwrapErr",
				Receiver: model.Field{
					Name: "r",
					Type: "Result[T, R]",
				},
				Return: []model.Field{
					{
						Type: "R",
					},
				},
				Code: []string{
					"if r.isOk {",
					"\tpanic(\"unwrap a ok value\")",
					"}",
					"v, ok := r.value.(R)",
					"if !ok {",
					"\tpanic(\"unwrap a invalid error\")",
					"}",
					"return v",
				},
			},
			{
				Name: "UnwrapErrOr",
				Receiver: model.Field{
					Name: "r",
					Type: "Result[T, R]",
				},
				Params: []model.Field{
					{
						Name: "defaultValue",
						Type: "R",
					},
				},
				Return: []model.Field{
					{
						Type: "R",
					},
				},
				Code: []string{
					"if r.isOk {",
					"\treturn defaultValue",
					"}",
					"v, ok := r.value.(R)",
					"if !ok {",
					"\treturn defaultValue",
					"}",
					"return v",
				},
			},
		},
		Functions: []model.Function{
			{
				Name: "Ok[T, R any]",
				Params: []model.Field{
					{
						Name: "value",
						Type: "T",
					},
				},
				Return: []model.Field{
					{
						Type: "Result[T, R]",
					},
				},
				Code: []string{
					"return Result[T, R]{",
					"\tvalue: value,",
					"\tisOk: true,",
					"}",
				},
			},
			{
				Name: "Err[T, R any]",
				Params: []model.Field{
					{
						Name: "err",
						Type: "R",
					},
				},
				Return: []model.Field{
					{
						Type: "Result[T, R]",
					},
				},
				Code: []string{
					"return Result[T, R]{",
					"\tvalue: err,",
					"\tisOk: false,",
					"}",
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
