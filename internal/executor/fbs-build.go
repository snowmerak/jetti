package executor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func FbsBuild() {
	if _, err := os.Stat(fbsFolder); os.IsNotExist(err) {
		panic(err)
	}

	if err := filepath.Walk(fbsFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".fbs" {
			return nil
		}

		slashPath := filepath.ToSlash(path)
		slashPath = slashPath[len(fbsFolder)+1:]
		genPath := filepath.Join(generated, "fbs", filepath.Dir(filepath.Dir(slashPath)))

		fmt.Println(genPath)

		cmd := exec.Command("flatc", "--go", "-o", genPath, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	switch output, err := exec.Command("go", "get", "-u", "github.com/google/flatbuffers").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}
}
