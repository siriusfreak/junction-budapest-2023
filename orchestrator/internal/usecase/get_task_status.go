package usecase

import (
	"context"
	"orchestrator/internal/domain"
	"orchestrator/internal/interfaces"
)

type GetTaskStatusUseCase struct {
	taskStorage interfaces.TaskStorage
}

type TaskStatusResponse struct {
	CompletionPercentage float64 `json:"completion_percentage"`
	ConfidenceLevel      bool    `json:"confidence_level"`
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
	totalFields := 1
	filledFields := 0
	confidenceLevel := false

	if video.OnePersonDetectResult != nil {
		filledFields++
		confidenceLevel = confidenceLevel || *video.OnePersonDetectResult
	}

	if totalFields == filledFields{
		return TaskStatusResponse{
			CompletionPercentage: 100,
			ConfidenceLevel:      confidenceLevel,
		}
	}

	return TaskStatusResponse{
		CompletionPercentage: (float64(filledFields) / float64(totalFields)) * 100,
		ConfidenceLevel:      false,
	}
}
