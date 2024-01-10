package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Run(root string, name string, args ...string) error {
	cmdArgs := []string{"go", "run", root + "/cmd/" + name + "/."}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if _, err := os.ReadDir(filepath.Join(root, "clib")); !os.IsNotExist(err) {
		cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, fmt.Sprintf("CGO_LDFLAGS=-L%s/clib/%s-%s/lib", root, runtime.GOOS, runtime.GOARCH))
		cmd.Env = append(cmd.Env, fmt.Sprintf("CGO_CFLAGS=-I%s/clib/%s-%s/include", root, runtime.GOOS, runtime.GOARCH))
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
