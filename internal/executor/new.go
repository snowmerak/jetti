package executor

import (
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"os"
	"os/exec"
	"path/filepath"
)

const commandScaffold = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

func New(root string, moduleName string, isCmd bool) error {
	switch isCmd {
	case true:
		if err := os.MkdirAll(filepath.Join(root, "cmd", moduleName), os.ModePerm); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(root, "cmd", moduleName, "main.go"))
		if err != nil {
			return err
		}

		if _, err := f.WriteString(commandScaffold); err != nil {
			return err
		}
	case false:
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
	}

	return nil
}
