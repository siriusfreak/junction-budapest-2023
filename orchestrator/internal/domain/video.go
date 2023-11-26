package domain

type VideoFakeCandidat struct {
	Format                      string `json:"Format"`
	UID                         string `json:"UID"`
	AudioFakeDetectionResult    *bool `json:"AudioFakeDetectionResult"`
	DeepfakeDetectResult        *bool `json:"DeepfakeDetectResult"`
	LipsMovementDetectionResult *bool `json:"LipsMovementDetectionResult"`
	OpenClosedEyeDetect         *bool `json:"OpenClosedEyeDetect"`
	WhisperLargeV3Result        *bool `json:"WhisperLargeV3Result"`
	OnePersonDetectResult       *bool `json:"OnePersonDetectResult"`
}

// if it is deepfake: bool = true, if it is not deepfake: bool = false
