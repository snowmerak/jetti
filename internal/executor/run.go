package executor

import (
	"os"
	"os/exec"
)

func Run(root string, name string, args ...string) error {
	cmdArgs := []string{"go", "run", root + "/cmd/" + name + "/."}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
