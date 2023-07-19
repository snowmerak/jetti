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

func Getter(path string, getters []check.Getter) error {
	dir := filepath.Dir(path)

	makeFilename := func(getter check.Getter) string {
		return filepath.Join(dir, fmt.Sprintf("%s.getter.go", getter.StructName))
	}

	fileMap := map[string]*model.Package{}

	for _, getter := range getters {
		filename := makeFilename(getter)

		pkg, ok := fileMap[filename]
		if !ok {
			pkg = &model.Package{
				Name:    getter.PackageName,
				Imports: getter.Imports,
			}
			fileMap[filename] = pkg
		}

		fieldName := strings.ToUpper(getter.FieldName[:1]) + getter.FieldName[1:]
		pkg.Methods = append(pkg.Methods, model.Method{
			Name: fmt.Sprintf("Get%s", fieldName),
			Receiver: model.Field{
				Name: getter.StructName[:1],
				Type: getter.StructName,
			},
			Return: []model.Field{
				{
					Type: getter.FieldType,
				},
			},
			Code: []string{
				fmt.Sprintf("return $RECEIVER$.%s", getter.FieldName),
			},
		})
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
