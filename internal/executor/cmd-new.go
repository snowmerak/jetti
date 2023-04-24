package executor

import "os"

func CmdNew(fileName string) {
	if err := os.MkdirAll("cmd/"+fileName, os.ModePerm); err != nil {
		panic(err)
	}

	if _, err := os.Stat("cmd/" + fileName + "/main.go"); !os.IsNotExist(err) {
		panic("file already exists")
	}

	f, err := os.Create("cmd/" + fileName + "/main.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.Write([]byte("package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n")); err != nil {
		panic(err)
	}
}
