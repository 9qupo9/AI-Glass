@echo off
cd /d "%~dp0"
echo ========================================================
echo JudgeAI-LutzEdge Setup and Training
echo ========================================================
echo.
echo Installing requirements...
pip install gdown pandas scikit-learn numpy opencv-python torch torchvision torchaudio tqdm matplotlib
if %errorlevel% neq 0 (
    echo Error: pip failed. Please make sure Python and pip are installed and in your PATH.
    pause
    exit /b %errorlevel%
)

echo.
echo Downloading dataset from Google Drive...
python setup_and_download.py
if %errorlevel% neq 0 (
    echo Error: Download script failed.
    pause
    exit /b %errorlevel%
)

echo.
echo Training Logistic Regression Model and saving weights...
python train_and_save_model.py
if %errorlevel% neq 0 (
    echo Error: Training failed.
    pause
    exit /b %errorlevel%
)

echo.
echo Setup and Training Completed Successfully!
echo The new classifier_wrapper.py is now ready to use.
pause
