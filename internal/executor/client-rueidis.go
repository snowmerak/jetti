package executor

import (
	"os"
	"path/filepath"
	"strings"
)

func ClientRueidis(path string) {
	dep := getDependency(path)
	if dep == nil {
		return
	}

	folder := makeSubPath(clientFolder, dep.Import)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	lowerName := strings.ToLower(dep.Type)
	f, err := os.Create(filepath.Join(folder, lowerName+".go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	buffer := strings.Builder{}
	buffer.WriteString("package " + lowerName + "\n\n")
	buffer.WriteString("import \"github.com/rueian/rueidis\"\n\n")
	buffer.WriteString("type " + dep.Type + " struct {\n")
	buffer.WriteString("\tclient rueidis.Client\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func New" + dep.Type + "(client rueidis.Client) *" + dep.Type + " {\n")
	buffer.WriteString("\treturn &" + dep.Type + "{\n\t\tclient: client,\n\t}\n")
	buffer.WriteString("}\n\n")

	if _, err := f.WriteString(buffer.String()); err != nil {
		panic(err)
	}

	if err := goGet("github.com/rueian/rueidis"); err != nil {
		panic(err)
	}
}
