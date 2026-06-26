import gdown
import os
import sys

# URL and ID for Google Drive
DRIVE_FOLDER_ID = '1WzERNs04uo_5xjybfKcXYOC9v8KL6Hk2'
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUTPUT_DIR = os.path.join(BASE_DIR, 'data', 'Video_data', 'dataset_videos')

def main():
    print("Starting download of training dataset from Google Drive...")
    
    if not os.path.exists(OUTPUT_DIR):
        os.makedirs(OUTPUT_DIR, exist_ok=True)
        
    try:
        # Download folder from Google Drive
        gdown.download_folder(id=DRIVE_FOLDER_ID, output=OUTPUT_DIR, quiet=False, use_cookies=False)
        print("Download finished successfully!")
    except Exception as e:
        print(f"Error during download: {e}")
        print("Please check your internet connection or install gdown properly.")
        sys.exit(1)

if __name__ == '__main__':
    main()
