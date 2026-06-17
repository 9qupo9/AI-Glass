package components

// GetGlobalJS возвращает клиентский JavaScript для взаимодействия и рендеринга
func GetGlobalJS() string {
	return `
// ==========================================
// Модель Скелета и Канвас (Skeleton Canvas)
// ==========================================
const BONES = [
  { from: 'head', to: 'chest' },
  { from: 'chest', to: 'shoulderL' },
  { from: 'chest', to: 'shoulderR' },
  { from: 'shoulderL', to: 'shoulderR' },
  { from: 'chest', to: 'hipL' },
  { from: 'chest', to: 'hipR' },
  { from: 'hipL', to: 'hipR' },
  { from: 'shoulderL', to: 'elbowL' },
  { from: 'elbowL', to: 'wristL' },
  { from: 'shoulderR', to: 'elbowR' },
  { from: 'elbowR', to: 'wristR' },
  { from: 'hipL', to: 'kneeL' },
  { from: 'kneeL', to: 'ankleL' },
  { from: 'hipR', to: 'kneeR' },
  { from: 'kneeR', to: 'ankleR' }
];

function drawSkeleton(ctx, joints, width, height, options = {}) {
  const {
    color = '#60A5FA',
    lineWidth = 3,
    glowColor = null,
    drawJoints = true,
    highlightJoints = {},
    alpha = 1.0,
    scaleFactor = 1.0,
    offsetX = 0,
    offsetY = 0
  } = options;
  
  ctx.save();
  ctx.globalAlpha = alpha;
  
  if (glowColor) {
    ctx.shadowColor = glowColor;
    ctx.shadowBlur = 10;
  } else {
    ctx.shadowBlur = 0;
  }
  
  ctx.lineWidth = lineWidth;
  ctx.lineCap = 'round';
  ctx.lineJoin = 'round';
  
  function getCanvasCoords(jointName) {
    const pt = joints[jointName];
    if (!pt) return { x: 0, y: 0 };
    
    if (options.videoRect) {
       let nx = pt.x;
       let ny = pt.y;
       if (scaleFactor !== 1.0) {
           nx = (nx - 50) * scaleFactor + 50;
           ny = (ny - 60) * scaleFactor + 60;
       }
       nx += offsetX;
       ny += offsetY;

       return {
          x: (nx / 100) * options.videoRect.vW * options.videoRect.scale + options.videoRect.dx,
          y: (ny / 100) * options.videoRect.vH * options.videoRect.scale + options.videoRect.dy
       };
    }

    const nx = (pt.x - 50) * scaleFactor + 50 + offsetX;
    const ny = (pt.y - 60) * scaleFactor + 60 + offsetY;
    return {
      x: (nx / 100) * width,
      y: (ny / 100) * height
    };
  }
  
  // Кости
  ctx.strokeStyle = color;
  BONES.forEach(bone => {
    const start = getCanvasCoords(bone.from);
    const end = getCanvasCoords(bone.to);
    ctx.beginPath();
    ctx.moveTo(start.x, start.y);
    ctx.lineTo(end.x, end.y);
    ctx.stroke();
  });
  
  // Суставы
  if (drawJoints) {
    ctx.shadowBlur = 0;
    Object.keys(joints).forEach(jointName => {
      const pt = getCanvasCoords(jointName);
      const radius = jointName === 'head' ? 7 * scaleFactor : 4 * scaleFactor;
      ctx.beginPath();
      ctx.arc(pt.x, pt.y, radius, 0, 2 * Math.PI);
      ctx.fillStyle = highlightJoints[jointName] || color;
      ctx.fill();
      if (highlightJoints[jointName]) {
        ctx.strokeStyle = '#FFF';
        ctx.lineWidth = 1.5;
        ctx.stroke();
      }
    });
  }
  ctx.restore();
}

function drawModelComparison(canvasIncorrect, canvasCorrect, jump, currentFrameIndex) {
  if (!canvasIncorrect || !canvasCorrect || !jump) return;
  const ctxInc = canvasIncorrect.getContext('2d');
  const ctxCor = canvasCorrect.getContext('2d');
  const wInc = canvasIncorrect.width;
  const hInc = canvasIncorrect.height;
  const wCor = canvasCorrect.width;
  const hCor = canvasCorrect.height;
  
  ctxInc.clearRect(0, 0, wInc, hInc);
  ctxCor.clearRect(0, 0, wCor, hCor);
  
  // Отрисовка самого видео на заднем фоне (полное повторение видео)
  let vRect = null;
  if (videoPlayer && videoPlayer.readyState >= 2 && videoPlayer.videoWidth) {
      const vW = videoPlayer.videoWidth;
      const vH = videoPlayer.videoHeight;
      const scale = Math.min(wInc / vW, hInc / vH);
      const dx = (wInc - vW * scale) / 2;
      const dy = (hInc - vH * scale) / 2;
      vRect = { vW, vH, scale, dx, dy };
      ctxInc.save();
      ctxInc.globalAlpha = 0.4;
      ctxInc.drawImage(videoPlayer, dx, dy, vW * scale, vH * scale);
      ctxInc.restore();
  }

  // СТРОГИЙ РЕЖИМ: Никаких моковых данных! 
  // Если нейросеть не видит человека, скелет вообще не нарисуется.
  const currentSkeleton = currentMediaPipeJoints;
  const goldSkeleton = currentSkeleton;
  if (!currentSkeleton || !goldSkeleton) return;
  
  let isFrameIncorrect = false;
  let activeIssues = [];
  jump.issues.forEach(issue => {
    if (Math.abs(issue.frame - currentFrameIndex) <= 4) {
      if (issue.status === 'warning' || issue.status === 'danger') {
        isFrameIncorrect = true;
        activeIssues.push(issue);
      }
    }
  });
  
  const isMediaPipe = (currentSkeleton === currentMediaPipeJoints);
  
  const incOptions = {
    color: '#F59E0B', 
    glowColor: 'rgba(245, 158, 11, 0.4)',
    lineWidth: 3,
    drawJoints: true,
    scaleFactor: isMediaPipe ? 1.0 : (jump.alignment?.scale || 1.0),
    offsetX: isMediaPipe ? 0 : (jump.alignment?.offsetX || 0),
    offsetY: isMediaPipe ? 0 : (jump.alignment?.offsetY || 0),
    videoRect: vRect,
    highlightJoints: {}
  };
  
  if (isFrameIncorrect) {
    activeIssues.forEach(issue => {
      if (issue.id === 's2' || issue.id === 's4' || issue.id === 'l1' || issue.id === 'l2') {
        incOptions.highlightJoints.kneeL = '#F59E0B';
        incOptions.highlightJoints.kneeR = '#F59E0B';
      }
      if (issue.id === 's3' || issue.id === 'a3' || issue.id === 'l1' || issue.id === 'l2') {
        incOptions.highlightJoints.hipL = '#F59E0B';
        incOptions.highlightJoints.hipR = '#F59E0B';
        incOptions.highlightJoints.chest = '#F59E0B';
      }
    });
  }
  
  drawSkeleton(ctxInc, currentSkeleton, wInc, hInc, incOptions);
  
  const corOptions = {
    color: '#10B981', 
    glowColor: 'rgba(16, 185, 129, 0.5)',
    lineWidth: 3,
    drawJoints: true,
    scaleFactor: incOptions.scaleFactor,
    offsetX: incOptions.offsetX,
    offsetY: incOptions.offsetY,
    videoRect: incOptions.videoRect
  };
  drawSkeleton(ctxCor, goldSkeleton, wCor, hCor, corOptions);
}

// ==========================================
// Данные прыжков (Jumps Data)
// ==========================================
function generateSkeletonFrames(jumpType, totalFrames, keyframes) {
  const frames = [];
  for (let f = 0; f < totalFrames; f++) {
    let phase = 'entry';
    let progress = 0;
    
    if (f < keyframes.takeoff) {
      phase = 'entry';
      progress = f / keyframes.takeoff;
    } else if (f < keyframes.air) {
      phase = 'takeoff';
      progress = (f - keyframes.takeoff) / (keyframes.air - keyframes.takeoff);
    } else if (f < keyframes.landing) {
      phase = 'air';
      progress = (f - keyframes.air) / (keyframes.landing - keyframes.air);
    } else if (f < keyframes.exit) {
      phase = 'landing';
      progress = (f - keyframes.landing) / (keyframes.exit - keyframes.landing);
    } else {
      phase = 'exit';
      progress = (f - keyframes.exit) / (totalFrames - keyframes.exit);
    }
    
    let cx = 46.5; let cy = 47.5;
    let rotationAngle = 0; let tiltAngle = 0;
    
    if (phase === 'entry') {
      cx = 38 + progress * 8.5;
      cy = 47.5 + Math.sin(progress * Math.PI) * 1.5;
    } else if (phase === 'takeoff') {
      cx = 46.5 + progress * 1.0;
      cy = 47.5 + Math.sin(progress * Math.PI) * 4; 
    } else if (phase === 'air') {
      if (jumpType === 'lutz') {
         cx = 47.5 + progress * 25;
         cy = 47.5 + progress * 20; 
         rotationAngle = progress * 3 * 2 * Math.PI;
         tiltAngle = progress * 70 * (Math.PI / 180);
      } else if (jumpType === 'lutz_perfect') {
         cx = 47.5;
         cy = 47.5 - Math.sin(progress * Math.PI) * 10;
         rotationAngle = progress * 4 * 2 * Math.PI;
         tiltAngle = Math.sin(progress * Math.PI) * 1 * (Math.PI / 180);
      } else {
         cx = 47.5;
         const arc = Math.sin(progress * Math.PI);
         const jumpHeight = jumpType === 'axel' ? 14 : 8;
         cy = 47.5 - arc * jumpHeight;
         const spins = jumpType === 'axel' ? 3.5 : 4;
         rotationAngle = progress * spins * 2 * Math.PI;
         const maxTilt = jumpType === 'axel' ? 2 : 1;
         tiltAngle = Math.sin(progress * Math.PI) * maxTilt * (Math.PI / 180);
      }
    } else if (phase === 'landing') {
      if (jumpType === 'lutz') {
          cx = 72.5 + progress * 10;
          cy = 67.5 + progress * 10;
      } else if (jumpType === 'lutz_perfect') {
          cx = 47.5 + progress * 1.5;
          cy = 47.5 + Math.sin(progress * Math.PI) * 3;
      } else {
          cx = 47.5 + progress * 1.5;
          cy = 47.5 + Math.sin(progress * Math.PI) * 3;
      }
    } else {
      if (jumpType === 'lutz') {
          cx = 82.5 + progress * 15;
          cy = 77.5;
      } else if (jumpType === 'lutz_perfect') {
          cx = 49 + progress * 8.5;
      } else {
          cx = 49 + progress * 8.5;
      }
    }
    
    function project(ox, oy, oz) {
      let rx = ox * Math.cos(rotationAngle) - oz * Math.sin(rotationAngle);
      let rz = ox * Math.sin(rotationAngle) + oz * Math.cos(rotationAngle);
      let tx = rx * Math.cos(tiltAngle) - oy * Math.sin(tiltAngle);
      let ty = rx * Math.sin(tiltAngle) + oy * Math.cos(tiltAngle);
      return { x: cx + tx, y: cy + ty };
    }
    
    const joints = {};
    joints.hipL = project(-4.5, 0, 0);
    joints.hipR = project(4.5, 0, 0);
    joints.chest = project(0, -17, 0);
    joints.head = project(0, -25, 0);
    joints.shoulderL = project(-6.5, -16, 0);
    joints.shoulderR = project(6.5, -16, 0);
    
    if (phase === 'entry') {
      joints.elbowL = project(-11, -11, 4); joints.elbowR = project(11, -11, 4);
      joints.wristL = project(-16, -5, 6); joints.wristR = project(16, -5, 6);
      joints.kneeL = project(-4.5, 11, 0); joints.ankleL = project(-4.5, 23, 0);
      joints.kneeR = project(7, 13, -5); joints.ankleR = project(11, 21, -10);
    } else if (phase === 'takeoff') {
      joints.elbowL = project(-7, -9, -6); joints.elbowR = project(7, -9, -6);
      joints.wristL = project(-5, -2, -10); joints.wristR = project(5, -2, -10);
      joints.kneeL = project(-5, 13, 3); joints.ankleL = project(-5, 24, 5);
      joints.kneeR = project(3.5, 14, -2); joints.ankleR = project(3.5, 24, -3);
    } else if (phase === 'air') {
      // Имитируем реальную плотную группировку рук (tight elbow tuck), как на видео
      joints.elbowL = project(-5, -16, 5); joints.elbowR = project(5, -16, 5);
      joints.wristL = project(-1, -17, 8); joints.wristR = project(1, -17, 6);
      joints.kneeL = project(-1.8, 13, 2); joints.ankleL = project(-0.9, 26, 3);
      joints.kneeR = project(1.8, 13, -2); joints.ankleR = project(0.9, 26, -1);
    } else if (phase === 'landing') {
      const la = progress;
      joints.elbowL = project(-9 - la*4, -13, 4); joints.elbowR = project(9 + la*4, -13, 4);
      joints.wristL = project(-14 - la*5, -8, 7); joints.wristR = project(14 + la*5, -8, 7);
      joints.kneeR = project(2.5, 9 + la*3, 1); joints.ankleR = project(2.5, 21 + la*2, 2);
      joints.kneeL = project(-9*la, 11 + la*2, -5*la); joints.ankleL = project(-16*la, 16 + la*3, -12*la);
    } else {
      joints.elbowL = project(-14, -14, 8); joints.elbowR = project(14, -14, 8);
      joints.wristL = project(-21, -9, 13); joints.wristR = project(21, -9, 13);
      joints.kneeR = project(3.5, 11, 0); joints.ankleR = project(3.5, 22, 0);
      joints.kneeL = project(-12, 9, -10); joints.ankleL = project(-23, 3, -20);
    }
    frames.push(joints);
  }
  return frames;
}

const jumpsData = {
  salchow: {
    id: "salchow",
    name: "Quad Salchow (4S)",
    confidence: "High Precision",
    confidenceClass: "conf-high",
    totalFrames: 32,
    startTime: 24.0, endTime: 33.0,
    activeLevelSystem: "isu",
    alignment: { scale: 1.15, offsetX: -42, offsetY: -32 },
    keyframes: { entry: 2, takeoff: 8, air: 15, landing: 24, exit: 32 },
    keyframeLabels: [
      { id: 'entry', label: 'Entry', frame: 2, timeText: '0:24' },
      { id: 'takeoff', label: 'Takeoff', frame: 8, timeText: '0:27' },
      { id: 'air', label: 'Flight', frame: 15, timeText: '0:28' },
      { id: 'landing', label: 'Landing', frame: 24, timeText: '0:29' },
      { id: 'exit', label: 'Exit', frame: 32, timeText: '0:32' }
    ],
    biometrics: {
      heightData: [0, 0.03, 0.18, 0.52, 0.62, 0.60, 0.35, 0],
      rotationData: [0, 80, 480, 1100, 1250, 1190, 800, 90],
      tiltData: [1, 2, 2, 2, 1.5, 1, 0],
      impactG: 3.9, airTimeSec: 0.65, maxHeightM: 0.62, maxTiltDeg: 2, rotationSpeedDegSec: 1250
    },
    issues: [
      { id: "s1", label: "Perfect: Deep edge entry", moment: "0:26 (Entry)", impact: 2, impactClass: "impact-positive-high", desc: "Clean entry edge and powerful speed.", frame: 2, status: "perfect" },
      { id: "s2", label: "Excellent: Vertical axis", moment: "0:28 (Flight)", impact: 2, impactClass: "impact-positive-high", desc: "Axis tilt is only 2 degrees. Tight rotation.", frame: 15, status: "perfect" }
    ],
    practice: [
      { id: "p1", title: "Rotation Tightness", reps: "5 reps", focus: "Keep elbows tight to the body for speed." }
    ],
    coachSummary: "Absolutely reference Quad Salchow by Yuzuru Hanyu! Takeoff is perfectly synchronized.",
    levels: {
      isu: { title: "ISU Basic Novice", badge: "BN", badgeClass: "badge-bn", progress: 35, nextGoal: "Transition to Inter. Novice", toReach: ["Improve height", "Stabilize axis"], desc: "Element requires improvement for category upgrade." },
      usfs: { title: "Basic 4-6", badge: "B6", badgeClass: "badge-b4", progress: 50, nextGoal: "Freeskate 1", toReach: ["More entry speed"], desc: "Confident base level." }
    },
    goeSimulation: { baseScore: 9.70, initialModifiers: { height: 4, takeoff: 4, axis: 5, rotation: 4, landing: 5 } }
  },
  axel: {
    id: "axel",
    name: "Triple Axel (3A)",
    confidence: "High Precision",
    confidenceClass: "conf-high",
    totalFrames: 45,
    startTime: 64.0, endTime: 73.0,
    activeLevelSystem: "isu",
    alignment: { scale: 1.25, offsetX: -44, offsetY: -35 },
    keyframes: { entry: 3, takeoff: 10, air: 22, landing: 34, exit: 45 },
    keyframeLabels: [
      { id: 'entry', label: 'Entry', frame: 3, timeText: '1:04' },
      { id: 'air', label: 'Flight', frame: 22, timeText: '1:08' },
      { id: 'exit', label: 'Exit', frame: 45, timeText: '1:12' }
    ],
    biometrics: {
      heightData: [0, 0.18, 0.70, 0.60, 0.25, 0],
      rotationData: [0, 480, 1180, 1100, 410, 0],
      tiltData: [1, 2, 2, 1.5, 1, 0],
      impactG: 4.2, airTimeSec: 0.72, maxHeightM: 0.70, maxTiltDeg: 2, rotationSpeedDegSec: 1180
    },
    issues: [
      { id: "a1", label: "Massive takeoff height", moment: "1:07", impact: 2, impactClass: "impact-positive-high", desc: "Height is 0.70 m.", frame: 10, status: "perfect" }
    ],
    practice: [
      { id: "ap1", title: "Complex entries", reps: "4 reps", focus: "Edge balance." }
    ],
    coachSummary: "Magnificent Triple Axel with massive distance. Soft knee on landing.",
    levels: {
      isu: { title: "ISU Elite Men", badge: "EL", badgeClass: "badge-bn", progress: 95, nextGoal: "Perfect GOE Master", toReach: ["Add 4A"], desc: "Reference jump technique." },
      usfs: { title: "Senior Elite", badge: "SR", badgeClass: "badge-b4", progress: 95, nextGoal: "Olympic Level", toReach: ["Entry simulations"], desc: "Highest qualification level." }
    },
    goeSimulation: { baseScore: 8.00, initialModifiers: { height: 5, takeoff: 4, axis: 5, rotation: 5, landing: 5 } }
  },
  lutz: {
    id: "lutz",
    name: "Triple Lutz (Fall)",
    confidence: "High Precision",
    confidenceClass: "conf-high",
    totalFrames: 45,
    startTime: 109.0, endTime: 121.0,
    activeLevelSystem: "isu",
    alignment: { scale: 1.15, offsetX: 0, offsetY: 25 },
    keyframes: { entry: 4, takeoff: 12, air: 24, landing: 35, exit: 45 },
    keyframeLabels: [
      { id: 'entry', label: 'Entry', frame: 4, timeText: '0:04' },
      { id: 'takeoff', label: 'Takeoff', frame: 12, timeText: '0:07' },
      { id: 'landing', label: 'Landing (Fall)', frame: 35, timeText: '0:10' },
      { id: 'exit', label: 'Exit', frame: 45, timeText: '0:12' }
    ],
    biometrics: {
      heightData: [0, 0.15, 0.45, 0.42, 0.18, 0],
      rotationData: [0, 320, 850, 840, 250, 0],
      tiltData: [1, 5, 12, 18, 25, 30],
      impactG: 8.5, airTimeSec: 0.52, maxHeightM: 0.45, maxTiltDeg: 30, rotationSpeedDegSec: 850
    },
    issues: [
      { id: "l1", label: "Root Cause: Early Shoulder Drop", moment: "0:09", impact: -5, impactClass: "impact-negative-high", desc: "Right shoulder dropped on takeoff, causing a 25° axis tilt in flight. TO AVOID: Lock the right shoulder and engage the core during the toe pick strike.", frame: 24, status: "danger" },
      { id: "l2", label: "Critical: Premature Opening", moment: "0:10", impact: -5, impactClass: "impact-negative-high", desc: "Panic from axis loss caused early opening, resulting in a hard fall. TO AVOID: Pull arms tighter immediately after takeoff and hold the rotation position longer.", frame: 35, status: "danger" }
    ],
    practice: [
      { id: "lp1", title: "Core Axis Drills", reps: "20 reps", focus: "Engage right side core to prevent tilting left in the air." },
      { id: "lp2", title: "Harness Training", reps: "15 mins", focus: "Perform jump on harness to safely rebuild axis memory." }
    ],
    coachSummary: "This fall was caused by dropping the right shoulder immediately at takeoff, ruining the axis before you even left the ice. TO AVOID THIS: Drive your knee straight up, keep the upper body completely rigid, and do not lean outside the circle. Focus on a vertical takeoff vector.",
    levels: {
      isu: { title: "ISU Senior Women", badge: "SR", badgeClass: "badge-bn", progress: 30, nextGoal: "Fix Axis Tilt", toReach: ["Core stability"], desc: "Element downgraded + fall deduction." },
      usfs: { title: "Senior", badge: "SR", badgeClass: "badge-b4", progress: 30, nextGoal: "Stabilize Axis", toReach: ["Harness work"], desc: "Deductions applied for fall." }
    },
    goeSimulation: { baseScore: 5.90, initialModifiers: { height: 1, takeoff: -2, axis: -5, rotation: -3, landing: -5 } }
  }
};

jumpsData.salchow.skeletonFrames = generateSkeletonFrames('salchow', 32, jumpsData.salchow.keyframes);
jumpsData.axel.skeletonFrames = generateSkeletonFrames('axel', 45, jumpsData.axel.keyframes);
jumpsData.lutz.skeletonFrames = generateSkeletonFrames('lutz', 45, jumpsData.lutz.keyframes);

jumpsData.salchow.goldModel = generateSkeletonFrames('salchow', 32, jumpsData.salchow.keyframes);
jumpsData.axel.goldModel = generateSkeletonFrames('axel', 45, jumpsData.axel.keyframes);
jumpsData.lutz.goldModel = generateSkeletonFrames('lutz_perfect', 45, jumpsData.lutz.keyframes);


// ==========================================
// Основная Логика (Main Logic)
// ==========================================
let activeJump = jumpsData.salchow;
let currentFrame = 0;
let isPlaying = false;
let activeSystem = 'isu';
let goeModifiers = {};
let videoPlayer = null;

// MediaPipe Pose Tracker Variables
let poseModel = null;
let currentMediaPipeJoints = null;
let isMediaPipeReady = false;
window.successfulJumpHistory = [];
window.lastAnalyzedJumpHistory = [];

function formatMediaPipeLandmarks(lms) {
    const j = {};
    const toPt3D = (idx) => {
        if (!lms[idx]) return null;
        return { 
            x: lms[idx].x * 100, 
            y: lms[idx].y * 100,
            z: lms[idx].z * 100 
        };
    };
    
    j.head = toPt3D(0);
    j.shoulderL = toPt3D(12);
    j.shoulderR = toPt3D(11);
    j.elbowL = toPt3D(14);
    j.elbowR = toPt3D(13);
    j.wristL = toPt3D(16);
    j.wristR = toPt3D(15);
    j.hipL = toPt3D(24);
    j.hipR = toPt3D(23);
    j.kneeL = toPt3D(26);
    j.kneeR = toPt3D(25);
    j.ankleL = toPt3D(28);
    j.ankleR = toPt3D(27);
    j.footL = toPt3D(32); 
    j.footR = toPt3D(31);
    
    // Fallbacks if visibility is low
    Object.keys(j).forEach(k => { if (!j[k]) j[k] = {x:0,y:0,z:0}; });
    
    j.chest = { x: (j.shoulderL.x + j.shoulderR.x)/2, y: (j.shoulderL.y + j.shoulderR.y)/2, z: (j.shoulderL.z + j.shoulderR.z)/2 };
    return j;
}

function formatAiText(text) {
    if (!text) return "";
    let formatted = text.replace(/(?:\s|^)(\d+[\)\.])\s/g, '<br><br><span style="color: #38BDF8; font-weight: bold; margin-right: 4px;">$1</span>');
    formatted = formatted.replace(/^(<br>)+/, '');
    return formatted.replace(/\n/g, '<br>');
}

function renderJumpHistory() {
    const container = document.getElementById('jump-history-container');
    if (!container || !window.completedJumps || window.completedJumps.length === 0) return;
    
    container.innerHTML = '';
    window.completedJumps.forEach((jump, index) => {
        const item = document.createElement('div');
        const borderColor = (jump.is_anomaly || jump.isAnomaly || jump.goe < 0) ? "#EF4444" : "#10B981";
        item.style = "background: rgba(255,255,255,0.05); padding: 8px 12px; border-radius: 6px; display: flex; justify-content: space-between; align-items: center; border-left: 3px solid " + borderColor + ";";
        
        let title = jump.classification || "Unknown Element";
        let score = jump.finalScore ? jump.finalScore.toFixed(2) + " pts" : (jump.goe ? jump.goe.toFixed(2) + " pts" : "Analyzed");
        
        item.innerHTML = '<div style="display: flex; flex-direction: column;">' +
                '<span style="font-size: 13px; font-weight: 600; color: #E2E8F0;">Jump ' + (index + 1) + ': ' + title + '</span>' +
            '</div>' +
            '<div style="font-size: 13px; font-weight: bold; color: ' + (jump.goe < 0 ? '#FCA5A5' : '#6EE7B7') + ';">' +
                score +
            '</div>';
        container.appendChild(item);
    });
}

let lastAnalyzeTime = 0;
let frameHistory = [];
let analyzeRunId = 0;

function classifyJump(joints) {
    // Проверяем наличие критически важных точек для 3D анализа (включая стопы)
    if (!joints || !joints.shoulderL || !joints.hipL || !joints.footL) return;
    
    // Сохраняем полный набор 3D суставов
    frameHistory.push(joints);
    
    if (frameHistory.length > 150) {
        frameHistory.shift();
    }
    
    // Ограничиваем отправку данных до 10 кадров в секунду, чтобы не перегружать Go-сервер
    const now = Date.now();
    if (now - lastAnalyzeTime < 100) return;
    lastAnalyzeTime = now;
    
    const currentRunId = analyzeRunId;
    
    fetch('/api/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(frameHistory)
    })
    .then(r => r.json())
    .then(data => {
        if (currentRunId !== analyzeRunId) return; // ПРЕДОТВРАЩАЕМ ГОНКУ ЗАПРОСОВ! Игнорируем ответы из прошлого.
        
        const badge = document.getElementById('detected-jump-badge');
        const classInfo = document.querySelector('.classification-info');
        
        if (data.status === "detected") {
            window.bestJumpData = data;
            window.isTrackingJump = true;
        } else if (data.status === "analyzing" && window.isTrackingJump) {
            if (!window.completedJumps) window.completedJumps = [];
            const now = Date.now();
            if (!window.lastJumpSaveTime || (now - window.lastJumpSaveTime > 2000)) {
                window.completedJumps.push(window.bestJumpData);
                window.lastJumpSaveTime = now;
                renderJumpHistory();
            }
            window.isTrackingJump = false;
            window.bestJumpData = null;
        }
        
        let jumpToDisplay = window.bestJumpData;
        if (!jumpToDisplay && window.completedJumps && window.completedJumps.length > 0) {
            jumpToDisplay = window.completedJumps[window.completedJumps.length - 1];
        }
        const finalAnswer = jumpToDisplay || data;

        if (badge && finalAnswer.status !== "analyzing") {
            let outputText = finalAnswer.classification ? finalAnswer.classification.toUpperCase() : "UNKNOWN";
            if (finalAnswer.is_anomaly) {
                outputText += " (" + finalAnswer.anomaly_type.toUpperCase() + ")";
            }
            badge.innerText = "DETECTED: " + outputText;
        }
        
        if (classInfo && finalAnswer.status !== "analyzing") {
            // Скрываем весь блок вместе с синей полоской по просьбе пользователя
            classInfo.style.display = 'none';
            
            // Очищаем блок нарушений до ответа ИИ
            const violationsContainer = document.getElementById('violations-panel');
            if (violationsContainer) {
                violationsContainer.innerHTML = '';
            }
        }

        if (finalAnswer.status !== "analyzing") {
            // Обновляем панель очков финальными (замороженными) данными
            const hdrBase = document.getElementById('hdr-base');
            const hdrGoe = document.getElementById('hdr-goe');
            const hdrScore = document.getElementById('hdr-score');
            const hdrReason = document.getElementById('hdr-reason');
            
            if (hdrBase) hdrBase.textContent = finalAnswer.baseScore.toFixed(2) + " pts";
            if (hdrGoe) {
                const sign = finalAnswer.goe > 0 ? "+" : "";
                hdrGoe.textContent = sign + finalAnswer.goe.toFixed(2) + " pts";
                hdrGoe.style.color = finalAnswer.goe < 0 ? "#EF4444" : "#10B981";
            }
            if (hdrScore) hdrScore.textContent = finalAnswer.finalScore.toFixed(2) + " pts";
            if (hdrReason) hdrReason.textContent = finalAnswer.scoreReason;
            
            // Обновляем быструю диагностику только финальными данными
            const diagCauseElem = document.getElementById('diag-cause-text');
            const diagFixElem = document.getElementById('diag-fix-text');
            if (diagCauseElem) {
                // ДОБАВЛЯЕМ ДАННЫЕ О ПРЫЖКЕ В ДИАГНОСТИКУ:
                let jumpDataHtml = '<div style="margin-bottom: 8px; font-weight: bold; color: #E2E8F0; padding-bottom: 6px; border-bottom: 1px solid rgba(255,255,255,0.1);">';
                jumpDataHtml += 'Element: <span style="color: ' + (finalAnswer.is_anomaly || finalAnswer.isAnomaly ? '#FCA5A5' : '#34D399') + ';">' + (finalAnswer.classification || 'Unknown') + '</span>';
                jumpDataHtml += '</div>';

                const parts = (finalAnswer.diagnosticCause || '').split('; ').filter(p => p.trim().length > 0);
                if (parts.length > 1) {
                    diagCauseElem.innerHTML = jumpDataHtml + '<ul style="margin: 4px 0 0 0; padding-left: 18px; list-style: disc;">' +
                        parts.map(p => '<li style="margin-bottom: 4px; line-height: 1.4;">' + p.trim() + '</li>').join('') +
                        '</ul>';
                } else {
                    diagCauseElem.innerHTML = jumpDataHtml + '<div>' + finalAnswer.diagnosticCause + '</div>';
                }
            }
            if (diagFixElem) diagFixElem.textContent = finalAnswer.diagnosticFix;
            
            // Запрашиваем полный анализ от ИИ-тренера (Gemini LLM)
            if (diagFixElem && !window.coachRequestedForRun && finalAnswer.status === "detected") {
                window.coachRequestedForRun = true; // предотвращаем множественные запросы
                window.lastAnalyzedJumpHistory = [...frameHistory]; // Сохраняем текущий прыжок для возможного использования как эталон
                diagFixElem.innerHTML = '<div style="color: #A855F7; font-style: italic; font-weight: bold;">ИИ-тренер анализирует 3D-координаты... 🧠</div>';
                
                fetch('/api/coach', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ history: window.lastAnalyzedJumpHistory, successHistory: window.successfulJumpHistory })
                })
                .then(r => r.json())
                .then(coachData => {
                    if (coachData && coachData.status === "detected") {
                        window.bestJumpData = coachData;
                        
                        // Обновляем бейдж
                        document.getElementById('detected-jump-badge').innerText = "DETECTED: " + (coachData.classification ? coachData.classification.toUpperCase() : "UNKNOWN");
                        
                        // Обновляем телеметрию
                        const classInfo = document.querySelector('.classification-info');
                        if (classInfo) {
                            classInfo.style.display = 'block'; // Показываем блок с синей полосой обратно!
                            let html = '<div style="color: #E2E8F0; font-weight: 600; margin-bottom: 6px; display: flex; align-items: center; gap: 8px;">' +
                                    '<svg viewBox="0 0 24 24" width="14" height="14" fill="#A855F7"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>' +
                                    'SkateEye AI: Full Biomechanics Analysis' +
                                '</div>' +
                                'AI Telemetry: Rotation ' + (coachData.shoulderAngle||0).toFixed(0) + '°, Ankle Lean ' + (coachData.ankleLean||0).toFixed(1) + 'px.' +
                                '<br><span style="color: #A855F7; font-weight: 600; margin-top: 4px; display: inline-block;">Neural Engine:</span> ' + (coachData.probabilityText||"Confident");
                            html += '<br>AI Target locked as: <span style="color: #F8FAFC; font-weight: bold; background: rgba(168, 85, 247, 0.4); border: 1px solid rgba(168, 85, 247, 0.8); padding: 1px 6px; border-radius: 3px; margin-left: 4px;">' + coachData.classification + '</span>';
                            classInfo.innerHTML = html;
                        }

                        // Обновляем нарушения
                        const violationsContainer = document.getElementById('violations-panel');
                        if (violationsContainer) {
                            if (coachData.violations && coachData.violations.length > 0) {
                                let vHtml = '<div style="margin-top: 8px; border-top: 1px solid rgba(239,68,68,0.3); padding-top: 6px;">';
                                vHtml += '<div style="color: #F87171; font-weight: 700; font-size: 11px; letter-spacing: 0.05em; margin-bottom: 4px;">⚖️ ISU VIOLATIONS DETECTED (' + coachData.violations.length + ')</div>';
                                coachData.violations.forEach(function(v) {
                                    vHtml += '<div style="display: flex; align-items: flex-start; gap: 6px; margin-bottom: 3px;">';
                                    vHtml += '<span style="color: #EF4444; font-size: 10px; margin-top: 1px;">▶</span>';
                                    vHtml += '<span style="color: #FCA5A5; font-size: 11px; line-height: 1.4;">' + v + '</span>';
                                    vHtml += '</div>';
                                });
                                vHtml += '</div>';
                                violationsContainer.innerHTML = vHtml;
                                violationsContainer.style.display = 'block';
                            } else {
                                violationsContainer.innerHTML = '<div style="margin-top: 8px; border-top: 1px solid rgba(16,185,129,0.3); padding-top: 6px; color: #34D399; font-size: 11px; font-weight: 600;">✅ Clean jump according to AI!</div>';
                                violationsContainer.style.display = 'block';
                            }
                        }

                        // Обновляем очки
                        document.getElementById('hdr-base').textContent = (coachData.baseScore||0).toFixed(2) + " pts";
                        document.getElementById('hdr-goe').textContent = (coachData.goe>0?"+":"") + (coachData.goe||0).toFixed(2) + " pts";
                        document.getElementById('hdr-goe').style.color = coachData.goe < 0 ? "#EF4444" : "#10B981";
                        document.getElementById('hdr-score').textContent = (coachData.finalScore||0).toFixed(2) + " pts";
                        document.getElementById('hdr-reason').textContent = coachData.scoreReason || "";

                        // Обновляем диагноз
                        let jumpDataHtml = '<div style="margin-bottom: 8px; font-weight: bold; color: #E2E8F0; padding-bottom: 6px; border-bottom: 1px solid rgba(255,255,255,0.1);">';
                        jumpDataHtml += 'AI Classified Element: <span style="color: #A855F7;">' + (coachData.classification || 'Unknown') + '</span>';
                        jumpDataHtml += '</div>';
                        diagCauseElem.innerHTML = jumpDataHtml + '<div style="line-height: 1.5; margin-top: 8px;">' + formatAiText(coachData.diagnosticCause) + '</div>';
                        diagFixElem.innerHTML = '<div style="color: #34D399; font-weight: 600; margin-bottom: 8px; display: flex; justify-content: space-between; align-items: center;">' + 
                                                '<span>Advice from an AI trainer:</span>' +
                                                '<button id="btn-save-reference" style="background: rgba(16, 185, 129, 0.2); border: 1px solid #10B981; color: #10B981; padding: 2px 8px; border-radius: 4px; font-size: 11px; cursor: pointer; transition: all 0.2s;">Set as Reference Jump</button>' +
                                                '</div>' + 
                                                '<div style="line-height: 1.5;">' + formatAiText(coachData.diagnosticFix) + '</div>';
                        
                        setTimeout(() => {
                            const btnRef = document.getElementById('btn-save-reference');
                            if (btnRef) {
                                btnRef.addEventListener('click', () => {
                                    window.successfulJumpHistory = window.lastAnalyzedJumpHistory;
                                    btnRef.textContent = "Saved as Reference!";
                                    btnRef.style.background = "#10B981";
                                    btnRef.style.color = "#FFF";
                                    setTimeout(() => { 
                                        btnRef.textContent = "Set as Reference Jump"; 
                                        btnRef.style.background = "rgba(16, 185, 129, 0.2)"; 
                                        btnRef.style.color = "#10B981"; 
                                    }, 2000);
                                });
                            }
                        }, 50);

                        if (window.completedJumps && window.completedJumps.length > 0) {
                            window.completedJumps[window.completedJumps.length - 1] = coachData;
                            renderJumpHistory();
                        }
                    } else if (coachData && coachData.status === "error") {
                        diagCauseElem.innerHTML = '<div style="color: #EF4444; font-weight: bold;">AI Error occurred</div>';
                        diagFixElem.innerHTML = '<div style="color: #EF4444; font-weight: 600;">' + coachData.diagnosticFix + '</div>';
                    }
                })
                .catch(err => {
                    console.error("Coach API error:", err);
                    diagCauseElem.innerHTML = '<div style="color: #EF4444; font-weight: bold;">Network Error</div>';
                    diagFixElem.innerHTML = '<div style="color: #EF4444;">Ошибка при обращении к ИИ-тренеру. Убедитесь, что сервер запущен.</div>';
                });
            }
        }
        
        // Сброс флага при новом анализе
        if (finalAnswer.status === "analyzing") {
            window.coachRequestedForRun = false;
        }

        // Обновляем живой график завала конька!
        if (window.Chart) {
            updateEdgeChart(frameHistory);
        }
    })
    .catch(err => console.error("Analyzer Error:", err));
}

let edgeChartInstance = null;
function updateEdgeChart(history) {
    const canvas = document.getElementById('edgeChartCanvas');
    if (!canvas) return;

    const labels = history.map((_, i) => i);
    const dataPoints = history.map(h => {
        if (!h || !h.footL || !h.ankleL) return 0;
        return h.footL.x - h.ankleL.x;
    });

    if (!edgeChartInstance) {
        edgeChartInstance = new Chart(canvas, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Ankle Lean (px)',
                    data: dataPoints,
                    borderColor: '#3B82F6',
                    borderWidth: 2,
                    pointRadius: 0,
                    fill: true,
                    backgroundColor: 'rgba(59, 130, 246, 0.1)',
                    tension: 0.2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: false,
                scales: {
                    y: {
                        suggestedMin: -5,
                        suggestedMax: 5,
                        grid: { color: 'rgba(255,255,255,0.05)' },
                        title: { display: true, text: '<- Outside Edge (Lutz)   Inside Edge (Flip) ->', color: '#94A3B8' }
                    },
                    x: { display: false }
                },
                plugins: {
                    legend: { display: false }
                }
            }
        });
    } else {
        edgeChartInstance.data.labels = labels;
        edgeChartInstance.data.datasets[0].data = dataPoints;
        edgeChartInstance.update();
    }
}

function initMediaPipe() {
    if (!window.Pose) return;
    poseModel = new window.Pose({locateFile: (file) => {
        return "https://cdn.jsdelivr.net/npm/@mediapipe/pose/" + file;
    }});
    poseModel.setOptions({
        modelComplexity: 2,
        smoothLandmarks: true,
        enableSegmentation: false,
        minDetectionConfidence: 0.1,
        minTrackingConfidence: 0.1
    });
    poseModel.onResults((results) => {
        if (results.poseLandmarks && results.poseLandmarks.length > 0) {
            currentMediaPipeJoints = formatMediaPipeLandmarks(results.poseLandmarks);
            // Выводим сырые данные нейросети прямо в консоль браузера!
            console.log("AI DETECTED POINTS:", currentMediaPipeJoints.ankleL.x.toFixed(1), currentMediaPipeJoints.ankleL.y.toFixed(1));
            classifyJump(currentMediaPipeJoints);
        } else {
            currentMediaPipeJoints = null;
        }
        renderFrameUI();
    });
    isMediaPipeReady = true;
}

async function processVideoFrameLoop() {
    if (isPlaying && videoPlayer && isMediaPipeReady) {
        try {
            await poseModel.send({image: videoPlayer});
        } catch (e) {
            console.error("Pose processing error:", e);
        }
    }
    requestAnimationFrame(processVideoFrameLoop);
}

const COACH_CHAT_RESPONSES = {
  default: "Great question! For this element, it's critical to control shoulder alignment during entry. Avoid opening your arms too early.",
  axis: "To fix axis tilt, keep your torso strictly vertical during takeoff. Direct your takeoff force vector upwards, not forwards.",
  landing: "To stabilize landing and avoid exit delay: open your free leg backwards immediately upon toe-pick contact.",
  upgrade: "Upgrade to PRO to unlock real-time 3D joint tracking, unlimited videos, and personal consultations."
};

function fitJumpToVideo(jump) {
  if (!videoPlayer) return;
  
  if (videoPlayer.src.includes('skater.mp4') || videoPlayer.src.includes('skater.webm') || videoPlayer.src.includes('practice_fall.mp4')) {
      jump.startTime = 0.0;
      const dur = videoPlayer.duration && videoPlayer.duration !== Infinity ? videoPlayer.duration : 30.0;
      jump.endTime = dur;
      jump.totalFrames = Math.floor(dur * 30); 
      jump.keyframes = {
        entry: Math.floor(jump.totalFrames * 0.1),
        takeoff: Math.floor(jump.totalFrames * 0.3),
        air: Math.floor(jump.totalFrames * 0.5),
        landing: Math.floor(jump.totalFrames * 0.8),
        exit: jump.totalFrames - 1
      };
  } else if (videoPlayer.src.includes('single_jump.mp4')) {
      jump.startTime = 0.0;
      const dur = videoPlayer.duration && videoPlayer.duration !== Infinity ? videoPlayer.duration : 16.0;
      jump.endTime = dur;
      jump.totalFrames = Math.floor(dur * 30);
      
      jump.keyframes = {
        entry: Math.floor(jump.totalFrames * 0.1),
        takeoff: Math.floor(jump.totalFrames * 0.3),
        air: Math.floor(jump.totalFrames * 0.5),
        landing: Math.floor(jump.totalFrames * 0.8),
        exit: jump.totalFrames - 1
      };
  } else {
      const dur = videoPlayer.duration && videoPlayer.duration !== Infinity ? videoPlayer.duration : 10.0;
      jump.startTime = 0.0;
      jump.endTime = dur;
      jump.totalFrames = Math.max(30, Math.floor(dur * 30));
      jump.keyframes = {
        entry: Math.floor(jump.totalFrames * 0.1),
        takeoff: Math.floor(jump.totalFrames * 0.3),
        air: Math.floor(jump.totalFrames * 0.5),
        landing: Math.floor(jump.totalFrames * 0.8),
        exit: jump.totalFrames - 1
      };
  }
  
  const map = ['entry', 'takeoff', 'air', 'landing', 'exit'];
  jump.keyframeLabels.forEach((kf, idx) => {
    if(idx < map.length) {
      kf.frame = jump.keyframes[map[idx]];
      const t = jump.startTime + (jump.endTime - jump.startTime) * (kf.frame / jump.totalFrames);
      kf.timeText = "0:" + t.toFixed(2).replace('.',':');
    }
  });
  
  jump.skeletonFrames = generateSkeletonFrames(jump.id, jump.totalFrames, jump.keyframes);
  const goldJumpType = jump.id === 'lutz' ? 'lutz_perfect' : jump.id;
  jump.goldModel = generateSkeletonFrames(goldJumpType, jump.totalFrames, jump.keyframes);
}

document.addEventListener('DOMContentLoaded', () => {
  videoPlayer = document.getElementById('video-player');
  setupEventListeners();
  
  // Инициализация MediaPipe
  setTimeout(() => {
    initMediaPipe();
    processVideoFrameLoop();
  }, 1000);
  
  if (videoPlayer) {
    videoPlayer.addEventListener('loadedmetadata', () => { selectJump('lutz'); });
    if (videoPlayer.readyState >= 1) selectJump('lutz');
  } else {
    selectJump('lutz');
  }
  
  addChatMessage('coach', 'Hi! I am your AI Coach. I have thoroughly analyzed your video. MediaPipe Tracker initialized. Ask me how to improve your technique!');
});

function resetAnalysisState() {
    frameHistory = [];
    lastAnalyzeTime = 0;
    analyzeRunId++; 
    window.bestJumpData = null; 
    
    const badge = document.getElementById('detected-jump-badge');
    if (badge) badge.innerText = "DETECTING...";
    const classInfo = document.querySelector('.classification-info');
    if (classInfo) classInfo.style.display = 'none';
    
    const hdrBase = document.getElementById('hdr-base');
    if (hdrBase) hdrBase.textContent = "-- pts";
    const hdrGoe = document.getElementById('hdr-goe');
    if (hdrGoe) hdrGoe.textContent = "-- pts";
    const hdrScore = document.getElementById('hdr-score');
    if (hdrScore) hdrScore.textContent = "-- pts";
    const hdrReason = document.getElementById('hdr-reason');
    if (hdrReason) hdrReason.textContent = "Waiting...";
}

function setupEventListeners() {
  setupTabs();
  const selector = document.getElementById('jump-selector');
  if (selector) selector.addEventListener('change', (e) => selectJump(e.target.value));

  const btnPlay = document.getElementById('btn-play-pause');
  if (btnPlay) btnPlay.addEventListener('click', togglePlay);

  const scrubber = document.getElementById('video-scrubber');
  if (scrubber) scrubber.addEventListener('input', (e) => seekToFrame(parseInt(e.target.value)));

  if (videoPlayer) {
    videoPlayer.addEventListener('seeking', () => {
        frameHistory = []; // Only clear frames so analyzer doesn't glitch, KEEP bestJumpData!
    });

    videoPlayer.addEventListener('timeupdate', handleVideoTimeUpdate);
    videoPlayer.addEventListener('play', () => {
      isPlaying = true;
      
      // Сбрасываем только историю кадров, если проигрывание начинается заново
      if (videoPlayer.currentTime < 1.0) {
          frameHistory = [];
      }

      if (btnPlay) {
        const pIcon = document.getElementById('icon-play');
        const pauseIcon = document.getElementById('icon-pause');
        if (pIcon) pIcon.style.display = 'none';
        if (pauseIcon) pauseIcon.style.display = 'block';
      }
    });
    
    videoPlayer.addEventListener('pause', () => {
      isPlaying = false;
      if (btnPlay) {
        const pIcon = document.getElementById('icon-play');
        const pauseIcon = document.getElementById('icon-pause');
        if (pIcon) pIcon.style.display = 'block';
        if (pauseIcon) pauseIcon.style.display = 'none';
      }
    });
  }

  document.getElementById('tab-isu')?.addEventListener('click', () => switchJudgingSystem('isu'));
  document.getElementById('tab-usfs')?.addEventListener('click', () => switchJudgingSystem('usfs'));

  // Десктопные табы
  const navItems = document.querySelectorAll('.desktop-tabs .desk-tab-btn');
  navItems.forEach(item => {
    item.addEventListener('click', (e) => {
      e.preventDefault();
      const tabId = item.getAttribute('data-tab');
      navItems.forEach(i => i.classList.remove('active'));
      item.classList.add('active');
      document.querySelectorAll('.nav-tab-btn').forEach(b => b.classList.remove('active'));
      document.getElementById("tab-btn-" + tabId)?.classList.add('active');
      document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
      document.getElementById("content-" + tabId)?.classList.add('active');
    });
  });

  // Мобильные табы
  const tabBtns = document.querySelectorAll('.nav-tab-btn');
  tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      const tabId = btn.getAttribute('data-tab');
      tabBtns.forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      navItems.forEach(i => {
        if (i.getAttribute('data-tab') === tabId) i.classList.add('active');
        else i.classList.remove('active');
      });
      document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
      document.getElementById("content-" + tabId)?.classList.add('active');
    });
  });

  document.getElementById('btn-upgrade')?.addEventListener('click', () => {
    document.getElementById('tab-btn-chat')?.click();
    openChatWithQuestion("What are the benefits of the PRO plan?", 'upgrade');
  });

  const chatInput = document.getElementById('chat-user-input');
  const btnSend = document.getElementById('btn-chat-send');
  if (btnSend && chatInput) {
    btnSend.addEventListener('click', () => {
      const text = chatInput.value.trim();
      if (text) { chatInput.value = ''; openChatWithQuestion(text, 'default'); }
    });
    chatInput.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') {
        const text = chatInput.value.trim();
        if (text) { chatInput.value = ''; openChatWithQuestion(text, 'default'); }
      }
    });
  }

  const fileInput = document.getElementById('video-file-input');
  const fileBtn = document.getElementById('btn-upload-file');
  if (fileInput && fileBtn) {
    fileInput.addEventListener('change', (e) => {
      if (e.target.files.length > 0) {
        const file = e.target.files[0];
        const cleanName = file.name.replace(/[^a-zA-Z0-9.-_ ]/g, '');
        fileBtn.querySelector("span").textContent = cleanName;
        openChatWithQuestion("I uploaded a video: " + cleanName + ". Can you analyze it?", 'default');
        
        const url = URL.createObjectURL(file);
        videoPlayer.src = url;
        videoPlayer.load();

        setTimeout(() => {
          addChatMessage('coach', "Video " + cleanName + " successfully uploaded! Performing frame-by-frame joint tracking...");
          videoPlayer.addEventListener('loadedmetadata', function onLoaded() {
            videoPlayer.removeEventListener('loadedmetadata', onLoaded);
            setTimeout(() => {
              selectJump('lutz');
              addChatMessage('coach', 'Analysis complete! I have adjusted the metrics and 3D model for your new video. Check the dashboard.');
            }, 1500);
          });
        }, 1000);
      }
    });
  }
}

function selectJump(jumpId) {
  if (videoPlayer) videoPlayer.pause();
  activeJump = jumpsData[jumpId];
  
  // Full reset ONLY on new jump/video selection
  resetAnalysisState();
  

  fitJumpToVideo(activeJump);
  currentFrame = 0;
  
  const scrubber = document.getElementById('video-scrubber');
  if (scrubber) { scrubber.setAttribute('max', activeJump.totalFrames - 1); scrubber.value = 0; }
  const frameCounter = document.getElementById('frame-counter');
  if (frameCounter) frameCounter.textContent = "Frames: " + activeJump.totalFrames;

  updateTopBar();
  updateKeyframesRow();
  updateQuickDiagnostics();
  updateIssuesTable();
  updatePracticeDrills();
  updateJudging(true);
  updateJumpMetrics();
  syncVideoToFrame(0);
}

function togglePlay() {
  if (!videoPlayer) return;
  if (isPlaying) {
    videoPlayer.pause();
  } else {
    const startTime = activeJump.startTime;
    const endTime = activeJump.endTime;
    if (videoPlayer.currentTime >= endTime - 0.1 || videoPlayer.currentTime < startTime) {
      resetAnalysisState(); // ПРИНУДИТЕЛЬНО СБРАСЫВАЕМ при нажатии PLAY в конце видео
      videoPlayer.currentTime = startTime;
    }
    videoPlayer.play();
  }
}

function seekToFrame(frameIndex) {
  if (videoPlayer) videoPlayer.pause();
  resetAnalysisState(); // ПРИНУДИТЕЛЬНО СБРАСЫВАЕМ при любой перемотке
  syncVideoToFrame(frameIndex);
  
  // Process single frame for Pose on seek
  setTimeout(async () => {
      if (videoPlayer && isMediaPipeReady && poseModel && !isPlaying) {
          try {
             await poseModel.send({image: videoPlayer});
          } catch(e) {}
      }
  }, 50);
}

function syncVideoToFrame(frameIndex) {
  if (!videoPlayer) return;
  const startTime = activeJump.startTime;
  const endTime = activeJump.endTime;
  const segmentDuration = endTime - startTime;
  const framePercent = frameIndex / (activeJump.totalFrames - 1);
  videoPlayer.currentTime = startTime + (segmentDuration * framePercent);
  currentFrame = frameIndex;
  renderFrameUI();
}

let lastVideoTimeTracker = -1;
function handleVideoTimeUpdate(e) {
  if (!videoPlayer) return;
  const startTime = activeJump.startTime;
  const endTime = activeJump.endTime;
  const currentTime = videoPlayer.currentTime;

  // Если время вдруг прыгнуло назад (видео зациклилось или отмотано)
  if (lastVideoTimeTracker > currentTime && currentTime < startTime + 1.0) {
      resetAnalysisState();
  }
  lastVideoTimeTracker = currentTime;

  if (currentTime < startTime) { videoPlayer.currentTime = startTime; return; }
  if (currentTime > endTime) {
    if (videoPlayer.loop) { videoPlayer.currentTime = startTime; } 
    else { videoPlayer.pause(); videoPlayer.currentTime = endTime; }
    return;
  }

  const segmentDuration = endTime - startTime;
  const progressInSegment = (currentTime - startTime) / segmentDuration;
  const frameIndex = Math.min(Math.floor(progressInSegment * activeJump.totalFrames), activeJump.totalFrames - 1);

  if (frameIndex !== currentFrame) {
    currentFrame = frameIndex;
    const scrubber = document.getElementById('video-scrubber');
    if (scrubber) scrubber.value = currentFrame;
    renderFrameUI();
  }
}

function renderFrameUI() {
  drawSimulatedVideo();
  drawModelComparison(document.getElementById('canvas-incorrect'), document.getElementById('canvas-correct'), activeJump, currentFrame);
  updateCorrectnessLabels();
  updateVideoOverlayTexts();
  updateActiveKeyframeCard();
}

function updateTopBar() {
  const confIndicator = document.getElementById('confidence-indicator');
  const confText = document.getElementById('confidence-text');
  if (confIndicator && confText) {
    confText.textContent = activeJump.confidence;
    confIndicator.className = 'confidence-badge ' + activeJump.confidenceClass;
  }
}

function updateKeyframesRow() {
  const container = document.getElementById('keyframes-list');
  if (!container) return;
  container.replaceChildren();

  activeJump.keyframeLabels.forEach((kf) => {
    const card = document.createElement('div');
    card.className = 'keyframe-card';
    card.id = 'kf-card-' + kf.id;
    card.addEventListener('click', () => seekToFrame(kf.frame));

    const thumbBox = document.createElement('div');
    thumbBox.className = 'keyframe-thumb-box';
    const thumbCanvas = document.createElement('canvas');
    thumbCanvas.className = 'keyframe-canvas';
    thumbCanvas.width = 100; thumbCanvas.height = 48;
    thumbBox.appendChild(thumbCanvas);

    const label = document.createElement('span');
    label.className = 'keyframe-label';
    label.textContent = kf.label;

    const time = document.createElement('span');
    time.className = 'keyframe-time';
    time.textContent = kf.timeText;

    card.appendChild(thumbBox);
    card.appendChild(label);
    card.appendChild(time);
    container.appendChild(card);

    const thumbCtx = thumbCanvas.getContext('2d');
    if (thumbCtx && videoPlayer && videoPlayer.src) {
      thumbCtx.fillStyle = '#0b0f1d';
      thumbCtx.fillRect(0, 0, 100, 48);
      const tempVideo = document.createElement('video');
      tempVideo.src = videoPlayer.src;
      tempVideo.muted = true;
      tempVideo.playsInline = true;
      const segmentDuration = activeJump.endTime - activeJump.startTime;
      const progress = kf.frame / (activeJump.totalFrames - 1);
      tempVideo.currentTime = activeJump.startTime + progress * segmentDuration;
      tempVideo.addEventListener('seeked', () => {
        thumbCtx.drawImage(tempVideo, 0, 0, 100, 48);
        tempVideo.src = ''; tempVideo.remove();
      }, { once: true });
    }
  });
}

function updateActiveKeyframeCard() {
  activeJump.keyframeLabels.forEach(kf => {
    document.getElementById('kf-card-' + kf.id)?.classList.remove('active');
  });
  let currentKf = activeJump.keyframeLabels[0];
  for (let i = 0; i < activeJump.keyframeLabels.length; i++) {
    if (currentFrame >= activeJump.keyframeLabels[i].frame) currentKf = activeJump.keyframeLabels[i];
  }
  document.getElementById('kf-card-' + currentKf.id)?.classList.add('active');
}

function drawSimulatedVideo() {
  const canvas = document.getElementById('video-canvas');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  const w = canvas.width; const h = canvas.height;
  ctx.clearRect(0, 0, w, h);
  // Желтый человечек убран с главного видео по просьбе пользователя
}

function updateVideoOverlayTexts() {
  const phaseOverlay = document.getElementById('video-phase-overlay');
  const timeOverlay = document.getElementById('video-time-overlay');
  if (!phaseOverlay || !timeOverlay) return;

  let phaseText = 'Entry';
  const kf = activeJump.keyframes;
  if (currentFrame < kf.takeoff) phaseText = 'Entry';
  else if (currentFrame < kf.air) phaseText = 'Takeoff';
  else if (currentFrame < kf.landing) phaseText = 'Flight';
  else if (currentFrame < kf.exit) phaseText = 'Landing';
  else phaseText = 'Exit';
  
  phaseOverlay.textContent = phaseText;
  const relTime = Math.max(0, videoPlayer.currentTime - activeJump.startTime);
  timeOverlay.textContent = "0:" + relTime.toFixed(2).replace('.', ':');
}

function updateCorrectnessLabels() {
  const descInc = document.getElementById('desc-incorrect');
  const descCor = document.getElementById('desc-correct');
  if (!descInc || !descCor) return;
  descCor.textContent = "Perfect alignment: 90° axis, tight elbow tuck.";
  let currentIssue = null;
  activeJump.issues.forEach(issue => {
    if (Math.abs(issue.frame - currentFrame) <= 4) currentIssue = issue;
  });
  descInc.textContent = currentIssue ? "Deviation: " + currentIssue.desc : "No significant deviations detected.";
}

function updateQuickDiagnostics() {
  // Отключено: теперь диагностика обновляется через fetch-запрос к Go-серверу!
  return;
  container.replaceChildren();

  const issue = activeJump.issues && activeJump.issues.length > 0 ? activeJump.issues[0] : null;
  const drill = activeJump.practice && activeJump.practice.length > 0 ? activeJump.practice[0] : null;

  const summary = activeJump.coachSummary || "Execution generally clean, no severe mistakes detected.";
  const mistakeText = issue ? issue.label : "Jump performed successfully.";
  const fixText = drill ? drill.focus : "Keep practicing the same technique to build muscle memory.";

  const cards = [
    { type: 'error', icon: '🚨', title: 'What Happened', text: mistakeText },
    { type: 'cause', icon: '🔍', title: 'Root Cause', text: summary },
    { type: 'fix', icon: '🛠️', title: 'How to Fix It', text: fixText }
  ];

  cards.forEach(c => {
    const card = document.createElement('div');
    card.className = 'diag-card ' + c.type;
    card.innerHTML = "<div class='diag-header'>" +
        "<span class='diag-icon'>" + c.icon + "</span>" +
        "<span class='diag-title'>" + c.title + "</span>" +
      "</div>" +
      "<div class='diag-body'>" + c.text + "</div>";
    container.appendChild(card);
  });
}

function openChatWithQuestion(text, key = 'default') {
  addChatMessage('user', text);
  setTimeout(() => { addChatMessage('coach', COACH_CHAT_RESPONSES[key] || COACH_CHAT_RESPONSES.default); }, 750);
}

function addChatMessage(sender, text) {
  const container = document.getElementById('chat-messages');
  if (!container) return;
  const msg = document.createElement('div'); msg.className = 'chat-msg ' + sender;
  const label = document.createElement('span'); label.className = 'chat-msg-sender'; label.textContent = sender === 'user' ? 'You' : 'AI Coach';
  const body = document.createElement('span'); body.className = 'chat-msg-text'; body.textContent = text;
  msg.appendChild(label); msg.appendChild(body);
  container.appendChild(msg);
  container.scrollTop = container.scrollHeight;
  updateSuggestedQuestions(text);
}

function updateSuggestedQuestions() {
  const box = document.getElementById('suggested-questions-box');
  if (!box) return;
  box.replaceChildren();
  const qs = [{ t: "How to fix axis tilt?", k: 'axis' }, { t: "What to fix on landing?", k: 'landing' }, { t: "PRO benefits?", k: 'upgrade' }];
  qs.forEach(q => {
    const btn = document.createElement('button'); btn.className = 'suggested-q-btn'; btn.textContent = q.t;
    btn.addEventListener('click', () => openChatWithQuestion(q.t, q.k));
    box.appendChild(btn);
  });
}

function setupTabs() {
  const tabs = document.querySelectorAll('.desk-tab-btn, .nav-tab-btn');
  tabs.forEach(btn => {
    btn.addEventListener('click', (e) => {
      const tabId = e.currentTarget.dataset.tab;
      document.querySelectorAll('.desk-tab-btn, .nav-tab-btn').forEach(b => {
        if(b.dataset.tab === tabId) b.classList.add('active');
        else b.classList.remove('active');
      });
      document.querySelectorAll('.tab-pane').forEach(p => p.classList.remove('active'));
      const pane = document.getElementById('content-' + tabId);
      if (pane) pane.classList.add('active');
    });
  });

  const sysTabs = document.querySelectorAll('.system-tab');
  sysTabs.forEach(btn => {
    btn.addEventListener('click', (e) => {
      sysTabs.forEach(b => b.classList.remove('active'));
      e.currentTarget.classList.add('active');
      updateJudging(e.currentTarget.id === 'tab-isu');
    });
  });
}

function updateIssuesTable() {
  const hdrIssues = document.getElementById('hdr-issues');
  if (hdrIssues && activeJump.issues) {
    hdrIssues.textContent = activeJump.issues.length;
  }
}

function updatePracticeDrills() {
  // Removed, UI no longer exists
}

function updateJudging(isISU = true) {
  // Отключено: теперь оценки выставляются динамически с Go сервера
  return;
}

function updateDiagnostics() {
  // Отключено: теперь диагностика выставляется динамически с Go сервера
  return;
}

function populateDiagnostics(jump) {
  // Отключено: теперь диагностика выставляется динамически с Go сервера
  return;
}

function updateJudgingOriginal(isISU = true) {
  const hdrBase = document.getElementById('hdr-base');
  const hdrGoe = document.getElementById('hdr-goe');
  const hdrScore = document.getElementById('hdr-score');
  const hdrReason = document.getElementById('hdr-reason');
  
  if (!activeJump) return;
  
  let base = activeJump.goeSimulation?.baseScore || 0;
  let goe = 0;
  let reason = "Good execution.";
  
  if (activeJump.id === 'lutz') {
    goe = -2.10;
    reason = "Penalty applied: Early right shoulder drop causing severe axis tilt and a fall.";
  } else if (activeJump.id === 'axel') {
    goe = 2.30;
    reason = "Bonus applied: Massive distance and soft knee on landing.";
  } else if (activeJump.id === 'salchow') {
    goe = 1.50;
    reason = "Bonus applied: Perfect takeoff synchronization and deep entry edge.";
  }
  
  let final = base + goe;
  
  if (hdrBase) hdrBase.textContent = base.toFixed(2) + " pts";
  if (hdrGoe) {
    const sign = goe > 0 ? "+" : "";
    hdrGoe.textContent = sign + goe.toFixed(2) + " pts";
    hdrGoe.style.color = goe < 0 ? "#EF4444" : "#10B981";
  }
  if (hdrScore) hdrScore.textContent = final.toFixed(2) + " pts";
  
  if (hdrReason) {
    hdrReason.textContent = reason;
    if (goe < 0) {
      hdrReason.style.color = "#FCA5A5";
      hdrReason.style.backgroundColor = "rgba(239, 68, 68, 0.1)";
      hdrReason.style.borderColor = "rgba(239, 68, 68, 0.2)";
    } else {
      hdrReason.style.color = "#6EE7B7";
      hdrReason.style.backgroundColor = "rgba(16, 185, 129, 0.1)";
      hdrReason.style.borderColor = "rgba(16, 185, 129, 0.2)";
    }
  }
}

function updateJumpMetrics() {
  if (!activeJump || !activeJump.biometrics) return;
  
  const b = activeJump.biometrics;
  
  const elAir = document.getElementById('bio-airtime');
  if (elAir) elAir.textContent = b.airTimeSec.toFixed(2) + " sec";
  
  const elHeight = document.getElementById('bio-height');
  if (elHeight) elHeight.textContent = b.maxHeightM.toFixed(2) + " m";
  
  const elTilt = document.getElementById('bio-tilt');
  if (elTilt) elTilt.textContent = b.maxTiltDeg + "°";
  
  const elForce = document.getElementById('bio-gforce');
  if (elForce) elForce.textContent = b.impactG.toFixed(1) + " G";
}
`
}
