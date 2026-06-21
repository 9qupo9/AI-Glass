package components

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AdvancedBiometrics struct {
	TakeoffFrame int     `json:"takeoff_frame"`
	LandingFrame int     `json:"landing_frame"`
	Rotations    float64 `json:"rotations"`
	TotalDegrees float64 `json:"total_degrees"`
	Direction    string  `json:"direction"`
	TakeoffFoot  string  `json:"support_foot"`
	TakeoffType  string  `json:"takeoff_type"`
	EdgeType     string  `json:"edge_type"`
	LTilt        float64 `json:"l_tilt"`
	RTilt        float64 `json:"r_tilt"`
}

type ClassifierResult struct {
	Verdict            string             `json:"verdict"`
	Confidence         float64            `json:"confidence"`
	Model              string             `json:"model"`
	AdvancedBiometrics AdvancedBiometrics `json:"advanced_biometrics"`
	Timeline           []interface{}      `json:"timeline"`
}

type AnalyzerResponse struct {
	Status             string             `json:"status"`
	AdvancedBiometrics AdvancedBiometrics `json:"advanced_biometrics"`
	Timeline           []interface{}      `json:"timeline"`
}

// classifyJump применяет Таблицу Истинности для определения прыжка
func classifyJump(b AdvancedBiometrics) string {
	fmt.Printf("[DEBUG] Takeoff=%s, Foot=%s, Dir=%s, Edge=%s | L_Tilt=%.4f, R_Tilt=%.4f\n", 
		b.TakeoffType, b.TakeoffFoot, b.Direction, b.EdgeType, b.LTilt, b.RTilt)

	// Приоритет физики: ребро и тип захода
	if b.TakeoffType == "Edge" {
		if b.Direction == "Forward" && b.TakeoffFoot == "Left Foot" && b.EdgeType == "Outside" {
			return "Axel"
		}
		if b.Direction == "Backward" && b.TakeoffFoot == "Left Foot" && b.EdgeType == "Inside" {
			return "Salchow"
		}
		if b.Direction == "Backward" && b.TakeoffFoot == "Right Foot" && b.EdgeType == "Outside" {
			return "Loop"
		}
	}

	if b.TakeoffType == "Toe" && b.Direction == "Backward" {
		if b.TakeoffFoot == "Right Foot" {
			return "Toe Loop"
		}
		if b.EdgeType == "Inside" {
			return "Flip"
		}
		if b.EdgeType == "Outside" {
			return "Lutz"
		}
	}

	return "Unknown Jump"
}

// RunFigureJumpsClassifier запускает python-анализатор
func RunFigureJumpsClassifier(videoPath string) (ClassifierResult, error) {
	cwd, _ := os.Getwd()
	workDir := filepath.Join(cwd, "Neiron", "data", "Video_data")
	wrapperPath := filepath.Join(workDir, "jump3d_analyzer.py")
	pythonExe := findPython()

	cmd := exec.Command(pythonExe, wrapperPath, videoPath)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[Python Error] %s\n", string(out))
		return ClassifierResult{}, fmt.Errorf("analyzer execution failed: %v", err)
	}

	// Ищем последний JSON объект в выводе (защита от мусорных логов)
	jsonLine := ""
	lines := strings.Split(string(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}") {
			jsonLine = line
			break
		}
	}

	if jsonLine == "" {
		return ClassifierResult{}, fmt.Errorf("no valid JSON found in analyzer output")
	}

	var response AnalyzerResponse
	if err := json.Unmarshal([]byte(jsonLine), &response); err != nil {
		return ClassifierResult{}, fmt.Errorf("json unmarshal error: %w", err)
	}

	return ClassifierResult{
		Verdict:            classifyJump(response.AdvancedBiometrics),
		Confidence:         1.0,
		Model:              "Go-TruthTable",
		AdvancedBiometrics: response.AdvancedBiometrics,
		Timeline:           response.Timeline,
	}, nil
}

func findPython() string {
	minicondaPath := `C:\Users\9qupo\Miniconda3\python.exe`
	if _, err := os.Stat(minicondaPath); err == nil {
		return minicondaPath
	}
	return "python"
}
