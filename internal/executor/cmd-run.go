package executor

import (
	"os"
	"os/exec"
)

func CmdRun(file string, args ...string) {
	execArgs := append([]string{"run", "cmd/" + file + "/main.go"}, args...)
	cmd := exec.Command("go", execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
