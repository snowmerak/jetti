package executor

import (
	"os"
	"path/filepath"
	"strings"
)

func ClientGoRedis(path string) {
	dep := getDependency(path)
	if dep == nil {
		return
	}

	lowerName := strings.ToLower(dep.Type)
	folder := makeSubPath(clientFolder, lowerName)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	buffer := strings.Builder{}
	buffer.WriteString("package " + lowerName + "\n\n")
	buffer.WriteString("import \"github.com/go-redis/redis/v8\"\n\n")
	buffer.WriteString("type " + dep.Type + " struct {\n")
	buffer.WriteString("\tclient redis.Cmdable\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func New" + dep.Type + "(client redis.Cmdable) *" + dep.Type + " {\n")
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

	if err := goGet("github.com/go-redis/redis/v8"); err != nil {
		panic(err)
	}
}
