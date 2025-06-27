package routes

import (
	"net/http"

	"github.com/JuBushThoo/gowebapp/internal/app/handlers"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", handlers.UploadFormHandler)
	router.HandleFunc("/upload", handlers.UploadHandler)
	router.HandleFunc("/downloads", handlers.ProxyDownloadHandler)

	// Обслуживание статики (если понадобится)
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	return router
}
