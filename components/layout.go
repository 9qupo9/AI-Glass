package components

// HeaderBlock
type HeaderBlock struct{}

func (h *HeaderBlock) HTML() string {
	return `
	<header class="top-bar">
      <div class="brand-logo-container">
        <div class="brand-name">SkateEye AI</div>
        <div id="confidence-indicator" class="confidence-badge conf-high">
          <span class="conf-dot"></span>
          <span id="confidence-text">High Precision</span>
        </div>
      </div>

      <div style="flex: 1; display: flex; justify-content: center; align-items: center;">
          <div style="color: #EF4444; font-weight: 800; font-size: 16px; text-transform: uppercase; text-decoration: underline; letter-spacing: 0.5px;">
              MAX 5 SECONDS VIDEO LIMIT. TRIM BEFORE UPLOADING.
          </div>
      </div>      <div class="user-profile" style="margin-bottom: 0;">
        <div class="user-avatar" style="width: 36px; height: 36px; font-size: 14px;">US</div>
        <div>
          <div class="user-name" style="font-size: 14px;">User (Guest)</div>
          <div class="user-plan" style="font-size: 11px;">Basic Access</div>
        </div>
      </div>
    </header>
	`
}

// VideoPanelBlock
type VideoPanelBlock struct{}

func (v *VideoPanelBlock) HTML() string {
	return `
	<div class="video-panel-wrapper">
		<div class="panel-card" style="width: 140px; display: flex; flex-direction: column; padding: 16px; border-radius: 6px;">
			<div class="panel-title" style="font-size: 13px; margin-bottom: 16px;">Frames</div>
			<div class="keyframes-scroll-container" style="flex: 1; overflow-y: auto; overflow-x: hidden; padding-right: 4px; padding-bottom: 0;">
				<div id="keyframes-list" class="keyframes-column" style="display: flex; flex-direction: column; gap: 12px; width: 100%;">
					<!-- Loaded via JS -->
				</div>
			</div>
		</div>
		<div class="panel-card" style="flex: 1; border-radius: 6px;">
			<div class="panel-title" style="display: flex; justify-content: space-between; align-items: center;">
				3D Biomechanics Analysis
				<label for="video-file-input" class="action-btn upload-btn" style="cursor: pointer; font-size: 14px; padding: 6px 12px; border-radius: 6px; background: rgba(255, 255, 255, 0.1); color: #fff; display: flex; align-items: center; gap: 8px;">
					<span>Upload Video</span>
				</label>
				<input type="file" id="video-file-input" accept="video/mp4,video/webm" style="display: none;">
			</div>
			
			<div class="video-container">
				<video id="video-player" class="video-player" src="" muted playsinline preload="metadata"></video>
				<canvas id="video-canvas" class="video-canvas" width="800" height="450"></canvas>
				
				<div class="video-overlay-ui">
					<div id="video-phase-overlay" class="overlay-tag">Entry</div>
					<div id="video-time-overlay" class="overlay-tag">0:00</div>
				</div>
			</div>

			<div class="video-controls">
				<div class="scrubber-row">
					<button id="btn-play-pause" class="play-pause-btn">
						<svg id="icon-play" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M8 5v14l11-7z"/></svg>
						<svg id="icon-pause" viewBox="0 0 24 24" width="20" height="20" fill="currentColor" style="display: none;"><path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/></svg>
					</button>
					<div class="timeline-slider-container">
						<input type="range" id="video-scrubber" class="timeline-slider" min="0" max="100" value="0">
					</div>
					<div id="frame-counter" class="frame-counter">Frame: 0</div>
				</div>
			</div>

			<div id="header-compressed-stats" style="display: flex; flex-direction: column; background: rgba(255,255,255,0.03); padding: 16px 32px; border-radius: 6px; border: 1px solid rgba(255,255,255,0.08); margin-top: 16px;">
				<div style="display: flex; justify-content: center; gap: 40px; align-items: center;">
					<div style="display:flex; flex-direction:column; align-items:center;">
						<span style="font-size: 11px; color: #94A3B8; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px;">Base Score</span>
						<span id="hdr-base" style="font-family: var(--font-mono); font-size: 18px; font-weight: 700; color: #F8FAFC; margin-top: 4px;">-- pts</span>
					</div>
					<div style="width: 1px; height: 32px; background: rgba(255,255,255,0.1);"></div>
					<div style="display:flex; flex-direction:column; align-items:center;">
						<span style="font-size: 11px; color: #94A3B8; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px;">Penalty (GOE)</span>
						<span id="hdr-goe" style="font-family: var(--font-mono); font-size: 18px; font-weight: 700; color: #EF4444; margin-top: 4px;">-- pts</span>
					</div>
					<div style="width: 1px; height: 32px; background: rgba(255,255,255,0.1);"></div>
					<div style="display:flex; flex-direction:column; align-items:center;">
						<span style="font-size: 11px; color: #94A3B8; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px;">Final Score</span>
						<span id="hdr-score" style="font-family: var(--font-mono); font-size: 18px; font-weight: 700; color: #10B981; margin-top: 4px;">-- pts</span>
					</div>
				</div>
			</div>

			<!-- COACH SUMMARY -->
			<div class="panel-section-title" style="margin-top: 24px; font-size: 14px; font-weight: 700; color: #F8FAFC; text-transform: uppercase; letter-spacing: 0.5px;">Coach Summary</div>
			<div class="coach-summary-container" style="background: rgba(59, 130, 246, 0.05); border: 1px solid rgba(59, 130, 246, 0.1); border-radius: 8px; margin-top: 12px; padding: 16px; display: flex; gap: 16px; flex-direction: column;">
				<div style="display: flex; gap: 16px;">
					<div class="ai-icon" style="width: 32px; height: 32px; background: rgba(59, 130, 246, 0.1); color: #60A5FA; border-radius: 6px; display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 12px;">AI</div>
					<div id="coach-summary-text" style="color: #E2E8F0; font-size: 14px; line-height: 1.5; flex: 1;">
						Awaiting analysis...
					</div>
				</div>
			</div>

			<!-- DETECTED ISSUES -->
			<div class="panel-section-title" style="margin-top: 24px; font-size: 14px; font-weight: 700; color: #F8FAFC; text-transform: uppercase; letter-spacing: 0.5px;">Detected Issues</div>
			<div class="issues-container" style="background: rgba(255, 255, 255, 0.02); border: 1px solid rgba(255, 255, 255, 0.05); border-radius: 8px; margin-top: 12px; padding: 16px;">
				<table class="issues-table" style="width: 100%; border-collapse: collapse; text-align: left;">
					<thead>
						<tr style="border-bottom: 1px solid rgba(255,255,255,0.1); color: #94A3B8; font-size: 11px; text-transform: uppercase;">
							<th style="padding-bottom: 8px; width: 60%;">Issue</th>
							<th style="padding-bottom: 8px; width: 20%;">Moment</th>
							<th style="padding-bottom: 8px; width: 20%;">ISU Impact</th>
						</tr>
					</thead>
					<tbody id="detected-issues-tbody">
						<!-- Loaded via JS -->
						<tr><td colspan="3" style="padding-top: 12px; color: #64748B; font-size: 13px;">Awaiting analysis...</td></tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
	`
}

// SidebarBlock
type SidebarBlock struct{}

func (s *SidebarBlock) HTML() string {
	return `
	<div class="scrollable-content-wrapper">
		<div id="content-dashboard" class="tab-pane active">
			<div class="panel-card" style="background: transparent; border: none; box-shadow: none; padding: 0; margin-top: 24px;">
				<div class="panel-title" style="display: flex; justify-content: space-between; align-items: center;">
					<span>AI Analysis & Diagnostics</span>
				</div>
				<div id="diagnostics-container" class="diagnostics-grid" style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-top: 16px;">
					


					<!-- ISU JUDGING SYSTEM -->
					<div class="diag-card" style="grid-column: span 1; display: flex; flex-direction: column;">
						<div class="diag-header" style="font-size: 13px; font-weight: 700; color: #F8FAFC; text-transform: uppercase; margin-bottom: 12px;">ISU Judging System</div>
						<div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px;">
							<div>
								<div style="font-size: 11px; color: #64748B; margin-bottom: 4px;">Current level</div>
								<div id="isu-current-level" style="font-size: 18px; font-weight: 700; color: #F8FAFC;">--</div>
							</div>
							<div id="isu-badge" style="width: 48px; height: 48px; border-radius: 24px; background: rgba(16, 185, 129, 0.1); border: 2px solid #10B981; color: #10B981; display: flex; flex-direction: column; align-items: center; justify-content: center; font-size: 14px; font-weight: 800; text-align: center; line-height: 1.1;">
								BN<br><span style="font-size: 9px; font-weight: 600;">Review</span>
							</div>
						</div>
						<div style="font-size: 12px; font-weight: 700; color: #E2E8F0; margin-bottom: 6px;">Next goal: <span id="isu-next-goal">--</span></div>
						<div style="font-size: 12px; color: #94A3B8; margin-bottom: 6px;">To reach this level:</div>
						<ul id="isu-requirements" style="font-size: 12px; color: #94A3B8; margin: 0 0 12px 0; padding-left: 16px; line-height: 1.5;">
							<li>Waiting for analysis...</li>
						</ul>
						<div id="isu-explanation" style="font-size: 12px; color: #64748B; line-height: 1.5; margin-bottom: 16px; flex: 1;">--</div>
						<div style="text-align: center;">
							<a href="#" style="color: #60A5FA; font-size: 13px; text-decoration: none; font-weight: 600;">Learn more about levels</a>
						</div>
					</div>

					<!-- WHAT TO PRACTICE -->
					<div class="diag-card" style="grid-column: span 1; display: flex; flex-direction: column;">
						<div class="diag-header" style="font-size: 13px; font-weight: 700; color: #F8FAFC; text-transform: uppercase; margin-bottom: 16px;">What to practice</div>
						<div id="practice-list" style="display: flex; flex-direction: column; gap: 16px; flex: 1;">
							<div style="color: #64748B; font-size: 13px;">Awaiting analysis...</div>
						</div>
						<div style="margin-top: 16px;">
							<a href="#" style="color: #60A5FA; font-size: 13px; text-decoration: none; font-weight: 600;">Go to training plan -></a>
						</div>
					</div>

					<!-- US FIGURE SKATING SYSTEM -->
					<div class="diag-card" style="grid-column: span 1; display: flex; flex-direction: column;">
						<div class="diag-header" style="font-size: 13px; font-weight: 700; color: #F8FAFC; text-transform: uppercase; margin-bottom: 12px;">US Figure Skating System</div>
						<div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px;">
							<div>
								<div style="font-size: 11px; color: #64748B; margin-bottom: 4px;">Current level</div>
								<div id="usfsa-current-level" style="font-size: 18px; font-weight: 700; color: #F8FAFC;">--</div>
							</div>
							<div id="usfsa-badge" style="width: 48px; height: 48px; border-radius: 24px; background: rgba(59, 130, 246, 0.1); border: 2px solid #3B82F6; color: #3B82F6; display: flex; flex-direction: column; align-items: center; justify-content: center; font-size: 16px; font-weight: 800; text-align: center; line-height: 1.1;">
								4<br><span style="font-size: 9px; font-weight: 600;">Good</span>
							</div>
						</div>
						<div style="font-size: 12px; font-weight: 700; color: #E2E8F0; margin-bottom: 6px;">Next goal: <span id="usfsa-next-goal">--</span></div>
						<div id="usfsa-explanation" style="font-size: 12px; color: #94A3B8; line-height: 1.5; margin-bottom: 16px; flex: 1;">--</div>
					</div>

					<!-- Biometrics Stats -->
					<div class="diag-card" style="grid-column: span 1; display: flex; flex-direction: column; gap: 16px;">
						<div class="diag-header" style="font-size: 13px; font-weight: 700; color: #F8FAFC; text-transform: uppercase; margin-bottom: 4px;">Biomechanics</div>
						
						<!-- Air Time Block -->
						<div style="display: flex; align-items: center; gap: 16px;">
							<div style="width: 48px; height: 48px; border-radius: 12px; background: rgba(56, 189, 248, 0.1); display: flex; align-items: center; justify-content: center; font-size: 24px;">⏱️</div>
							<div>
								<div style="font-size: 11px; text-transform: uppercase; color: #71717A; font-weight: 700; margin-bottom: 4px;">Air Time</div>
								<div id="air-time-val" style="font-family: var(--font-mono); font-size: 22px; font-weight: 800; color: #F8FAFC; letter-spacing: 0.5px;">-- s</div>
							</div>
						</div>

						<!-- Height Block -->
						<div style="display: flex; align-items: center; gap: 16px;">
							<div style="width: 48px; height: 48px; border-radius: 12px; background: rgba(16, 185, 129, 0.1); display: flex; align-items: center; justify-content: center; font-size: 24px;">📏</div>
							<div>
								<div style="font-size: 11px; text-transform: uppercase; color: #71717A; font-weight: 700; margin-bottom: 4px;">Max Height</div>
								<div id="height-val" style="font-family: var(--font-mono); font-size: 22px; font-weight: 800; color: #F8FAFC; letter-spacing: 0.5px;">-- m</div>
							</div>
						</div>

						<!-- Rotations Block -->
						<div style="display: flex; align-items: center; gap: 16px;">
							<div style="position: relative; width: 48px; height: 48px;">
								<svg width="48" height="48" viewBox="0 0 48 48" style="transform: rotate(-90deg);">
									<circle cx="24" cy="24" r="20" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="4"></circle>
									<circle id="rotations-ring" cx="24" cy="24" r="20" fill="none" stroke="#A855F7" stroke-width="4" stroke-dasharray="125.6" stroke-dashoffset="125.6" stroke-linecap="round" style="transition: stroke-dashoffset 1s ease-out;"></circle>
								</svg>
								<div id="rotations-ring-text" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 800; color: #F8FAFC; font-family: var(--font-mono);">0</div>
							</div>
							<div>
								<div style="font-size: 11px; text-transform: uppercase; color: #71717A; font-weight: 700; margin-bottom: 4px;">Rotations</div>
								<div id="rotations-val" style="font-family: var(--font-mono); font-size: 22px; font-weight: 800; color: #F8FAFC; letter-spacing: 0.5px;">--</div>
							</div>
						</div>
					</div>
				</div>
			</div>

		</div>
	</div>
	`
}

// MobileNavBlock
type MobileNavBlock struct{}

func (m *MobileNavBlock) HTML() string {
	return `
	<!-- Bottom Nav Bar Removed completely per user request -->
	`
}

// MainLayoutBlock
type MainLayoutBlock struct{}

func (m *MainLayoutBlock) HTML() string {
	videoPanel := &VideoPanelBlock{}
	sidebar := &SidebarBlock{}

	return `
	<div style="display: flex; flex-direction: column; height: 100%; overflow: hidden;">
		<main class="main-content" style="flex: 1; padding-bottom: 10px;">
			` + videoPanel.HTML() + `
			` + sidebar.HTML() + `
		</main>
	</div>
	`
}
