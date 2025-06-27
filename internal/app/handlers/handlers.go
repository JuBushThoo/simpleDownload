package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	uploadDir      = "static/uploads"
	uploadFormPath = "templates/down_page.html"
)

// Обработчик для отображения формы
func UploadFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(uploadFormPath)
	if err != nil {
		fmt.Printf("Ошибка загрузки шаблона: %v\n", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Printf("Ошибка выполнения шаблона: %v\n", err)
		http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
	}
}

// Обработчик для загрузки файла
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос на загрузку файла получен")
	if r.Method != http.MethodPost {
		log.Println("Неподдерживаемый метод:", r.Method)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		log.Println("Ошибка при разборе формы:", err)
		http.Error(w, "Не удалось разобрать форму", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Ошибка при получении файла:", err)
		http.Error(w, "Не удалось получить файл", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Println("Файл получен:", handler.Filename)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Println("Ошибка при создании директории для загрузок:", err)
		http.Error(w, "Не удалось создать папку для загрузок", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("Ошибка при создании файла на сервере:", err)
		http.Error(w, "Не удалось создать файл на сервере", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("Ошибка при сохранении файла:", err)
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}

	log.Println("Файл успешно загружен:", handler.Filename)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "success", "message": "Файл '%s' успешно загружен!"}`, handler.Filename)
}

// Прокси-скачивание файла по URL
func ProxyDownloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Не указан url", http.StatusBadRequest)
		return
	}

	// Проверка валидности URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		http.Error(w, "Неподдерживаемый протокол", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка при запросе файла:", err)
		http.Error(w, "Ошибка скачивания файла", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Ошибка статуса при скачивании:", resp.Status)
		http.Error(w, "Ошибка скачивания файла: "+resp.Status, http.StatusBadGateway)
		return
	}

	// Получаем имя файла из URL или заголовков
	filename := "downloaded_file"
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		if strings.Contains(contentDisposition, "filename=") {
			filename = strings.Split(contentDisposition, "filename=")[1]
			filename = strings.Trim(filename, `"`)
		}
	} else {
		if lastSlash := strings.LastIndex(url, "/"); lastSlash != -1 {
			filename = url[lastSlash+1:]
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("Ошибка при копировании файла:", err)
		http.Error(w, "Ошибка при передаче файла", http.StatusInternalServerError)
	}
}
