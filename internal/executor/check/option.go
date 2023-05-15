package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

func HasOptional(pkg *model.Package) ([]string, error) {
	optional := []string(nil)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:optional") {
			optional = append(optional, st.Name)
		}
	}

	for _, it := range pkg.Interfaces {
		if strings.Contains(it.Doc, "jetti:optional") {
			optional = append(optional, it.Name)
		}
	}

	for _, ali := range pkg.Aliases {
		if strings.Contains(ali.Doc, "jetti:optional") {
			optional = append(optional, ali.Name)
		}
	}

	return optional, nil
}
