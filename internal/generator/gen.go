package generator

import (
	"bytes"
	"github.com/snowmerak/go-bean/internal/model"
	"strings"
)

func Generate(models ...model.Model) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	buffer.Write([]byte("package bean\n\n"))
	buffer.Write([]byte("type Bean struct {\n"))
	for _, m := range models {
		for _, structName := range m.StructNames {
			buffer.Write([]byte("\t" + strings.ToLower(structName) + " *" + m.PackageName + "." + structName + "\n"))
		}
	}
	buffer.Write([]byte("}\n\n"))

	buffer.Write([]byte("type Builder struct {\n"))
	buffer.Write([]byte("\tbean *Bean\n"))
	buffer.Write([]byte("}\n\n"))

	buffer.Write([]byte("func New() *Builder {\n"))
	buffer.Write([]byte("\treturn &Builder{}\n"))
	buffer.Write([]byte("}\n\n"))

	buffer.Write([]byte("func (b *Builder) Build() *Bean {\n"))
	buffer.Write([]byte("\treturn b.bean\n"))
	buffer.Write([]byte("}\n\n"))

	for _, m := range models {
		for _, structName := range m.StructNames {
			buffer.Write([]byte("func (b *Builder) " + structName + "(" + strings.ToLower(structName) + " *" + m.PackageName + "." + structName + ") *Builder {\n"))
			buffer.Write([]byte("\tb.bean." + strings.ToLower(structName) + " = " + strings.ToLower(structName) + "\n"))
			buffer.Write([]byte("\treturn b\n"))
			buffer.Write([]byte("}\n\n"))
		}
	}

	for _, m := range models {
		for _, structName := range m.StructNames {
			buffer.Write([]byte("func (b *Bean) " + structName + "() *" + m.PackageName + "." + structName + " {\n"))
			buffer.Write([]byte("\treturn b." + strings.ToLower(structName) + "\n"))
			buffer.Write([]byte("}\n\n"))
		}
	}

	return buffer.Bytes(), nil
}
