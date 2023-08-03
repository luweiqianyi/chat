package store

import "fmt"

const (
	FileSavePathPrefix = "file-server-127.0.0.1/images/avatar"
)

func SavedPath(folder string, originFileName string) string {
	return fmt.Sprintf("%v/%v", folder, originFileName)
}
