package executor

import (
	"os"
	"os/exec"
	"path/filepath"
)

func CmdBuild(name string, args ...string) {
	arr := append([]string{"build", "-o", filepath.Join("bin", name)}, args...)
	arr = append(arr, "./cmd/"+name+"/.")
	cmd := exec.Command("go", arr...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
