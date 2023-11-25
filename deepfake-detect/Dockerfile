FROM nvidia/cuda:11.8.0-runtime-ubuntu20.04

RUN apt-get update && \
    apt-get install --no-install-recommends -y curl && \
    apt-get install -y libsm6 libxext6 libxrender-dev

ENV CONDA_AUTO_UPDATE_CONDA=false \
    PATH=/opt/miniconda/bin:$PATH

RUN curl -sLo ~/miniconda.sh https://repo.anaconda.com/miniconda/Miniconda3-py38_4.9.2-Linux-x86_64.sh \
    && chmod +x ~/miniconda.sh \
    && ~/miniconda.sh -b -p /opt/miniconda \
    && rm ~/miniconda.sh \
    && sed -i "$ a PATH=/opt/miniconda/bin:\$PATH" /etc/environment

RUN python3 -m pip --no-cache-dir install --upgrade pip && \
    python3 --version && \
    pip3 --version

RUN pip3 install poetry

WORKDIR /usr/src/app

COPY pyproject.toml poetry.lock* /usr/src/app/

RUN poetry config virtualenvs.create false \
    && poetry install --no-interaction --no-ansi

RUN apt-get update \
    && apt-get install -y libgl1-mesa-glx \
    && rm -rf /var/lib/apt/lists/*

COPY . /usr/src/app

EXPOSE 8000

# Define environment variable
ENV MODULE_NAME="app"
ENV VARIABLE_NAME="app"
ENV PORT=8000

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
