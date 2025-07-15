package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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
				w.WriteHeader(http.StatusOK)
				return
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
			Status: "created",
		}

		data, err := store.CreateTask(newTask)
		if err != nil {
			if errors.Is(err, response.ErrorTaskLimit) {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		} else {
			response.JSON(w, data, http.StatusCreated)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AddLinksInTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/tasks/"), "/")
		if len(pathSegments) < 1 || pathSegments[0] == "" {
			response.JSON(w, map[string]string{"error": "Отсутствует идентификатор в URL"}, http.StatusBadRequest)
			return
		}
		id := pathSegments[0]
		if id != "" {
		idInt, _ := strconv.Atoi(id)
		

		var AddLinksBody struct {
			Links []string `json:"links"`
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := decoder.Decode(&AddLinksBody); err != nil {
			response.JSON(w, map[string]string{"error": "Отсутствуют links в body"}, http.StatusBadRequest)
			return
		}
		task, err := store.UpdateTask(idInt, AddLinksBody.Links)
		if err != nil {
			if errors.Is(err, response.TaskNoFound) {
				response.JSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			}
		} else {
			response.JSON(w, task, http.StatusOK)
		}
	}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}