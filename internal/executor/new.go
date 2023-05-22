package executor

import (
	"fmt"
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

const protoScaffold = `syntax = "proto3";

package %s;

option go_package = "%s";

`

const (
	NewKindModule = iota
	NewKindCmd
	NewKindProto
)

func New(root string, moduleName string, kind int) error {
	switch kind {
	case NewKindCmd:
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
	case NewKindModule:
		cmd := exec.Command("go", "mod", "init", moduleName)
		cmd.Dir = root
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}

		subDirs := []string{"lib", "cmd", "internal", "model"}
		for _, subDir := range subDirs {
			if err := generate.MakeDocGo(filepath.Join(root, subDir)); err != nil {
				return err
			}
		}

		if err := generate.MakeReadme(root, moduleName); err != nil {
			return err
		}
	case NewKindProto:
		path := filepath.Join(filepath.ToSlash(root), filepath.ToSlash(moduleName))
		dir := filepath.Dir(path)
		base := filepath.Base(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		f, err := os.Create(path + ".proto")
		if err != nil {
			return err
		}

		if _, err := f.WriteString(fmt.Sprintf(protoScaffold, base, moduleName)); err != nil {
			return err
		}
	}

	return nil
}
