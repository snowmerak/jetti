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

	f, err := os.Create(filepath.Join(dir, "errface.go"))
	if err != nil {
		return err
	}

	pkg := &model.Package{
		Name: "errface",
	}

	importsSet := map[string]struct{}{}
	for _, getter := range getters {
		for _, imp := range getter.Imports {
			if _, ok := importsSet[imp.Path]; ok {
				continue
			}
			importsSet[imp.Path] = struct{}{}
			pkg.Imports = append(pkg.Imports, imp)
		}

		fieldName := strings.ToUpper(getter.FieldName[:1]) + getter.FieldName[1:]
		typeName := strings.TrimPrefix(getter.FieldType, "*")
		typeName = strings.TrimPrefix(typeName, "[]")
		typeName = strings.TrimPrefix(typeName, "...")
		typeName = strings.ToUpper(typeName[:1]) + typeName[1:]
		switch {
		case strings.HasPrefix(getter.FieldType, "*"):
			typeName = typeName + "Pointer"
		case strings.HasPrefix(getter.FieldType, "[]"):
			typeName = typeName + "Slice"
		case strings.HasPrefix(getter.FieldType, "..."):
			typeName = typeName + "Slice"
		}
		pkg.Interfaces = append(pkg.Interfaces, model.Interface{
			Name: fmt.Sprintf("%s%s", fieldName, typeName),
			Methods: []model.Method{
				{
					Name: fmt.Sprintf("Get%s", fieldName),
					Return: []model.Field{
						{
							Type: getter.FieldType,
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
}
