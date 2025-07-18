package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bogomollov/15.07.2025/internal/config"
	"github.com/bogomollov/15.07.2025/internal/taskstore"
)

var store = taskstore.GlobalStore

type File struct {
	Name string
	Data io.ReadCloser
}

func WriteArchive(id int) {
	task, _ := store.GetTask(id)

	archiveName := fmt.Sprintf("%d.zip", id)
	archivePath := filepath.Join("assets", archiveName)

	archiveFile, err := os.Create(archivePath)
	if err != nil {
		_, err2 := store.UpdateTaskStatus(id, "failed", "")
		if err2 != nil {
			fmt.Printf("Ошибка при обновлении статуса у %d задачи: %v\n", id, err2)
		}
		return
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	client := &http.Client{Timeout: 30 * time.Second}

	filesChan := make(chan File)
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 3)

	for i, link := range task.Links {
		wg.Add(1)
		limiter <- struct{}{}
		go func(index int, fileLink string) {
			defer wg.Done()
			defer func() { <-limiter }()

			cleanLink := strings.Split(fileLink, "?")[0]
			cleanLink = strings.Split(cleanLink, "#")[0]

			res, _ := client.Get(cleanLink)
			if res.StatusCode != http.StatusOK {
				res.Body.Close()
				return
			}

			fileName := filepath.Base(cleanLink)
			if fileName == "." || fileName == "/" || fileName == "" {
				fileName = fmt.Sprintf("file_%d%s", index+1, filepath.Ext(cleanLink))
			}

			filesChan <- File{Name: fileName, Data: res.Body}
		}(i, link)
	}

	go func() {
		wg.Wait()
		close(filesChan)
	}()

	var archiveError error
	for file := range filesChan {
		fileInZip, err := zipWriter.Create(file.Name)
		if err != nil {
			file.Data.Close()
			archiveError = err
			continue
		}

		if _, err := io.Copy(fileInZip, file.Data); err != nil {
			archiveError = err
		}
		file.Data.Close()
	}
	if archiveError != nil {
		fmt.Printf("Ошибка при создании архива у %d задачи: %v\n", id, archiveError)
		return
	}

	URL := config.GetURL() + ":" + config.GetPort() + "/assets/" + archiveName
	_, err = store.UpdateTaskURL(id, URL)
	if err != nil {
		_, updateErr := store.UpdateTaskStatus(id, "failed", "")
		if updateErr != nil {
			fmt.Printf("Ошибка обновления статуса у %d задачи: %v\n", id, updateErr)
		}
	}
}
