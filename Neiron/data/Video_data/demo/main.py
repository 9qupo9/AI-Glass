import sys
import argparse
import cv2
import os
import numpy as np
import torch
import glob
from tqdm import tqdm
sys.path.append(os.getcwd())
from common.camera import *
from model.strided_transformer import Model
import copy

import matplotlib
import matplotlib.pyplot as plt
import matplotlib.gridspec as gridspec

import warnings
warnings.simplefilter('ignore', RuntimeWarning)
warnings.simplefilter('ignore', UserWarning)
plt.switch_backend('agg')
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42

import itertools
from sklearn.preprocessing import StandardScaler
import pickle


# whole body 17keys
KEY = ["Hips", "R_UpLeg", "R_Leg", "R_Foot", "L_UpLeg", "L_Leg", "L_Foot", "Spine",
       "Neck", "Neck1", "Head", "R_Arm", "R_ForeArm", "R_Hand", "L_Arm", "L_ForeArm", "L_Hand"]
# lower body 7keys
NOT_KEY = ["Hips", "R_UpLeg", "R_Leg", "R_Foot", "L_UpLeg", "L_Leg", "L_Foot"]
# number of keypoints
N_KEY = len(KEY)
# number of skaters
N_SKATER = 6


def feature_extraction(path_list, fps=60, n_key=N_KEY):
    step = int(240 / fps)
    keypoints = []
    for path in tqdm(path_list, leave=False):
        keypoint = list(np.load(path)["reconstruction"])
        keypoint = keypoint[::step]
        if n_key == 7:
            # 下半身のkeypoint、7点のみ使用の場合
            for f in range(len(keypoint)):
                keypoint[f] = keypoint[f][0:7]
        keypoint_1d = list(itertools.chain.from_iterable(list(itertools.chain.from_iterable(keypoint))))
        
        # Replace NaNs with 0.0 to prevent StandardScaler from crashing
        arr_1d = np.array(keypoint_1d, dtype=np.float32)
        arr_1d = np.nan_to_num(arr_1d, nan=0.0, posinf=0.0, neginf=0.0)
            
        keypoints.append(arr_1d)
    keypoints_array = np.array(keypoints)
    return keypoints_array


def show2Dpose(kps, img):
    connections = [[0, 1], [1, 2], [2, 3], [0, 4], [4, 5],
                   [5, 6], [0, 7], [7, 8], [8, 9], [9, 10],
                   [8, 11], [11, 12], [12, 13], [8, 14], [14, 15], [15, 16]]

    LR = np.array([0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0], dtype=bool)

    lcolor = (255, 0, 0)
    rcolor = (0, 0, 255)
    thickness = 5

    for j, c in enumerate(connections):
        start = map(int, kps[c[0]])
        end = map(int, kps[c[1]])
        start = list(start)
        end = list(end)
        cv2.line(img, (start[0], start[1]), (end[0], end[1]), lcolor if LR[j] else rcolor, thickness)
        cv2.circle(img, (start[0], start[1]), thickness=-1, color=(0, 255, 0), radius=5)
        cv2.circle(img, (end[0], end[1]), thickness=-1, color=(0, 255, 0), radius=5)

    return img


def show3Dpose(vals, ax):
    ax.view_init(elev=15., azim=70)

    I = np.array([0, 0, 1, 4, 2, 5, 0, 7,  8,  8, 14, 15, 11, 12, 8,  9])
    J = np.array([1, 4, 2, 5, 3, 6, 7, 8, 14, 11, 15, 16, 12, 13, 9, 10])

    LR = np.array([0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0], dtype=bool)

    for i in np.arange(len(I)):
        x, y, z = [np.array([vals[I[i], j], vals[J[i], j]]) for j in range(3)]
        ax.plot(x, y, z, lw=2)
        ax.scatter(x, y, z)

    if vals[3][2] == 0:
        ax.plot(vals[3][0], vals[3][1], vals[3][2], marker='^', markersize=10)
    if vals[6][2] == 0:
        ax.plot(vals[6][0], vals[6][1], vals[6][2], marker='^', markersize=10)

    RADIUS = 0.8

    ax.set_xlim3d([-RADIUS, RADIUS])
    ax.set_ylim3d([-RADIUS, RADIUS])
    ax.set_aspect('auto')

    white = (1.0, 1.0, 1.0, 0.0)
    ax.xaxis.set_pane_color(white)
    ax.yaxis.set_pane_color(white)
    ax.zaxis.set_pane_color(white)

    ax.tick_params('x', labelbottom=False)
    ax.tick_params('y', labelleft=False)
    ax.tick_params('z', labelleft=False)


def h36m_coco_format(keypoints, scores):
    assert len(keypoints.shape) == 4
    _, T, _, _ = keypoints.shape
    new_kpts = np.zeros((1, T, 17, 2))
    new_scores = np.zeros((1, T, 17))
    
    # 0: Pelvis = (LHip + RHip) / 2
    new_kpts[:, :, 0] = (keypoints[:, :, 11] + keypoints[:, :, 12]) / 2
    new_scores[:, :, 0] = (scores[:, :, 11] + scores[:, :, 12]) / 2
    
    # 1: RHip, 2: RKnee, 3: RAnkle
    new_kpts[:, :, 1:4] = keypoints[:, :, 12:15]
    new_scores[:, :, 1:4] = scores[:, :, 12:15]
    
    # 4: LHip, 5: LKnee, 6: LAnkle
    new_kpts[:, :, 4] = keypoints[:, :, 11]
    new_kpts[:, :, 5] = keypoints[:, :, 13]
    new_kpts[:, :, 6] = keypoints[:, :, 15]
    new_scores[:, :, 4] = scores[:, :, 11]
    new_scores[:, :, 5] = scores[:, :, 13]
    new_scores[:, :, 6] = scores[:, :, 15]
    
    # 8: Neck = (LShoulder + RShoulder) / 2
    new_kpts[:, :, 8] = (keypoints[:, :, 5] + keypoints[:, :, 6]) / 2
    new_scores[:, :, 8] = (scores[:, :, 5] + scores[:, :, 6]) / 2
    
    # 7: Spine = (Pelvis + Neck) / 2
    new_kpts[:, :, 7] = (new_kpts[:, :, 0] + new_kpts[:, :, 8]) / 2
    new_scores[:, :, 7] = (new_scores[:, :, 0] + new_scores[:, :, 8]) / 2
    
    # 9: Nose
    new_kpts[:, :, 9] = keypoints[:, :, 0]
    new_scores[:, :, 9] = scores[:, :, 0]
    
    # 10: Head = Nose + (Nose - Neck)
    new_kpts[:, :, 10] = keypoints[:, :, 0] + (keypoints[:, :, 0] - new_kpts[:, :, 8])
    new_scores[:, :, 10] = scores[:, :, 0]
    
    # 11: LShoulder, 12: LElbow, 13: LWrist
    new_kpts[:, :, 11] = keypoints[:, :, 5]
    new_kpts[:, :, 12] = keypoints[:, :, 7]
    new_kpts[:, :, 13] = keypoints[:, :, 9]
    new_scores[:, :, 11] = scores[:, :, 5]
    new_scores[:, :, 12] = scores[:, :, 7]
    new_scores[:, :, 13] = scores[:, :, 9]
    
    # 14: RShoulder, 15: RElbow, 16: RWrist
    new_kpts[:, :, 14] = keypoints[:, :, 6]
    new_kpts[:, :, 15] = keypoints[:, :, 8]
    new_kpts[:, :, 16] = keypoints[:, :, 10]
    new_scores[:, :, 14] = scores[:, :, 6]
    new_scores[:, :, 15] = scores[:, :, 8]
    new_scores[:, :, 16] = scores[:, :, 10]
    
    return new_kpts, new_scores

def get_pose2D(video_path, output_dir, cut_frame=160):
    from ultralytics import YOLO
    
    model_path = "yolo11n-pose.pt"
    print(f'Loading {model_path} and generating 2D pose...')
    
    # Load model
    model = YOLO(model_path)
    
    # Run tracking
    results = model.track(video_path, persist=True, tracker="bytetrack.yaml", verbose=False)
    
    track_dict = {}
    total_frames = len(results)
    
    for frame_idx, r in enumerate(results):
        if r.boxes is None or r.boxes.id is None:
            continue
        ids = r.boxes.id.int().cpu().tolist()
        kpts = r.keypoints.xy.cpu().numpy()
        confs = r.keypoints.conf.cpu().numpy()
        boxes = r.boxes.xyxy.cpu().numpy()
        for i, track_id in enumerate(ids):
            if track_id not in track_dict:
                track_dict[track_id] = {'frames': [], 'kpts': [], 'confs': [], 'y1': []}
            track_dict[track_id]['frames'].append(frame_idx)
            track_dict[track_id]['kpts'].append(kpts[i])
            track_dict[track_id]['confs'].append(confs[i])
            track_dict[track_id]['y1'].append(boxes[i][1])
            
    if not track_dict:
        raise Exception("YOLOv11 не нашел фигуриста на видео!")

    # Find jumper by looking for the highest jump apex (min Y coordinate)
    best_id = max(track_dict.keys(), key=lambda k: len(track_dict[k]['frames']))
    jump_apex_frame = None
    
    # Evaluate locus for each track to find a jump
    for t_id, data in track_dict.items():
        if len(data['frames']) < 30: continue
        y1 = np.array(data['y1'])
        apex_idx = np.argmin(y1)
        if 0 < apex_idx < len(y1)-1:
            jump_apex_frame = data['frames'][apex_idx]
            best_id = t_id
            break
            
    if jump_apex_frame is None:
        jump_apex_frame = total_frames // 2
        
    start_frame = max(0, jump_apex_frame - 80)
    end_frame = start_frame + cut_frame
    
    # Reconstruct 160 frame window
    final_kpts = np.zeros((1, cut_frame, 17, 2))
    final_confs = np.zeros((1, cut_frame, 17))
    
    data = track_dict[best_id]
    for i, frame_idx in enumerate(data['frames']):
        if start_frame <= frame_idx < end_frame:
            idx = frame_idx - start_frame
            final_kpts[0, idx] = data['kpts'][i]
            final_confs[0, idx] = data['confs'][i]
            
    # Convert COCO to H36M
    re_kpts, _ = h36m_coco_format(final_kpts, final_confs)

    output_dir += 'input_2D/'
    os.makedirs(output_dir, exist_ok=True)
    output_npz = output_dir + 'keypoints.npz'
    np.savez_compressed(output_npz, reconstruction=re_kpts)
    
    print('Generating 2D pose successful!')
    return start_frame

def img2video(video_path, output_dir):
    cap = cv2.VideoCapture(video_path)
    fps = int(cap.get(cv2.CAP_PROP_FPS)) + 5

    fourcc = cv2.VideoWriter_fourcc(*"mp4v")

    names = sorted(glob.glob(os.path.join(output_dir + 'pose/', '*.png')))
    img = cv2.imread(names[0])
    size = (img.shape[1], img.shape[0])

    videoWrite = cv2.VideoWriter(output_dir + video_name + '.mp4', fourcc, fps, size)

    for name in names:
        img = cv2.imread(name)
        videoWrite.write(img)

    videoWrite.release()


def showimage(ax, img):
    ax.set_xticks([])
    ax.set_yticks([])
    plt.axis('off')
    ax.imshow(img)


def get_pose3D(video_path, output_dir, start_frame=0, cut_frame=160):
    args, _ = argparse.ArgumentParser().parse_known_args()
    args.layers, args.channel, args.d_hid, args.frames = 3, 256, 512, 351
    args.stride_num = [3, 9, 13]
    args.pad = (args.frames - 1) // 2
    args.previous_dir = 'checkpoint/pretrained'
    args.n_joints, args.out_joints = 17, 17

    # Reload
    device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
    model = Model(args).to(device)

    model_dict = model.state_dict()
    model_paths = sorted(glob.glob(os.path.join(args.previous_dir, '*.pth')))
    model_path = None
    for path in model_paths:
        if os.path.basename(path).startswith('n'):
            model_path = path

    pre_dict = torch.load(model_path, map_location=device)
    for name, key in model_dict.items():
        model_dict[name] = pre_dict[name]
    model.load_state_dict(model_dict)

    model.eval()

    # input
    keypoints = np.load(output_dir + 'input_2D/keypoints.npz', allow_pickle=True)['reconstruction']

    cap = cv2.VideoCapture(video_path)
    video_length = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    cap.set(cv2.CAP_PROP_POS_FRAMES, start_frame)

    # 3D
    post_outs = []
    print('\nGenerating 3D pose...')
    actual_frames = min(cut_frame, len(keypoints[0]))
    for i in tqdm(range(actual_frames), leave=False):
        ret, img = cap.read()
        if not ret:
            break
        img_size = img.shape

        # input frames
        start = max(0, i - args.pad)
        end = min(i + args.pad, len(keypoints[0])-1)

        input_2D_no = keypoints[0][start:end+1]

        left_pad, right_pad = 0, 0
        if input_2D_no.shape[0] != args.frames:
            if i < args.pad:
                left_pad = args.pad - i
            if i > len(keypoints[0]) - args.pad - 1:
                right_pad = i + args.pad - (len(keypoints[0]) - 1)

            input_2D_no = np.pad(input_2D_no, ((left_pad, right_pad), (0, 0), (0, 0)), 'edge')

        joints_left = [4, 5, 6, 11, 12, 13]
        joints_right = [1, 2, 3, 14, 15, 16]

        input_2D = normalize_screen_coordinates(input_2D_no, w=img_size[1], h=img_size[0])

        input_2D_aug = copy.deepcopy(input_2D)
        input_2D_aug[:, :, 0] *= -1
        input_2D_aug[:, joints_left + joints_right] = input_2D_aug[:, joints_right + joints_left]
        input_2D = np.concatenate((np.expand_dims(input_2D, axis=0), np.expand_dims(input_2D_aug, axis=0)), 0)

        input_2D = input_2D[np.newaxis, :, :, :, :]

        input_2D = torch.from_numpy(input_2D.astype('float32')).to(device)

        N = input_2D.size(0)

        # estimation
        output_3D_non_flip, _ = model(input_2D[:, 0])
        output_3D_flip, _ = model(input_2D[:, 1])

        output_3D_flip[:, :, :, 0] *= -1
        output_3D_flip[:, :, joints_left + joints_right, :] = output_3D_flip[:, :, joints_right + joints_left, :]

        output_3D = (output_3D_non_flip + output_3D_flip) / 2

        output_3D[:, :, 0, :] = 0
        post_out = output_3D[0, 0].cpu().detach().numpy()

        rot = [0.1407056450843811, -0.1500701755285263, -0.755240797996521, 0.6223280429840088]
        rot = np.array(rot, dtype='float32')
        post_out = camera_to_world(post_out, R=rot, t=0)
        post_out[:, 2] -= np.min(post_out[:, 2])
        post_outs.append(post_out)

        input_2D_no = input_2D_no[args.pad]

        # 2D
        # image = show2Dpose(input_2D_no, copy.deepcopy(img))
        # output_dir_2D = output_dir + 'pose2D/'
        # os.makedirs(output_dir_2D, exist_ok=True)
        # cv2.imwrite(output_dir_2D + str(('%04d' % i)) + '_2D.png', image)

        # 3D
        # fig = plt.figure(figsize=(9.6, 5.4))
        # gs = gridspec.GridSpec(1, 1)
        # gs.update(wspace=-0.00, hspace=0.05)
        # ax = plt.subplot(gs[0], projection='3d')
        # show3Dpose(post_out, ax)
        # output_dir_3D = output_dir + 'pose3D/'
        # os.makedirs(output_dir_3D, exist_ok=True)
        # plt.savefig(output_dir_3D + str(('%04d' % i)) + '_3D.png', dpi=200, format='png', bbox_inches='tight')
        # plt.close(fig) # Prevent Matplotlib memory leak!
        
    cap.release() # Prevent OpenCV memory leak!
    keypoints_3d = np.array(post_outs)
    os.makedirs(output_dir, exist_ok=True)
    output_npz = output_dir + 'keypoints.npz'
    np.savez_compressed(output_npz, reconstruction=keypoints_3d)
    print('Generating 3D pose successful!')

    # Rendering removed for speed
    
    return keypoints_3d


# Legacy __main__ execution removed as we now use jump3d_analyzer.py
