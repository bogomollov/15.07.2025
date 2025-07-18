package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bogomollov/15.07.2025/internal/config"
	"github.com/bogomollov/15.07.2025/internal/transport/rest"
)

func main() {
	cfg := config.Load()

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs)) // статические файлы
	http.HandleFunc("/api/tasks", rest.CreateHandler)         // POST /api/tasks
	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/links") {
			rest.PatchHandler(w, r) // PATCH /api/tasks/{id}/links
		} else if strings.Count(path, "/") == 3 {
			rest.GetHandler(w, r) // GET /api/tasks/{id}
		} else {
			http.NotFound(w, r)
		}
	})

	fmt.Println("Сервер запущен на " + cfg.APP_URL + ":" + cfg.APP_PORT)
	http.ListenAndServe(":"+cfg.APP_PORT, nil)
}
