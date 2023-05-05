package executor

import (
	"os"
	"path/filepath"
	"strings"
)

func ServerTcp(path string) {
	dep := getDependency(path)
	if dep == nil {
		return
	}

	lowerName := strings.ToLower(dep.Type)
	folder := makeSubPath(serverFolder, lowerName)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	buffer := strings.Builder{}
	buffer.WriteString("package " + lowerName + "\n\n")
	buffer.WriteString("import \"net\"\n\n")
	buffer.WriteString("type " + dep.Type + " struct {\n")
	buffer.WriteString("\tlis " + "net.Listener\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func New" + dep.Type + "(addr string) *" + dep.Type + " {\n")
	buffer.WriteString("\tlis, err := net.Listen(\"tcp\", addr)\n")
	buffer.WriteString("\tif err != nil {\n\t\treturn nil\n\t}\n")
	buffer.WriteString("\treturn &" + dep.Type + "{\n\t\tlis: lis,\n\t}\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") Handler(conn net.Conn) {\n")
	buffer.WriteString("\tdefer conn.Close()\n")
	buffer.WriteString("\tbuff := make([]byte, 1024)\n")
	buffer.WriteString("\tfor {\n")
	buffer.WriteString("\t\t_, err := conn.Read(buff)\n")
	buffer.WriteString("\t\tif err != nil {\n\t\t\treturn\n\t\t}\n")
	buffer.WriteString("\t\t// TODO: handle message\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") Serve() {\n")
	buffer.WriteString("\tfor {\n")
	buffer.WriteString("\t\tconn, err := s.lis.Accept()\n")
	buffer.WriteString("\t\tif err != nil {\n\t\t\treturn\n\t\t}\n")
	buffer.WriteString("\t\tgo s.Handler(conn)\n")
	buffer.WriteString("\t}\n")
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
}
