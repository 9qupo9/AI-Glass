package components

// GetGlobalCSS возвращает весь глобальный CSS для приложения
func GetGlobalCSS() string {
	return `
/* SkateEye AI — Премиальный дизайн-код */

:root {
  /* Strict Professional Palette */
  --bg-main: #0A0A0A;
  --bg-card: #111111;
  --bg-card-hover: #1A1A1A;
  
  --border-color: rgba(255, 255, 255, 0.1);
  --border-glow: transparent;
  
  --color-primary: #3b82f6; /* Blue 500 */
  --color-primary-hover: #2563eb;
  --color-accent: #60a5fa; 
  --color-success: #059669; /* Emerald 600 */
  --color-warning: #d97706; /* Amber 600 */
  --color-danger: #dc2626; /* Red 600 */
  
  --text-main: #f8fafc;
  --text-muted: #94a3b8;

  --font-primary: 'Inter', system-ui, sans-serif;
  --font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  --font-display: 'Inter', system-ui, sans-serif;

  --transition-fast: 0.1s ease-out;
  --transition-smooth: 0.15s ease-out;
}

/* Сброс и базовые настройки */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  -webkit-tap-highlight-color: transparent;
}

body {
  font-family: var(--font-primary);
  background-color: var(--bg-main);
  color: var(--text-main);
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: hidden;
  background-image: none;
}

/* Кастомизация скроллбара */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 6px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.25);
}

/* Основной контейнер-сетка */
.app-container {
  display: grid;
  grid-template-rows: auto 1fr;
  height: 100vh;
  height: 100dvh;
}

/* ==========================================
   Шапка (Header)
   ========================================== */
.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 32px;
  background: #09090b;
  border-bottom: 1px solid var(--border-color);
  z-index: 10;
  position: sticky;
  top: 0;
}

.brand-logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-name {
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.5px;
  color: var(--text-main);
}

.desktop-tabs {
  display: flex;
  gap: 8px;
  background: rgba(255, 255, 255, 0.03);
  padding: 6px;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.desk-tab-btn {
  background: transparent;
  color: var(--text-muted);
  border: none;
  padding: 8px 20px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 6px;
  cursor: pointer;
  transition: var(--transition-smooth);
}

.desk-tab-btn:hover {
  color: var(--text-main);
  background: rgba(255, 255, 255, 0.05);
}

.desk-tab-btn.active {
  color: #fff;
  background: var(--color-primary);
  box-shadow: none;
}

.jump-select {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  color: #fff;
  padding: 10px 16px;
  border-radius: 6px;
  font-family: var(--font-primary);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  outline: none;
  transition: var(--transition-fast);
}

.jump-select:hover, .jump-select:focus {
  border-color: var(--color-primary);
  background: rgba(255, 255, 255, 0.08);
}

.confidence-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.conf-high {
  background: rgba(16, 185, 129, 0.1);
  color: var(--color-success);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.conf-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: currentColor;
}

/* ==========================================
   Основной макет (Main Layout)
   ========================================== */
.main-content {
  display: flex;
  flex-direction: row;
  overflow: hidden;
  padding: 24px 32px;
  gap: 24px;
}

/* ==========================================
   Универсальные Карточки (Panel Cards)
   ========================================== */
.diagnostics-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.diag-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  box-shadow: none;
  
  transition: transform 0.2s ease, border-color 0.2s ease;
}

.diag-card:hover {
  transform: translateY(-2px);
}

.diag-card.error {
  border-left: 5px solid #EF4444;
}
.diag-card.error:hover { border-color: #EF4444; }

.diag-card.cause {
  border-left: 5px solid #F59E0B;
}
.diag-card.cause:hover { border-color: #F59E0B; }

.diag-card.fix {
  border-left: 5px solid #10B981;
}
.diag-card.fix:hover { border-color: #10B981; }

.diag-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.diag-icon {
  font-size: 22px;
  background: rgba(255, 255, 255, 0.1);
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.diag-title {
  color: #F8FAFC;
  font-weight: 700;
  font-size: 16px;
  letter-spacing: 0.3px;
}

.diag-body {
  color: #94A3B8;
  font-size: 14.5px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.scrollable-content-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  height: 100%;
}

.panel-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  transition: var(--transition-smooth);
}

.panel-card:hover {
  border-color: rgba(255, 255, 255, 0.15);
}

.panel-title {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 20px;
  letter-spacing: 0.5px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.panel-title::before {
  content: "";
  display: inline-block;
  width: 4px;
  height: 16px;
  background: var(--color-primary);
  border-radius: 6px;
}

/* ==========================================
   Видеоплеер (Video Panel)
   ========================================== */
.video-panel-wrapper {
  flex: 1.2;
  min-width: 0;
  display: flex;
  flex-direction: row;
  gap: 24px;
}

.video-container {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 9;
  min-height: 300px;
  background: #000;
  border-radius: 6px;
  overflow: hidden;
  box-shadow: none;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.video-player {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.video-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 5;
}

.video-overlay-ui {
  position: absolute;
  top: 16px;
  left: 16px;
  right: 16px;
  display: flex;
  justify-content: space-between;
  z-index: 10;
  pointer-events: none;
}

.overlay-tag {
  background: rgba(0, 0, 0, 0.6);
  
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.video-controls {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.scrubber-row {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(0, 0, 0, 0.2);
  padding: 12px 16px;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.play-pause-btn {
  background: var(--color-primary);
  border: none;
  color: #fff;
  width: 36px;
  height: 36px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  transition: var(--transition-fast);
}

.play-pause-btn:hover {
  background: var(--color-primary-hover);
  transform: scale(1.05);
}

.timeline-slider-container {
  flex-grow: 1;
}

.timeline-slider {
  width: 100%;
  -webkit-appearance: none;
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  outline: none;
}

.timeline-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #fff;
  border: 3px solid var(--color-primary);
  cursor: pointer;
  transition: transform var(--transition-fast);
}

.timeline-slider::-webkit-slider-thumb:hover {
  transform: scale(1.2);
}

.frame-counter {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  min-width: 60px;
  text-align: right;
}

/* Ключевые кадры */
.keyframes-scroll-container {
  width: 100%;
  overflow-x: auto;
  padding-bottom: 8px;
}

.keyframes-row {
  display: flex;
  gap: 12px;
  width: max-content;
}

.keyframe-card {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 6px;
  text-align: center;
  cursor: pointer;
  transition: var(--transition-smooth);
}

.keyframe-card:hover, .keyframe-card.active {
  background: rgba(99, 102, 241, 0.1);
  border-color: var(--color-primary);
}

.keyframe-thumb-box {
  width: 100%;
  height: 48px;
  background: #000;
  border-radius: 6px;
  margin-bottom: 8px;
  overflow: hidden;
}

.keyframe-canvas {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.keyframe-label {
  font-size: 10px;
  font-weight: 700;
  color: #fff;
  display: block;
}

.keyframe-time {
  font-size: 9px;
  color: var(--text-muted);
}

/* ==========================================
   Скроллируемый Контент (Scrollable Content)
   ========================================== */
.scrollable-content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
  overflow-y: auto;
  padding-right: 8px;
}

.tab-pane {
  display: none;
  animation: fadeIn 0.4s ease forwards;
}

.tab-pane.active {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Биометрия (Metrics Grid) */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.bio-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 20px;
  text-align: center;
  transition: var(--transition-fast);
}

.bio-card:hover {
  background: rgba(255, 255, 255, 0.03);
  border-color: rgba(99, 102, 241, 0.3);
}

.bio-card-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.bio-card-val {
  font-family: var(--font-mono);
  font-size: 26px;
  font-weight: 800;
  color: #fff;
  text-shadow: none;
}

.bio-card-sub {
  font-size: 10px;
  color: var(--text-muted);
  margin-top: 6px;
}

/* Графики и Canvas */
.chart-container {
  width: 100%;
  height: 140px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
  border: 1px inset rgba(255, 255, 255, 0.05);
  overflow: hidden;
  padding: 10px;
}

.comparative-skeletons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.skeleton-box {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.skeleton-box-title {
  font-size: 10px;
  font-weight: 800;
  text-transform: uppercase;
  padding: 6px 12px;
  border-radius: 6px;
  margin-bottom: 12px;
  letter-spacing: 0.5px;
}

.skeleton-box.incorrect .skeleton-box-title {
  background: rgba(239, 68, 68, 0.1);
  color: var(--color-danger);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.skeleton-box.correct .skeleton-box-title {
  background: rgba(16, 185, 129, 0.1);
  color: var(--color-success);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.compare-canvas {
  width: 100%;
  aspect-ratio: 16/9;
  background: #000;
  border-radius: 6px;
  box-shadow: none;
}

.correctness-desc {
  font-size: 11px;
  color: var(--text-muted);
  text-align: center;
  margin-top: 12px;
  line-height: 1.5;
}

/* ==========================================
   Судейство и GOE (Judging)
   ========================================== */
.system-tabs {
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  padding: 4px;
  border-radius: 6px;
  margin-bottom: 20px;
}

.system-tab {
  flex: 1;
  text-align: center;
  padding: 10px;
  font-size: 12px;
  font-weight: 700;
  color: var(--text-muted);
  border-radius: 6px;
  cursor: pointer;
  transition: var(--transition-fast);
  background: transparent;
  border: none;
}

.system-tab.active {
  background: rgba(99, 102, 241, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #fff;
  box-shadow: none;
}

.judging-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.judging-level-name {
  font-family: var(--font-mono);
  font-size: 22px;
  font-weight: 800;
  color: #fff;
}

.judging-level-sub {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.judging-badge {
  background: rgba(255, 255, 255, 0.03);
  color: #fff;
  font-family: var(--font-mono);
  font-weight: 800;
  font-size: 16px;
  width: 48px;
  height: 48px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: none;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.level-progress-tracker {
  margin: 24px 20px 40px 20px;
  position: relative;
}

.level-line-container {
  height: 6px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  position: relative;
}

.level-line-fill {
  height: 100%;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 6px;
  box-shadow: none;
  transition: width 0.5s ease-out;
}

.level-markers {
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  transform: translateY(-50%);
  display: flex;
  justify-content: space-between;
  z-index: 10;
}

.level-dot-marker {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #1a1c23;
  border: 2px solid rgba(255, 255, 255, 0.2);
  position: relative;
  z-index: 2;
  box-sizing: border-box;
}

.level-dot-marker.passed {
  background: #8b5cf6;
  border-color: #fff;
  box-shadow: none;
}

.level-dot-marker.active {
  background: #6366f1;
  border-color: #fff;
  transform: scale(1.3);
  box-shadow: none;
}

.level-marker-label {
  position: absolute;
  top: 18px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  font-weight: 600;
  color: var(--text-muted);
  white-space: nowrap;
}

.level-dot-marker.active .level-marker-label {
  color: #fff;
}

.judging-next-goals {
  background: rgba(0, 0, 0, 0.2);
  border: 1px dashed rgba(255, 255, 255, 0.15);
  padding: 16px;
  border-radius: 6px;
  margin-top: 16px;
}

.goal-title {
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 10px;
}

.goals-list {
  list-style: none;
  font-size: 11px;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.goals-list li::before {
  content: "→";
  color: var(--color-accent);
  margin-right: 8px;
  font-weight: 800;
}

.judging-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 16px;
  line-height: 1.5;
}

/* GOE Simulator */
.goe-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.goe-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-main);
}

.goe-ctrls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.goe-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #fff;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-fast);
}

.goe-btn:hover {
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.goe-val {
  font-family: var(--font-mono);
  font-size: 15px;
  font-weight: 800;
  width: 24px;
  text-align: center;
}

.goe-val.positive { color: var(--color-success); }
.goe-val.negative { color: var(--color-danger); }

.goe-result-box {
  margin-top: 24px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 20px;
  border-radius: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: none;
}

.result-label {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.result-formula {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
}

.result-score {
  font-family: var(--font-mono);
  font-size: 28px;
  font-weight: 800;
  text-shadow: none;
}

/* ==========================================
   Ошибки и таблица (Issues)
   ========================================== */
.issues-table {
  width: 100%;
  border-collapse: collapse;
}

.issues-table th {
  text-align: left;
  padding: 12px;
  font-size: 10px;
  font-weight: 800;
  color: var(--text-muted);
  text-transform: uppercase;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  letter-spacing: 0.5px;
}

.issues-table td {
  padding: 16px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.issue-row {
  cursor: pointer;
  transition: background 0.2s;
}

.issue-row:hover {
  background: rgba(255, 255, 255, 0.02);
}

.issue-label-cell {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.issue-dot {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
  margin-top: 2px;
  box-shadow: none;
}

.dot-ok { background: var(--color-success); }
.dot-warning { background: var(--color-warning); }
.dot-danger { background: var(--color-danger); }
.dot-perfect { background: var(--color-primary); }

.issue-desc-text {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 6px;
  line-height: 1.4;
}

.issue-moment {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
}

.impact-badge {
  font-family: var(--font-mono);
  font-weight: 800;
  font-size: 13px;
  padding: 4px 8px;
  border-radius: 6px;
  background: rgba(0,0,0,0.2);
}

.impact-positive-high { color: var(--color-success); border: 1px solid rgba(255, 255, 255, 0.08); }
.impact-negative-high { color: var(--color-danger); border: 1px solid rgba(255, 255, 255, 0.08); }
.impact-positive { color: var(--color-success); }
.impact-negative { color: var(--color-danger); }

/* Рекомендации Тренера */
.coach-summary-box {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 20px;
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.coach-avatar-badge {
  background: rgba(255, 255, 255, 0.03);
  color: #fff;
  font-weight: 800;
  font-size: 11px;
  padding: 6px 12px;
  border-radius: 6px;
  box-shadow: none;
  align-self: flex-start;
}

.coach-paragraph {
  font-size: 13px;
  line-height: 1.6;
  color: #e2e8f0;
}

.practice-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.practice-item {
  display: flex;
  gap: 16px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 16px;
  border-radius: 6px;
  transition: var(--transition-fast);
}

.practice-item:hover {
  background: rgba(255, 255, 255, 0.03);
  border-color: rgba(255, 255, 255, 0.1);
}

.practice-num {
  font-family: var(--font-mono);
  font-size: 24px;
  font-weight: 800;
  color: var(--color-primary);
  opacity: 0.8;
}

.practice-title {
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 4px;
}

.practice-reps {
  font-size: 11px;
  color: var(--color-accent);
  background: rgba(0, 240, 255, 0.1);
  padding: 2px 8px;
  border-radius: 6px;
  display: inline-block;
  margin-bottom: 8px;
}

.practice-focus {
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.5;
}

/* ==========================================
   Профиль и Чат (Profile & Chat)
   ========================================== */
.user-profile {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.user-avatar {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  box-shadow: none;
}

.user-name {
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.user-plan {
  font-size: 12px;
  color: var(--color-warning);
  font-weight: 600;
  margin-top: 4px;
}

.limit-tracker {
  background: rgba(0, 0, 0, 0.2);
  padding: 16px;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.limit-text {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 12px;
}

.limit-progress-bar {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  margin-bottom: 16px;
  overflow: hidden;
}

.limit-progress {
  height: 100%;
  background: var(--color-primary);
  border-radius: 6px;
}

.upgrade-button {
  width: 100%;
  background: rgba(255, 255, 255, 0.03);
  color: #000;
  border: none;
  padding: 12px;
  font-size: 13px;
  font-weight: 800;
  border-radius: 6px;
  cursor: pointer;
  box-shadow: none;
  transition: transform var(--transition-fast);
}

.upgrade-button:hover {
  transform: translateY(-2px);
}

.chat-card-full {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  border: 1px inset rgba(255, 255, 255, 0.02);
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 16px;
}

.chat-msg {
  max-width: 85%;
  display: flex;
  flex-direction: column;
  animation: fadeIn 0.3s ease forwards;
}

.chat-msg.user {
  align-self: flex-end;
}

.chat-msg.coach {
  align-self: flex-start;
}

.chat-msg-sender {
  font-size: 10px;
  color: var(--text-muted);
  margin-bottom: 4px;
  margin-left: 8px;
}

.chat-msg-text {
  padding: 12px 16px;
  font-size: 13px;
  line-height: 1.5;
  border-radius: 6px;
  color: #fff;
}

.chat-msg.user .chat-msg-text {
  background: var(--color-primary);
  border-bottom-right-radius: 4px;
}

.chat-msg.coach .chat-msg-text {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-bottom-left-radius: 4px;
}

.chat-suggested-questions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.suggested-q-btn {
  background: rgba(0, 240, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--color-accent);
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-fast);
}

.suggested-q-btn:hover {
  background: rgba(0, 240, 255, 0.15);
}

.chat-input-area {
  display: flex;
  gap: 12px;
}

.file-upload-btn-wrapper {
  position: relative;
  overflow: hidden;
  display: inline-block;
}

.file-upload-visual-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #fff;
  padding: 0 16px;
  height: 44px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.file-upload-input {
  position: absolute;
  left: 0;
  top: 0;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.chat-input-wrapper {
  flex: 1;
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  overflow: hidden;
}

.chat-text-input {
  flex: 1;
  background: transparent;
  border: none;
  padding: 12px 16px;
  color: #fff;
  font-size: 13px;
  outline: none;
}

.chat-text-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.chat-send-btn {
  background: var(--color-primary);
  border: none;
  color: #fff;
  padding: 0 20px;
  font-weight: 700;
  cursor: pointer;
  transition: var(--transition-fast);
}

.chat-send-btn:hover {
  background: var(--color-primary-hover);
}

/* ==========================================
   Адаптивность (Mobile & Tablets)
   ========================================== */
.bottom-nav-bar {
  display: none;
}

@media (max-width: 1024px) {
  .main-content {
    flex-direction: column;
    padding: 16px;
    gap: 16px;
    height: calc(100vh - 140px); /* Учитываем шапку и мобильное меню */
  }

  .desktop-tabs {
    display: none;
  }

  .bottom-nav-bar {
    display: flex;
    justify-content: space-around;
    align-items: center;
    background: rgba(10, 11, 16, 0.9);
    
    border-top: 1px solid var(--border-color);
    padding: 12px 8px 24px; /* Учет безопасных зон iOS */
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    z-index: 100;
  }

  .nav-tab-btn {
    background: transparent;
    border: none;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    color: var(--text-muted);
    cursor: pointer;
    transition: var(--transition-fast);
  }

  .nav-tab-btn.active {
    color: var(--color-primary);
  }

  .nav-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
  }

  .nav-text {
    font-size: 10px;
    font-weight: 600;
  }

  .metrics-grid, .comparative-skeletons {
    grid-template-columns: 1fr;
  }
}

/* Модальное окно чата (Резерв) */
.chat-dialog-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.chat-dialog-window {
  width: 100%;
  max-width: 500px;
  background: var(--bg-main);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.chat-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid var(--border-color);
}

.chat-dialog-title {
  font-weight: 700;
  font-size: 15px;
}

.chat-dialog-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 24px;
  cursor: pointer;
}

.chat-dialog-close:hover {
  color: #fff;
}
`
}

