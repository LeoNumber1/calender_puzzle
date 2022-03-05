package shape

import (
	"testing"
)

func TestEqual(t *testing.T) {
	a := NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	b := NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	if a.Equal(b) {
		t.Log("pass")
		return
	}
	t.Error("not pass")
}
