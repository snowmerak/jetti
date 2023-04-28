package executor

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Dependency struct {
	Import string
	Type   string
}

func getDependencies(dataList []string) map[string]*Dependency {
	dependencies := make(map[string]*Dependency)
	for _, data := range dataList {
		split := strings.Split(data, " ")
		if len(split) != 3 {
			return nil
		}

		typ := ""
		if strings.HasSuffix(split[0], "]") {
			sp := strings.SplitAfterN(split[0], "[", 2)
			typ = strings.TrimSuffix(sp[1], "]")
		}

		if typ == "" {
			return nil
		}

		dep := getDependency(typ)

		if dep == nil {
			continue
		}
		dependencies[data] = dep
	}

	return dependencies
}

func getDependency(data string) *Dependency {
	split := strings.Split(data, "/")

	lastPackageNames := strings.Split(split[len(split)-1], ".")
	lastPackageName := lastPackageNames[0]
	typeName := split[len(split)-1]
	split[len(split)-1] = lastPackageName

	dep := Dependency{
		Import: strings.Join(split, "/"),
		Type:   typeName,
	}

	return &dep
}

func makeSubPath(sub, path string) string {
	folder := path
	if runtime.GOOS == "windows" {
		folder = strings.ReplaceAll(folder, "/", "\\")
	}
	return filepath.Join(sub, folder)
}

func goGet(dep string) error {
	switch _, err := exec.Command("go", "get", "-u", dep).Output(); err.(type) {
	case nil:
	default:
		return err
	}
	return nil
}
