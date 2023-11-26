import torch

from transformers import pipeline, AutoModel
from transformers.pipelines.audio_utils import ffmpeg_read

from fastapi import FastAPI, File, UploadFile
import shutil
import ffmpeg
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

app = FastAPI()

@app.post("/predict")
async def predict(video: UploadFile = File(...)):
    filename = video.filename
    with tempfile.NamedTemporaryFile(suffix=filename) as tmp:
        # Сохраняем файл
        shutil.copyfileobj(video.file, tmp)
        tmp_path = tmp.name
        
        # Конвертируем mp4 в wav используя ffmpeg
        process = (
            ffmpeg
            .input(tmp_path)
            .output('pipe:', format='wav')
            .run_async(pipe_stdout=True, pipe_stderr=True)
        )
        out, err = process.communicate()
        
        if process.returncode != 0:
            raise Exception(f"ffmpeg error: {err.decode()}")

        # Транскрибация аудио
        inputs = {"array": out, "sampling_rate": pipe.feature_extractor.sampling_rate}
        text = pipe(inputs, batch_size=BATCH_SIZE)["text"]
        
        return {"text": text}