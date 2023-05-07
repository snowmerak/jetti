package check

import "errors"

var errCannotFindModule = errors.New("cannot find module name")

func IsCannotFindModuleErr(err error) bool {
	return errors.Is(err, errCannotFindModule)
}
