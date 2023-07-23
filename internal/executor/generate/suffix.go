package generate

import "strings"

const Suffix = "jet.go"

func MakeGeneratedFileName(elem ...string) string {
	return strings.Join(elem, ".") + "." + Suffix
}
