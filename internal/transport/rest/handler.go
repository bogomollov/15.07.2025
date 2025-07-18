package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bogomollov/15.07.2025/internal/archive"
	"github.com/bogomollov/15.07.2025/internal/taskstore"
	"github.com/bogomollov/15.07.2025/internal/transport/response"
)

var store = taskstore.GlobalStore

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
		if id != "" {
			idInt, _ := strconv.Atoi(id)
			data, found := store.GetTask(idInt)
			if !found {
				response.JSONError(w, response.TaskNotFound, http.StatusNotFound)
			} else {
				response.JSON(w, data, http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newTask := taskstore.Task{
			Status: taskstore.TaskStatusCreated,
		}

		data, err := store.CreateTask(newTask)
		if err != nil {
			if errors.Is(err, response.ErrorTaskLimit) {
				response.JSONError(w, err, http.StatusServiceUnavailable)
			}
		} else {
			response.JSON(w, data, http.StatusCreated)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		path := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/tasks/"), "/")
		if len(path) < 1 || path[0] == "" {
			response.JSONError(w, response.IDURLNotFound, http.StatusBadRequest)
			return
		}
		idStr := path[0]

		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSONError(w, response.IDInvalid, http.StatusBadRequest)
			return
		}

		var addLinksBody struct {
			Links []string `json:"links"`
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&addLinksBody); err != nil {
			response.JSONError(w, response.AddLinksError, http.StatusBadRequest)
			return
		}

		var filteredLinks []string
		exts := map[string]struct{}{
			".pdf":  {},
			".jpeg": {},
			".jpg":  {},
		}
		for _, link := range addLinksBody.Links {
			cleanLink := strings.Split(link, "?")[0]
			cleanLink = strings.Split(cleanLink, "#")[0]

			ext := strings.ToLower(filepath.Ext(cleanLink))
			if _, ok := exts[ext]; ok {
				filteredLinks = append(filteredLinks, link)
			}
		}
		if len(filteredLinks) == 0 {
			response.JSONError(w, response.LinksNotFound, http.StatusBadRequest)
			return
		}

		task, err := store.UpdateTask(idInt, filteredLinks)
		if err != nil {
			if errors.Is(err, response.TaskNotFound) {
				response.JSONError(w, response.TaskNotFound, http.StatusNotFound)
			}
		} else {
			response.JSON(w, task, http.StatusOK)
			updatedTask, _ := store.GetTask(idInt)
			if len(updatedTask.Links) >= 3 {
				go archive.WriteArchive(task.ID)
			}
		}
	}
}
