package fp

import (
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
)

func LambdaCond(modulePath, genPath string) error {
	genPath = filepath.Join(genPath, "cond")
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(generate.MakeGeneratedFileName(genPath, "cond"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	pkg := &model.Package{
		Name: "cond",
		Imports: []model.Import{
			{
				Path: modulePath + "/result",
			},
		},
		Functions: []model.Function{
			{
				Name: "If[T, R any]",
				Params: []model.Field{
					{
						Name: "cond",
						Type: "bool",
					},
					{
						Name: "trueFn",
						Type: "func() T",
					},
					{
						Name: "falseFn",
						Type: "func() R",
					},
				},
				Return: []model.Field{
					{
						Type: "result.Result[T, R]",
					},
				},
				Code: []string{
					"if cond {",
					"return result.Ok[T, R](trueFn())",
					"}",
					"return result.Err[T, R](falseFn())",
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

func LambdaWhen(modulePath, genPath string) error {
	genPath = filepath.Join(genPath, "when")
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(generate.MakeGeneratedFileName(genPath, "when"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	pkg := &model.Package{
		Name: "when",
		Imports: []model.Import{
			{
				Path: modulePath + "/option",
			},
		},
		Structs: []model.Struct{
			{
				Name: "Condition[T, R any]",
				Fields: []model.Field{
					{
						Name: "criteria",
						Type: "func (T) bool",
					},
					{
						Name: "fn",
						Type: "func (T) R",
					},
				},
			},
		},
		Functions: []model.Function{
			{
				Name: "Cond[T, R any]",
				Params: []model.Field{
					{
						Name: "criteria",
						Type: "func (T) bool",
					},
					{
						Name: "fn",
						Type: "func (T) R",
					},
				},
				Return: []model.Field{
					{
						Type: "Condition[T, R]",
					},
				},
				Code: []string{
					"return Condition[T, R]{",
					"\tcriteria: criteria,",
					"\tfn: fn,",
					"}",
				},
			},
			{
				Name: "When[T, R any]",
				Params: []model.Field{
					{
						Name: "criteria",
						Type: "T",
					},
					{
						Name: "cond",
						Type: "...Condition[T, R]",
					},
				},
				Return: []model.Field{
					{
						Type: "option.Option[R]",
					},
				},
				Code: []string{
					"for _, c := range cond {",
					"\tif c.criteria(criteria) {",
					"\t\treturn option.Some[R](c.fn(criteria))",
					"\t}",
					"}",
					"return option.None[R]()",
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
