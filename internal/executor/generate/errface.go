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

func ErrFace(root string, getters []check.Getter) error {
	dir := filepath.Join(root, "gen", "errface")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, getter := range getters {
		if err := func() error {
			for structName, structData := range getter.StructMap {
				if err := func() error {
					f, err := os.Create(filepath.Join(dir, getter.PackageName+"."+structName+".errface.go"))
					if err != nil {
						return err
					}
					defer f.Close()

					pkg := &model.Package{
						Name:       getter.PackageName,
						Imports:    getter.Imports,
						Interfaces: []model.Interface{},
					}

					for i, fieldName := range structData.FieldNames {
						fieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]
						pkg.Interfaces = append(pkg.Interfaces, model.Interface{
							Name: fmt.Sprintf("Get%s", fieldName),
							Methods: []model.Method{
								{
									Name: "Get" + fieldName,
									Return: []model.Field{
										{
											Type: structData.FieldTypes[i],
										},
									},
								},
							},
						})
					}

					value, err := generator.GenerateFile(pkg)
					if err != nil {
						return err
					}

					if _, err := f.Write(value); err != nil {
						return err
					}

					return nil
				}(); err != nil {
					return err
				}
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}
