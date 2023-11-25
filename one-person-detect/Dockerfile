FROM python:3.9-slim

WORKDIR /app

RUN apt-get update && apt-get install -y libgl1-mesa-glx ffmpeg libsm6 libxext6
RUN apt-get install -y --no-install-recommends gcc libc6-dev

ENV CONF_THRESHOLD=0.25
ENV IOU_THRESHOLD=0.45
ENV AGNOSTIC_NMS=False
ENV MAX_DET=1000
ENV MAX_VIDEO_LEN_SECOND=30
ENV MAX_FPS=200

RUN pip3 install poetry

COPY pyproject.toml poetry.lock* /app/

RUN poetry config virtualenvs.create false \
    && poetry install --no-interaction --no-ansi

COPY main.py .
COPY ultralyticsplus /app/ultralyticsplus

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
