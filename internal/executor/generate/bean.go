package generate

import (
	"fmt"
	"github.com/snowmerak/jetti/internal/executor/check"
	"github.com/snowmerak/jetti/lib/generator"
	"github.com/snowmerak/jetti/lib/model"
	"os"
	"path/filepath"
	"strings"
)

func Bean(path string, beans []check.Bean) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)

	for _, bean := range beans {
		for _, alias := range bean.Aliases {
			alias = strings.ToUpper(alias[:1]) + alias[1:]
			filePath := filepath.Join(dir, strings.ToLower(alias)+".context.go")
			pkg := &model.Package{
				Name: packageName,
				Imports: []model.Import{
					{
						Path: "context",
					},
				},
				Aliases: []model.Alias{
					{
						Name: alias + "ContextKey",
						Type: "string",
					},
				},
				Functions: []model.Function{
					{
						Name: "Push" + alias,
						Params: []model.Field{
							{
								Name: "ctx",
								Type: "context.Context",
							},
							{
								Name: "v",
								Type: "*" + bean.StructName,
							},
						},
						Return: []model.Field{
							{
								Type: "context.Context",
							},
						},
						Code: []string{
							fmt.Sprintf("return context.WithValue(ctx, %s(\"%s\"), v)", alias+"ContextKey", alias),
						},
					},
					{
						Name: "Get" + alias,
						Params: []model.Field{
							{
								Name: "ctx",
								Type: "context.Context",
							},
						},
						Return: []model.Field{
							{
								Type: "*" + bean.StructName,
							},
							{
								Type: "bool",
							},
						},
						Code: []string{
							fmt.Sprintf("v, ok := ctx.Value(%s(\"%s\")).(*%s)", alias+"ContextKey", alias, bean.StructName),
							"return v, ok",
						},
					},
				},
			}

			data, err := generator.GenerateFile(pkg)
			if err != nil {
				return err
			}

			if err := func() error {
				f, err := os.Create(filePath)
				if err != nil {
					return err
				}
				defer func() {
					_ = f.Close()
				}()

				if _, err := f.Write(data); err != nil {
					return err
				}

				return nil
			}(); err != nil {
				return err
			}
		}
	}

	return nil
}
