package helpers

import (
	"fmt"
	"testing"
)

func TestWalkDir(t *testing.T) {
	files, err := WalkDir("../static")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	fmt.Println(files)
}
