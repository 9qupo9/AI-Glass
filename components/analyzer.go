package components

import (
	"encoding/json"

	"math"
	"net/http"
)

const (
	// Градусы вращения (900° — это планка для тройных/четверных с недокрутом,
	// но если меньше 540° (1.5 оборота) — это гарантированный срыв)
	MinRotationComplete = 540.0

	// Пороги для наклона лодыжки (Ankle Lean) в пикселях или градусах
	// Положительный наклон — внутреннее ребро (Флип), отрицательный — наружное (Лутц)
	EdgeThreshold = 1.5
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
	WristL    Point `json:"wristL"` // Запястье левой руки — для детекции касания льда
	WristR    Point `json:"wristR"` // Запястье правой руки
}

type AnalysisResult struct {
	Classification  string  `json:"classification"`
	ProbabilityText string  `json:"probabilityText"`
	ShoulderAngle   float64 `json:"shoulderAngle"`
	AnkleLean       float64 `json:"ankleLean"`
	Status          string  `json:"status"`

	BaseScore       float64 `json:"baseScore"`
	GOE             float64 `json:"goe"`
	GOEDeductions   float64 `json:"goe_deductions"` // Суммарный штраф по всем нарушениям
	FinalScore      float64 `json:"finalScore"`
	ScoreReason     string  `json:"scoreReason"`
	DiagnosticWhat  string  `json:"diagnosticWhat"`
	DiagnosticCause string  `json:"diagnosticCause"`
	DiagnosticFix   string  `json:"diagnosticFix"`

	IsAnomaly   bool   `json:"is_anomaly"`
	AnomalyType string `json:"anomaly_type"`

	// Детализированные нарушения по ISU
	Violations  []string `json:"violations"`    // Список всех нарушений
	PreRotation bool     `json:"pre_rotation"`  // Пре-ротация
	AxisTilt    float64  `json:"axis_tilt"`     // Максимальный угол наклона оси в воздухе (градусы)
	StepOut     bool     `json:"step_out"`      // Степаут при приземлении
	HandTouch   bool     `json:"hand_touch"`    // Касание льда рукой
	ComboMissed bool     `json:"combo_missed"`  // Незавершенный каскад
	HasFlipTurn bool     `json:"has_flip_turn"` // Вектор движения (разворот) перед прыжком
}

// =============================================================================
// ДЕТЕКТОРЫ НАРУШЕНИЙ (ISU Judge Module)
// =============================================================================

// detectPreRotation: Плечи начали вращаться ДО взлёта
// Сравниваем направление плеч в первых кадрах (разгон) и кадрах взлёта.
// Если разворот плеч в фазе разгона уже > порога — это пре-ротация.
func detectPreRotation(history []FrameData, takeoffIdx int) bool {
	if takeoffIdx < 10 {
		return false
	}
	// Угол плеч за 10 кадров до взлёта
	startDx := history[takeoffIdx-10].ShoulderL.X - history[takeoffIdx-10].ShoulderR.X
	// Угол плеч в момент взлёта
	takeoffDx := history[takeoffIdx].ShoulderL.X - history[takeoffIdx].ShoulderR.X
	return math.Abs(takeoffDx-startDx) > 0.05
}

// detectAxisTilt: Наклон тела в воздухе
// Вычисляем угол вектора (бедро → плечо) относительно вертикали.
// Если скейтер в воздухе (ankleY высоко), и ось тела отклонена > 15° — аномалия.
func detectAxisTilt(history []FrameData) float64 {
	var maxTilt float64
	for _, frame := range history {
		ankleY := (frame.AnkleL.Y + frame.AnkleR.Y) / 2
		hipY := (frame.HipL.Y + frame.HipR.Y) / 2
		// Только если фигурист в воздухе
		if ankleY-hipY < 0.10 {
			continue
		}
		// Вектор позвоночника: от бёдер к плечам
		shoulderX := (frame.ShoulderL.X + frame.ShoulderR.X) / 2
		shoulderY := (frame.ShoulderL.Y + frame.ShoulderR.Y) / 2
		hipX := (frame.HipL.X + frame.HipR.X) / 2
		dx := shoulderX - hipX
		dy := shoulderY - hipY
		// Угол к вертикали (dy строго < 0 когда голова вверху)
		if dy == 0 {
			continue
		}
		// atan2 даёт угол отклонения от вертикали в радианах
		tiltRad := math.Atan2(math.Abs(dx), math.Abs(dy))
		tiltDeg := tiltRad * 180.0 / math.Pi
		if tiltDeg > maxTilt {
			maxTilt = tiltDeg
		}
	}
	return maxTilt
}

// detectStepOut: Степаут при приземлении
// Смотрим на последние кадры. Если стопы резко разъехались по оси X
// в момент, когда фигурист снова на льду — степаут.
func detectStepOut(history []FrameData, landingIdx int) bool {
	if landingIdx == -1 || landingIdx >= len(history)-5 {
		return false
	}
	// Анализируем кадры после приземления
	landingPhase := history[landingIdx:]
	var prevFootSpread float64 = -1
	for _, frame := range landingPhase {
		ankleY := (frame.AnkleL.Y + frame.AnkleR.Y) / 2
		hipY := (frame.HipL.Y + frame.HipR.Y) / 2
		// Только когда фигурист уже на льду
		if ankleY-hipY < 0.15 {
			spread := math.Abs(frame.FootL.X - frame.FootR.X)
			if prevFootSpread >= 0 && spread-prevFootSpread > 0.08 {
				return true
			}
			prevFootSpread = spread
		}
	}
	return false
}

// detectHandTouch: Запястье опустилось до уровня льда (касание рукой)
func detectHandTouch(history []FrameData, landingIdx int) bool {
	if landingIdx == -1 {
		return false
	}
	for i := landingIdx; i < len(history); i++ {
		frame := history[i]
		ankleY := (frame.AnkleL.Y + frame.AnkleR.Y) / 2
		if math.Abs(frame.WristL.Y-ankleY) < 0.08 || math.Abs(frame.WristR.Y-ankleY) < 0.08 {
			hipY := (frame.HipL.Y + frame.HipR.Y) / 2
			if ankleY-hipY < 0.15 {
				return true
			}
		}
	}
	return false
}

// detectComboMissed: После приземления первого прыжка ожидаем второй взлёт.
// Если ankleY не поднимается снова в течение следующих 30 кадров — каскад не завершён.
func detectComboMissed(history []FrameData, landingIdx int) bool {
	if landingIdx == -1 || landingIdx > len(history)-15 {
		return false
	}
	// Смотрим: поднялся ли фигурист снова после приземления?
	for i := landingIdx + 5; i < len(history); i++ {
		ankleY := (history[i].AnkleL.Y + history[i].AnkleR.Y) / 2
		hipY := (history[i].HipL.Y + history[i].HipR.Y) / 2
		if ankleY-hipY > 0.20 { // Снова в воздухе — второй прыжок есть!
			return false
		}
	}
	return true
}

// detectFlipTurn анализирует смену знака движения по оси X (поиск тройки/разворота)
// Это вектор движения для более точного отличия Флипа от Лутца
func detectFlipTurn(history []FrameData) bool {
	if len(history) < 15 {
		return false
	}

	firstFrame := history[0]
	midFrame := history[len(history)/2]
	lastFrame := history[len(history)-1]

	// Векторы движения таза на первой и второй половине буфера захода
	// Используем среднее между HipL и HipR для большей стабильности
	dx1 := ((midFrame.HipL.X + midFrame.HipR.X) / 2) - ((firstFrame.HipL.X + firstFrame.HipR.X) / 2)
	dx2 := ((lastFrame.HipL.X + lastFrame.HipR.X) / 2) - ((midFrame.HipL.X + midFrame.HipR.X) / 2)

	// Если знак движения по X изменился (двигался влево, стал вправо или наоборот) — был разворот
	if (dx1 > 0 && dx2 < -0.01) || (dx1 < 0 && dx2 > 0.01) {
		return true
	}
	return false
}

func AnalyzeJump(history []FrameData) AnalysisResult {
	if len(history) < 15 {
		return AnalysisResult{Status: "analyzing", Classification: "Approach phase...", ProbabilityText: "Analyzing approach..."}
	}

	// 0. Ищем фазы прыжка по высоте таза (hipY)
	peakIdx := -1
	minHipY := 9999.0
	for i, frame := range history {
		hipY := (frame.HipL.Y + frame.HipR.Y) / 2
		if hipY < minHipY {
			minHipY = hipY
			peakIdx = i
		}
	}

	if peakIdx == -1 {
		return AnalysisResult{Status: "analyzing", Classification: "Preparing...", ProbabilityText: "Waiting..."}
	}

	// Ищем отрыв (takeoff) — самый глубокий присед ДО пика
	takeoffIdx := peakIdx
	maxHipYBefore := minHipY
	for i := peakIdx; i >= 0; i-- {
		hipY := (history[i].HipL.Y + history[i].HipR.Y) / 2
		if hipY >= maxHipYBefore {
			maxHipYBefore = hipY
			takeoffIdx = i
		} else if maxHipYBefore-hipY > 3.0 {
			break
		}
	}

	// Ищем приземление (landing) — локальный максимум ПОСЛЕ пика
	landingIdx := peakIdx
	maxHipYAfter := minHipY
	for i := peakIdx; i < len(history); i++ {
		hipY := (history[i].HipL.Y + history[i].HipR.Y) / 2
		if hipY >= maxHipYAfter {
			maxHipYAfter = hipY
			landingIdx = i
		} else if maxHipYAfter-hipY > 3.0 {
			// Нормальное приземление: таз начал подниматься (амортизация завершена)
			break
		}

		// Падение: таз опустился ниже первоначального приседа (уже летит на лед)
		if hipY > maxHipYBefore+2.0 {
			landingIdx = i
			break
		}
	}

	// Если высота прыжка меньше 1.5 единиц — это не прыжок
	if maxHipYBefore-minHipY < 1.5 {
		return AnalysisResult{Status: "analyzing", Classification: "Gliding...", ProbabilityText: "Waiting for jump..."}
	}

	// ПРЕДОХРАНИТЕЛЬ ОТ АМНЕЗИИ:
	if takeoffIdx < 5 {
		return AnalysisResult{Status: "analyzing", Classification: "Jump completed", ProbabilityText: "Waiting for next jump..."}
	}

	// ПРЕДОХРАНИТЕЛЬ ОТ "ВСТАВАНИЯ СО ЛЬДА":
	if peakIdx-takeoffIdx > 60 {
		return AnalysisResult{Status: "analyzing", Classification: "Recovering...", ProbabilityText: "Getting up from ice"}
	}

	// Второй предохранитель: если в момент "отрыва" таз и лодыжки находятся почти на одной высоте,
	// значит фигурист лежит на льду (горизонтально), а не стоит в приседе!
	takeoffAnkleY := (history[takeoffIdx].AnkleL.Y + history[takeoffIdx].AnkleR.Y) / 2
	takeoffHipY := (history[takeoffIdx].HipL.Y + history[takeoffIdx].HipR.Y) / 2
	if math.Abs(takeoffAnkleY-takeoffHipY) < 0.12 {
		return AnalysisResult{Status: "analyzing", Classification: "Fallen...", ProbabilityText: "Skater is on the ice"}
	}

	// Рассчитываем вектор движения до отрыва (заход на прыжок)
	approachFrames := history
	if takeoffIdx > 0 && takeoffIdx < len(history) {
		approachFrames = history[:takeoffIdx]
	}
	hasFlipTurn := detectFlipTurn(approachFrames)

	// Вычисляем показатели, используя написанные детекторы
	axisTilt := detectAxisTilt(history)
	stepOut := detectStepOut(history, landingIdx)
	handTouch := detectHandTouch(history, landingIdx)
	comboMissed := detectComboMissed(history, landingIdx)
	preRotation := detectPreRotation(history, takeoffIdx)

	// Наклон конька в момент отрыва (для определения ребра)
	takeoffFrame := history[takeoffIdx]
	ankleLean := takeoffFrame.FootL.X - takeoffFrame.AnkleL.X

	var violations []string
	if axisTilt > 15.0 {
		violations = append(violations, "Severe Axis Tilt")
	}
	if stepOut {
		violations = append(violations, "Step Out")
	}
	if handTouch {
		violations = append(violations, "Hand Touch")
	}
	if preRotation {
		violations = append(violations, "Pre-rotation")
	}

	return AnalysisResult{
		Status:          "detected",
		Classification:  "Jump detected",
		ShoulderAngle:   0,
		AnkleLean:       ankleLean,
		ProbabilityText: "Biometric recording finished. Sending to AI...",
		BaseScore:       0,
		GOE:             0,
		GOEDeductions:   0,
		FinalScore:      0,
		ScoreReason:     "Awaiting AI analysis.",
		DiagnosticCause: "Processing frame data with SkateEye AI...",
		DiagnosticFix:   "Please wait for the AI coach...",
		IsAnomaly:       len(violations) > 0,
		AnomalyType:     "",
		Violations:      violations,
		PreRotation:     preRotation,
		AxisTilt:        axisTilt,
		StepOut:         stepOut,
		HandTouch:       handTouch,
		ComboMissed:     comboMissed,
		HasFlipTurn:     hasFlipTurn,
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
