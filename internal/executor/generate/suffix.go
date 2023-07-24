package generate

import (
	"path/filepath"
	"strings"
)

const Suffix = "jet.go"

func MakeGeneratedFileName(path string, elem ...string) string {
	dir := filepath.Dir(path)
	return filepath.Join(dir, strings.Join(elem, ".")+"."+Suffix)
}
