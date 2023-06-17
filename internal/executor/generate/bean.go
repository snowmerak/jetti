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

const BeanDirective = "bean"

func BeanContainer(root string) error {
	genPath := filepath.Join(root, "gen", "bean")
	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(genPath, "bean.container.go"))
	if err != nil {
		return err
	}

	pkg := &model.Package{
		Name: "bean",
		Imports: []model.Import{
			{
				Path: "sync",
			},
		},
		Structs: []model.Struct{
			{
				Name: "Container",
				Fields: []model.Field{
					{
						Name: "beans",
						Type: "map[any]any",
					},
					{
						Name: "lock",
						Type: "sync.RWMutex",
					},
				},
				Methods: []model.Method{
					{
						Name: "Get",
						Params: []model.Field{
							{
								Name: "key",
								Type: "any",
							},
						},
						Return: []model.Field{
							{
								Name: "value",
								Type: "any",
							},
							{
								Name: "ok",
								Type: "bool",
							},
						},
						Code: []string{
							"$RECEIVER$.lock.RLock()",
							"value, ok = $RECEIVER$.beans[key]",
							"$RECEIVER$.lock.RUnlock()",
							"return",
						},
					},
					{
						Name: "Set",
						Params: []model.Field{
							{
								Name: "key",
								Type: "any",
							},
							{
								Name: "value",
								Type: "any",
							},
						},
						Code: []string{
							"$RECEIVER$.lock.Lock()",
							"$RECEIVER$.beans[key] = value",
							"$RECEIVER$.lock.Unlock()",
						},
					},
					{
						Name: "Delete",
						Params: []model.Field{
							{
								Name: "key",
								Type: "any",
							},
						},
						Code: []string{
							"$RECEIVER$.lock.Lock()",
							"delete($RECEIVER$.beans, key)",
							"$RECEIVER$.lock.Unlock()",
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

func Bean(moduleName string, path string, beans []check.Bean) error {
	dir := filepath.Dir(path)
	packageName := filepath.Base(dir)
	beanPath := filepath.Join(moduleName, "gen", "bean")

	for _, bean := range beans {
		for _, alias := range bean.Aliases {
			alias = strings.ToUpper(alias[:1]) + alias[1:]
			filePath := filepath.Join(dir, strings.ToLower(alias)+".bean.go")
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
						Path: beanPath,
					},
					{
						Path: "errors",
					},
				},
				Aliases: []model.Alias{
					{
						Name: alias + "BeanKey",
						Type: "string",
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
								Name: "beanContainer",
								Type: "*bean.Container",
							},
							{
								Name: "value",
								Type: typ,
							},
						},
						Code: []string{
							fmt.Sprintf("beanContainer.Set(%sBeanKey(\"%skey\"), value)", alias, alias),
						},
					},
					{
						Name: "Get" + alias,
						Params: []model.Field{
							{
								Name: "beanContainer",
								Type: "*bean.Container",
							},
						},
						Return: []model.Field{
							{
								Name: "value",
								Type: typ,
							},
							{
								Name: "err",
								Type: "error",
							},
						},
						Code: []string{
							fmt.Sprintf("maybe, ok := beanContainer.Get(%sBeanKey(\"%skey\"))", alias, alias),
							"if !ok {",
							fmt.Sprintf("return nil, err%sNotFound", alias),
							"}",
							fmt.Sprintf("value, ok = maybe.(%s)", typ),
							"if !ok {",
							fmt.Sprintf("return nil, err%sNotFound", alias),
							"}",
							"return value, nil",
						},
					},
					{
						Name: "IsErr" + alias + "NotFound",
						Params: []model.Field{
							{
								Name: "err",
								Type: "error",
							},
						},
						Return: []model.Field{
							{
								Name: "ok",
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
