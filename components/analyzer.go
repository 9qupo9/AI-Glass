package components

import (
	"encoding/json"
	"math"
	"net/http"
)

const (
	MinRotationComplete = 540.0
	EdgeThreshold       = 1.5
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type FrameData struct {
	ShoulderL Point `json:"shoulderL"`
	ShoulderR Point `json:"shoulderR"`
	HipL      Point `json:"hipL"`
	HipR      Point `json:"hipR"`
	KneeL     Point `json:"kneeL"`
	KneeR     Point `json:"kneeR"`
	AnkleL    Point `json:"ankleL"`
	AnkleR    Point `json:"ankleR"`
	FootL     Point `json:"footL"`
	FootR     Point `json:"footR"`
	WristL    Point `json:"wristL"`
	WristR    Point `json:"wristR"`
}

type AnalysisResult struct {
	Classification  string   `json:"classification"`
	ProbabilityText string   `json:"probabilityText"`
	ShoulderAngle   float64  `json:"shoulderAngle"`
	AnkleLean       float64  `json:"ankleLean"`
	Status          string   `json:"status"`
	BaseScore       float64  `json:"baseScore"`
	GOE             float64  `json:"goe"`
	GOEDeductions   float64  `json:"goe_deductions"`
	FinalScore      float64  `json:"finalScore"`
	ScoreReason     string   `json:"scoreReason"`
	DiagnosticWhat  string   `json:"diagnosticWhat"`
	DiagnosticCause string   `json:"diagnosticCause"`
	DiagnosticFix   string   `json:"diagnosticFix"`
	IsAnomaly       bool     `json:"is_anomaly"`
	AnomalyType     string   `json:"anomaly_type"`
	Violations      []string `json:"violations"`
	PreRotation     bool     `json:"pre_rotation"`
	AxisTilt        float64  `json:"axis_tilt"`
	StepOut         bool     `json:"step_out"`
	HandTouch       bool     `json:"hand_touch"`
	ComboMissed     bool     `json:"combo_missed"`
	HasFlipTurn     bool     `json:"has_flip_turn"`
	IsAxelSetup     bool     `json:"is_axel_setup"`
}

func detectPreRotation(history []FrameData, takeoffIdx int) bool {
	if takeoffIdx < 10 || takeoffIdx >= len(history) {
		return false
	}
	startDx := history[takeoffIdx-10].ShoulderL.X - history[takeoffIdx-10].ShoulderR.X
	takeoffDx := history[takeoffIdx].ShoulderL.X - history[takeoffIdx].ShoulderR.X
	return math.Abs(takeoffDx-startDx) > 0.15
}

func detectAxisTilt(history []FrameData) float64 {
	var maxTilt float64
	for _, frame := range history {
		ankleY := (frame.AnkleL.Y + frame.AnkleR.Y) / 2
		hipY := (frame.HipL.Y + frame.HipR.Y) / 2
		if ankleY > hipY { 
			continue
		}
		shoulderX := (frame.ShoulderL.X + frame.ShoulderR.X) / 2
		shoulderY := (frame.ShoulderL.Y + frame.ShoulderR.Y) / 2
		hipX := (frame.HipL.X + frame.HipR.X) / 2
		dx := shoulderX - hipX
		dy := shoulderY - hipY
		if dy == 0 {
			continue
		}
		tiltDeg := math.Atan2(math.Abs(dx), math.Abs(dy)) * 180.0 / math.Pi
		if tiltDeg > maxTilt {
			maxTilt = tiltDeg
		}
	}
	return maxTilt
}

func detectStepOut(history []FrameData, landingIdx int) bool {
	if landingIdx == -1 || landingIdx >= len(history)-5 {
		return false
	}
	var prevSpread float64 = -1
	for i := landingIdx; i < len(history); i++ {
		frame := history[i]
		spread := math.Abs(frame.FootL.X - frame.FootR.X)
		if prevSpread >= 0 && (spread-prevSpread) > 0.12 {
			return true
		}
		prevSpread = spread
	}
	return false
}

func detectHandTouch(history []FrameData, landingIdx int) bool {
	if landingIdx == -1 || landingIdx >= len(history) {
		return false
	}
	for i := landingIdx; i < len(history); i++ {
		frame := history[i]
		ankleY := (frame.AnkleL.Y + frame.AnkleR.Y) / 2
		if frame.WristL.Y > ankleY-0.05 || frame.WristR.Y > ankleY-0.05 {
			return true
		}
	}
	return false
}

func detectComboMissed(history []FrameData, landingIdx int) bool {
	if landingIdx == -1 || landingIdx >= len(history)-10 {
		return false
	}
	return true
}

func detectFlipTurn(history []FrameData) bool {
	if len(history) < 10 {
		return false
	}
	first := history[0]
	last := history[len(history)-1]
	dx1 := ((first.HipL.X + first.HipR.X) / 2)
	dx2 := ((last.HipL.X + last.HipR.X) / 2)
	return math.Abs(dx2-dx1) > 0.2
}

func calculateBodyDirection(f FrameData) (float64, float64) {
	dx := f.ShoulderR.X - f.ShoulderL.X
	dy := f.ShoulderR.Y - f.ShoulderL.Y
	return -dy, dx
}

func detectAxelSetup(history []FrameData, takeoffIdx int) bool {
	if takeoffIdx < 5 || takeoffIdx >= len(history) {
		return false
	}
	start := history[takeoffIdx-5]
	end := history[takeoffIdx]
	moveX := end.HipL.X - start.HipL.X
	moveY := end.HipL.Y - start.HipL.Y
	bodyX, bodyY := calculateBodyDirection(end)
	return (moveX*bodyX + moveY*bodyY) > 0
}

func detectRotation(history []FrameData, takeoffIdx int, landingIdx int) float64 {
	if takeoffIdx < 0 || landingIdx <= takeoffIdx || landingIdx >= len(history) {
		return 0
	}
	airPhase := history[takeoffIdx:landingIdx]
	if len(airPhase) < 3 {
		return 180.0
	}
	var totalRotation float64
	for i := 1; i < len(airPhase); i++ {
		v1x := airPhase[i-1].ShoulderR.X - airPhase[i-1].ShoulderL.X
		v1y := airPhase[i-1].ShoulderR.Y - airPhase[i-1].ShoulderL.Y
		v2x := airPhase[i].ShoulderR.X - airPhase[i].ShoulderL.X
		v2y := airPhase[i].ShoulderR.Y - airPhase[i].ShoulderL.Y
		
		dot := v1x*v2x + v1y*v2y
		cross := v1x*v2y - v1y*v2x
		angle := math.Atan2(cross, dot)
		totalRotation += math.Abs(angle)
	}
	return totalRotation * 180.0 / math.Pi
}

func AnalyzeJump(history []FrameData) AnalysisResult {
	if len(history) < 15 {
		return AnalysisResult{Status: "analyzing", Classification: "Approach phase...", ProbabilityText: "Analyzing approach..."}
	}

	peakIdx := -1
	minHipY := 9999.0
	for i, frame := range history {
		hipY := (frame.HipL.Y + frame.HipR.Y) / 2
		if hipY < minHipY {
			minHipY = hipY
			peakIdx = i
		}
	}

	if peakIdx == -1 || peakIdx >= len(history) {
		return AnalysisResult{Status: "analyzing", Classification: "Preparing...", ProbabilityText: "Waiting..."}
	}

	takeoffIdx := 0
	maxHipYBefore := -1.0
	for i := 0; i < peakIdx; i++ {
		hipY := (history[i].HipL.Y + history[i].HipR.Y) / 2
		if hipY > maxHipYBefore {
			maxHipYBefore = hipY
			takeoffIdx = i
		}
	}

	landingIdx := len(history) - 1
	maxHipYAfter := -1.0
	for i := peakIdx; i < len(history); i++ {
		hipY := (history[i].HipL.Y + history[i].HipR.Y) / 2
		if hipY > maxHipYAfter {
			maxHipYAfter = hipY
			landingIdx = i
		}
	}

	if takeoffIdx < 2 {
		return AnalysisResult{Status: "analyzing", Classification: "Gliding...", ProbabilityText: "Waiting for jump..."}
	}

	approachFrames := history[:takeoffIdx]
	hasFlipTurn := detectFlipTurn(approachFrames)
	isAxelSetup := detectAxelSetup(history, takeoffIdx)
	rotationDegrees := detectRotation(history, takeoffIdx, landingIdx)
	
	classification := "Unclassified Jump"
	baseScore := 0.0

	if isAxelSetup {
		if rotationDegrees >= 360 && rotationDegrees < 680 {
			classification = "1A (Single Axel)"
			baseScore = 1.10
		} else {
			classification = "2A (Double Axel)"
			baseScore = 3.30
		}
	} else {
		if rotationDegrees >= 540 {
			classification = "2x Double Jump"
			baseScore = 1.30
		}
	}

	axisTilt := detectAxisTilt(history)
	stepOut := detectStepOut(history, landingIdx)
	handTouch := detectHandTouch(history, landingIdx)
	comboMissed := detectComboMissed(history, landingIdx)
	preRotation := detectPreRotation(history, takeoffIdx)

	takeoffFrame := history[takeoffIdx]
	ankleLean := takeoffFrame.FootL.X - takeoffFrame.AnkleL.X

	var violations []string
	var goe float64 = 0.0

	if axisTilt > 15.0 {
		violations = append(violations, "Severe Axis Tilt")
		goe -= 1.0
	}
	if stepOut {
		violations = append(violations, "Step Out")
		goe -= 2.0
	}
	if handTouch {
		violations = append(violations, "Hand Touch")
		goe -= 1.0
	}
	if preRotation {
		violations = append(violations, "Pre-rotation")
		goe -= 0.5
	}

	finalScore := baseScore + goe
	if finalScore < 0 {
		finalScore = 0
	}

	return AnalysisResult{
		Status:          "detected",
		Classification:  classification,
		BaseScore:       baseScore,
		GOE:             goe,
		FinalScore:      finalScore,
		ProbabilityText: "Biometrics verified locally.",
		ShoulderAngle:   rotationDegrees,
		AnkleLean:       ankleLean,
		IsAnomaly:       len(violations) > 0,
		Violations:      violations,
		PreRotation:     preRotation,
		AxisTilt:        axisTilt,
		StepOut:         stepOut,
		HandTouch:       handTouch,
		ComboMissed:     comboMissed,
		HasFlipTurn:     hasFlipTurn,
		IsAxelSetup:     isAxelSetup,
		ScoreReason:     "Deductions based on execution metrics.",
	}
}

func HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var history []FrameData
	if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	result := AnalyzeJump(history)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}