package usecase

import (
	"context"
	"log"
	"encoding/json"
	"orchestrator/internal/domain"
	"orchestrator/internal/interfaces"
)

type GetTaskStatusUseCase struct {
	taskStorage interfaces.TaskStorage
}

type TaskStatusResponse struct {
	CompletionPercentage float64                  `json:"completion_percentage"`
	ConfidenceLevel      bool         		      `json:"confidence_level"`
	VideoFakeCandidat 	 domain.VideoFakeCandidat `json:"VideoFakeCandidat"`
}

func NewGetTaskStatusUseCase(taskStorage interfaces.TaskStorage) *GetTaskStatusUseCase {
	return &GetTaskStatusUseCase{
		taskStorage: taskStorage,
	}
}

func (uc *GetTaskStatusUseCase) GetTaskStatus(ctx context.Context, uid string) (TaskStatusResponse, error) {
	video, err := uc.taskStorage.GetTask(ctx, uid)
	if err != nil {
		return TaskStatusResponse{}, err
	}

	return calculateCompletion(video), nil
}

func calculateCompletion(video *domain.VideoFakeCandidat) TaskStatusResponse {
	totalFields := 4
	filledFields := 0
	confidenceLevel := true

	if video.OnePersonDetectResult != nil {
		filledFields++
		confidenceLevel = confidenceLevel && *video.OnePersonDetectResult
	}

	if video.DeepfakeDetectResult != nil {
		filledFields++
		confidenceLevel = confidenceLevel && *video.DeepfakeDetectResult
	}

	//if video.LipsMovementDetectionResult != nil {
	//	filledFields++
	//	confidenceLevel = confidenceLevel && *video.LipsMovementDetectionResult
	//}

	if video.WhisperLargeV3Result != nil {
		filledFields++
		confidenceLevel = confidenceLevel && *video.WhisperLargeV3Result
	}

	if video.AudioFakeDetectionResult != nil {
		filledFields++
		confidenceLevel = confidenceLevel && *video.AudioFakeDetectionResult
	}

	//if video.OpenClosedEyeDetect != nil {
	//	filledFields++
	//	confidenceLevel = confidenceLevel && *video.OpenClosedEyeDetect
	//}
	var response TaskStatusResponse
	if totalFields == filledFields {
		response = TaskStatusResponse{
			CompletionPercentage: 100,
			ConfidenceLevel:      confidenceLevel,
			VideoFakeCandidat:    *video, 
		}
		logTaskStatusResponse(response)
	}else{
		response = TaskStatusResponse{
			CompletionPercentage: (float64(filledFields) / float64(totalFields)) * 100,
			ConfidenceLevel:      false,
			VideoFakeCandidat:    *video, 
		}
	}

	return response
}


func logTaskStatusResponse(response TaskStatusResponse) {
    responseJSON, err := json.MarshalIndent(response, "", "  ")
    if err != nil {
        log.Printf("Error marshaling response: %v", err)
        return
    }
    log.Printf("TaskStatusResponse: %s", responseJSON)
}