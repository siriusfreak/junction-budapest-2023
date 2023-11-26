import os
import sys

from fastapi import FastAPI, File, UploadFile
from fastapi.responses import PlainTextResponse
import shutil
import tempfile

sys.path.append('./av_hubert/avhubert')

sys.argv.append('dummy')


import dlib, cv2, os
import numpy as np
import skvideo
import skvideo.io
from tqdm import tqdm
from preparation.align_mouth import landmarks_interpolate, crop_patch, write_video_ffmpeg
from base64 import b64encode
import torch
import cv2
import tempfile
from argparse import Namespace
import fairseq
from fairseq import checkpoint_utils, options, tasks, utils
from fairseq.dataclass.configs import GenerationConfig
from huggingface_hub import hf_hub_download
import gradio as gr
from pytube import YouTube

# os.chdir('/home/user/app/av_hubert/avhubert')

user_dir = "./av_hubert/avhubert"
utils.import_user_module(Namespace(user_dir=user_dir))
data_dir = "./video"

ckpt_path = './model.pt'
face_detector_path = "./mmod_human_face_detector.dat"
face_predictor_path = "./shape_predictor_68_face_landmarks.dat"
mean_face_path = "./20words_mean_face.npy"
mouth_roi_path = "./roi.mp4"
modalities = ["video"]
gen_subset = "test"
gen_cfg = GenerationConfig(beam=20)
models, saved_cfg, task = checkpoint_utils.load_model_ensemble_and_task([ckpt_path])
models = [model.eval().cuda() if torch.cuda.is_available() else model.eval() for model in models]


def detect_landmark(image, detector, predictor):
    gray = cv2.cvtColor(image, cv2.COLOR_RGB2GRAY)
    face_locations  = detector(gray, 1)
    coords = None
    for (_, face_location) in enumerate(face_locations):
        if torch.cuda.is_available():
            rect = face_location.rect
        else:
            rect = face_location
        shape = predictor(gray, rect)
        coords = np.zeros((68, 2), dtype=np.int32)
        for i in range(0, 68):
            coords[i] = (shape.part(i).x, shape.part(i).y)
    return coords

def preprocess_video(input_video_path, out_file):
    if torch.cuda.is_available():
        print("Using CUDA")
        detector = dlib.cnn_face_detection_model_v1(face_detector_path)
    else:
        detector = dlib.get_frontal_face_detector()
    
    predictor = dlib.shape_predictor(face_predictor_path)
    STD_SIZE = (256, 256)
    mean_face_landmarks = np.load(mean_face_path)
    stablePntsIDs = [33, 36, 39, 42, 45]
    videogen = skvideo.io.vread(input_video_path)
    frames = np.array([frame for frame in videogen])
    landmarks = []
    for frame in tqdm(frames):
        landmark = detect_landmark(frame, detector, predictor)
        landmarks.append(landmark)
    preprocessed_landmarks = landmarks_interpolate(landmarks)
    rois = crop_patch(input_video_path, preprocessed_landmarks, mean_face_landmarks, stablePntsIDs, STD_SIZE, 
                          window_margin=12, start_idx=48, stop_idx=68, crop_height=96, crop_width=96)
    write_video_ffmpeg(rois, out_file, "ffmpeg")
    return out_file

def model_predict(process_video):
    with tempfile.TemporaryDirectory() as temp_dir:
        num_frames = int(cv2.VideoCapture(process_video).get(cv2.CAP_PROP_FRAME_COUNT))

        saved_cfg.task.modalities = modalities
        saved_cfg.task.data = temp_dir
        saved_cfg.task.label_dir = temp_dir
        task = tasks.setup_task(saved_cfg.task)
        generator = task.build_generator(models, gen_cfg)

        tsv_cont = ["/\n", f"test-0\t{process_video}\t{None}\t{num_frames}\t{int(16_000*num_frames/25)}\n"]
        label_cont = ["DUMMY\n"]
        with open(f"{temp_dir}/test.tsv", "w") as fo:
          fo.write("".join(tsv_cont))
        with open(f"{temp_dir}/test.wrd", "w") as fo:
          fo.write("".join(label_cont))
        task.load_dataset(gen_subset, task_cfg=saved_cfg.task)

        def decode_fn(x):
            dictionary = task.target_dictionary
            symbols_ignore = generator.symbols_to_strip_from_output
            symbols_ignore.add(dictionary.pad())
            return task.datasets[gen_subset].label_processors[0].decode(x, symbols_ignore)

        itr = task.get_batch_iterator(dataset=task.dataset(gen_subset)).next_epoch_itr(shuffle=False)
        sample = next(itr)
        if torch.cuda.is_available():
            sample = utils.move_to_cuda(sample)
        hypos = task.inference_step(generator, models, sample)
        ref = decode_fn(sample['target'][0].int().cpu())
        hypo = hypos[0][0]['tokens'].int().cpu()
        hypo = decode_fn(hypo)
        return hypo


app = FastAPI()


@app.post("/predict")
async def predict(video: UploadFile = File(...)):
    filename = video.filename
    print("DLIB_CUDA: ", dlib.DLIB_USE_CUDA, flush=True)
    with tempfile.NamedTemporaryFile(suffix=filename) as temp_file:
        # Copy the contents of the uploaded file to the temporary file
        shutil.copyfileobj(video.file, temp_file)
        # Get the path of the temporary file
        temp_file_path = temp_file.name

        with tempfile.NamedTemporaryFile(suffix=filename) as out_file:
            preprocess_video(temp_file_path, out_file.name)

            res = model_predict(out_file.name)

            return PlainTextResponse(content=str(res))
