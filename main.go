package main

import (
	"github.com/sabahtalateh/toolboxgen/internal/discovery"
)

var (
	root = "/Users/kravtsov777/Code/go/src/github.com/sabahtalateh/ex88"
)

func main() {
	conf, err := discovery.ConfigBuilder().Root(root).Build()
	if err != nil {
		panic(err)
	}

	if err = discovery.Discover(conf); err != nil {
		panic(err)
	}
}
