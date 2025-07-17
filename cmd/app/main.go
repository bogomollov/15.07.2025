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
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/api/tasks", rest.CreateHandler)
	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
    	path := r.URL.Path
    	if strings.HasSuffix(path, "/links") {
        	rest.AddLinksHandler(w, r)
    	} else if strings.Count(path, "/") == 3 {
        	rest.GetHandler(w, r)
    	} else {
			http.NotFound(w, r)
		}
	})

    fmt.Println("Сервер запущен на "+cfg.APP_URL+":"+cfg.APP_PORT)
    http.ListenAndServe(":"+cfg.APP_PORT, nil)
}