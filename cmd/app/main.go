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

	go http.HandleFunc("/api/tasks", rest.CreateHandler)
	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
    	path := r.URL.Path
    	if strings.HasSuffix(path, "/links") {
        rest.AddLinksInTask(w, r)
    	} else if strings.Count(path, "/") == 3 {
        	rest.GetHandler(w, r)
    	}
	})

    fmt.Println("Сервер запущен на порту", cfg.PORT)
    http.ListenAndServe(":"+cfg.PORT, nil)

	// url := "https://i.pinimg.com/originals/b8/07/c2/b807c2282ab0a491bd5c5c1051c6d312.jpg"
	// response, e := http.Get(url)
	// if e != nil {
	//     log.Fatal(e)
	// }
	// defer response.Body.Close()

	// file, err := os.Create("./assets/test.jpg")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer file.Close()

	// _, err = io.Copy(file, response.Body)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// fmt.Println("Success!")
}