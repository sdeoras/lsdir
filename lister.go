package lsdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// Lister defines the interface of a file lister
type Lister interface {
	// lister lists files in paths
	List(paths ...string) ([]string, error)
}

// lister implements Lister
type lister struct {
}

// New creates a new instance of lister
func NewLister() Lister {
	return &lister{}
}

func (l *lister) List(paths ...string) ([]string, error) {
	out := make(map[string]struct{})
	for _, path := range paths {
		files, err := list(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			out[file] = struct{}{}
		}
	}

	files := make([]string, 0, len(out))
	for file := range out {
		files = append(files, file)
	}

	sort.Strings(files)

	return files, nil
}

// list is an internal function that is called recursively
func list(path string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var out []string
	if !fileInfo.IsDir() {
		out = append(out, path)
		return out, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		file := file
		if file.IsDir() {
			outSub, err := list(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}

			out = append(out, outSub...)
		} else {
			out = append(out, filepath.Join(path, file.Name()))
		}
	}

	return out, nil
}
