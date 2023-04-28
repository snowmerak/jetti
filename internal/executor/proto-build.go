package executor

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ProtoBuild() {
	if _, err := os.Stat(protoFolder); os.IsNotExist(err) {
		panic(err)
	}

	protoFiles := []string(nil)
	if err := filepath.Walk(protoFolder, func(path string, info os.FileInfo, err error) error {
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

	if err := os.MkdirAll(filepath.Join(generated, "proto"), os.ModePerm); err != nil {
		panic(err)
	}

	cmd := exec.Command("protoc", append([]string{"--proto_path=template/proto", "--go_out=gen/proto", "--go_opt=paths=source_relative", "--go-grpc_out=gen/proto", "--go-grpc_opt=paths=source_relative"}, protoFiles...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	switch output, err := exec.Command("go", "get", "-u", "google.golang.org/grpc").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}

	switch output, err := exec.Command("go", "get", "-u", "google.golang.org/protobuf").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}
}
