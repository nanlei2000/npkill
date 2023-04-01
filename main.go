package main

import (
	"log"
	"os"

	"github.com/nanlei2000/npkill/pkg/file"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	f := file.New(cwd)
	f.FindNodeModules()
}
