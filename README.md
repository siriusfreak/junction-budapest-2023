# FelineGuard: Whisker Wise Deepfake Detector
## Presentation

## GPU requirement
You need TeslaV100 or higher. Preferred GPU is L40
[Our presentation](https://docs.google.com/presentation/d/1BfM6qaNuU4hnDgIUuWwPxeUEIKIMI6CNc9MwmLWAYjs/edit#slide=id.ged86f964b3_0_0)
## Installation

1. Install docker following [official guide](https://docs.docker.com/engine/install/)

2. Configure Nvidia repository:
  ~~~shell
  curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
    && curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
      sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
      sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list \
    && \
      sudo apt-get update
  ~~~

3. Install the NVIDIA Container Toolkit packages:
  ~~~shell
  sudo apt-get install -y nvidia-container-toolkit
  ~~~

4. Configure the container runtime by using the nvidia-ctk command:
  ~~~shell
  sudo nvidia-ctk runtime configure --runtime=docker
  ~~~

## Compiling images
System relies on precimpiled images hosted in GCP registry. If you want to change something you can always run:
~~~shell
make all
~~~

## Running network cluster
To run networks type:
~~~shell
make run
~~~

If you want to stop them, type:
~~~shell
make down
~~~

## Run orchestation service
~~~shell
go run cmd/main.go
~~~


