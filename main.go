package main

import (
	"github.com/nanlei2000/npkill/pkg/file"
)

func main() {
	f := file.New("/Users/xx/Project")
	f.FindNodeModules()
}