package fiber

import (
	"fmt"
	"testing"
)

func TestFiber(t *testing.T) {
	fi := New[string](func(suspend SuspendFunc[string]) {
		fmt.Println("started")
		v := suspend("first suspend")
		if v != "second resume" {
			fmt.Println(v)
			t.Error("Resume error")
		}
	})

	first := fi.Start()
	if first != "first suspend" {
		fmt.Println(first)
		t.Error("Suspend error")
	}
	end := fi.Resume("second resume")
	if end != "" {
		fmt.Println(end)
		t.Error("Terminated error")
	}
}
