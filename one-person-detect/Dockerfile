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

COPY requirements.txt requirements.txt

RUN pip3 install --no-cache-dir -r requirements.txt
RUN pip3 install python-multipart
RUN pip3 install uvicorn

COPY main.py .
COPY requirements.txt .
COPY ultralyticsplus /app/ultralyticsplus

ENV PATH="${PATH}:/path/to/uvicorn"

EXPOSE 8000

ENTRYPOINT ["uvicorn"]
CMD ["main:app", "--host", "0.0.0.0", "--port", "8000"]
