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

type DetectedIssue struct {
	Title       string  `json:"title"`
	Moment      string  `json:"moment"`
	ISUImpact   float64 `json:"isu_impact"`
	Description string  `json:"description"`
}

type JudgingSystem struct {
	CurrentLevel string   `json:"current_level"`
	NextGoal     string   `json:"next_goal"`
	Requirements []string `json:"requirements"`
	Explanation  string   `json:"explanation"`
	Badge        string   `json:"badge"`
}

type HowToCorrect struct {
	IncorrectText string `json:"incorrect_text"`
	CorrectText   string `json:"correct_text"`
}

type PracticeAdvice struct {
	Title       string `json:"title"`
	Repetitions string `json:"repetitions"`
	Focus       string `json:"focus"`
}

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

	// Новые поля для UI
	CoachSummary   string           `json:"coach_summary"`
	DetectedIssues []DetectedIssue  `json:"detected_issues"`
	WhatToPractice []PracticeAdvice `json:"what_to_practice"`
	ISUJudging     JudgingSystem    `json:"isu_judging"`
	USFSAJudging   JudgingSystem    `json:"usfsa_judging"`
	HowToCorrect   HowToCorrect     `json:"how_to_correct"`
}

// GenerateCoachAdvice отправляет данные прыжка в Mistral для анализа
func GenerateCoachAdvice(ctx context.Context, classResult ClassifierResult) (AnalysisResult, error) {
	// Формируем промпт, который включает все детали биометрии и определенный классификатором тип прыжка
	prompt := fmt.Sprintf(`You are an elite figure skating coach and ISU judge. Please analyze this jump: %s.
Technical data:
- Takeoff Type: %s
- Edge Type: %s
- Rotations in the air: %.2f
- Total Degrees rotated: %.2f
- Left leg tilt (LTilt): %.2f degrees
- Right leg tilt (RTilt): %.2f degrees
- Takeoff Foot: %s
- Direction: %s

Based on this data, provide a comprehensive analysis for the skater.
IMPORTANT: You MUST ALWAYS answer in English and return ONLY a valid JSON object.

The JSON format must be exactly:
{
  "classification": "string",
  "probabilityText": "string",
  "baseScore": 3.3,
  "goe": -1.5,
  "finalScore": 1.8,
  "scoreReason": "string",
  "diagnosticCause": "string",
  "diagnosticFix": "string",
  "coach_summary": "A warm, encouraging 2-3 sentence summary of the jump, noting what was good and the main thing to fix.",
  "detected_issues": [
    {
      "title": "Short title (e.g., 'ok: clear jump rise' or 'Late landing check')",
      "moment": "focus window",
      "isu_impact": 0,
      "description": "Detailed explanation of the issue or good point and its effect."
    }
  ],
  "what_to_practice": [
    {
      "title": "Exercise name (e.g. 'Salchow entry walk-throughs')",
      "repetitions": "e.g. '8 repetitions' or '5 sets of 3 seconds'",
      "focus": "What to focus on during the exercise"
    }
  ],
  "isu_judging": {
    "current_level": "e.g. ISU Basic Novice",
    "next_goal": "e.g. Intermediate Novice",
    "requirements": ["improve centering", "stabilize landings", "increase spin quality"],
    "explanation": "Explanation of the jump's level under ISU.",
    "badge": "e.g. BN Review"
  },
  "usfsa_judging": {
    "current_level": "e.g. Basic 4",
    "next_goal": "e.g. Basic 5",
    "requirements": ["stabilize basic elements", "improve rotation quality"],
    "explanation": "Explanation of the jump's level under US Figure Skating.",
    "badge": "e.g. 4 Good"
  },
  "how_to_correct": {
    "incorrect_text": "e.g. The skater leaves the ice, keeps turning in the air, and lands without getting the free leg back.",
    "correct_text": "e.g. The skater should leave the ice, finish the rotation in the air, land with a soft knee, and open the free leg back right away."
  }
}

Make sure to provide 3-4 detected issues (at least one positive 'ok:' and some areas to improve with negative ISU impact) and 2-3 specific practice exercises.`,
		classResult.Verdict,
		classResult.AdvancedBiometrics.TakeoffType,
		classResult.AdvancedBiometrics.EdgeType,
		classResult.AdvancedBiometrics.Rotations,
		classResult.AdvancedBiometrics.TotalDegrees,
		classResult.AdvancedBiometrics.LTilt,
		classResult.AdvancedBiometrics.RTilt,
		classResult.AdvancedBiometrics.TakeoffFoot,
		classResult.AdvancedBiometrics.Direction)

	reqBody := mistralRequest{
		Model: "mistral-small-latest",
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

	// Строго используем нашу детерминированную классификацию, чтобы ИИ не дописывал лишние слова (типа "прыжок")
	finalResult.Classification = classResult.Verdict
	
	// Если Mistral не вернул причину или фикс, добавим дефолтный текст
	if finalResult.DiagnosticCause == "" {
		finalResult.DiagnosticCause = "Analysis completed, but no specific errors were highlighted by the AI."
	}
	if finalResult.DiagnosticFix == "" {
		finalResult.DiagnosticFix = "Keep practicing and focus on maintaining a strong core and proper takeoff edge."
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
