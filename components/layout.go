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
		<div class="panel-card" style="width: 140px; display: flex; flex-direction: column; padding: 16px;">
			<div class="panel-title" style="font-size: 13px; margin-bottom: 16px;">Frames</div>
			<div class="keyframes-scroll-container" style="flex: 1; overflow-y: auto; overflow-x: hidden; padding-right: 4px; padding-bottom: 0;">
				<div id="keyframes-list" class="keyframes-column" style="display: flex; flex-direction: column; gap: 12px; width: 100%;">
					<!-- Loaded via JS -->
				</div>
			</div>
		</div>
		<div class="panel-card" style="flex: 1;">
			<div class="panel-title" style="display: flex; justify-content: space-between; align-items: center;">
				3D Biomechanics Analysis
				<label for="video-file-input" class="action-btn upload-btn" style="cursor: pointer; font-size: 14px; padding: 6px 12px; border-radius: 4px; background: rgba(255, 255, 255, 0.1); color: #fff; display: flex; align-items: center; gap: 8px;">
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

			<div id="header-compressed-stats" style="display: flex; flex-direction: column; background: rgba(30,41,59,0.5); padding: 16px 32px; border-radius: 12px; border: 1px solid rgba(255,255,255,0.05); margin-top: 16px;">
				<div style="display: flex; justify-content: center; gap: 40px; align-items: center;">
					<div style="display:flex; flex-direction:column; align-items:center;">
						<span style="font-size: 11px; color: #94A3B8; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px;">Base Score</span>
						<span id="hdr-base" style="font-size: 18px; font-weight: 700; color: #F8FAFC; margin-top: 4px;">-- pts</span>
					</div>
					<div style="width: 1px; height: 32px; background: rgba(255,255,255,0.1);"></div>
					<div style="display:flex; flex-direction:column; align-items:center;">
						<span style="font-size: 11px; color: #94A3B8; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px;">Penalty (GOE)</span>
						<span id="hdr-goe" style="font-size: 18px; font-weight: 700; color: #EF4444; margin-top: 4px;">-- pts</span>
					</div>
					<div style="width: 1px; height: 32px; background: rgba(255,255,255,0.1);"></div>
					<div style="display:flex; flex-direction:column; align-items:center;">
						<span style="font-size: 11px; color: #94A3B8; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px;">Final Score</span>
						<span id="hdr-score" style="font-size: 18px; font-weight: 700; color: #10B981; margin-top: 4px;">-- pts</span>
					</div>
				</div>
				<div id="hdr-reason" style="display: none; margin-top: 16px; font-size: 13px; color: #FCA5A5; text-align: center; background: rgba(239, 68, 68, 0.1); padding: 8px; border-radius: 6px; border: 1px solid rgba(239, 68, 68, 0.2);">
					--
				</div>
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
				<div id="diagnostics-container" class="diagnostics-grid">
					<div class="diag-card cause">
						<div class="diag-header"><span class="diag-icon">🔍</span><span class="diag-title">Why It Failed (or Succeeded)</span></div>
						<div class="diag-body" id="diag-cause-text">Awaiting analysis...</div>
					</div>
					<div class="diag-card fix">
						<div class="diag-header"><span class="diag-icon">🛠️</span><span class="diag-title">How to Fix It</span></div>
						<div class="diag-body" id="diag-fix-text">Awaiting analysis...</div>
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
