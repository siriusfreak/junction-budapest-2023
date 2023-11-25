from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import JSONResponse
from tempfile import TemporaryDirectory
from ultralyticsplus import YOLO
import cv2
import torch
import base64
import shutil
import os

if torch.cuda.is_available():
    DEVICE = 'cuda'
elif torch.backends.mps.is_built():
    DEVICE = 'mps'  # Для Mac с M1 chip
else:
    DEVICE = 'cpu'
print("model device: ", DEVICE)

MAX_VIDEO_LENGHT_SECOND = os.getenv('MAX_VIDEO_LENGHT_SECOND', 30)
MAX_FPS = os.getenv('MAX_FPS', 30)

model = YOLO('ultralyticsplus/yolov8l.pt')
model.overrides['conf'] = float(os.getenv('CONF_THRESHOLD', 0.25))
model.overrides['iou'] = float(os.getenv('IOU_THRESHOLD', 0.45))
model.overrides['agnostic_nms'] = os.getenv('AGNOSTIC_NMS', 'False') == 'True'
model.overrides['max_det'] = int(os.getenv('MAX_DET', 1000))
model.to(DEVICE)
app = FastAPI()

@app.post("/upload/")
async def upload_file(
    video: UploadFile = File(...),
    processed_percent: int = 1,
    confidence_threshold: float = 0.3,
    skip_milliseconds: int = 1,
):
    """
    This endpoint receives a video file and processes it to detect frames with exactly one person.

    Args:
    - video: A video file to be uploaded.
    - processed_percent: The percent of frames to process.
    - confidence_threshold: Confidence threshold for person detection.
    - skip_milliseconds: Milliseconds to skip after a non-compliant frame is found.

    Returns:
    - frames: Base64 encoded frames where exactly one person was detected.
    - total_frames: The total number of frames in the video.
    - processed_frames: The total number of frames processed.
    """
    if processed_percent > 100 or processed_percent <= 0:
        raise HTTPException(status_code=400, detail="Processed percent must be between 0 and 100.")

    with TemporaryDirectory() as temp_dir:
        tmp_path = os.path.join(temp_dir, 'temp_video.mp4')
        with open(tmp_path, 'wb') as tmp_file:
            shutil.copyfileobj(video.file, tmp_file)
        cap = cv2.VideoCapture(tmp_path)
    try:
        if not cap.isOpened():
            raise HTTPException(status_code=400, detail="Unable to read video file.")
        fps = cap.get(cv2.CAP_PROP_FPS)
        if fps == 0:
            raise HTTPException(status_code=400, detail="FPS of video is zero, which indicates a problem with the video file.")

        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        video_length = total_frames / fps

        if video_length > MAX_VIDEO_LENGHT_SECOND:
            raise HTTPException(status_code=400, detail=f"Video is too long. Maximum length allowed is {MAX_VIDEO_LENGHT_SECOND} seconds.")

        if fps > MAX_FPS:
            raise HTTPException(status_code=400, detail=f"Video FPS is too high. Maximum FPS allowed is {MAX_FPS}.")

        frame_step = max(int(100 / processed_percent), 1)
        skip_frames = max(int(fps * skip_milliseconds / 1000), 1)

        frames_to_return = []
        frame_count = 0
        processed_count = 0

        while True:
            ret, frame = cap.read()
            if not ret:
                break

            if frame_count % frame_step == 0:
                processed_count += 1

                result = model.predict(frame)[0]  # Для Mac с M1 chip: model.predict(frame, device=DEVICE)[0]
                confidences = result.boxes.conf
                categories = result.boxes.cls

                n_people = sum(1 for i in range(len(confidences)) if confidences[i] >= confidence_threshold and categories[i] == 0)

                if n_people != 1:
                    ret, buffer = cv2.imencode('.jpg', frame)
                    frame_base64 = base64.b64encode(buffer).decode('utf-8')
                    frames_to_return.append(frame_base64)

                    new_position = frame_count + skip_frames
                    if new_position < total_frames:
                        cap.set(cv2.CAP_PROP_POS_FRAMES, new_position)
                        frame_count = new_position - 1
                    else:
                        break

            frame_count += 1
    finally:
        cap.release()

    return JSONResponse(content={
        "frames": frames_to_return,
        "total_frames": total_frames,
        "processed_frames": processed_count,
    })
