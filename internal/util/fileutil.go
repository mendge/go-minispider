package util

import (
	"os"
)

// BuildDir build directory if dirPath is not existed.
func BuildDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	needToBuild := os.IsNotExist(err)
	if needToBuild {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
