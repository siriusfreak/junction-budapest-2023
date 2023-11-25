FROM python:3.8-bookworm

ENV TZ=Etc/UTC
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get install -y ffmpeg gcc-11 make cmake unzip wget git

RUN python3.8 -m pip install --upgrade pip setuptools wheel

RUN sed -i -e's/ main/ main contrib non-free/g' /etc/apt/sources.list.d/debian.sources
RUN apt-get update
RUN apt-get install -y nvidia-cuda-toolkit
RUN apt-get install -y nvidia-cudnn
RUN /usr/sbin/update-nvidia-cudnn -u

RUN rm -rf /usr/bin/gcc /usr/bin/g++
RUN ln -s /usr/bin/gcc-11 /usr/bin/gcc
RUN ln -s /usr/bin/g++-11 /usr/bin/g++

RUN python3.8 -m pip install dlib==19.24.2 --config-settings="DLIB_USE_CUDA=1 --allow-unsupported-compiler" -vvv

COPY requirements.txt /requirements.txt
RUN python3.8 -m pip  install -r /requirements.txt

RUN mkdir /app
WORKDIR /app
COPY av_hubert av_hubert
RUN python3.8 -m pip install file:///app/av_hubert/fairseq

RUN wget https://huggingface.co/vumichien/AV-HuBERT/resolve/main/model.pt?download=true -O model.pt
RUN wget https://huggingface.co/spaces/vumichien/Lip_movement_reading/resolve/main/mmod_human_face_detector.dat?download=true -O mmod_human_face_detector.dat
RUN wget https://huggingface.co/spaces/vumichien/Lip_movement_reading/resolve/main/20words_mean_face.npy?download=true -O 20words_mean_face.npy
RUN wget https://huggingface.co/spaces/vumichien/Lip_movement_reading/resolve/main/shape_predictor_68_face_landmarks.dat?download=true -O shape_predictor_68_face_landmarks.dat

COPY app.py app.py

RUN rm -rf /usr/bin/nvidia-smi

ENTRYPOINT ["python3.8", "-m", "uvicorn", "app:app", "--host", "0.0.0.0"]
