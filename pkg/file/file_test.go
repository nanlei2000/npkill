package file

import "testing"

func TestFindNodeModules(t *testing.T) {
	f := New("/Users/bytedance/Project")
	f.FindNodeModules()
}
