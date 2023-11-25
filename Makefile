deepfake:
	nvidia-docker build -t deepfake-detect:v1 deepfake-detect/.

eyes:
	nvidia-docker build -t open-closed-eye-detect:v1 open-closed-eye-detect/.


run:
	nvidia-docker compose up -d

down:
	nvidia-docker compose down

