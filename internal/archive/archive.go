package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bogomollov/15.07.2025/internal/config"
	"github.com/bogomollov/15.07.2025/internal/taskstore"
)

var store = taskstore.GlobalStore

func WriteArchive(id int) error {
	task, _ := store.GetTask(id)

	archive := fmt.Sprintf("%d.zip", id)
	archivePath := filepath.Join("assets", archive)
	archiveFile, _ := os.Create(archivePath)
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	client := &http.Client{Timeout: 30 * time.Second}

	for i, link := range task.Links {
		cleanLink := strings.Split(link, "?")[0]
		cleanLink = strings.Split(cleanLink, "#")[0]

		resp, err := client.Get(cleanLink)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		fileName := filepath.Base(cleanLink)
		if fileName == "." || fileName == "/" || fileName == "" {
			fileName = fmt.Sprintf("file_%d", i+1)
		}

		fileInZip, err := zipWriter.Create(fileName)
		if err != nil {
			resp.Body.Close()
			continue
		}

		if _, err := io.Copy(fileInZip, resp.Body); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
	}

	URL := config.GetURL()+":"+config.GetPort()+"/assets/" + archive
	store.UpdateTaskURL(id, URL)

	return nil
}
