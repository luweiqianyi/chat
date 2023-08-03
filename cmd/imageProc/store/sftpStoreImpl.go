package store

type SftpFileStoreImpl struct {
}

func (impl SftpFileStoreImpl) SaveBytesToRemoteServer(fileBytes []byte, fileName string) (string, error) {
	//remoteSavedPath := SavedPath(FileSavePathPrefix, fileName)
	//dstFile, err := global.SftpClient().Create(remoteSavedPath)
	//if err != nil {
	//	return "", err
	//}
	//defer dstFile.Close()
	//
	//_, err = dstFile.Write(fileBytes)
	//if err != nil {
	//	return "", err
	//}

	// TODO implement me
	return "", nil
}
