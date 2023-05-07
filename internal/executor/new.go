package executor

import (
	"github.com/snowmerak/jetti/internal/executor/generate"
	"os"
	"os/exec"
	"path/filepath"
)

func New(root string, moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	subDirs := []string{"lib", "cmd", "internal"}
	for _, subDir := range subDirs {
		if err := generate.MakeDocGo(filepath.Join(root, subDir)); err != nil {
			return err
		}
	}

	if err := generate.MakeReadme(root, moduleName); err != nil {
		return err
	}

	return nil
}
