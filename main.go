package main

import (
	"github.com/nanlei2000/npkill/pkg/file"
)

func main() {
	f := file.New("/Users/bytedance/Project")
	f.FindNodeModules()
}
