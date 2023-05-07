package check

import (
	"os"
	"strings"
)

func GetModuleName(root string) (string, error) {
	data, err := os.ReadFile(root + "/go.mod")
	if err != nil {
		return "", err
	}

	split := strings.SplitN(string(data), "\n", 2)
	if strings.HasPrefix(split[0], "module ") {
		return strings.TrimSpace(strings.TrimPrefix(split[0], "module ")), nil
	}

	return "", errCannotFindModule
}
