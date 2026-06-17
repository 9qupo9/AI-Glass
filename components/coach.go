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
`
	}

	prompt += `Рассчитай:
1. Какой это прыжок (Lutz, Flip, Axel, Loop, Toe Loop, Salchow)? Учитывай ребро отрыва (Ankle Lean) и заход.
2. Сколько было полных оборотов? (плечи X перекрещиваются, 1 оборот = 360 градусов). Умножь на 180 за каждое пересечение (за заход Акселя +180).
3. Наклон оси (Axis Tilt) в градусах.
4. Было ли падение (таз упал на лед после приземления)?
5. Напиши вердикт тренера.

ВАЖНО: ОТВЕЧАЙ СТРОГО НА РУССКОМ ЯЗЫКЕ! Все текстовые поля (scoreReason, diagnosticCause, diagnosticFix, violations, classification) должны быть на русском.

Требуемый формат ответа (строгий JSON, ключи должны совпадать точно):
{
  "classification": "Название прыжка и нарушения в скобках, например Triple Lutz (Fall)",
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
	axelStr := "СПИНОЙ ВПЕРЕД (Не Аксель)"
	if preAnalysis.IsAxelSetup {
		axelStr = "ЛИЦОМ ВПЕРЕД (ЭТО СТРОГО АКСЕЛЬ!)"
	}

	prompt += fmt.Sprintf(`
ВНИМАНИЕ! Бэкенд уже рассчитал для тебя ключевые параметры:
- Направление захода: %s
- Траектория захода: %s
- Ребро в момент отрыва: %s
- Угол наклона оси в воздухе: %.1f градусов
- Статус приземления: %s
- Пре-ротация (ранний разворот плеч): %v

ПРАВИЛА ISU:
0. ЕСЛИ ЗАХОД "ЛИЦОМ ВПЕРЕД" — ЭТО 100%% АКСЕЛЬ (Axel). ТЕБЕ ЗАПРЕЩЕНО НАЗЫВАТЬ ЕГО ФЛИПОМ, ЛУТЦЕМ ИЛИ САЛЬХОВОМ. Выбирай только между Single Axel, Double Axel или Triple Axel (в зависимости от оборотов).
1. FLIP (Флип): правильный прыжок исполняется строго с ВНУТРЕННЕГО ребра + заход с разворота (тройки). Если ребро наружное — это ошибка [e] (Lip).
2. LUTZ (Лутц): правильный прыжок исполняется строго с НАРУЖНОГО ребра + заход по прямой/дуге без разворота. Если ребро внутреннее — это ошибка [e] (Flutz).
3. Наклон оси более 15-20 градусов вызывает падение и огромные штрафы GOE.

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
