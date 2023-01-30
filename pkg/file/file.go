package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/dustin/go-humanize"
)

type Fs struct {
	cwd string
}

func New(cwd string) *Fs {
	return &Fs{
		cwd: cwd,
	}
}

// 找到当前目录下所有的 node_modules
func (f *Fs) FindNodeModulesInner(cwd string) {
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		log.Fatal(err)
	}

	var needFurther []string
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		absPath := path.Join(cwd, file.Name())
		if file.Name() == "node_modules" {
			dirSize, err := DirSize(absPath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Path: %s, Size: %s \n", absPath, humanize.Bytes(uint64(dirSize)))
			needFurther = []string{}
			break
		}
		needFurther = append(needFurther, absPath)
	}

	if len(needFurther) == 0 {
		return
	}
	maxGoroutines := 3
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup
	for _, v := range needFurther {
		v := v
		guard <- struct{}{} // would block if guard channel is already filled
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.FindNodeModulesInner(v)
			<-guard
		}()
	}
	wg.Wait()
}

func (f *Fs) FindNodeModules() {
	f.FindNodeModulesInner(f.cwd)
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
