package generate

import (
	"path/filepath"
	"strings"
)

const Suffix = "jet"

func MakeGeneratedFileName(dir string, elem ...string) string {
	return filepath.Join(dir, strings.Join(elem, ".")+"."+Suffix+".go")
}

func MakeGeneratedTestFileName(dir string, elem ...string) string {
	return filepath.Join(dir, strings.Join(elem, ".")+"."+Suffix+"_test.go")
}
