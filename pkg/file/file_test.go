package file

import "testing"

func TestFindNodeModules(t *testing.T) {
	f := New("/Users/xx/Project")
	f.FindNodeModules()
}
