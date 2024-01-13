package main

import (
	"fmt"
	"github.com/snowmerak/jetti/v2/lib/tools"
)

func main() {
	if err := tools.CloneIfNotExists(); err != nil {
		panic(err)
	}

	registries, err := tools.GetRegistries()
	if err != nil {
		panic(err)
	}

	for _, registry := range registries {
		reg, err := tools.GetRegistryInfo(registry)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", reg)
	}

}
