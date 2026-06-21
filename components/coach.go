package components

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// AnalysisResult хранит итоговый ответ для фронтенда
type AnalysisResult struct {
	Classification  string             `json:"classification"`
	ProbabilityText string             `json:"probabilityText"`
	BaseScore       float64            `json:"baseScore"`
	GOE             float64            `json:"goe"`
	FinalScore      float64            `json:"finalScore"`
	ScoreReason     string             `json:"scoreReason"`
	DiagnosticCause string             `json:"diagnosticCause"`
	DiagnosticFix   string             `json:"diagnosticFix"`
	Status          string             `json:"status"`
	IsAnomaly       bool               `json:"is_anomaly"`
	Violations      []string           `json:"violations"`
	Biometrics      AdvancedBiometrics `json:"advanced_biometrics"` // Данные от нашего классификатора
	Timeline        []interface{}      `json:"timeline"`
}

// GenerateCoachAdvice отправляет данные прыжка в Mistral для анализа
func GenerateCoachAdvice(ctx context.Context, classResult ClassifierResult) (AnalysisResult, error) {
	// Формируем промпт, который включает все детали биометрии
	prompt := fmt.Sprintf(`Analyze this figure skating jump:
Jump: %s
Biometrics: %+v
Provide analysis in JSON format with fields: classification, probabilityText, baseScore, goe, finalScore, scoreReason, diagnosticCause, diagnosticFix.
Make sure to explain why it is classified as this jump based on takeoff foot and edge type.`,
		classResult.Verdict, classResult.AdvancedBiometrics)

	reqBody := mistralRequest{
		Model: "mistral-tiny",
		Messages: []mistralMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
		ResponseFmt: mistralFmt{Type: "json_object"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("MISTRAL_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AnalysisResult{}, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return AnalysisResult{}, fmt.Errorf("Mistral API error: %s", string(respBytes))
	}

	var apiResp mistralAPIResponse
	json.Unmarshal(respBytes, &apiResp)

	// Очистка от markdown-блоков
	raw := strings.Trim(apiResp.Choices[0].Message.Content, "`json \n")

	var finalResult AnalysisResult
	if err := json.Unmarshal([]byte(raw), &finalResult); err != nil {
		return AnalysisResult{}, err
	}

	// Важно: всегда привязываем биометрию обратно к результату
	finalResult.Biometrics = classResult.AdvancedBiometrics
	finalResult.Timeline = classResult.Timeline
	finalResult.Status = "complete"

	return finalResult, nil
}

// --- Вспомогательные структуры (оставляем без изменений) ---
type mistralRequest struct {
	Model       string           `json:"model"`
	Messages    []mistralMessage `json:"messages"`
	Temperature float64          `json:"temperature"`
	ResponseFmt mistralFmt       `json:"response_format"`
}
type mistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type mistralFmt struct {
	Type string `json:"type"`
}
type mistralAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
