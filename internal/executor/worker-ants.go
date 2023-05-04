package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func WorkerAnts(path string, params ...string) {
	dep := getDependency(path)
	if dep == nil {
		return
	}

	lowerName := strings.ToLower(dep.Type)

	folder := makeSubPath(workerFolder, lowerName)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	arguments := strings.Join(params, ", ")
	argumentsWithName := strings.Builder{}
	for i, param := range params {
		argumentsWithName.WriteString(fmt.Sprintf("p%d %s", i, param))
		if i < len(params)-1 {
			argumentsWithName.WriteString(", ")
		}
	}
	parameters := strings.Builder{}
	for i := range params {
		parameters.WriteString(fmt.Sprintf("p%d", i))
		if i < len(params)-1 {
			parameters.WriteString(", ")
		}
	}
	parametersWithArgs := strings.Builder{}
	for i := range params {
		parametersWithArgs.WriteString(fmt.Sprintf("args.p%d", i))
		if i < len(params)-1 {
			parametersWithArgs.WriteString(", ")
		}
	}

	f, err := os.Create(filepath.Join(folder, "worker_queue.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsWorkerQueue, lowerName, arguments, arguments)); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "worker_stack.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsWorkerStack, lowerName)); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "worker_loop_queue.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsWorkerLoopQueue, lowerName)); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "spin_lock.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsSpinLock, lowerName)); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "options.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsOptions, lowerName)); err != nil {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}

	f, err = os.Create(filepath.Join(folder, "parameter.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	fields := strings.Builder{}
	for i, param := range params {
		fields.WriteString(fmt.Sprintf("\t%s %s\n", "p"+strconv.FormatInt(int64(i), 16), param))
	}
	if _, err := f.WriteString(fmt.Sprintf(antsParameter, lowerName, fields.String())); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "worker_func.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsWorkerFunc, lowerName, parametersWithArgs.String(), arguments, argumentsWithName.String(), parameters.String())); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "pool_func.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsPoolFunc, lowerName, arguments, arguments, argumentsWithName.String(), parameters.String())); err != nil {
		panic(err)
	}

	f, err = os.Create(filepath.Join(folder, "ants.go"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)

	if _, err := f.WriteString(fmt.Sprintf(antsMain, lowerName)); err != nil {
		panic(err)
	}
}
