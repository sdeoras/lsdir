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
	pattern string
	recurse bool
}

// New creates a new instance of lister using a pattern that matches using
// shell wildcards (and not regexp). recurse indicates if recursive search is on.
func NewLister(recurse bool, pattern string) Lister {
	l := new(lister)
	l.pattern = pattern
	l.recurse = recurse

	return l
}

func (l *lister) List(paths ...string) ([]string, error) {
	out := make(map[string]struct{})
	var listFunc func(path string) ([]string, error)

	if l.recurse {
		listFunc = list
	} else {
		listFunc = listNonRecursive
	}

	for _, path := range paths {
		files, err := listFunc(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			out[file] = struct{}{}
		}
	}

	files := make([]string, 0, len(out))
	for file := range out {
		if matched, err := filepath.Match(l.pattern, filepath.Base(file)); matched && err == nil {
			files = append(files, file)
		} else if err != nil {
			return nil, err
		}
	}

	sort.Strings(files)

	return files, nil
}

func listNonRecursive(path string) ([]string, error) {
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
		if !file.IsDir() {
			out = append(out, filepath.Join(path, file.Name()))
		}
	}

	return out, nil
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
