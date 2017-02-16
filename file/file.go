package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Mv(src, dst string) error {
	return os.Rename(src, dst)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	if Exists(path) {
		info, _ := os.Stat(path)
		return info.IsDir()
	}
	return false
}

func MkPath(path string, perm os.FileMode) error {
	if !Exists(path) {
		return os.MkdirAll(path, perm)
	}
	return nil
}

func GetPath(file string) (string, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	base := filepath.Base(abs)
	return abs[:len(abs)-len(base)-1], nil
}

func GetFilesList(path string) ([]string, error) {
	if !IsDir(path) {
		return nil, fmt.Errorf("file or directory not found: %s", path)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, f.Name())
		}
	}
	return names, nil
}

func Size(path string) int64 {
	if Exists(path) {
		info, _ := os.Stat(path)
		return info.Size()
	}
	return 0
}
