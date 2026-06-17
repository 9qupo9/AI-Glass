package components

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CoachRequest struct {
	History        []FrameData `json:"history"`
	SuccessHistory []FrameData `json:"successHistory,omitempty"`
}

type CoachResponse struct {
	Advice string `json:"advice"`
}

type MistralRequest struct {
	Model       string                  `json:"model"`
	Messages    []MistralMessage        `json:"messages"`
	Temperature float64                 `json:"temperature"`
	ResponseFmt MistralResponseFormat   `json:"response_format"`
}

type MistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MistralResponseFormat struct {
	Type string `json:"type"`
}

type MistralAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateCoachAdvice(ctx context.Context, history []FrameData, successHistory []FrameData) (AnalysisResult, error) {
	apiKey := "ydIANPrsuxbcZyO7qI6t2KjM156IlFSc"
	if apiKey == "" {
		return AnalysisResult{Status: "error", DiagnosticFix: "Пожалуйста, установите ключ Mistral API."}, nil
	}

	prompt := `Ты - профессиональный тренер по фигурному катанию мирового класса и эксперт по биомеханике ISU. 
Твоя задача - проанализировать сырые координаты (JSON) скелета фигуриста (30 fps) во время текущего прыжка и вернуть СТРОГИЙ JSON.
В массиве history лежат кадры текущего (возможно неудачного) прыжка (от начала приседа до приземления).
Y растет вниз. Z - глубина. X - горизонталь.
`
	if len(successHistory) > 0 {
		prompt += `Также тебе передан массив successHistory - это эталонный (удачный) прыжок этого же фигуриста.
ВАЖНО: Основной анализ и оценки ( classification, score, violations ) делай по текущему прыжку (history).
Однако в поле "diagnosticFix" и "diagnosticCause" обязательно сравни текущий прыжок с эталонным (successHistory) и укажи конкретные отличия, которые привели к ошибке в текущем прыжке.
		return AnalysisResult{Status: "error", DiagnosticFix: "Please set the Mistral API key."}, nil
	}

	prompt := `You are SkateEye, an elite AI figure skating coach. Analyze the biomechanical data provided by the MediaPipe tracker and give short, actionable feedback to the skater. IMPORTANT: ANSWER STRICTLY IN ENGLISH. Explain what failed and how to fix it based on the ISU rules.
In the history array, you have the frames of the current (possibly failed) jump (from start of squat to landing).
Y grows downwards. Z is depth. X is horizontal.
`
	if len(successHistory) > 0 {
		prompt += `You have also been provided with a successHistory array - this is the reference (successful) jump by the same skater.
IMPORTANT: Do your main analysis and evaluations (classification, score, violations) based on the current jump (history).
However, in the "diagnosticFix" and "diagnosticCause" fields, you must compare the current jump with the reference one (successHistory) and point out the specific differences that led to the error in the current jump.
`
	}

	prompt += `Calculate:
1. What is the jump (Lutz, Flip, Axel, Loop, Toe Loop, Salchow)? Consider the edge at takeoff (Ankle Lean) and entry.
2. How many full rotations? (shoulders X cross, 1 rotation = 360 degrees. Multiply by 180 for each crossing; +180 for Axel setup).
3. Axis tilt in degrees.
4. Was there a fall (pelvis touched the ice after landing)?
5. Write a coach's verdict.

IMPORTANT: ANSWER STRICTLY IN ENGLISH. All text fields (scoreReason, diagnosticCause, diagnosticFix, violations, classification) must be in English.

Required response format (strict JSON, keys must match exactly):
{
  "classification": "Name of the jump and violations in parentheses, e.g., Triple Lutz (Fall)",
  "shoulderAngle": 1080,
  "ankleLean": 15.5,
  "probabilityText": "AI Matrix: Confidence 99%",
  "baseScore": 5.9,
  "goe": -1.5,
  "finalScore": 4.4,
  "scoreReason": "Краткое объяснение оценки",
  "diagnosticCause": "Причина ошибки с точки зрения физики (сравнить с эталонным, если передан)",
  "diagnosticFix": "Совет тренера по исправлению (сравнить с эталонным, если передан)",
  "isAnomaly": true,
  "violations": ["список", "штрафов", "текстом"]
}

Вот данные текущего прыжка (history):
`

	// Вычисляем физические параметры нашим локальным бэкендом, чтобы помочь ИИ
	preAnalysis := AnalyzeJump(history)

	edgeType := "Внутреннее (Inside Edge)"
	if preAnalysis.AnkleLean < 0 {
		edgeType = "Наружное (Outside Edge)"
	}

	turnType := "Заход по прямой / дуге (без резкого разворота)"
	if preAnalysis.HasFlipTurn {
		turnType = "Резкий разворот таза перед прыжком (смена направления)"
	}

	landingStatus := "Нормальное"
	if preAnalysis.StepOut {
		landingStatus = "Степаут (Step-out)"
	} else if preAnalysis.HandTouch {
		landingStatus = "Касание льда рукой"
	}

	// Вставляем наш текст в промпт для Mistral
	axelStr := "BACKWARD TAKEOFF (Not an Axel)"
	if preAnalysis.IsAxelSetup {
		axelStr = "FORWARD TAKEOFF (THIS IS STRICTLY AN AXEL!)"
	}

	prompt += fmt.Sprintf(`
ATTENTION! The backend has already calculated key parameters for you:
- Takeoff direction: %s
- Entry trajectory: %s
- Edge at takeoff: %s
- Axis tilt in the air: %.1f degrees
- Landing status: %s
- Pre-rotation (early shoulder rotation): %v

ISU RULES:
0. IF TAKEOFF IS "FORWARD" — IT IS 100%% AN AXEL. YOU ARE FORBIDDEN TO CALL IT A FLIP, LUTZ, OR SALCHOW. Choose only between Single Axel, Double Axel, or Triple Axel (depending on rotations).
1. FLIP: Correct jump is strictly from an INSIDE edge + turn entry (three-turn). If the edge is outside — it's an [e] error (Lip).
2. LUTZ: Correct jump is strictly from an OUTSIDE edge + straight/curve entry without a turn. If the edge is inside — it's an [e] error (Flutz).
3. Axis tilt over 15-20 degrees causes falls and massive GOE deductions.

`, axelStr, turnType, edgeType, preAnalysis.AxisTilt, landingStatus, preAnalysis.PreRotation)


	historyJson, _ := json.Marshal(history)
	prompt += string(historyJson)

	if len(successHistory) > 0 {
		successJson, _ := json.Marshal(successHistory)
		prompt += "\n\nВот данные эталонного (удачного) прыжка (successHistory):\n" + string(successJson)
	}

	reqBody := MistralRequest{
		Model: "mistral-large-latest",
		Messages: []MistralMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.2,
		ResponseFmt: MistralResponseFormat{Type: "json_object"},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("ошибка сборки запроса: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewReader(reqBytes))
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("ошибка создания запроса к Mistral API: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("ошибка отправки запроса к Mistral API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("ошибка чтения ответа от Mistral: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return AnalysisResult{}, fmt.Errorf("Mistral API вернул ошибку %d: %s", resp.StatusCode, string(respBytes))
	}

	var orResp MistralAPIResponse
	if err := json.Unmarshal(respBytes, &orResp); err != nil {
		return AnalysisResult{}, fmt.Errorf("ошибка парсинга ответа Mistral: %w", err)
	}

	if len(orResp.Choices) == 0 {
		return AnalysisResult{}, fmt.Errorf("пустой ответ от Mistral")
	}

	rawJson := orResp.Choices[0].Message.Content
	rawJson = strings.TrimPrefix(rawJson, "```json")
	rawJson = strings.TrimPrefix(rawJson, "```")
	rawJson = strings.TrimSuffix(rawJson, "```")
	rawJson = strings.TrimSpace(rawJson)

	var aiResult AnalysisResult
	if err := json.Unmarshal([]byte(rawJson), &aiResult); err != nil {
		return AnalysisResult{}, fmt.Errorf("ошибка парсинга JSON от ИИ: %w\nRaw string: %s", err, rawJson)
	}
	aiResult.Status = "detected"
	return aiResult, nil
}

func HandleCoach(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CoachRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := GenerateCoachAdvice(r.Context(), req.History, req.SuccessHistory)
	if err != nil {
		result = AnalysisResult{
			Status: "error",
			Classification: "AI Error",
			DiagnosticFix: fmt.Sprintf("Ошибка ИИ-тренера: %v", err),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
