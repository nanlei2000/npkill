package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
)

type Fs struct {
	cwd string
}

type FolderInfo struct {
	AbsPath string
	Bytes   int64
	Count   int64
}

func New(cwd string) *Fs {
	return &Fs{
		cwd: cwd,
	}
}

// 找到当前目录下所有的 node_modules
func (f *Fs) FindNodeModulesInner(cwd string, folderInfo *[]FolderInfo) {
	isHidden, _ := IsHiddenFile(cwd)
	if isHidden {
		return
	}
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		return
	}

	var needFurther []string
	for _, file := range files {
		isHidden, _ := IsHiddenFile(file.Name())
		if !file.IsDir() || isHidden {
			continue
		}
		absPath := path.Join(cwd, file.Name())
		if file.Name() == "node_modules" {
			size, count, err := DirSize(absPath)
			if err != nil {
				continue
			}
			*folderInfo = append(*folderInfo, FolderInfo{
				AbsPath: absPath,
				Bytes:   size,
				Count:   count,
			})
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
			f.FindNodeModulesInner(v, folderInfo)
			<-guard
		}()
	}
	wg.Wait()
}

func (f *Fs) FindNodeModules() {
	var folderInfo []FolderInfo
	f.FindNodeModulesInner(f.cwd, &folderInfo)
	sort.Slice(folderInfo, func(i, j int) bool {
		return folderInfo[i].Bytes-folderInfo[j].Bytes > 0
	})
	if len(folderInfo) == 0 {
		fmt.Println("No node_modules found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 4, 4, ' ', 0)
	fmt.Fprintln(w, "path\tfile_count\tsize\t")
	for _, info := range folderInfo {
		fmt.Fprintf(w, "%s\t%d\t%s\t\n", info.AbsPath, info.Count, humanize.Bytes(uint64(info.Bytes)))
	}
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
