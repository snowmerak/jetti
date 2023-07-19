package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

type Getter struct {
	PackageName string
	Imports     []model.Import
	StructName  string
	FieldName   string
	FieldType   string
}

func HasGetter(pkg *model.Package) ([]Getter, error) {
	list := make([]Getter, 0)

	for _, s := range pkg.Structs {
		if !strings.Contains(s.Doc, "jetti:getter") {
			continue
		}
		for _, f := range s.Fields {
			list = append(list, Getter{
				PackageName: pkg.Name,
				Imports:     pkg.Imports,
				StructName:  s.Name,
				FieldName:   f.Name,
				FieldType:   f.Type,
			})
		}
	}

	return list, nil
}
