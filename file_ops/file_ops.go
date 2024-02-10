package fileops

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"time"

	"github.com/docker/go-units"
	"github.com/zakisk/redhat/server/models"
)

type FileOps struct {
	Hasher hash.Hash
}

func NewFileOps(hasher hash.Hash) *FileOps {
	return &FileOps{Hasher: hasher}
}

func (fo *FileOps) FileChecksum(info fs.FileInfo) (string, error) {
	f, err := os.Open("./assets/" + info.Name())
	if err != nil {
		return "", fmt.Errorf("Failed to open file\nerror: %s", err.Error())
	}

	if _, err = io.Copy(fo.Hasher, f); err != nil {
		return "", fmt.Errorf("Failed to copy file content\nerror: %s", err.Error())
	}
	defer fo.Hasher.Reset()

	hash := fo.Hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func (fo *FileOps) CreateFile(fileName string, file multipart.File) error {
	newFile, err := os.Create("./assets/" + fileName)
	if err != nil {
		return fmt.Errorf("Failed to create new file\nerror: %s", err.Error())
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return fmt.Errorf("Failed to copy file data\nerror: %s", err.Error())
	}

	return nil
}

func (fo *FileOps) UpdateFile(fileName string, file multipart.File) error {
	currentFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open destination file\nerror: %s", err.Error())
	}
	defer currentFile.Close()

	_, err = io.Copy(currentFile, file)
	if err != nil {
		return fmt.Errorf("Failed to copy data into file\nerror: %s", err.Error())
	}

	return nil
}

func (fo *FileOps) RemoveFile(fileName string) error {
	_, err := os.Stat("./assets/" + fileName)
	if err != nil {
		return err
	}

	err = os.Remove("./assets/" + fileName)
	if err != nil {
		return fmt.Errorf("Failed to delete file\nerror: %s", err.Error())
	}

	return nil
}

func (fo *FileOps) ListFile(dir string) ([]*models.File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Failed to read directory\nerror: %s", err.Error())
	}

	files := []*models.File{}
	for _, e := range entries {
		info, _ := e.Info()
		file := &models.File{
			Name:       info.Name(),
			Mode:       info.Mode().String(),
			ModifiedAt: info.ModTime().Format(time.DateTime),
			Size:       units.HumanSizeWithPrecision(float64(info.Size()), 3),
		}

		files = append(files, file)
	}

	return files, nil
}
