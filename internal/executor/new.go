package executor

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/internal/executor/generate"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const commandScaffold = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

const generateScaffold = `package main

//go:generate jetti generate

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

		for _, subDir := range subDirectories {
			if err := generate.MakeDocGo(filepath.Join(root, subDir)); err != nil {
				return err
			}
		}

		if err := generate.MakeReadme(root, moduleName); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(root, "generate.go"))
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			f.Close()
		}(f)
		if _, err := f.WriteString(generateScaffold); err != nil {
			return err
		}

		gitIgnores := []string{".jetti-cache"}
		f, err = os.Create(filepath.Join(root, ".gitignore"))
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			f.Close()
		}(f)
		if _, err := f.WriteString(strings.Join(gitIgnores, "\n")); err != nil {
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
