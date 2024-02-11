package fileops

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"strings"
	"sync"
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

func (fo *FileOps) FileChecksum(fileName string) (string, error) {
	f, err := os.Open("./assets/" + fileName)
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

func (fo *FileOps) countWords(fileName string) (*models.WordCount, error) {
	f, err := os.Open("./assets/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file\nerror: %s", err.Error())
	}
	defer f.Close()

	wc := &models.WordCount{
		WordsCountMap: map[string]int{},
	}
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		wg.Add(1)
		go func() {
			words := strings.Split(text, " ")
			for _, word := range words {
				wc.Mu.Lock()
				word = strings.Trim(word, " :;.,-*")
				word = strings.ToLower(word)
				if wc.WordsCountMap[word] == 0 {
					wc.TotalWordsCount++
				}
				wc.WordsCountMap[word]++
				wc.Mu.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return wc, nil
}

func (fo *FileOps) CountAllWords() (*models.WordCount, error) {
	entries, err := os.ReadDir("./assets")
	if err != nil {
		return nil, fmt.Errorf("Failed to open file\nerror: %s", err.Error())
	}

	var wg sync.WaitGroup
	wc := &models.WordCount{
		WordsCountMap: map[string]int{},
	}
	// var totalWordsCount :=
	for _, e := range entries {
		wg.Add(1)
		go func(de fs.DirEntry) {
			wordCount, _ := fo.countWords(de.Name())
			wc.TotalFileCount++
			wc.TotalWordsCount += wordCount.TotalWordsCount
			for k, v := range wordCount.WordsCountMap {
				wc.Mu.Lock()
				wc.WordsCountMap[k] += v
				wc.Mu.Unlock()
			}
			wg.Done()
		}(e)
	}

	wg.Wait()
	return wc, nil
}
