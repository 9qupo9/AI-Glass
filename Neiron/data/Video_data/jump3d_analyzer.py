import sys
import os
import json
import numpy as np
import warnings
warnings.filterwarnings("ignore")

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))
os.chdir(CURRENT_DIR)
sys.path.append(CURRENT_DIR)
sys.path.append(os.path.join(CURRENT_DIR, "demo"))

from demo.main import get_pose2D, get_pose3D
from demo.biometrics import determine_biometrics

def analyze_video(video_path):
    if not os.path.exists(video_path):
        print(json.dumps({"error": "Video file not found"}))
        return
        
    video_name = os.path.basename(video_path).split('.')[0]
    output_dir = os.path.join(CURRENT_DIR, 'demo', 'output', video_name) + os.sep
    os.makedirs(output_dir, exist_ok=True)
    
    # Очищаем sys.argv, чтобы внутренние парсеры (yacs/argparse) в hrnet не пытались прочитать путь к видео как аргумент командной строки
    sys.argv = [sys.argv[0]]
    
    # 1. 2D Pose Estimation (YOLO11)
    start_frame = get_pose2D(video_path, output_dir)
    
    # Load the generated 2D keypoints
    import numpy as np
    kpts_2d = np.load(output_dir + 'input_2D/keypoints.npz', allow_pickle=True)['reconstruction']
    
    # kpts_2d shape is usually (1, frames, 17, 2). Squeeze batch dim.
    if len(kpts_2d.shape) == 4:
        kpts_2d = kpts_2d[0]
        
    # 2. 3D Pose Estimation (VideoPose3D)
    kpts_3d = get_pose3D(video_path, output_dir, start_frame, 160)
        
    # 3. Extract Biometrics using 2D + 3D math
    features = determine_biometrics(kpts_2d, kpts_3d, start_frame)

    # Формируем timeline для отрисовки зеленого человечка
    timeline = []
    joint_names = [
        "MidHip", "RHip", "RKnee", "RAnkle",
        "LHip", "LKnee", "LAnkle",
        "Spine", "Neck", "Nose", "Head",
        "LShoulder", "LElbow", "LWrist",
        "RShoulder", "RElbow", "RWrist"
    ]
    for frame_idx in range(len(kpts_2d)):
        joints_dict = {}
        for i, name in enumerate(joint_names):
            joints_dict[name] = {"x": float(kpts_2d[frame_idx, i, 0]), "y": float(kpts_2d[frame_idx, i, 1])}
        timeline.append({"joints": joints_dict})

    # Печатаем чистый JSON в stdout для парсера на Go
    print(json.dumps({
        "status": "success",
        "advanced_biometrics": features,
        "timeline": timeline
    }))

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "Video path required"}))
        sys.exit(1)
        
    analyze_video(sys.argv[1])
