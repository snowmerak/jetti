package executor

import (
	"log"
	"os"
	"os/exec"
)

func Init(projectName string) {
	if err := os.MkdirAll("lib", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("lib/worker", os.ModePerm); err != nil {
		panic(err)
	}

	if err := makeDocFile("lib/worker"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("lib/client", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("lib/client", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("lib/server", os.ModePerm); err != nil {
		panic(err)
	}

	if err := makeDocFile("lib/server"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("lib/service", os.ModePerm); err != nil {
		panic(err)
	}

	if err := makeDocFile("lib/service"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("internal", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("cmd", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("statics", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(protoFolder, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(umlFolder, os.ModePerm); err != nil {
		panic(err)
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		output, err := exec.Command("go", "mod", "init", projectName).Output()
		if err != nil {
			panic(err)
		}
		log.Println(string(output))
	}

	if _, err := os.Stat("README.md"); os.IsNotExist(err) {
		if err := os.WriteFile("README.md", []byte("# "+projectName+"\n"), os.ModePerm); err != nil {
			panic(err)
		}
	}

	//switch output, err := exec.Command("go", "get", "github.com/goccy/go-json").Output(); err.(type) {
	//case nil:
	//	log.Println(string(output))
	//default:
	//	panic(err)
	//}
	//
	//switch output, err := exec.Command("go", "get", "github.com/goccy/go-yaml").Output(); err.(type) {
	//case nil:
	//	log.Println(string(output))
	//default:
	//	panic(err)
	//}
}
