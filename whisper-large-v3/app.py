import torch

from transformers import pipeline, AutoModel
from transformers.pipelines.audio_utils import ffmpeg_read
from audio_extract import extract_audio
from fastapi import FastAPI, File, UploadFile
import shutil
import subprocess
import tempfile

BATCH_SIZE = 8
FILE_LIMIT_MB = 1000

device = 0 if torch.cuda.is_available() else "cpu"

pipe = pipeline(
    task="automatic-speech-recognition",
    model="./model",
    chunk_length_s=30,
    device=device,
)

def transcribe(audio_path, task):
    text = pipe({'audio': audio_path}, batch_size=BATCH_SIZE, generate_kwargs={"task": task}, return_timestamps=True)["text"]
    return text


def file_transcribe(filepath, task):
    with tempfile.NamedTemporaryFile(suffix="audio.mp3") as temp_file:
        extract_audio(input_path=filepath, output_path="1"+temp_file.name)

        with open("1"+temp_file.name, "rb") as f:
            audio = f.read()


        inputs = ffmpeg_read(audio, pipe.feature_extractor.sampling_rate)
        inputs = {"array": inputs, "sampling_rate": pipe.feature_extractor.sampling_rate}

        text = pipe(inputs, batch_size=BATCH_SIZE, generate_kwargs={"task": task}, return_timestamps=True)["text"]

        return text


app = FastAPI()


@app.post("/predict")
async def predict(video: UploadFile = File(...)):
    filename = video.filename
    with tempfile.NamedTemporaryFile(suffix=filename) as temp_file:
        # Copy the contents of the uploaded file to the temporary file
        shutil.copyfileobj(video.file, temp_file)
        # Get the path of the temporary file
        temp_file_path = temp_file.name
        return file_transcribe(temp_file_path, "transcribe")