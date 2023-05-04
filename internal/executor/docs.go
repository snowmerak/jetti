package executor

import (
	"os"
	"path/filepath"
)

var umlFolder = filepath.Join("docs", "uml")

func makeDocFile(path string) error {
	f, err := os.Create(filepath.Join(path, "docs.go"))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.WriteString("package " + filepath.Base(path) + "\n"); err != nil {
		return err
	}

	return nil
}
