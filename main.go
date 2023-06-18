package main

import (
	"github.com/sabahtalateh/toolboxgen/internal/discovery"
)

var (
	root = "/Users/kravtsov777/Code/go/src/github.com/sabahtalateh/ex88"
)

func main() {
	tools, err := discovery.Discover(root)
	if err != nil {
		panic(err)
	}
	println(tools)
}
