package check

import (
	"path/filepath"
	"strings"

	"github.com/snowmerak/jetti/lib/model"
)

func HasBean(path string, pkg *model.Package) ([]string, error) {
	rs := make([]string, 0)
	path = filepath.Dir(path)

	for _, st := range pkg.Structs {
		if strings.Contains(st.Doc, "jetti:bean") {
			rs = append(rs, path+"/"+st.Name)
		}
	}

	return rs, nil
}
