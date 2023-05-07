package generate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func BuildFlatbuffers(root, path string) error {
	genPath := filepath.Join(root, "gen", "flatbuffers")
	slashPath := filepath.ToSlash(path)
	slashRoot := filepath.ToSlash(root)
	slashPath = strings.TrimPrefix(slashPath, slashRoot)
	slashPath = filepath.Dir(slashPath)
	genPath = filepath.Join(genPath, slashPath)

	cmd := exec.Command("flatc", "--go", "-o", genPath, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := goGet(googleFlatbuffers); err != nil {
		return err
	}

	return nil
}
