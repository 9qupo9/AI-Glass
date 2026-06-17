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

type CoachRequest struct {
	History        []FrameData `json:"history"`
	SuccessHistory []FrameData `json:"successHistory,omitempty"`
}

type MistralRequest struct {
	Model       string                `json:"model"`
	Messages    []MistralMessage      `json:"messages"`
	Temperature float64               `json:"temperature"`
	ResponseFmt MistralResponseFormat `json:"response_format"`
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
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return AnalysisResult{Status: "error", DiagnosticFix: "Please configure MISTRAL_API_KEY env variable."}, nil
	}

	preAnalysis := AnalyzeJump(history)

	prompt := `You are SkateEye, an elite AI figure skating coach. Analyze raw 3D coordinates.
Return STRICTLY a JSON object matching the schema below. No markdown wrappers. No backticks. 

Schema:
{
  "classification": "Jump Name",
  "shoulderAngle": 720.0,
  "ankleLean": 1.2,
  "probabilityText": "Confidence score description",
  "baseScore": 3.3,
  "goe": -1.0,
  "finalScore": 2.3,
  "scoreReason": "Reason for score in English",
  "diagnosticCause": "Biomechanical error cause",
  "diagnosticFix": "Actionable instructions to repair technique",
  "isAnomaly": true,
  "violations": ["List of errors"]
}`

	reqBody := MistralRequest{
		Model: "mistral-large-latest",
		Messages: []MistralMessage{
			{Role: "user", Content: prompt + fmt.Sprintf("\nData: %+v", preAnalysis)},
		},
		Temperature: 0.1,
		ResponseFmt: MistralResponseFormat{Type: "json_object"},
	}

	reqBytes, _ := json.Marshal(reqBody)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "[https://api.mistral.ai/v1/chat/completions](https://api.mistral.ai/v1/chat/completions)", bytes.NewReader(reqBytes))
	if err != nil {
		return AnalysisResult{}, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return AnalysisResult{}, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var mResp MistralAPIResponse
	if err := json.Unmarshal(respBytes, &mResp); err != nil || len(mResp.Choices) == 0 {
		return AnalysisResult{}, fmt.Errorf("invalid response from LLM layer")
	}

	rawJson := mResp.Choices[0].Message.Content
	rawJson = strings.TrimPrefix(rawJson, "```json")
	rawJson = strings.TrimSuffix(rawJson, "```")
	rawJson = strings.TrimSpace(rawJson)

	var aiResult AnalysisResult
	if err := json.Unmarshal([]byte(rawJson), &aiResult); err != nil {
		return AnalysisResult{}, err
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
		result = AnalysisResult{Status: "error", DiagnosticFix: err.Error()}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}