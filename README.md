# My Upload Server

Простой сервер для загрузки файлов на Go.

## Запуск

```bash
go run cmd/myapp/main.go
```

## Структура проекта

- `cmd/myapp/main.go`: Основной файл для запуска сервера.
- `internal/handlers/handlers.go`: Обработчики HTTP-запросов.
- `internal/routes/routes.go`: Настройка маршрутов.
- `static/uploads/`: Папка для хранения загруженных файлов.
- `templates/down_page.html`: HTML-форма для загрузки файла.