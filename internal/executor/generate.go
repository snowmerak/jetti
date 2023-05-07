package executor

import (
	"os"
	"path/filepath"
)

func Generate(root string) error {
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".go":
		case ".json":
			if err := convertJson(path); err != nil {
				return err
			}
		case ".yml":
			fallthrough
		case ".yaml":
		default:
			return nil
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
