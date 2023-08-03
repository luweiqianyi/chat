package util

import "os"

func isPathExist(folder string) bool {
	_, err := os.Stat(folder)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func EnsureFolderExist(folder string) error {
	if !isPathExist(folder) {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
