FROM python:3.11-slim-bookworm

ENV TZ=Etc/UTC
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get install -y ffmpeg gcc-11 make cmake unzip wget git

RUN mkdir -p /app/model

WORKDIR /app/model
RUN wget https://huggingface.co/openai/whisper-large-v3/resolve/main/model.safetensors?download=true -O model.safetensors
RUN wget https://huggingface.co/openai/whisper-large-v3/resolve/main/config.json?download=true -O config.json
RUN wget https://huggingface.co/openai/whisper-large-v3/resolve/main/tokenizer.json?download=true -O tokenizer.json
RUN wget https://huggingface.co/openai/whisper-large-v3/resolve/main/tokenizer_config.json?download=true -O tokenizer_config.json
RUN wget https://huggingface.co/openai/whisper-large-v3/resolve/main/preprocessor_config.json?download=true -O preprocessor_config.json
RUN wget https://huggingface.co/openai/whisper-large-v3/resolve/main/generation_config.json?download=true -O generation_config.json

# Poetry installing
RUN pip3 install poetry

WORKDIR /app

COPY pyproject.toml poetry.lock* /app/

# Dependencies installing
RUN poetry config virtualenvs.create false \
    && poetry install --no-interaction --no-ansi


WORKDIR /app
COPY app.py app.py

ENTRYPOINT ["python3.11", "-m", "uvicorn", "app:app", "--host", "0.0.0.0"]
