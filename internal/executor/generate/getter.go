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

func Getter(path string, getter check.Getter) error {
	dir := filepath.Dir(path)

	makeFilename := func(pkgName string, structName string) string {
		return filepath.Join(dir, fmt.Sprintf("%s.%s.getter.go", pkgName, structName))
	}

	fileMap := map[string]*model.Package{}

	for structName, structData := range getter.StructMap {
		filename := makeFilename(getter.PackageName, structName)

		pkg, ok := fileMap[filename]
		if !ok {
			pkg = &model.Package{
				Name:    getter.PackageName,
				Imports: getter.Imports,
			}
			fileMap[filename] = pkg
		}

		for i, fieldName := range structData.FieldNames {
			fieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]
			pkg.Methods = append(pkg.Methods, model.Method{
				Name: fmt.Sprintf("Get%s", fieldName),
				Receiver: model.Field{
					Name: structName[:1],
					Type: structName,
				},
				Return: []model.Field{
					{
						Type: structData.FieldTypes[i],
					},
				},
				Code: []string{
					fmt.Sprintf("return $RECEIVER$.%s", structData.FieldNames[i]),
				},
			})
		}
	}

	for filename, pkg := range fileMap {
		if err := func() error {
			f, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer f.Close()

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
}
