package fileStorage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type LocalFileStorage struct {
	StoragePath string
}

func NewLocalFileStorage(storagePath string) *LocalFileStorage {
	return &LocalFileStorage{StoragePath: storagePath}
}

func (storage *LocalFileStorage) UploadFile(fileFromRequest multipart.File, fileName string, ext string) error {
	fileToSave, err := os.OpenFile(
		filepath.Join(storage.StoragePath, strings.Join([]string{fileName, ext}, ".")),
		os.O_WRONLY|os.O_CREATE,
		0666,
	)
	if err != nil {
		return err
	}
	defer fileToSave.Close()
	_, err = io.Copy(fileToSave, fileFromRequest)
	if err != nil {
		return err
	}
	return nil
}
