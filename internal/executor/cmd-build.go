package executor

import (
	"os"
	"os/exec"
)

func CmdBuild(name string) {
	cmd := exec.Command("go", "build", "-o", name, "./cmd/"+name+"/.")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
