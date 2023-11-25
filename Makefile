deepfake:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/deepfake-detect:v1 deepfake-detect/.
	docker push europe-west3-docker.pkg.dev/junction-budapest-2023/main/deepfake-detect:v1

eyes:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/open-closed-eye-detect:v1 open-closed-eye-detect/.
	docker push europe-west3-docker.pkg.dev/junction-budapest-2023/main/open-closed-eye-detect:v1

audio:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/audio-fake-detection:v1 audio-fake-detection/.
	docker push europe-west3-docker.pkg.dev/junction-budapest-2023/main/audio-fake-detection:v1

lips:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/lips-movement-detection:v1 lips-movement-detection/.
	docker push europe-west3-docker.pkg.dev/junction-budapest-2023/main/lips-movement-detection:v1

whisper:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/whisper-large-v3:v1 whisper-large-v3/.
	docker push europe-west3-docker.pkg.dev/junction-budapest-2023/main/whisper-large-v3:v1

person:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/one-person-detect:v1 one-person-detect/.
	docker push europe-west3-docker.pkg.dev/junction-budapest-2023/main/one-person-detect:v1

all:
	make deepfake
	make eyes
	make audio
	make lips
	make whisper
	make person

run:
	nvidia-docker compose up -d

down:
	nvidia-docker compose down

