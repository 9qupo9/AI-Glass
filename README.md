# SkateEye AI (AI-Glass)

SkateEye AI is an advanced biometrics and pose estimation system for figure skating jump analysis. By processing video input, the system automatically detects the skater's takeoff, flight phase, and landing, and classifies the type of jump (e.g., Axel, Lutz, Flip, Loop, Salchow, Toe Loop).

## Features

*   **Video Processing Pipeline**: Extracts 2D and 3D keypoints from figure skating videos.
*   **Biometrics Extraction**: Analyzes body posture, rotations, takeoff foot, and edge types (Inside/Outside).
*   **Jump Classification**: Implements a physics-based logic (in Go) to determine the exact jump performed based on the takeoff parameters (Edge/Toe, Forward/Backward, Left/Right foot, Inside/Outside edge).
*   **Hybrid Architecture**: Uses Python for heavy Machine Learning inference (Pose Estimation, 3D Lifting) and Go for the fast and reliable classification logic and web server.

## Technologies Used

*   **Python**: Neural network inference, 2D/3D pose extraction, and biomechanical mathematical analysis (`numpy`, etc.).
*   **Go**: High-performance backend server, rule-based jump classification (`Truth Table`), and REST API.

## Project Structure

*   `components/classifier.go`: Go module containing the truth table for jump classification based on `AdvancedBiometrics`.
*   `Neiron/data/Video_data/demo/biometrics.py`: Python module that calculates physical attributes from keypoints (support foot, edge angle, takeoff type, rotation).

## How it works

1.  **Pose Estimation**: A video is processed to extract frame-by-frame 2D and 3D skeleton keypoints.
2.  **Biometrics Calculation**: Python scripts analyze the trajectory, calculating the knee angles, hip vectors, and ankle tilts to precisely determine the take-off edge (Inside/Outside) and direction (Forward/Backward).
3.  **Classification**: The extracted metrics are sent to the Go backend, which runs through a definitive truth table to output the final jump verdict (e.g., Axel, Flip, Lutz).