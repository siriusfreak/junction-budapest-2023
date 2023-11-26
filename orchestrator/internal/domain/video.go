package domain

type VideoFakeCandidat struct {
	Format                      string
	UID                         string
	AudioFakeDetectionResult    *bool
	DeepfakeDetectResult        *bool
	LipsMovementDetectionResult *bool
	OpenClosedEyeDetect         *bool
	WhisperLargeV3Result        *bool
	OnePersonDetectResult       *bool
}

// if it is deepfake: bool = true, if it is not deepfake: bool = false
