package main

import (
	"fmt"
	"log"
	"net/http"

	"skateeye/components" // Замените на правильный модуль, если используете go mod
)

func main() {
	// Создаем страницу и добавляем блоки
	page := &components.Page{
		Title: "SkateEye AI - Премиальный анализ фигурного катания",
		Components: []components.UIBlock{
			&components.HeaderBlock{},
			&components.MainLayoutBlock{},
			&components.MobileNavBlock{},
		},
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/api/analyze", components.HandleAnalyze)
	http.HandleFunc("/api/coach", components.HandleCoach)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Обязательные заголовки безопасности (Security Skills)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://fonts.gstatic.com https://storage.googleapis.com https://cdn.jsdelivr.net blob: data:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; connect-src 'self' https://cdn.jsdelivr.net blob: data:; object-src 'none';")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Рендерим страницу
		html := page.Render()
		fmt.Fprint(w, html)
	})

	// Запуск сервера исключительно на localhost для безопасности
	port := "8080"
	addr := "127.0.0.1:" + port
	log.Printf("Сервер SkateEye AI запущен по адресу http://%s", addr)
	
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
