package generate

import (
	"path/filepath"
	"strings"
)

const Suffix = "jet.go"

func MakeGeneratedFileName(dir string, elem ...string) string {
	return filepath.Join(dir, strings.Join(elem, ".")+"."+Suffix)
}
