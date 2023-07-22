package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

type Getter struct {
	PackageName string
	Imports     []model.Import
	StructMap   map[string]GetterStruct
}

type GetterStruct struct {
	FieldNames []string
	FieldTypes []string
}

func HasGetter(pkg *model.Package) (Getter, error) {
	list := Getter{
		PackageName: pkg.Name,
		Imports:     pkg.Imports,
		StructMap:   map[string]GetterStruct{},
	}

	for _, s := range pkg.Structs {
		if !strings.Contains(s.Doc, "jetti:getter") {
			continue
		}
		gs, ok := list.StructMap[s.Name]
		if !ok {
			gs = GetterStruct{}
		}
		for _, f := range s.Fields {
			gs.FieldNames = append(gs.FieldNames, f.Name)
			gs.FieldTypes = append(gs.FieldTypes, f.Type)
		}
		list.StructMap[s.Name] = gs
	}

	return list, nil
}
