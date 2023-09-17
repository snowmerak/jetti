package executor

import (
	"github.com/snowmerak/jetti/v2/lib/parser"
	"os"
	"path/filepath"
	"strings"
)

func Check(root string) error {
	for _, subDir := range subDirectories {
		if err := filepath.Walk(filepath.Join(root, subDir), func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			if !strings.HasSuffix(info.Name(), ".go") {
				return nil
			}

			sp := strings.Split(filepath.ToSlash(path), "/")
			for _, s := range sp {
				if strings.HasPrefix(s, ".") {
					return nil
				}
			}

			pkg, err := parser.ParseFile(path)
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
