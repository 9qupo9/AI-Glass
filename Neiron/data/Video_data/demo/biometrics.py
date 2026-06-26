import numpy as np
from scipy.signal import savgol_filter
import json
import os
import math
import statistics
from collections import deque

class BiomechProfiler:
    def __init__(self, history_file="profiler_history.json", max_history=15):
        self.history_file = history_file
        self.max_history = max_history
        self.history = deque(maxlen=max_history)
        self.base_threshold = 0.20
        self.load_history()

    def load_history(self):
        if os.path.exists(self.history_file):
            try:
                with open(self.history_file, 'r') as f:
                    data = json.load(f)
                    for val in data.get("history", []):
                        self.history.append(float(val))
            except Exception:
                pass

    def save_history(self):
        try:
            with open(self.history_file, 'w') as f:
                json.dump({"history": list(self.history)}, f)
        except Exception:
            pass

    def update_and_calculate(self, current_tilt):
        if math.isnan(current_tilt):
            return self.base_threshold

        self.history.append(current_tilt)
        self.save_history()
        
        if len(self.history) < 3:
            return self.base_threshold
            
        return statistics.median(self.history)

    def score_confidence(self, current_tilt, threshold):
        if math.isnan(current_tilt):
            return 0.0, 0.0

        delta = abs(current_tilt - threshold)
        if delta < 0.05:
            confidence = max(0.50, delta / 0.05 * 0.80) 
        else:
            confidence = min(0.99, 0.80 + (delta / 0.5) * 0.19)
            
        return round(confidence, 2), round(delta, 4)

# Файл сохраняем в папке Video_data, на уровень выше демо
history_path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "profiler_history.json")
profiler = BiomechProfiler(history_file=history_path)

# Индексы ключевых точек (YOLO11 Pose format)
IDX_HIPS, IDX_R_HIP, IDX_R_KNEE, IDX_R_ANKLE = 0, 1, 2, 3
IDX_L_HIP, IDX_L_KNEE, IDX_L_ANKLE = 4, 5, 6
IDX_NECK, IDX_NOSE, IDX_HEAD = 8, 9, 10
IDX_L_SHOULDER, IDX_R_SHOULDER = 11, 14

def find_flight_phase(keypoints_2d):
    feet_y = (keypoints_2d[:, IDX_R_ANKLE, 1] + keypoints_2d[:, IDX_L_ANKLE, 1]) / 2.0
    window_len = min(7, len(feet_y))
    if window_len % 2 == 0: window_len -= 1
    smoothed = savgol_filter(feet_y, window_length=window_len, polyorder=2) if window_len >= 3 else feet_y
    peak = np.argmin(smoothed)
    height = np.abs(np.mean(keypoints_2d[:, IDX_HIPS, 1]) - np.mean(keypoints_2d[:, IDX_HEAD, 1]))
    threshold = np.max(smoothed) - height * 0.05
    takeoff = next((i for i in range(peak, 0, -1) if smoothed[i] >= threshold), 0)
    landing = next((i for i in range(peak, len(smoothed)) if smoothed[i] >= threshold), len(smoothed)-1)
    return takeoff, landing

def smooth_keypoints_xy(keypoints_2d):
    smoothed = keypoints_2d.copy()
    window_len = 7
    for idx in range(keypoints_2d.shape[1]):
        if window_len < len(smoothed):
            smoothed[:, idx, 0] = savgol_filter(smoothed[:, idx, 0], window_length=window_len, polyorder=2)
            smoothed[:, idx, 1] = savgol_filter(smoothed[:, idx, 1], window_length=window_len, polyorder=2)
    return smoothed

def calculate_rotations(keypoints_2d, takeoff, landing):
    total_degrees = 0.0
    last_angle = None
    for i in range(takeoff, landing):
        dx, dy = keypoints_2d[i, IDX_L_SHOULDER, 0] - keypoints_2d[i, IDX_R_SHOULDER, 0], \
                 keypoints_2d[i, IDX_L_SHOULDER, 1] - keypoints_2d[i, IDX_R_SHOULDER, 1]
        current_angle = np.arctan2(dy, dx)
        if last_angle is not None:
            diff = (current_angle - last_angle + np.pi) % (2 * np.pi) - np.pi
            total_degrees += np.abs(np.degrees(diff))
        last_angle = current_angle
    rotations = round((total_degrees / 360.0) * 4) / 4
    return float(rotations), float(rotations * 360.0)

def calculate_joint_angle(p1, p2, p3):
    v1, v2 = np.array([p1[0]-p2[0], p1[1]-p2[1]]), np.array([p3[0]-p2[0], p3[1]-p2[1]])
    return np.degrees(np.arccos(np.clip(np.dot(v1, v2) / (np.linalg.norm(v1)*np.linalg.norm(v2) + 1e-6), -1.0, 1.0)))

def determine_takeoff_type(keypoints_2d, takeoff):
    start = max(0, takeoff - 15)
    max_l_diff, max_r_diff = 0, 0
    for i in range(start, takeoff):
        l_knee = calculate_joint_angle(keypoints_2d[i, IDX_L_HIP], keypoints_2d[i, IDX_L_KNEE], keypoints_2d[i, IDX_L_ANKLE])
        r_knee = calculate_joint_angle(keypoints_2d[i, IDX_R_HIP], keypoints_2d[i, IDX_R_KNEE], keypoints_2d[i, IDX_R_ANKLE])
        if l_knee - r_knee > max_l_diff: max_l_diff = l_knee - r_knee
        if r_knee - l_knee > max_r_diff: max_r_diff = r_knee - l_knee
    return ("Toe", "Left Foot") if max_l_diff > 15 else ("Toe", "Right Foot") if max_r_diff > 15 else ("Edge", "None")

def determine_edge_math(keypoints_2d, keypoints_3d, takeoff):
    l_ankle_y = np.mean(keypoints_2d[max(0, takeoff-5):takeoff, IDX_L_ANKLE, 1])
    r_ankle_y = np.mean(keypoints_2d[max(0, takeoff-5):takeoff, IDX_R_ANKLE, 1])
    support_foot = "Left Foot" if l_ankle_y > r_ankle_y - 5 else "Right Foot"
    
    takeoff_type, _ = determine_takeoff_type(keypoints_2d, takeoff)
    
    window = slice(max(0, takeoff - 5), takeoff)
    hip_vec = keypoints_3d[takeoff, IDX_L_HIP] - keypoints_3d[takeoff, IDX_R_HIP]
    hip_dir = hip_vec / (np.linalg.norm(hip_vec) + 1e-6)
    
    l_tilt = np.mean([np.dot(keypoints_3d[i, IDX_L_ANKLE] - keypoints_3d[i, IDX_L_HIP], hip_dir) for i in range(window.start, window.stop)])
    r_tilt = np.mean([np.dot(keypoints_3d[i, IDX_R_ANKLE] - keypoints_3d[i, IDX_R_HIP], hip_dir) for i in range(window.start, window.stop)])
    
    # ====== АДАПТИВНАЯ ЛОГИКА ======
    current_tilt = l_tilt if support_foot == "Left Foot" else r_tilt
    adaptive_threshold = profiler.update_and_calculate(current_tilt)
    confidence, delta = profiler.score_confidence(current_tilt, adaptive_threshold)
    
    if support_foot == "Left Foot":
        edge = "Outside" if l_tilt < adaptive_threshold else "Inside"
    else:
        edge = "Outside" if r_tilt < adaptive_threshold else "Inside"
    
    return support_foot, edge, "Backward" if takeoff_type == "Toe" else "Forward", l_tilt, r_tilt, confidence, delta

def determine_biometrics(keypoints_2d, keypoints_3d, start_frame=0):
    k2d = smooth_keypoints_xy(keypoints_2d)
    takeoff, landing = find_flight_phase(k2d)
    rot, total_deg = calculate_rotations(k2d, takeoff, landing)
    supp, edge, direction, l_tilt, r_tilt, conf, delta = determine_edge_math(k2d, keypoints_3d, takeoff)
    t_type, _ = determine_takeoff_type(k2d, takeoff)
    
    def safe_float(val):
        return 0.0 if math.isnan(val) else float(val)

    return {
        "takeoff_frame": int(takeoff + start_frame),
        "landing_frame": int(landing + start_frame),
        "rotations": safe_float(rot),
        "total_degrees": safe_float(total_deg),
        "direction": str(direction),
        "takeoff_type": str(t_type),
        "support_foot": str(supp),
        "edge_type": str(edge),
        "l_tilt": safe_float(l_tilt),
        "r_tilt": safe_float(r_tilt),
        "edge_confidence": safe_float(conf),
        "biomech_delta": safe_float(delta)
    }