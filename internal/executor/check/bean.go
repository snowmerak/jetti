package check

import (
	"github.com/snowmerak/jetti/lib/model"
	"strings"
)

func HasBean(pkg *model.Package) ([]model.Struct, error) {
	rs := make([]model.Struct, 0)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:bean") {
			rs = append(rs, st)
		}
	}

	return rs, nil
}
