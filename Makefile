deepfake:
	nvidia-docker build -t deepfake-detect:v1 deepfake-detect/.

eyes:
	nvidia-docker build -t eyes-detect:v1 open_closed_eye_detection/.