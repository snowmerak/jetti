package generate

import (
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
	"strings"
)

func OptionalParameter(path string, opts []check.OptionalParameter) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)

	for _, opt := range opts {
		err := func() error {
			f, err := os.Create(filepath.Join(dir, strings.ToLower(opt.Name)+".optional.go"))
			if err != nil {
				return err
			}
			defer f.Close()

			pkg := &model.Package{
				Name: packageName,
				Aliases: []model.Alias{
					{
						Name: opt.Name + "Optional",
						Type: "func(*" + opt.Name + ") *" + opt.Name,
					},
				},
				Functions: []model.Function{
					{
						Name: "Apply" + opt.Name,
						Params: []model.Field{
							{
								Name: "defaultValue",
								Type: opt.Name,
							},
							{
								Name: "fn",
								Type: "..." + opt.Name + "Optional",
							},
						},
						Return: []model.Field{
							{
								Type: "*" + opt.Name,
							},
						},
						Code: []string{
							"param := &defaultValue",
							"for _, f := range fn {",
							"\tparam = f(param)",
							"}",
							"return param",
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
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
