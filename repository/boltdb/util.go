package boltdb

import (
	"os"
	"path/filepath"
)

func getOrCreateDir(path string) (os.FileInfo, error) {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			cwd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			newDir := filepath.Join(cwd, path)
			err = os.Mkdir(filepath.Join(cwd, path), os.ModePerm)
			if err != nil {
				return nil, err
			}
			return os.Stat(newDir)
		}
		return f, err
	}
	return f, nil
}
