//go:build !windows
// +build !windows

package file

func IsHiddenFile(filename string) (bool, error) {
	return filename[0:1] == ".", nil
}
