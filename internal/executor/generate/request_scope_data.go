package generate

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/generator"
	"github.com/snowmerak/jetti/v2/lib/model"
	"os"
	"path/filepath"
	"strings"
)

const RequestDirective = "request"

func RequestScopeData(path string, beans []check.Bean) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)

	for _, bean := range beans {
		for _, alias := range bean.Aliases {
			alias = strings.ToUpper(alias[:1]) + alias[1:]
			filePath := filepath.Join(MakeGeneratedFileName(dir, strings.ToLower(alias), "context"))
			typ := bean.Name
			switch bean.Type {
			case check.TypeStruct:
				fallthrough
			case check.TypeAlias:
				typ = "*" + bean.Name
			}
			pkg := &model.Package{
				Name: packageName,
				Imports: []model.Import{
					{
						Path: "context",
					},
					{
						Path: "errors",
					},
				},
				Aliases: []model.Alias{
					{
						Name: alias + "ContextKey",
						Type: "struct{}",
					},
				},
				GlobalVariables: []model.GlobalVariable{
					{
						Name:  "err" + alias + "NotFound",
						Type:  "error",
						Value: fmt.Sprintf("errors.New(\"%s not found\")", strings.ToLower(alias)),
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
								Type: typ,
							},
						},
						Return: []model.Field{
							{
								Type: "context.Context",
							},
						},
						Code: []string{
							fmt.Sprintf("return context.WithValue(ctx, %s{}, v)", alias+"ContextKey"),
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
								Type: typ,
							},
							{
								Type: "bool",
							},
						},
						Code: []string{
							fmt.Sprintf("v, ok := ctx.Value(%s{}).(%s)", alias+"ContextKey", typ),
							"return v, ok",
						},
					},
					{
						Name: "Err" + alias + "NotFound",
						Return: []model.Field{
							{
								Type: "error",
							},
						},
						Code: []string{
							fmt.Sprintf("return err%sNotFound", alias),
						},
					},
					{
						Name: "Is" + alias + "NotFoundErr",
						Params: []model.Field{
							{
								Name: "err",
								Type: "error",
							},
						},
						Return: []model.Field{
							{
								Type: "bool",
							},
						},
						Code: []string{
							fmt.Sprintf("return errors.Is(err, err%sNotFound)", alias),
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
