import io
import os
import base64
import torch
import torch.nn.functional as F
from facenet_pytorch import MTCNN, InceptionResnetV1
from PIL import Image
import cv2
from pytorch_grad_cam import GradCAM
from pytorch_grad_cam.utils.model_targets import ClassifierOutputTarget
from pytorch_grad_cam.utils.image import show_cam_on_image
from fastapi import FastAPI, UploadFile, File, HTTPException
from tempfile import TemporaryDirectory
import shutil

app = FastAPI()

DEVICE = 'cuda:0' if torch.cuda.is_available() else 'cpu'

mtcnn = MTCNN(
    select_largest=False,
    post_process=False,
    device=DEVICE
).to(DEVICE).eval()

model = InceptionResnetV1(
    pretrained="vggface2",
    classify=True,
    num_classes=1,
    device=DEVICE
)

checkpoint = torch.load("resnetinceptionv1_epoch_32.pth", map_location=torch.device('cpu'))
model.load_state_dict(checkpoint['model_state_dict'])
model.to(DEVICE)
model.eval()


@app.post("/predict")
async def predict(video: UploadFile = File(...)):    
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

        if video_length > 20:
            raise HTTPException(status_code=400, detail=f"Video is too long. Maximum length allowed is {20} seconds.")

        if fps > 60:
            raise HTTPException(status_code=400, detail=f"Video FPS is too high. Maximum FPS allowed is {60}.")

        frame_count = 0
        frames = 0
        processed_count = 0
        fake = 0
        
        while True:
            percent = (frames/total_frames * 100)
            print(f"{percent:.2f}")

            ret, frame = cap.read()
            if not ret:
                break
            
            
            frames += 1
            if frames % 2 == 0:
                continue

            frame_count += 1
            face = mtcnn(frame)
            if face is None:
                continue

            face = face.unsqueeze(0)  # add the batch dimension
            face = F.interpolate(face, size=(256, 256), mode='bilinear', align_corners=False)

            # convert the face into a numpy array to be able to plot it
            prev_face = face.squeeze(0).permute(1, 2, 0).cpu().detach().int().numpy()
            prev_face = prev_face.astype('uint8')

            face = face.to(DEVICE)
            face = face.to(torch.float32)
            face = face / 255.0
            face_image_to_plot = face.squeeze(0).permute(1, 2, 0).cpu().detach().int().numpy()

            target_layers = [model.block8.branch1[-1]]
            use_cuda = True if torch.cuda.is_available() else False
            cam = GradCAM(model=model, target_layers=target_layers, use_cuda=use_cuda)
            targets = [ClassifierOutputTarget(0)]

            grayscale_cam = cam(input_tensor=face, targets=targets, eigen_smooth=True)
            grayscale_cam = grayscale_cam[0, :]
            visualization = show_cam_on_image(face_image_to_plot, grayscale_cam, use_rgb=True)
            face_with_mask = cv2.addWeighted(prev_face, 1, visualization, 0.5, 0)

            with torch.no_grad():
                output = torch.sigmoid(model(face).squeeze(0))
                prediction = "real" if output.item() < 0.5 else "fake"

                real_prediction = 1 - output.item()
                fake_prediction = output.item()

                confidences = {
                    'real': real_prediction,
                    'fake': fake_prediction
                }

                if fake_prediction > 0.5:
                    fake+= 1
            
            processed_count += 1
            
        if processed_count / frame_count <= 0.80:
            raise HTTPException(status_code=400, detail=f"Face must be seen throughout the entire video.")

        return {"fake": fake/processed_count}

    finally:
        cap.release()

    return {"fake": False}
