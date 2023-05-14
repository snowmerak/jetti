package check

import (
	"github.com/snowmerak/jetti/v2/lib/model"
	"path/filepath"
)

type DependencyLink struct {
	To   string
	From []string
}

func GetImports(root string, moduleName string, path string, pkg *model.Package) (*DependencyLink, error) {
	link := new(DependencyLink)
	packagePath := filepath.Join(moduleName, filepath.ToSlash(path)[len(filepath.ToSlash(root))+1:])

	link.To = packagePath
	link.From = make([]string, len(pkg.Imports))

	for i, imp := range pkg.Imports {
		link.From[i] = imp.Path
	}

	return link, nil
}
