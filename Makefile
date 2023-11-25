deepfake:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/deepfake-detect:v1 deepfake-detect/.

eyes:
	nvidia-docker build -t europe-west3-docker.pkg.dev/junction-budapest-2023/main/open-closed-eye-detect:v1 open-closed-eye-detect/.




run:
	nvidia-docker compose up -d

down:
	nvidia-docker compose down

