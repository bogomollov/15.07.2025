package response

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
    ErrorTaskLimit = errors.New("Нельзя создать новую задачу т.к превышен лимит")
    TaskNoFound = errors.New("Задача с таким идентификатором не найдена")
)

func JSON(w http.ResponseWriter, data interface{}, status int) {
    if status == 0 {
        status = http.StatusOK
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}