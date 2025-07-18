package response

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrorTaskLimit = errors.New("Нельзя создать новую задачу т.к превышен лимит")
	TaskNotFound   = errors.New("Задача с таким идентификатором не найдена")
	IDURLNotFound  = errors.New("Отсутствует идентификатор в URL запроса")
	IDInvalid      = errors.New("Неверный идентификатор задачи")
	AddLinksError  = errors.New("Неверный формат JSON в теле запроса или отсутствует массив links")
	LinksNotFound  = errors.New("Отсутствуют ссылки с расширениями .pdf или .jpeg и .jpg")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// ответ в формате JSON
func JSON(w http.ResponseWriter, data interface{}, status int) {
	if status == 0 {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ошибка в формате JSON
func JSONError(w http.ResponseWriter, err error, status int) {
	jsonError := ErrorResponse{Error: err.Error()}
	JSON(w, jsonError, status)
}
