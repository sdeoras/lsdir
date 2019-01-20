package lsdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestLister_List(t *testing.T) {
	id := uuid.New().String()

	if err := os.MkdirAll(filepath.Join("/tmp", id, "a", "b", "c", "d"), 0755); err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.RemoveAll(filepath.Join("/tmp", id))
	}()

	b := []byte("this is a test")

	expectedPaths := make(map[string]struct{})

	paths := []string{
		filepath.Join("/tmp", id),
		filepath.Join("/tmp", id, "a"),
		filepath.Join("/tmp", id, "a", "b"),
		filepath.Join("/tmp", id, "a", "b", "c"),
		filepath.Join("/tmp", id, "a", "b", "c", "d"),
	}

	for _, path := range paths {
		for _, file := range []string{"x.txt", "y.txt", "z.txt"} {
			fileName := filepath.Join(path, file)
			if err := ioutil.WriteFile(fileName, b, 0644); err != nil {
				t.Fatal(err)
			}
			expectedPaths[fileName] = struct{}{}
		}
	}

	lister := NewLister(true, "*")

	files, err := lister.List(filepath.Join("/tmp", id))
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 15 {
		t.Fatal("expected 15 files, found:", len(files))
	}

	for _, file := range files {
		if _, ok := expectedPaths[file]; !ok {
			t.Fatal(file, "not found in list")
		}
	}

	// test with a pattern
	lister = NewLister(true, "x*")

	files, err = lister.List(filepath.Join("/tmp", id))
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 5 {
		t.Fatal("expected 5 files, found:", len(files))
	}

	for _, file := range files {
		if _, ok := expectedPaths[file]; !ok {
			t.Fatal(file, "not found in list")
		}
	}

	// test with recursion off
	lister = NewLister(false, "*")

	files, err = lister.List(filepath.Join("/tmp", id))
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatal("expected 3 files, found:", len(files))
	}

	for _, file := range files {
		if _, ok := expectedPaths[file]; !ok {
			t.Fatal(file, "not found in list")
		}
	}

	// test with recursion off and a pattern
	lister = NewLister(false, "x*")

	files, err = lister.List(filepath.Join("/tmp", id))
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatal("expected 1 files, found:", len(files))
	}

	for _, file := range files {
		if _, ok := expectedPaths[file]; !ok {
			t.Fatal(file, "not found in list")
		}
	}
}
