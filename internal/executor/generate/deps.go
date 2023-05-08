package generate

import (
	"os/exec"
)

const (
	goccyJson = "github.com/goccy/go-json"
	goccyYaml = "github.com/goccy/go-yaml"
)

const (
	googleGrpc     = "google.golang.org/grpc"
	googleProtobuf = "google.golang.org/protobuf"

	googleFlatbuffers = "github.com/google/flatbuffers"
)

func goGet(dep string) error {
	switch _, err := exec.Command("go", "get", "-u", dep).Output(); err.(type) {
	case nil:
	default:
		return err
	}
	return nil
}
