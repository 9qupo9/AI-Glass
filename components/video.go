package components

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// HandleAnalyzeVideo принимает MP4, прогоняет через классификатор и LLM.
func HandleAnalyzeVideo(w http.ResponseWriter, r *http.Request) {
	// 0. Валидация метода
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Прием файла с ограничением 200 MB
	if err := r.ParseMultipartForm(200 << 20); err != nil {
		http.Error(w, "Файл слишком большой", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Ошибка получения файла из формы", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Создание временного файла
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%d.mp4", time.Now().UnixNano()))
	tmpFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Ошибка сервера при создании временного файла", http.StatusInternalServerError)
		return
	}

	// Копируем данные в файл
	_, err = io.Copy(tmpFile, file)
	tmpFile.Close() // Закрываем файл перед передачей в функции анализа
	if err != nil {
		os.Remove(tempPath)
		http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
		return
	}

	// Гарантированное удаление файла после завершения всех операций
	defer os.Remove(tempPath)

	// 3. Классификация и анализ (Шаг 1 и 2) - используем контекст с увеличенным таймаутом
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	classResult, err := RunFigureJumpsClassifier(tempPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка классификации прыжка: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Генерация советов через LLM
	coachResult, err := GenerateCoachAdvice(ctx, classResult)
	if err != nil {
		// Fallback: возвращаем то, что удалось получить от классификатора
		coachResult = AnalysisResult{
			Status:          "partial",
			Classification:  classResult.Verdict,
			Biometrics:      classResult.AdvancedBiometrics,
			Timeline:        classResult.Timeline,
			DiagnosticCause: "Анализ выполнен, но LLM-тренер недоступен. " + err.Error(),
			DiagnosticFix:   "Попробуйте повторить запрос позже.",
		}
	}

	// 5. Ответ фронтенду
	w.Header().Set("Content-Type", "application/json")
	
	if coachResult.Status == "partial" {
		fmt.Printf("[LLM Status] %s\n", coachResult.DiagnosticCause)
	}

	if err := json.NewEncoder(w).Encode(coachResult); err != nil {
		// Логируем ошибку, так как заголовки уже могли быть отправлены
		fmt.Printf("Error encoding JSON: %v\n", err)
	}
}
