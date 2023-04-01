package file

import "testing"

func TestFindNodeModules(t *testing.T) {
	f := New("/Users/lielienan/Project")
	f.FindNodeModules()
}
