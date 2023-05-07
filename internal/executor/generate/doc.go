package generate

import (
	"github.com/snowmerak/jetti/lib/generator"
	"github.com/snowmerak/jetti/lib/model"
	"os"
	"path/filepath"
)

func MakeDocGo(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	pkg := &model.Package{
		Name: filepath.Base(path),
	}

	data, err := generator.GenerateFile(pkg)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, "doc.go"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func MakeReadme(root string, moduleName string) error {
	f, err := os.Create(filepath.Join(root, "README.md"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString("# " + moduleName + "\n"); err != nil {
		return err
	}

	return nil
}
