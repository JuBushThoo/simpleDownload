package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JuBushThoo/gowebapp/internal/app/routes"
)

func main() {
	router := routes.SetupRouter()
	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
