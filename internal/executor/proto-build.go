package executor

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ProtoBuild() {
	if _, err := os.Stat("proto"); os.IsNotExist(err) {
		panic(err)
	}

	protoFiles := []string(nil)
	if err := filepath.Walk("proto", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".proto" {
			return nil
		}
		protoFiles = append(protoFiles, path)
		return nil
	}); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("src/model", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("src/service", os.ModePerm); err != nil {
		panic(err)
	}

	cmd := exec.Command("protoc", append([]string{"--proto_path=proto", "--go_out=src/model", "--go_opt=paths=source_relative", "--go-grpc_out=src/model", "--go-grpc_opt=paths=source_relative"}, protoFiles...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("src/server", os.ModePerm); err != nil {
		panic(err)
	}
}
