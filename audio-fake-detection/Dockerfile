FROM python:3.11-slim-bookworm

ENV TZ=Etc/UTC
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get install -y ffmpeg gcc-11 make cmake unzip wget git

RUN mkdir -p /app/checkpoints

WORKDIR /app/checkpoints
RUN wget https://huggingface.co/spaces/Kimata/multimodal_deepfake_detection/resolve/main/checkpoints/model.pth?download=true -O model.pth

# Dependencies installing
RUN pip3 install poetry

COPY pyproject.toml poetry.lock* /app/

RUN poetry config virtualenvs.create false \
    && poetry install --no-interaction --no-ansi


WORKDIR /app
COPY models models
COPY app.py app.py

ENTRYPOINT ["python3.11", "-m", "uvicorn", "app:app", "--host", "0.0.0.0"]
