package utils

import (
	"io/ioutil"
	"path/filepath"
)

func InEachRepo(dir string, fn func(path string) error) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range files {
		if !v.IsDir() {
			continue // skip non dirs
		}
		err := fn(filepath.Join(dir, v.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func InSlice(files []string, filename string) bool {
	for _, file := range files {
		if file == filename {
			return true
		}
	}
	return false
}
