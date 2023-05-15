package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

type OptionalParameter struct {
	Name string
}

func HasOptionalParameter(pkg *model.Package) ([]OptionalParameter, error) {
	params := []OptionalParameter(nil)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:parameter") {
			param := OptionalParameter{
				Name: st.Name,
			}
			params = append(params, param)
		}
	}

	for _, it := range pkg.Interfaces {
		if strings.Contains(it.Doc, "jetti:parameter") {
			param := OptionalParameter{
				Name: it.Name,
			}
			params = append(params, param)
		}
	}

	return params, nil
}
