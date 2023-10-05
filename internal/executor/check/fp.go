package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"strings"
)

func HasFp(pkg *model.Package) (bool, error) {
	directive := "jetti:fp"

	for _, fn := range pkg.Functions {
		if strings.Contains(fn.Doc, directive) {
			return true, nil
		}
	}

	return false, nil
}
