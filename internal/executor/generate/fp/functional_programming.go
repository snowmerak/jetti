package fp

import (
	"os"
	"path/filepath"
)

func FunctionalProgramming(moduleName, root string) error {
	genPath, err := InitFunctionalProgramming(root)
	if err != nil {
		return err
	}

	if err := MonadOption(genPath); err != nil {
		return err
	}

	if err := MonadResult(genPath); err != nil {
		return err
	}

	moduleName = moduleName + "/gen/fp"

	if err := LambdaCond(moduleName, genPath); err != nil {
		return err
	}

	if err := LambdaWhen(moduleName, genPath); err != nil {
		return err
	}

	return nil
}

func InitFunctionalProgramming(root string) (string, error) {
	functionalProgrammingPath := filepath.Join(root, "gen", "fp")
	if err := os.MkdirAll(functionalProgrammingPath, os.ModePerm); err != nil {
		return "", err
	}
	return functionalProgrammingPath, nil
}
