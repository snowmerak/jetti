package generate

import (
	"os"
	"os/exec"
	"path/filepath"
)

func BuildProtobuf(root, path string) error {
	genPath := filepath.Join(root, "gen", "grpc")

	if err := os.MkdirAll(genPath, os.ModePerm); err != nil {
		return err
	}

	cmd := exec.Command("protoc", "--go_out="+genPath, "--go_opt=paths=source_relative", "--go-grpc_out="+genPath, "--go-grpc_opt=paths=source_relative", "--proto_path="+root, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := goGet(googleProtobuf); err != nil {
		return err
	}

	if err := goGet(googleGrpc); err != nil {
		return err
	}

	return nil
}
