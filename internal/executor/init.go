package executor

import (
	"log"
	"os"
	"os/exec"
)

func Init(projectName string) {
	if err := os.MkdirAll("src", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("internal", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("cmd", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("proto", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("configs", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("uml", os.ModePerm); err != nil {
		panic(err)
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		output, err := exec.Command("go", "mod", "init", projectName).Output()
		if err != nil {
			panic(err)
		}
		log.Println(string(output))
	}

	switch output, err := exec.Command("go", "get", "github.com/goccy/go-json").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}

	switch output, err := exec.Command("go", "get", "github.com/goccy/go-yaml").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}

	switch output, err := exec.Command("go", "get", "google.golang.org/grpc").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}

	switch output, err := exec.Command("go", "get", "google.golang.org/protobuf").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}
}
