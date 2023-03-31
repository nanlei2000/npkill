package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"text/tabwriter"

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
func (f *Fs) FindNodeModulesInner(cwd string, w *tabwriter.Writer) {
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
			size, count, err := DirSize(absPath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t\n", absPath, humanize.Bytes(uint64(size)), count)
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
			f.FindNodeModulesInner(v, w)
			<-guard
		}()
	}
	wg.Wait()
}

func (f *Fs) FindNodeModules() {
	w := tabwriter.NewWriter(os.Stdout, 1, 4, 4, ' ', 0)
	fmt.Fprintln(w, "PATH\tSIZE\tCOUNT\t")
	f.FindNodeModulesInner(f.cwd, w)
	w.Flush()
}

func DirSize(path string) (int64, int64, error) {
	var size int64
	var count int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		} else {
			count++
		}
		return err
	})
	return size, count, err
}
