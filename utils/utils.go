package utils

import (
	"os"
	"path/filepath"
)

func GetRecursiveDirSize(folder string) (int64, error) {
	var size int64
	err := filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}
