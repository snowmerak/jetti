package executor

import (
	"os"
	"path/filepath"
	"strings"
)

func ClientNats(path string) {
	dep := getDependency(path)
	if dep == nil {
		return
	}

	folder := makeSubPath(clientFolder, dep.Import)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	lowerName := strings.ToLower(dep.Type)

	buffer := strings.Builder{}
	buffer.WriteString("package " + lowerName + "\n\n")
	buffer.WriteString("import \"github.com/nats-io/nats.go\"\n\n")
	buffer.WriteString("type " + dep.Type + " struct {\n")
	buffer.WriteString("\tclient *nats.Conn\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func New" + dep.Type + "(client *nats.Conn) *" + dep.Type + " {\n")
	buffer.WriteString("\treturn &" + dep.Type + "{\n\t\tclient: client,\n\t}\n")
	buffer.WriteString("}\n\n")

	f, err := os.Create(filepath.Join(folder, lowerName+".go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.WriteString(buffer.String()); err != nil {
		panic(err)
	}

	if err := goGet("github.com/nats-io/nats.go"); err != nil {
		panic(err)
	}
}
